package service

import (
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/ciliverse/cilikube/configs"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/ciliverse/cilikube/pkg/k8s"
	ldap "github.com/go-ldap/ldap/v3"
)

func (s *AuthService) ldapConfig() *configs.LDAPConfig {
	if s.config == nil || !s.config.LDAP.Enabled || strings.TrimSpace(s.config.LDAP.URL) == "" {
		return nil
	}
	return &s.config.LDAP
}

// TestLDAP checks the directory bind and an optional user search. It does not create a local account.
func (s *AuthService) TestLDAP(username string) (string, error) {
	cfg := s.ldapConfig()
	if cfg == nil {
		if k8s.IsShowcase() {
			if strings.TrimSpace(username) == "" || username == k8s.ShowcaseDirectoryUser {
				return "simulated directory: bind ok, 1 entries", nil
			}
			return "simulated directory: bind ok, 0 entries", nil
		}
		return "", errors.New("ldap is disabled")
	}
	conn, err := s.ldapConn()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if cfg.BindDN != "" {
		if err := conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			return "", fmt.Errorf("bind failed: %w", err)
		}
	}
	filter := "(objectClass=*)"
	if strings.TrimSpace(username) != "" {
		pattern := cfg.UserFilter
		if pattern == "" {
			pattern = "(uid=%s)"
		}
		filter = fmt.Sprintf(pattern, ldap.EscapeFilter(username))
	}
	req := ldap.NewSearchRequest(
		cfg.UserBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 5, 0, false,
		filter, []string{"dn"}, nil,
	)
	res, err := conn.Search(req)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}
	return fmt.Sprintf("bind ok, %d entries", len(res.Entries)), nil
}

func (s *AuthService) ldapPasswordOK(user *store.User, password string) bool {
	cfg := s.ldapConfig()
	if cfg == nil || strings.TrimSpace(user.LDAPDN) == "" {
		return false
	}
	return s.ldapBind(user.LDAPDN, password) == nil
}

func (s *AuthService) provisionLDAPUser(username, password string) (*store.User, error) {
	cfg := s.ldapConfig()
	if cfg == nil {
		return nil, errors.New("ldap disabled")
	}
	entry, err := s.ldapSearch(username, password)
	if err != nil {
		return nil, err
	}
	email := entry.email
	if email == "" {
		email = username + "@ldap.local"
	}
	if existing, err := s.store.GetUserByEmail(email); err == nil && existing != nil {
		return nil, errors.New("email already exists")
	}
	raw := make([]byte, 24)
	_, _ = rand.Read(raw)
	user := &store.User{
		Username:    username,
		Email:       email,
		DisplayName: entry.display,
		IsActive:    true,
		LDAPDN:      entry.dn,
	}
	if err := user.HashPassword(fmt.Sprintf("%x", raw)); err != nil {
		return nil, err
	}
	if err := s.store.CreateUser(user); err != nil {
		return nil, err
	}
	role, err := s.store.GetRoleByName("viewer")
	if err != nil {
		return nil, err
	}
	if err := s.store.AssignRole(user.ID, role.ID); err != nil {
		return nil, err
	}
	if s.syncRoles != nil {
		_ = s.syncRoles(user.ID)
	}
	return user, nil
}

type ldapEntry struct {
	dn      string
	email   string
	display string
}

func (s *AuthService) ldapSearch(username, password string) (*ldapEntry, error) {
	cfg := s.ldapConfig()
	conn, err := s.ldapConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if cfg.BindDN != "" {
		if err := conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			return nil, err
		}
	}
	filter := cfg.UserFilter
	if filter == "" {
		filter = "(uid=%s)"
	}
	if !strings.Contains(filter, "%s") {
		return nil, errors.New("ldap user_filter must contain one %s")
	}
	req := ldap.NewSearchRequest(
		cfg.UserBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 2, 0, false,
		fmt.Sprintf(filter, ldap.EscapeFilter(username)),
		[]string{"dn", firstNonEmpty(cfg.EmailAttr, "mail"), firstNonEmpty(cfg.DisplayAttr, "cn")},
		nil,
	)
	res, err := conn.Search(req)
	if err != nil || len(res.Entries) != 1 {
		return nil, errors.New("ldap user not found")
	}
	entry := res.Entries[0]
	if err := s.ldapBind(entry.DN, password); err != nil {
		return nil, err
	}
	return &ldapEntry{
		dn:      entry.DN,
		email:   entry.GetAttributeValue(firstNonEmpty(cfg.EmailAttr, "mail")),
		display: entry.GetAttributeValue(firstNonEmpty(cfg.DisplayAttr, "cn")),
	}, nil
}

func (s *AuthService) ldapBind(dn, password string) error {
	conn, err := s.ldapConn()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Bind(dn, password)
}

func (s *AuthService) ldapConn() (*ldap.Conn, error) {
	cfg := s.ldapConfig()
	conn, err := ldap.DialURL(cfg.URL)
	if err != nil {
		return nil, err
	}
	if cfg.StartTLS {
		if err := conn.StartTLS(&tls.Config{InsecureSkipVerify: cfg.InsecureSkip}); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
