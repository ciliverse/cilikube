package store

import (
	"fmt"
	"sync"
	"time"
)

// MemoryAuthStore is an in-memory implementation of all auth-related stores
type MemoryAuthStore struct {
	// Data storage
	users          map[uint]*User
	usersByName    map[string]*User
	usersByEmail   map[string]*User
	roles          map[uint]*Role
	rolesByName    map[string]*Role
	userRoles      map[uint][]uint           // userID -> roleIDs
	oauthProviders map[string]*OAuthProvider // key: userID_provider
	auditLogs      []*AuditLog

	// ID generators
	nextUserID     uint
	nextRoleID     uint
	nextUserRoleID uint
	nextOAuthID    uint
	nextAuditLogID uint

	// Mutex for thread safety
	mutex sync.RWMutex
}

// NewMemoryAuthStore creates a new in-memory auth store
func NewMemoryAuthStore() *MemoryAuthStore {
	return &MemoryAuthStore{
		users:          make(map[uint]*User),
		usersByName:    make(map[string]*User),
		usersByEmail:   make(map[string]*User),
		roles:          make(map[uint]*Role),
		rolesByName:    make(map[string]*Role),
		userRoles:      make(map[uint][]uint),
		oauthProviders: make(map[string]*OAuthProvider),
		auditLogs:      make([]*AuditLog, 0),
		nextUserID:     1,
		nextRoleID:     1,
		nextUserRoleID: 1,
		nextOAuthID:    1,
		nextAuditLogID: 1,
	}
}

// Initialize initializes the memory store with default data
func (s *MemoryAuthStore) Initialize() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create default roles
	adminRole := &Role{
		ID:          s.nextRoleID,
		Name:        "admin",
		DisplayName: "Administrator",
		Description: "System administrator with all permissions",
		IsSystem:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.nextRoleID++
	s.roles[adminRole.ID] = adminRole
	s.rolesByName[adminRole.Name] = adminRole

	editorRole := &Role{
		ID:          s.nextRoleID,
		Name:        "editor",
		DisplayName: "Editor",
		Description: "Can edit most resources, but cannot manage clusters and users",
		IsSystem:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.nextRoleID++
	s.roles[editorRole.ID] = editorRole
	s.rolesByName[editorRole.Name] = editorRole

	viewerRole := &Role{
		ID:          s.nextRoleID,
		Name:        "viewer",
		DisplayName: "Viewer",
		Description: "Can only view resources, cannot perform modification operations",
		IsSystem:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.nextRoleID++
	s.roles[viewerRole.ID] = viewerRole
	s.rolesByName[viewerRole.Name] = viewerRole

	// Create default admin user
	adminUser := &User{
		ID:            s.nextUserID,
		Username:      "admin",
		Email:         "admin@cilikube.com",
		DisplayName:   "System Administrator",
		IsActive:      true,
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Set password
	err := adminUser.HashPassword("12345678")
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	s.nextUserID++
	s.users[adminUser.ID] = adminUser
	s.usersByName[adminUser.Username] = adminUser
	s.usersByEmail[adminUser.Email] = adminUser

	// Assign admin role to admin user
	s.userRoles[adminUser.ID] = []uint{adminRole.ID}

	return nil
}

// Close closes the memory store (no-op for memory store)
func (s *MemoryAuthStore) Close() error {
	return nil
}

// UserStore implementation

func (s *MemoryAuthStore) CreateUser(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if username already exists
	if _, exists := s.usersByName[user.Username]; exists {
		return fmt.Errorf("username '%s' already exists", user.Username)
	}

	// Check if email already exists
	if _, exists := s.usersByEmail[user.Email]; exists {
		return fmt.Errorf("email '%s' already exists", user.Email)
	}

	// Assign ID and timestamps
	user.ID = s.nextUserID
	s.nextUserID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Store user
	userCopy := *user
	s.users[user.ID] = &userCopy
	s.usersByName[user.Username] = &userCopy
	s.usersByEmail[user.Email] = &userCopy

	return nil
}

func (s *MemoryAuthStore) GetUserByID(id uint) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", id)
	}

	userCopy := *user
	return &userCopy, nil
}

func (s *MemoryAuthStore) GetUserByUsername(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.usersByName[username]
	if !exists {
		return nil, fmt.Errorf("user with username '%s' not found", username)
	}

	userCopy := *user
	return &userCopy, nil
}

func (s *MemoryAuthStore) GetUserByEmail(email string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.usersByEmail[email]
	if !exists {
		return nil, fmt.Errorf("user with email '%s' not found", email)
	}

	userCopy := *user
	return &userCopy, nil
}

func (s *MemoryAuthStore) UpdateUser(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existingUser, exists := s.users[user.ID]
	if !exists {
		return fmt.Errorf("user with ID %d not found", user.ID)
	}

	// Check if username is being changed and if it conflicts
	if existingUser.Username != user.Username {
		if _, exists := s.usersByName[user.Username]; exists {
			return fmt.Errorf("username '%s' already exists", user.Username)
		}
		delete(s.usersByName, existingUser.Username)
		s.usersByName[user.Username] = user
	}

	// Check if email is being changed and if it conflicts
	if existingUser.Email != user.Email {
		if _, exists := s.usersByEmail[user.Email]; exists {
			return fmt.Errorf("email '%s' already exists", user.Email)
		}
		delete(s.usersByEmail, existingUser.Email)
		s.usersByEmail[user.Email] = user
	}

	user.UpdatedAt = time.Now()
	userCopy := *user
	s.users[user.ID] = &userCopy

	return nil
}

func (s *MemoryAuthStore) DeleteUser(id uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("user with ID %d not found", id)
	}

	// Remove from all indexes
	delete(s.users, id)
	delete(s.usersByName, user.Username)
	delete(s.usersByEmail, user.Email)
	delete(s.userRoles, id)

	// Remove OAuth providers
	for key := range s.oauthProviders {
		if s.oauthProviders[key].UserID == id {
			delete(s.oauthProviders, key)
		}
	}

	return nil
}

