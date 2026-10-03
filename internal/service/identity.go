package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	idauth "github.com/ciliverse/cilikube/internal/auth"
	"github.com/ciliverse/cilikube/internal/models"
	"github.com/ciliverse/cilikube/internal/store"
	appauth "github.com/ciliverse/cilikube/pkg/auth"
)

type mfaChallenge struct {
	userID  uint
	expires time.Time
}

var mfaChallenges = struct {
	sync.Mutex
	items map[string]mfaChallenge
}{items: map[string]mfaChallenge{}}

func (s *AuthService) beginMFA(userID uint) (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := hex.EncodeToString(buf)
	mfaChallenges.Lock()
	mfaChallenges.items[token] = mfaChallenge{userID: userID, expires: time.Now().Add(5 * time.Minute)}
	mfaChallenges.Unlock()
	return token, nil
}

func takeMFA(token string) (uint, error) {
	mfaChallenges.Lock()
	defer mfaChallenges.Unlock()
	item, ok := mfaChallenges.items[token]
	if !ok || time.Now().After(item.expires) {
		delete(mfaChallenges.items, token)
		return 0, errors.New("mfa challenge expired")
	}
	delete(mfaChallenges.items, token)
	return item.userID, nil
}

// VerifyMFA completes a password login that is waiting on a TOTP code.
func (s *AuthService) VerifyMFA(mfaToken, code, ipAddress, userAgent string) (*models.LoginResponse, error) {
	userID, err := takeMFA(mfaToken)
	if err != nil {
		return nil, err
	}
	user, err := s.store.GetUserByID(userID)
	if err != nil || !user.TOTPEnabled || !idauth.ValidateTOTP(user.TOTPSecret, code, time.Now()) {
		return nil, errors.New("invalid authenticator code")
	}
	return s.finishLogin(user, ipAddress, userAgent, "totp")
}

// BeginTOTPSetup returns a secret that is not enabled until ConfirmTOTP.
func (s *AuthService) BeginTOTPSetup(userID uint) (secret, url string, err error) {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return "", "", errors.New("user not found")
	}
	secret, err = idauth.GenerateTOTPSecret()
	if err != nil {
		return "", "", err
	}
	user.TOTPSecret = secret
	user.TOTPEnabled = false
	if err := s.store.UpdateUser(user); err != nil {
		return "", "", err
	}
	return secret, idauth.TOTPAuthURL("CiliKube", user.Username, secret), nil
}

// ConfirmTOTP turns on MFA after the user proves they can generate a code.
func (s *AuthService) ConfirmTOTP(userID uint, code string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if user.TOTPSecret == "" || !idauth.ValidateTOTP(user.TOTPSecret, code, time.Now()) {
		return errors.New("invalid authenticator code")
	}
	user.TOTPEnabled = true
	return s.store.UpdateUser(user)
}

// DisableTOTP turns MFA off when the current code matches.
func (s *AuthService) DisableTOTP(userID uint, code string) error {
	user, err := s.store.GetUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if !user.TOTPEnabled || !idauth.ValidateTOTP(user.TOTPSecret, code, time.Now()) {
		return errors.New("invalid authenticator code")
	}
	user.TOTPEnabled = false
	user.TOTPSecret = ""
	return s.store.UpdateUser(user)
}

func (s *AuthService) finishLogin(storeUser *store.User, ipAddress, userAgent, method string) (*models.LoginResponse, error) {
	if err := s.securityService.RecordSuccessfulLogin(storeUser.ID, ipAddress, userAgent); err != nil {
		fmt.Printf("Failed to record successful login: %v\n", err)
	}
	now := time.Now()
	storeUser.LastLoginAt = &now
	if err := s.store.UpdateUser(storeUser); err != nil {
		fmt.Printf("Failed to update last login time: %v\n", err)
	}
	if s.syncRoles != nil {
		_ = s.syncRoles(storeUser.ID)
	}
	user := s.convertStoreUserToModelsUser(storeUser)
	roles, err := s.store.GetUserRoles(storeUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	if len(roles) > 0 {
		user.Role = primaryRoleName(roles)
	} else {
		user.Role = "viewer"
	}
	sessionID, err := s.securityService.CreateSession(storeUser.ID, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	token, expiresAt, err := appauth.GenerateToken(&user, sessionID)
	if err != nil {
		_ = s.securityService.InvalidateSession(sessionID)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	loginDetails, _ := json.Marshal(map[string]interface{}{
		"username": storeUser.Username, "ip": ipAddress, "user_agent": userAgent,
		"session_id": sessionID, "method": method, "result": "success",
	})
	s.createAuditLog(&storeUser.ID, "login", "user", fmt.Sprintf("%d", storeUser.ID), ipAddress, userAgent, string(loginDetails))
	return &models.LoginResponse{Token: token, ExpiresAt: expiresAt, User: user.ToResponse()}, nil
}
