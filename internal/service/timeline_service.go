package service

import (
	"context"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/database"
	"github.com/ciliverse/cilikube/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"gorm.io/gorm"
)

const (
	timelineScanSeconds      = 15
	timelineHeartbeatSeconds = 30
	timelineRetentionDays    = 7
	timelineMaxPods          = 2000
	timelineMaxEvents        = 500
)

type timelineCacheEntry struct {
	Status    string
	Reason    string
	SampledAt time.Time
}

// TimelineService samples resource health and builds the Timeline API payload.
type TimelineService struct {
	cm     *k8s.ClusterManager
	events *EventService
	db     *gorm.DB

	mu    sync.Mutex
	cache map[string]timelineCacheEntry // key: cluster|ns|kind|name

	lastOK   map[string]time.Time
	lastErr  map[string]string
	started  bool
	stopCh   chan struct{}
}

func NewTimelineService(cm *k8s.ClusterManager, events *EventService) *TimelineService {
	s := &TimelineService{
		cm:     cm,
		events: events,
		cache:  map[string]timelineCacheEntry{},
		lastOK: map[string]time.Time{},
		lastErr: map[string]string{},
		stopCh: make(chan struct{}),
	}
	if db, err := database.GetDB(); err == nil && db != nil {
		s.db = db
	}
	return s
}

func (s *TimelineService) Start() {
	if s == nil || s.db == nil || s.started {
		return
	}
	s.started = true
	go s.loop()
	log.Println("timeline sampler started (scan=15s heartbeat=30s retention=7d)")
}

func (s *TimelineService) Stop() {
	if s == nil || !s.started {
		return
	}
	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
}

func (s *TimelineService) loop() {
	// Initial delay so cluster clients finish warming.
	select {
	case <-time.After(8 * time.Second):
	case <-s.stopCh:
		return
	}
	s.scanAll()
	scan := time.NewTicker(timelineScanSeconds * time.Second)
	prune := time.NewTicker(time.Hour)
	defer scan.Stop()
	defer prune.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-scan.C:
			s.scanAll()
		case <-prune.C:
			s.prune()
		}
	}
}

func (s *TimelineService) prune() {
	if s.db == nil {
		return
	}
	cutoff := time.Now().UTC().Add(-time.Duration(timelineRetentionDays) * 24 * time.Hour)
	if err := s.db.Where("sampled_at < ?", cutoff).Delete(&store.TimelineStatusSample{}).Error; err != nil {
		log.Printf("timeline prune: %v", err)
	}
}

func cacheKey(clusterID, ns, kind, name string) string {
	return clusterID + "|" + ns + "|" + kind + "|" + name
}

type liveResource struct {
	Namespace string
	Kind      string
	Name      string
	UID       string
	AppGroup  string
	Status    string
	Reason    string
	Href      string
}

func (s *TimelineService) scanAll() {
	if s.cm == nil || s.db == nil {
		return
	}
	if k8s.IsShowcase() {
		s.seedShowcase()
		return
	}
	clients := s.cm.SnapshotClients()
	for id, client := range clients {
		if client == nil || client.Clientset == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		resources, err := listTimelineResources(ctx, client.Clientset, "")
		cancel()
		if err != nil {
			s.mu.Lock()
			s.lastErr[id] = err.Error()
			s.mu.Unlock()
			continue
		}
		now := time.Now().UTC()
		for _, r := range resources {
			s.maybeWrite(id, r, now)
		}
		s.mu.Lock()
		s.lastOK[id] = now
		delete(s.lastErr, id)
		s.mu.Unlock()
	}
}

