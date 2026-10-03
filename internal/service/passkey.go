package service

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type passkeyUser struct {
	id          []byte
	name        string
	display     string
	credentials []webauthn.Credential
}

func (u passkeyUser) WebAuthnID() []byte                         { return u.id }
func (u passkeyUser) WebAuthnName() string                       { return u.name }
func (u passkeyUser) WebAuthnDisplayName() string                { return u.display }
func (u passkeyUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

type passkeySessions struct {
	sync.Mutex
	items map[string]webauthn.SessionData
}

func (p *passkeySessions) put(data webauthn.SessionData) string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	id := hex.EncodeToString(buf)
	p.Lock()
	if p.items == nil {
		p.items = map[string]webauthn.SessionData{}
	}
	p.items[id] = data
	p.Unlock()
	return id
}

func (p *passkeySessions) take(id string) (webauthn.SessionData, bool) {
	p.Lock()
	defer p.Unlock()
	data, ok := p.items[id]
	delete(p.items, id)
	return data, ok
}

var passkeySessionStore passkeySessions

func (s *AuthService) webAuthn(origin string) (*webauthn.WebAuthn, error) {
	rpid := "localhost"
	origins := []string{"http://localhost:8888", "http://127.0.0.1:8888"}
	if s.config != nil {
		if s.config.WebAuthn.RPID != "" {
			rpid = s.config.WebAuthn.RPID
		}
		if len(s.config.WebAuthn.Origins) > 0 {
			origins = s.config.WebAuthn.Origins
		}
	}
	if origin != "" {
		if parsed, err := url.Parse(origin); err == nil && parsed.Hostname() != "" && s.config != nil && s.config.WebAuthn.RPID == "" {
			rpid = parsed.Hostname()
		}
		found := false
		for _, item := range origins {
			if item == origin {
				found = true
			}
		}
		if !found {
			origins = append(origins, origin)
		}
	}
	return webauthn.New(&webauthn.Config{
		RPID:          rpid,
		RPDisplayName: "CiliKube",
		RPOrigins:     origins,
	})
}

func (s *AuthService) passkeyUser(user *store.User) (passkeyUser, error) {
	if user.WebAuthnID == "" {
		buf := make([]byte, 32)
		if _, err := rand.Read(buf); err != nil {
			return passkeyUser{}, err
		}
		user.WebAuthnID = hex.EncodeToString(buf)
		if err := s.store.UpdateUser(user); err != nil {
			return passkeyUser{}, err
		}
	}
	id, err := hex.DecodeString(user.WebAuthnID)
	if err != nil {
		return passkeyUser{}, err
	}
	rows, err := s.store.ListWebAuthnCredentials(user.ID)
	if err != nil {
		return passkeyUser{}, err
	}
	creds := make([]webauthn.Credential, 0, len(rows))
	for _, row := range rows {
		var cred webauthn.Credential
		if err := json.Unmarshal([]byte(row.CredentialJSON), &cred); err != nil {
			continue
		}
		creds = append(creds, cred)
	}
	display := user.DisplayName
	if display == "" {
		display = user.Username
	}
	return passkeyUser{id: id, name: user.Username, display: display, credentials: creds}, nil
}

// BeginPasskeyRegister starts enrollment for the signed-in user.
func (s *AuthService) BeginPasskeyRegister(userID uint, origin string) (string, *protocol.CredentialCreation, error) {
	wa, err := s.webAuthn(origin)
	if err != nil {
		return "", nil, err
	}
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return "", nil, errors.New("user not found")
	}
	account, err := s.passkeyUser(user)
	if err != nil {
		return "", nil, err
	}
	creation, session, err := wa.BeginRegistration(account)
	if err != nil {
		return "", nil, err
	}
	return passkeySessionStore.put(*session), creation, nil
}

// FinishPasskeyRegister stores a credential created by the browser.
func (s *AuthService) FinishPasskeyRegister(userID uint, sessionID string, req *http.Request) error {
	session, ok := passkeySessionStore.take(sessionID)
	if !ok {
		return errors.New("passkey registration expired")
	}
	wa, err := s.webAuthn(req.Header.Get("Origin"))
	if err != nil {
		return err
	}
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	account, err := s.passkeyUser(user)
	if err != nil {
		return err
	}
	cred, err := wa.FinishRegistration(account, session, req)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	if err := s.store.SaveWebAuthnCredential(&store.WebAuthnCredential{
		UserID:         userID,
		CredentialID:   cred.ID,
		CredentialJSON: string(raw),
	}); err != nil {
		return err
	}
	s.createAuditLog(&userID, "passkey_register", "user", fmt.Sprintf("%d", userID), "", req.UserAgent(), "passkey registered")
	return nil
}

// BeginPasskeyLogin starts a discoverable passkey assertion.
func (s *AuthService) BeginPasskeyLogin(origin string) (string, *protocol.CredentialAssertion, error) {
	wa, err := s.webAuthn(origin)
	if err != nil {
		return "", nil, err
	}
	assertion, session, err := wa.BeginDiscoverableLogin()
	if err != nil {
		return "", nil, err
	}
	return passkeySessionStore.put(*session), assertion, nil
}

// FinishPasskeyLogin verifies an assertion and issues the normal session JWT.
func (s *AuthService) FinishPasskeyLogin(sessionID string, req *http.Request, ip, ua string) (*models.LoginResponse, error) {
	session, ok := passkeySessionStore.take(sessionID)
	if !ok {
		return nil, errors.New("passkey login expired")
	}
	wa, err := s.webAuthn(req.Header.Get("Origin"))
	if err != nil {
		return nil, err
	}
	var matched *store.User
	_, err = wa.FinishDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		row, findErr := s.store.FindWebAuthnByCredentialID(rawID)
		if findErr != nil {
			return nil, findErr
		}
		user, userErr := s.store.GetUserByID(row.UserID)
		if userErr != nil {
			return nil, userErr
		}
		account, accountErr := s.passkeyUser(user)
		if accountErr != nil {
			return nil, accountErr
		}
		if !bytes.Equal(account.id, userHandle) {
			return nil, errors.New("passkey user handle mismatch")
		}
		matched = user
		return account, nil
	}, session, req)
	if err != nil {
		return nil, err
	}
	if matched == nil || !matched.IsActive {
		return nil, errors.New("passkey user is not active")
	}
	return s.finishLogin(matched, ip, ua, "passkey")
}

// ListPasskeys returns enrolled credentials for the profile page.
func (s *AuthService) ListPasskeys(userID uint) ([]store.WebAuthnCredential, error) {
	return s.store.ListWebAuthnCredentials(userID)
}

// DeletePasskey removes one credential owned by the user.
func (s *AuthService) DeletePasskey(userID uint, credentialID []byte) error {
	if err := s.store.DeleteWebAuthnCredential(userID, credentialID); err != nil {
		return err
	}
	s.createAuditLog(&userID, "passkey_delete", "user", fmt.Sprintf("%d", userID), "", "", "passkey removed")
	return nil
}
