package store

import (
	"bytes"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (s *MemoryStore) SaveWebAuthnCredential(c *WebAuthnCredential) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	c.ID = uint(len(s.webauthnCreds[c.UserID]) + 1)
	copied := *c
	s.webauthnCreds[c.UserID] = append(s.webauthnCreds[c.UserID], &copied)
	return nil
}

func (s *MemoryStore) ListWebAuthnCredentials(userID uint) ([]WebAuthnCredential, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	src := s.webauthnCreds[userID]
	out := make([]WebAuthnCredential, 0, len(src))
	for _, c := range src {
		out = append(out, *c)
	}
	return out, nil
}

func (s *MemoryStore) FindWebAuthnByCredentialID(credentialID []byte) (*WebAuthnCredential, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, list := range s.webauthnCreds {
		for _, c := range list {
			if bytes.Equal(c.CredentialID, credentialID) {
				copied := *c
				return &copied, nil
			}
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *MemoryStore) DeleteWebAuthnCredential(userID uint, credentialID []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	list := s.webauthnCreds[userID]
	next := list[:0]
	found := false
	for _, c := range list {
		if bytes.Equal(c.CredentialID, credentialID) {
			found = true
			continue
		}
		next = append(next, c)
	}
	if !found {
		return errors.New("passkey not found")
	}
	s.webauthnCreds[userID] = next
	return nil
}

func (s *DatabaseStore) SaveWebAuthnCredential(c *WebAuthnCredential) error {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
	}
	return s.db.Create(c).Error
}

func (s *DatabaseStore) ListWebAuthnCredentials(userID uint) ([]WebAuthnCredential, error) {
	var list []WebAuthnCredential
	err := s.db.Where("user_id = ?", userID).Order("id ASC").Find(&list).Error
	return list, err
}

func (s *DatabaseStore) FindWebAuthnByCredentialID(credentialID []byte) (*WebAuthnCredential, error) {
	var c WebAuthnCredential
	err := s.db.Where("credential_id = ?", credentialID).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *DatabaseStore) DeleteWebAuthnCredential(userID uint, credentialID []byte) error {
	res := s.db.Where("user_id = ? AND credential_id = ?", userID, credentialID).Delete(&WebAuthnCredential{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("passkey not found")
	}
	return nil
}