func (s *TimelineService) maybeWrite(clusterID string, r liveResource, now time.Time) {
	key := cacheKey(clusterID, r.Namespace, r.Kind, r.Name)
	s.mu.Lock()
	prev, ok := s.cache[key]
	write := !ok || prev.Status != r.Status || prev.Reason != r.Reason ||
		now.Sub(prev.SampledAt) >= time.Duration(timelineHeartbeatSeconds)*time.Second
	if write {
		s.cache[key] = timelineCacheEntry{Status: r.Status, Reason: r.Reason, SampledAt: now}
	}
	s.mu.Unlock()
	if !write {
		return
	}
	row := store.TimelineStatusSample{
		ClusterID: clusterID,
		Namespace: r.Namespace,
		Kind:      r.Kind,
		Name:      r.Name,
		UID:       r.UID,
		AppGroup:  r.AppGroup,
		Status:    r.Status,
		Reason:    r.Reason,
		SampledAt: now,
	}
	if err := s.db.Create(&row).Error; err != nil {
		log.Printf("timeline sample write: %v", err)
	}
}

func (s *TimelineService) seedShowcase() {
	now := time.Now().UTC()
	// Seed every fleet cluster so switching demo/prod-east/staging-lab still has bars.
	for _, clusterID := range []string{
		k8s.ShowcaseClusterID,
		k8s.ShowcaseProdClusterID,
		k8s.ShowcaseStagingClusterID,
	} {
		s.ensureShowcaseTimelineHistory(clusterID, now)
		for _, r := range showcaseTimelineLive(clusterID, now) {
			s.maybeWrite(clusterID, r, now)
		}
		s.mu.Lock()
		s.lastOK[clusterID] = now
		s.mu.Unlock()
	}
}

// showcaseTimelineLive matches real showcase inventory names (not fictional demo/ns).
func showcaseTimelineLive(clusterID string, now time.Time) []liveResource {
	tick := now.Unix() / 90
	switch clusterID {
	case k8s.ShowcaseProdClusterID:
		orders := liveResource{
			Namespace: "production", Kind: "deployment", Name: "orders-api", AppGroup: "orders-api",
			Status: TLDegraded, Reason: "4/5 ready", Href: "/deployments/production/orders-api",
		}
		if tick%3 == 0 {
			orders.Status = TLRolling
			orders.Reason = "rolling update"
		}
		return []liveResource{
			orders,
			{Namespace: "production", Kind: "pod", Name: "orders-api-a1b2c-55555", AppGroup: "orders-api", Status: TLUnhealthy, Reason: "0/1 CrashLoopBackOff", Href: "/pods/production/orders-api-a1b2c-55555"},
			{Namespace: "production", Kind: "deployment", Name: "checkout", AppGroup: "checkout", Status: TLDegraded, Reason: "2/3 ready", Href: "/deployments/production/checkout"},
			{Namespace: "payments", Kind: "deployment", Name: "pay-gateway", AppGroup: "pay-gateway", Status: TLDegraded, Reason: "2/3 ready", Href: "/deployments/payments/pay-gateway"},
			{Namespace: "payments", Kind: "pod", Name: "pay-gateway-9x8y7-p3333", AppGroup: "pay-gateway", Status: TLUnhealthy, Reason: "0/1 CrashLoopBackOff", Href: "/pods/payments/pay-gateway-9x8y7-p3333"},
			{Namespace: "kube-system", Kind: "deployment", Name: "coredns", AppGroup: "coredns", Status: TLHealthy, Reason: "2/2 ready", Href: "/deployments/kube-system/coredns"},
		}
	case k8s.ShowcaseStagingClusterID:
		return []liveResource{
			{Namespace: "staging", Kind: "deployment", Name: "web-frontend", AppGroup: "web-frontend", Status: TLHealthy, Reason: "2/2 ready", Href: "/deployments/staging/web-frontend"},
			{Namespace: "staging", Kind: "deployment", Name: "api-gateway", AppGroup: "api-gateway", Status: TLHealthy, Reason: "1/1 ready", Href: "/deployments/staging/api-gateway"},
			{Namespace: "preview", Kind: "pod", Name: "feature-x-pr42-eee05", AppGroup: "feature-x", Status: TLDegraded, Reason: "Pending", Href: "/pods/preview/feature-x-pr42-eee05"},
			{Namespace: "kube-system", Kind: "deployment", Name: "coredns", AppGroup: "coredns", Status: TLHealthy, Reason: "1/1 ready", Href: "/deployments/kube-system/coredns"},
		}
	default:
		api := liveResource{
			Namespace: "default", Kind: "deployment", Name: "api-gateway", AppGroup: "api-gateway",
			Status: TLHealthy, Reason: "2/2 ready", Href: "/deployments/default/api-gateway",
		}
		if tick%2 == 0 {
			api.Status = TLDegraded
			api.Reason = "1/2 ready"
		} else if tick%5 == 0 {
			api.Status = TLRolling
			api.Reason = "updating"
		}
		return []liveResource{
			{Namespace: "default", Kind: "deployment", Name: "web-frontend", AppGroup: "web-frontend", Status: TLHealthy, Reason: "3/3 ready", Href: "/deployments/default/web-frontend"},
			api,
			{Namespace: "default", Kind: "pod", Name: "api-gateway-6c4d5-jkl78", AppGroup: "api-gateway", Status: TLDegraded, Reason: "Readiness probe failing", Href: "/pods/default/api-gateway-6c4d5-jkl78"},
			{Namespace: "default", Kind: "service", Name: "web-frontend", AppGroup: "web-frontend", Status: TLHealthy, Reason: "ClusterIP", Href: "/services/default/web-frontend"},
			{Namespace: "production", Kind: "deployment", Name: "orders-api", AppGroup: "orders-api", Status: TLHealthy, Reason: "4/4 ready", Href: "/deployments/production/orders-api"},
			{Namespace: "kube-system", Kind: "deployment", Name: "coredns", AppGroup: "coredns", Status: TLHealthy, Reason: "2/2 ready", Href: "/deployments/kube-system/coredns"},
			{Namespace: "cilibase", Kind: "deployment", Name: "cilikube-demo", AppGroup: "cilikube-demo", Status: TLHealthy, Reason: "2/2 ready", Href: "/deployments/cilibase/cilikube-demo"},
		}
	}
}

