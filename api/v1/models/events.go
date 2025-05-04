package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Event 表示Kubernetes中的事件
// 事件通常表示集群中某处发生的状态变化
type Event struct {
	Metadata       metav1.ObjectMeta `json:"metadata"`
	Type           string            `json:"type"`                     // 事件类型：Normal 或 Warning
	Reason         string            `json:"reason"`                   // 事件原因，人类可读的简短描述
	Message        string            `json:"message"`                  // 事件详细信息
	Source         EventSource       `json:"source"`                   // 事件来源
	InvolvedObject ObjectReference   `json:"involvedObject"`           // 事件相关对象
	Count          int32             `json:"count,omitempty"`          // 事件发生次数
	FirstTimestamp string            `json:"firstTimestamp,omitempty"` // 首次发生时间
	LastTimestamp  string            `json:"lastTimestamp,omitempty"`  // 最后发生时间
	EventTime      string            `json:"eventTime,omitempty"`      // 事件时间（微秒级精度）
}

// EventSource 表示事件的来源
type EventSource struct {
	Component string `json:"component,omitempty"` // 产生事件的组件
	Host      string `json:"host,omitempty"`      // 产生事件的节点名称
}

// ObjectReference 包含对另一个Kubernetes对象的引用
type ObjectReference struct {
	Kind            string `json:"kind,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	Name            string `json:"name,omitempty"`
	UID             string `json:"uid,omitempty"`
	APIVersion      string `json:"apiVersion,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	FieldPath       string `json:"fieldPath,omitempty"`
}

// EventList 表示事件列表
type EventList struct {
	Items []Event `json:"items"`
	Total int     `json:"total"`
}

// K8sEventToEvent 将Kubernetes Event转换为应用模型
func K8sEventToEvent(k8sEvent *corev1.Event) Event {
	event := Event{
		Metadata: k8sEvent.ObjectMeta,
		Type:     k8sEvent.Type,
		Reason:   k8sEvent.Reason,
		Message:  k8sEvent.Message,
		Count:    k8sEvent.Count,
		InvolvedObject: ObjectReference{
			Kind:            k8sEvent.InvolvedObject.Kind,
			Namespace:       k8sEvent.InvolvedObject.Namespace,
			Name:            k8sEvent.InvolvedObject.Name,
			UID:             string(k8sEvent.InvolvedObject.UID),
			APIVersion:      k8sEvent.InvolvedObject.APIVersion,
			ResourceVersion: k8sEvent.InvolvedObject.ResourceVersion,
			FieldPath:       k8sEvent.InvolvedObject.FieldPath,
		},
		Source: EventSource{
			Component: k8sEvent.Source.Component,
			Host:      k8sEvent.Source.Host,
		},
	}

	// 处理时间戳
	if !k8sEvent.FirstTimestamp.IsZero() {
		event.FirstTimestamp = k8sEvent.FirstTimestamp.Format(metav1.RFC3339Micro)
	}
	if !k8sEvent.LastTimestamp.IsZero() {
		event.LastTimestamp = k8sEvent.LastTimestamp.Format(metav1.RFC3339Micro)
	}
	if !k8sEvent.EventTime.IsZero() {
		event.EventTime = k8sEvent.EventTime.Format(metav1.RFC3339Micro)
	}

	return event
}