func (s *MemoryAuthStore) ListUsers(offset, limit int) ([]*User, int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	total := int64(len(s.users))
	users := make([]*User, 0)

	i := 0
	for _, user := range s.users {
		if i >= offset && len(users) < limit {
			userCopy := *user
			users = append(users, &userCopy)
		}
		i++
	}

	return users, total, nil
}

// RoleStore implementation

func (s *MemoryAuthStore) CreateRole(role *Role) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if role name already exists
	if _, exists := s.rolesByName[role.Name]; exists {
		return fmt.Errorf("role with name '%s' already exists", role.Name)
	}

	// Assign ID and timestamps
	role.ID = s.nextRoleID
	s.nextRoleID++
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	// Store role
	roleCopy := *role
	s.roles[role.ID] = &roleCopy
	s.rolesByName[role.Name] = &roleCopy

	return nil
}

func (s *MemoryAuthStore) GetRoleByID(id uint) (*Role, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	role, exists := s.roles[id]
	if !exists {
		return nil, fmt.Errorf("role with ID %d not found", id)
	}

	roleCopy := *role
	return &roleCopy, nil
}

func (s *MemoryAuthStore) GetRoleByName(name string) (*Role, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	role, exists := s.rolesByName[name]
	if !exists {
		return nil, fmt.Errorf("role with name '%s' not found", name)
	}

	roleCopy := *role
	return &roleCopy, nil
}

func (s *MemoryAuthStore) UpdateRole(role *Role) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	existingRole, exists := s.roles[role.ID]
	if !exists {
		return fmt.Errorf("role with ID %d not found", role.ID)
	}

	// Check if name is being changed and if it conflicts
	if existingRole.Name != role.Name {
		if _, exists := s.rolesByName[role.Name]; exists {
			return fmt.Errorf("role with name '%s' already exists", role.Name)
		}
		delete(s.rolesByName, existingRole.Name)
		s.rolesByName[role.Name] = role
	}

	role.UpdatedAt = time.Now()
	roleCopy := *role
	s.roles[role.ID] = &roleCopy

	return nil
}