func (s *TimelineService) ensureShowcaseTimelineHistory(clusterID string, now time.Time) {
	if s.db == nil {
		return
	}
	var count int64
	if err := s.db.Model(&store.TimelineStatusSample{}).Where("cluster_id = ?", clusterID).Count(&count).Error; err != nil {
		return
	}
	if count >= 48 {
		return
	}
	// Backfill ~6h of status segments so Timeline is not empty on first open.
	step := 3 * time.Minute
	start := now.Add(-6 * time.Hour)
	batch := make([]store.TimelineStatusSample, 0, 256)
	for t := start; !t.After(now); t = t.Add(step) {
		for _, r := range showcaseTimelineAt(clusterID, t) {
			batch = append(batch, store.TimelineStatusSample{
				ClusterID: clusterID,
				Namespace: r.Namespace,
				Kind:      r.Kind,
				Name:      r.Name,
				UID:       r.UID,
				AppGroup:  r.AppGroup,
				Status:    r.Status,
				Reason:    r.Reason,
				SampledAt: t,
			})
		}
	}
	if len(batch) == 0 {
		return
	}
	if err := s.db.CreateInBatches(batch, 100).Error; err != nil {
		log.Printf("timeline showcase backfill: %v", err)
	}
}

// showcaseTimelineAt builds a status story arc for a point in time (demo bars + markers).
func showcaseTimelineAt(clusterID string, at time.Time) []liveResource {
	mins := int(at.Unix()/60) % 360 // 6h cycle
	base := showcaseTimelineLive(clusterID, at)
	out := make([]liveResource, 0, len(base))
	for _, r := range base {
		rr := r
		switch {
		case r.Name == "api-gateway" && r.Kind == "deployment" && clusterID == k8s.ShowcaseClusterID:
			switch {
			case mins >= 40 && mins < 70:
				rr.Status, rr.Reason = TLRolling, "rolling update"
			case mins >= 70 && mins < 110:
				rr.Status, rr.Reason = TLDegraded, "1/2 ready"
			case mins >= 200 && mins < 230:
				rr.Status, rr.Reason = TLUnhealthy, "0/2 available"
			default:
				rr.Status, rr.Reason = TLHealthy, "2/2 ready"
			}
		case r.Name == "orders-api" && r.Kind == "deployment" && clusterID == k8s.ShowcaseProdClusterID:
			switch {
			case mins >= 60 && mins < 100:
				rr.Status, rr.Reason = TLRolling, "rolling update"
			case mins >= 100 && mins < 160:
				rr.Status, rr.Reason = TLDegraded, "4/5 ready"
			default:
				rr.Status, rr.Reason = TLDegraded, "4/5 ready"
			}
		case strings.Contains(r.Name, "CrashLoop") || strings.HasSuffix(r.Name, "55555") || strings.HasSuffix(r.Name, "p3333"):
			rr.Status, rr.Reason = TLUnhealthy, "0/1 CrashLoopBackOff"
		case r.Name == "api-gateway-6c4d5-jkl78":
			if mins >= 70 && mins < 130 {
				rr.Status, rr.Reason = TLUnhealthy, "0/1 CrashLoopBackOff"
			} else if mins >= 130 && mins < 160 {
				rr.Status, rr.Reason = TLDegraded, "Readiness probe failing"
			} else {
				rr.Status, rr.Reason = TLHealthy, "1/1 Running"
			}
		}
		out = append(out, rr)
	}
	return out
}