func (s *MemoryAuthStore) DeleteRole(id uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	role, exists := s.roles[id]
	if !exists {
		return fmt.Errorf("role with ID %d not found", id)
	}

	// Don't allow deleting system roles
	if role.IsSystem {
		return fmt.Errorf("cannot delete system role '%s'", role.Name)
	}

	// Remove from all indexes
	delete(s.roles, id)
	delete(s.rolesByName, role.Name)

	// Remove role assignments
	for userID, roleIDs := range s.userRoles {
		newRoleIDs := make([]uint, 0)
		for _, roleID := range roleIDs {
			if roleID != id {
				newRoleIDs = append(newRoleIDs, roleID)
			}
		}
		s.userRoles[userID] = newRoleIDs
	}

	return nil
}

func (s *MemoryAuthStore) ListRoles() ([]*Role, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roles := make([]*Role, 0, len(s.roles))
	for _, role := range s.roles {
		roleCopy := *role
		roles = append(roles, &roleCopy)
	}

	return roles, nil
}

// UserRoleStore implementation

func (s *MemoryAuthStore) AssignRole(userID, roleID uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if user exists
	if _, exists := s.users[userID]; !exists {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	// Check if role exists
	if _, exists := s.roles[roleID]; !exists {
		return fmt.Errorf("role with ID %d not found", roleID)
	}

	// Check if assignment already exists
	roleIDs := s.userRoles[userID]
	for _, existingRoleID := range roleIDs {
		if existingRoleID == roleID {
			return nil // Already assigned
		}
	}

	// Add role assignment
	s.userRoles[userID] = append(roleIDs, roleID)

	return nil
}

func (s *MemoryAuthStore) RemoveRole(userID, roleID uint) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	roleIDs := s.userRoles[userID]
	newRoleIDs := make([]uint, 0)
	found := false

	for _, existingRoleID := range roleIDs {
		if existingRoleID != roleID {
			newRoleIDs = append(newRoleIDs, existingRoleID)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("user %d does not have role %d", userID, roleID)
	}

	s.userRoles[userID] = newRoleIDs

	return nil
}

func (s *MemoryAuthStore) GetUserRoles(userID uint) ([]*Role, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roleIDs := s.userRoles[userID]
	roles := make([]*Role, 0, len(roleIDs))

	for _, roleID := range roleIDs {
		if role, exists := s.roles[roleID]; exists {
			roleCopy := *role
			roles = append(roles, &roleCopy)
		}
	}

	return roles, nil
}

func (s *MemoryAuthStore) GetRoleUsers(roleID uint) ([]*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]*User, 0)

	for userID, roleIDs := range s.userRoles {
		for _, assignedRoleID := range roleIDs {
			if assignedRoleID == roleID {
				if user, exists := s.users[userID]; exists {
					userCopy := *user
					users = append(users, &userCopy)
				}
				break
			}
		}
	}

	return users, nil
}

func (s *MemoryAuthStore) HasRole(userID, roleID uint) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	roleIDs := s.userRoles[userID]
	for _, assignedRoleID := range roleIDs {
		if assignedRoleID == roleID {
			return true, nil
		}
	}

	return false, nil
}

// OAuthStore implementation

func (s *MemoryAuthStore) CreateOAuthProvider(provider *OAuthProvider) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%d_%s", provider.UserID, provider.Provider)

	// Check if already exists
	if _, exists := s.oauthProviders[key]; exists {
		return fmt.Errorf("OAuth provider '%s' already exists for user %d", provider.Provider, provider.UserID)
	}

	// Assign ID and timestamps
	provider.ID = s.nextOAuthID
	s.nextOAuthID++
	provider.CreatedAt = time.Now()
	provider.UpdatedAt = time.Now()

	// Store provider
	providerCopy := *provider
	s.oauthProviders[key] = &providerCopy

	return nil
}

func (s *MemoryAuthStore) GetOAuthProvider(userID uint, provider string) (*OAuthProvider, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	key := fmt.Sprintf("%d_%s", userID, provider)
	oauthProvider, exists := s.oauthProviders[key]
	if !exists {
		return nil, fmt.Errorf("OAuth provider '%s' not found for user %d", provider, userID)
	}

	providerCopy := *oauthProvider
	return &providerCopy, nil
}