func listTimelineResources(ctx context.Context, cs kubernetes.Interface, namespace string) ([]liveResource, error) {
	ns := namespace
	out := make([]liveResource, 0, 256)

	deps, err := cs.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range deps.Items {
		d := &deps.Items[i]
		st, reason := ClassifyDeploymentStatus(d)
		out = append(out, liveResource{
			Namespace: d.Namespace, Kind: "deployment", Name: d.Name, UID: string(d.UID),
			AppGroup: TimelineAppGroup(d.Labels), Status: st, Reason: reason,
			Href: TimelineHref("deployment", d.Namespace, d.Name),
		})
	}

	sts, err := cs.AppsV1().StatefulSets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range sts.Items {
		d := &sts.Items[i]
		st, reason := ClassifyStatefulSetStatus(d)
		out = append(out, liveResource{
			Namespace: d.Namespace, Kind: "statefulset", Name: d.Name, UID: string(d.UID),
			AppGroup: TimelineAppGroup(d.Labels), Status: st, Reason: reason,
			Href: TimelineHref("statefulset", d.Namespace, d.Name),
		})
	}

	dss, err := cs.AppsV1().DaemonSets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range dss.Items {
		d := &dss.Items[i]
		st, reason := ClassifyDaemonSetStatus(d)
		out = append(out, liveResource{
			Namespace: d.Namespace, Kind: "daemonset", Name: d.Name, UID: string(d.UID),
			AppGroup: TimelineAppGroup(d.Labels), Status: st, Reason: reason,
			Href: TimelineHref("daemonset", d.Namespace, d.Name),
		})
	}

	jobs, err := cs.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range jobs.Items {
		j := &jobs.Items[i]
		st, reason := ClassifyJobStatus(j)
		out = append(out, liveResource{
			Namespace: j.Namespace, Kind: "job", Name: j.Name, UID: string(j.UID),
			AppGroup: TimelineAppGroup(j.Labels), Status: st, Reason: reason,
			Href: TimelineHref("job", j.Namespace, j.Name),
		})
	}

	cjs, err := cs.BatchV1().CronJobs(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range cjs.Items {
		c := &cjs.Items[i]
		st, reason := ClassifyCronJobStatus(c)
		out = append(out, liveResource{
			Namespace: c.Namespace, Kind: "cronjob", Name: c.Name, UID: string(c.UID),
			AppGroup: TimelineAppGroup(c.Labels), Status: st, Reason: reason,
			Href: TimelineHref("cronjob", c.Namespace, c.Name),
		})
	}

	svcs, err := cs.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for i := range svcs.Items {
		svc := &svcs.Items[i]
		st, reason := ClassifyServiceStatus(svc)
		out = append(out, liveResource{
			Namespace: svc.Namespace, Kind: "service", Name: svc.Name, UID: string(svc.UID),
			AppGroup: TimelineAppGroup(svc.Labels), Status: st, Reason: reason,
			Href: TimelineHref("service", svc.Namespace, svc.Name),
		})
	}

	pods, err := cs.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{Limit: timelineMaxPods})
	if err != nil {
		return nil, err
	}
	for i := range pods.Items {
		p := &pods.Items[i]
		st, reason := ClassifyPodStatus(p)
		out = append(out, liveResource{
			Namespace: p.Namespace, Kind: "pod", Name: p.Name, UID: string(p.UID),
			AppGroup: TimelineAppGroup(p.Labels), Status: st, Reason: reason,
			Href: TimelineHref("pod", p.Namespace, p.Name),
		})
	}

	return out, nil
}

// --- API types ---

type TimelineSegment struct {
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Status string    `json:"status"`
	Reason string    `json:"reason,omitempty"`
}

type TimelineEventMarker struct {
	At      time.Time `json:"at"`
	Type    string    `json:"type"`
	Reason  string    `json:"reason"`
	Message string    `json:"message"`
	Marker  string    `json:"marker"`
}

type TimelineRow struct {
	Kind        string                 `json:"kind"`
	Namespace   string                 `json:"namespace"`
	Name        string                 `json:"name"`
	UID         string                 `json:"uid,omitempty"`
	AppGroup    string                 `json:"appGroup"`
	Href        string                 `json:"href"`
	Segments    []TimelineSegment      `json:"segments"`
	Events      []TimelineEventMarker  `json:"events"`
	EventBadges map[string]int         `json:"eventBadges"`
}

type TimelineGroup struct {
	Name string        `json:"name"`
	Rows []TimelineRow `json:"rows"`
}

type TimelineSamplingInfo struct {
	Enabled           bool `json:"enabled"`
	HeartbeatSeconds  int  `json:"heartbeatSeconds"`
	RetentionDays     int  `json:"retentionDays"`
	Provisional       bool `json:"provisional"`
	ScanSeconds       int  `json:"scanSeconds"`
}

type TimelineResponse struct {
	ClusterID string               `json:"clusterId"`
	From      time.Time            `json:"from"`
	To        time.Time            `json:"to"`
	GroupBy   string               `json:"groupBy"`
	Groups    []TimelineGroup      `json:"groups"`
	Sampling  TimelineSamplingInfo `json:"sampling"`
}

type TimelineMetaResponse struct {
	Enabled          bool              `json:"enabled"`
	ScanSeconds      int               `json:"scanSeconds"`
	HeartbeatSeconds int               `json:"heartbeatSeconds"`
	RetentionDays    int               `json:"retentionDays"`
	SampleCount      int64             `json:"sampleCount"`
	LastOK           map[string]string `json:"lastOk"`
	LastError        map[string]string `json:"lastError"`
}

type TimelineQuery struct {
	ClusterID string
	Namespace string
	From      time.Time
	To        time.Time
	GroupBy   string
	Q         string
	Kinds     map[string]bool
}

func (s *TimelineService) Meta(clusterID string) *TimelineMetaResponse {
	out := &TimelineMetaResponse{
		Enabled:          s != nil && s.db != nil,
		ScanSeconds:      timelineScanSeconds,
		HeartbeatSeconds: timelineHeartbeatSeconds,
		RetentionDays:    timelineRetentionDays,
		LastOK:           map[string]string{},
		LastError:        map[string]string{},
	}
	if s == nil {
		return out
	}
	s.mu.Lock()
	for k, v := range s.lastOK {
		out.LastOK[k] = v.UTC().Format(time.RFC3339)
	}
	for k, v := range s.lastErr {
		out.LastError[k] = v
	}
	s.mu.Unlock()
	if s.db != nil {
		q := s.db.Model(&store.TimelineStatusSample{})
		if clusterID != "" {
			q = q.Where("cluster_id = ?", clusterID)
		}
		_ = q.Count(&out.SampleCount).Error
	}
	return out
}