func (s *MemoryAuthStore) GetOAuthProviderByProviderUserID(provider, providerUserID string) (*OAuthProvider, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, oauthProvider := range s.oauthProviders {
		if oauthProvider.Provider == provider && oauthProvider.ProviderUserID == providerUserID {
			providerCopy := *oauthProvider
			return &providerCopy, nil
		}
	}

	return nil, fmt.Errorf("OAuth provider '%s' with provider user ID '%s' not found", provider, providerUserID)
}

func (s *MemoryAuthStore) UpdateOAuthProvider(provider *OAuthProvider) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%d_%s", provider.UserID, provider.Provider)

	if _, exists := s.oauthProviders[key]; !exists {
		return fmt.Errorf("OAuth provider '%s' not found for user %d", provider.Provider, provider.UserID)
	}

	provider.UpdatedAt = time.Now()
	providerCopy := *provider
	s.oauthProviders[key] = &providerCopy

	return nil
}

func (s *MemoryAuthStore) DeleteOAuthProvider(userID uint, provider string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := fmt.Sprintf("%d_%s", userID, provider)

	if _, exists := s.oauthProviders[key]; !exists {
		return fmt.Errorf("OAuth provider '%s' not found for user %d", provider, userID)
	}

	delete(s.oauthProviders, key)

	return nil
}

func (s *MemoryAuthStore) ListUserOAuthProviders(userID uint) ([]*OAuthProvider, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	providers := make([]*OAuthProvider, 0)

	for _, provider := range s.oauthProviders {
		if provider.UserID == userID {
			providerCopy := *provider
			providers = append(providers, &providerCopy)
		}
	}

	return providers, nil
}

// AuditLogStore implementation

func (s *MemoryAuthStore) CreateAuditLog(log *AuditLog) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Assign ID and timestamp
	log.ID = s.nextAuditLogID
	s.nextAuditLogID++
	log.CreatedAt = time.Now()

	// Store log
	logCopy := *log
	s.auditLogs = append(s.auditLogs, &logCopy)

	return nil
}

func (s *MemoryAuthStore) GetAuditLogsByUserID(userID uint, offset, limit int) ([]*AuditLog, int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	filteredLogs := make([]*AuditLog, 0)

	for _, log := range s.auditLogs {
		if log.UserID != nil && *log.UserID == userID {
			filteredLogs = append(filteredLogs, log)
		}
	}

	total := int64(len(filteredLogs))

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(filteredLogs) {
		start = len(filteredLogs)
	}
	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}

	result := make([]*AuditLog, 0)
	for i := start; i < end; i++ {
		logCopy := *filteredLogs[i]
		result = append(result, &logCopy)
	}

	return result, total, nil
}

func (s *MemoryAuthStore) GetAuditLogsByAction(action string, offset, limit int) ([]*AuditLog, int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	filteredLogs := make([]*AuditLog, 0)

	for _, log := range s.auditLogs {
		if log.Action == action {
			filteredLogs = append(filteredLogs, log)
		}
	}

	total := int64(len(filteredLogs))

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(filteredLogs) {
		start = len(filteredLogs)
	}
	if end > len(filteredLogs) {
		end = len(filteredLogs)
	}

	result := make([]*AuditLog, 0)
	for i := start; i < end; i++ {
		logCopy := *filteredLogs[i]
		result = append(result, &logCopy)
	}

	return result, total, nil
}

func (s *MemoryAuthStore) ListAuditLogs(offset, limit int) ([]*AuditLog, int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	total := int64(len(s.auditLogs))

	// Apply pagination
	start := offset
	end := offset + limit
	if start > len(s.auditLogs) {
		start = len(s.auditLogs)
	}
	if end > len(s.auditLogs) {
		end = len(s.auditLogs)
	}

	result := make([]*AuditLog, 0)
	for i := start; i < end; i++ {
		logCopy := *s.auditLogs[i]
		result = append(result, &logCopy)
	}

	return result, total, nil
}