func (s *TimelineService) Build(ctx context.Context, cs kubernetes.Interface, q TimelineQuery) (*TimelineResponse, error) {
	if q.To.IsZero() {
		q.To = time.Now().UTC()
	}
	if q.From.IsZero() || !q.From.Before(q.To) {
		q.From = q.To.Add(-15 * time.Minute)
	}
	if q.GroupBy != "none" {
		q.GroupBy = "app"
	}

	live, err := listTimelineResources(ctx, cs, q.Namespace)
	if err != nil {
		return nil, err
	}
	live = filterLive(live, q)

	samplesByKey := map[string][]store.TimelineStatusSample{}
	provisional := true
	if s != nil && s.db != nil && q.ClusterID != "" {
		var rows []store.TimelineStatusSample
		dbq := s.db.Where("cluster_id = ? AND sampled_at >= ? AND sampled_at <= ?", q.ClusterID, q.From.Add(-timelineHeartbeatSeconds*time.Second), q.To)
		if q.Namespace != "" {
			dbq = dbq.Where("namespace = ?", q.Namespace)
		}
		if err := dbq.Order("sampled_at asc").Find(&rows).Error; err == nil && len(rows) > 0 {
			provisional = false
			for _, row := range rows {
				k := cacheKey(q.ClusterID, row.Namespace, row.Kind, row.Name)
				samplesByKey[k] = append(samplesByKey[k], row)
			}
		}
	}

	events := s.loadEvents(ctx, cs, q)

	rowsOut := make([]TimelineRow, 0, len(live))
	for _, r := range live {
		key := cacheKey(q.ClusterID, r.Namespace, r.Kind, r.Name)
		segs := buildSegments(samplesByKey[key], r, q.From, q.To)
		evs, badges := attachEvents(events, r.Namespace, r.Kind, r.Name, q.From, q.To)
		rowsOut = append(rowsOut, TimelineRow{
			Kind: r.Kind, Namespace: r.Namespace, Name: r.Name, UID: r.UID,
			AppGroup: r.AppGroup, Href: r.Href, Segments: segs, Events: evs, EventBadges: badges,
		})
	}

	sort.Slice(rowsOut, func(i, j int) bool {
		if rowsOut[i].AppGroup != rowsOut[j].AppGroup {
			return rowsOut[i].AppGroup < rowsOut[j].AppGroup
		}
		if rowsOut[i].Kind != rowsOut[j].Kind {
			return rowsOut[i].Kind < rowsOut[j].Kind
		}
		return rowsOut[i].Name < rowsOut[j].Name
	})

	groups := groupRows(rowsOut, q.GroupBy)
	return &TimelineResponse{
		ClusterID: q.ClusterID,
		From:      q.From,
		To:        q.To,
		GroupBy:   q.GroupBy,
		Groups:    groups,
		Sampling: TimelineSamplingInfo{
			Enabled:          s != nil && s.db != nil,
			HeartbeatSeconds: timelineHeartbeatSeconds,
			RetentionDays:    timelineRetentionDays,
			Provisional:      provisional,
			ScanSeconds:      timelineScanSeconds,
		},
	}, nil
}

func filterLive(live []liveResource, q TimelineQuery) []liveResource {
	qq := strings.ToLower(strings.TrimSpace(q.Q))
	out := live[:0]
	for _, r := range live {
		if len(q.Kinds) > 0 && !q.Kinds[r.Kind] {
			continue
		}
		if qq != "" {
			hay := strings.ToLower(r.Name + " " + r.AppGroup + " " + r.Namespace + " " + r.Kind)
			if !strings.Contains(hay, qq) {
				continue
			}
		}
		out = append(out, r)
	}
	return out
}

func buildSegments(samples []store.TimelineStatusSample, live liveResource, from, to time.Time) []TimelineSegment {
	if len(samples) == 0 {
		return []TimelineSegment{{From: from, To: to, Status: live.Status, Reason: live.Reason}}
	}
	status, reason := live.Status, live.Reason
	for _, s := range samples {
		if !s.SampledAt.After(from) {
			status, reason = s.Status, s.Reason
		}
	}
	type point struct {
		at             time.Time
		status, reason string
	}
	points := []point{{at: from, status: status, reason: reason}}
	for _, s := range samples {
		if s.SampledAt.Before(from) || s.SampledAt.After(to) {
			continue
		}
		last := points[len(points)-1]
		if s.Status == last.status && s.Reason == last.reason {
			continue
		}
		points = append(points, point{at: s.SampledAt, status: s.Status, reason: s.Reason})
	}
	segs := make([]TimelineSegment, 0, len(points))
	for i, p := range points {
		end := to
		if i+1 < len(points) {
			end = points[i+1].at
		}
		if !end.After(p.at) {
			continue
		}
		segs = append(segs, TimelineSegment{From: p.at, To: end, Status: p.status, Reason: p.reason})
	}
	if len(segs) == 0 {
		return []TimelineSegment{{From: from, To: to, Status: live.Status, Reason: live.Reason}}
	}
	return segs
}

func (s *TimelineService) loadEvents(ctx context.Context, cs kubernetes.Interface, q TimelineQuery) []models.ClusterEvent {
	if cs == nil {
		return nil
	}
	list, err := cs.CoreV1().Events(q.Namespace).List(ctx, metav1.ListOptions{Limit: timelineMaxEvents})
	if err != nil {
		// Fallback to EventService (active cluster) if direct list fails.
		if s != nil && s.events != nil {
			resp, e2 := s.events.ListEvents(models.EventListRequest{
				Namespace: q.Namespace, Limit: timelineMaxEvents, Since: q.From.Format(time.RFC3339),
			})
			if e2 == nil && resp != nil {
				return resp.Events
			}
		}
		return nil
	}
	out := make([]models.ClusterEvent, 0, len(list.Items))
	for i := range list.Items {
		ce := models.ConvertK8sEventToClusterEvent(&list.Items[i])
		at := ce.LastTime
		if at.IsZero() {
			at = ce.CreatedAt
		}
		if at.Before(q.From) || at.After(q.To) {
			continue
		}
		out = append(out, ce)
	}
	return out
}

func attachEvents(events []models.ClusterEvent, ns, kind, name string, from, to time.Time) ([]TimelineEventMarker, map[string]int) {
	badges := map[string]int{}
	out := []TimelineEventMarker{}
	kindL := strings.ToLower(kind)
	for _, e := range events {
		if e.Namespace != ns {
			continue
		}
		if !strings.EqualFold(e.ObjectKind, kind) && strings.ToLower(e.ObjectKind) != kindL {
			continue
		}
		if e.Object != name {
			continue
		}
		at := e.LastTime
		if at.IsZero() {
			at = e.CreatedAt
		}
		if at.Before(from) || at.After(to) {
			continue
		}
		m := MapEventMarker(e.Type, e.Reason)
		badges[m]++
		out = append(out, TimelineEventMarker{
			At: at, Type: e.Type, Reason: e.Reason, Message: e.Message, Marker: m,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out, badges
}

func groupRows(rows []TimelineRow, groupBy string) []TimelineGroup {
	if groupBy == "none" {
		return []TimelineGroup{{Name: "all", Rows: rows}}
	}
	order := []string{}
	m := map[string][]TimelineRow{}
	for _, r := range rows {
		g := r.AppGroup
		if g == "" {
			g = "_ungrouped"
		}
		if _, ok := m[g]; !ok {
			order = append(order, g)
		}
		m[g] = append(m[g], r)
	}
	sort.Strings(order)
	out := make([]TimelineGroup, 0, len(order))
	for _, g := range order {
		out = append(out, TimelineGroup{Name: g, Rows: m[g]})
	}
	return out
}
