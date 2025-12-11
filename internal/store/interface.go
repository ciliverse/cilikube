package store

import "time"

// ClusterStore defines all methods required for interacting with cluster data persistent storage.
type ClusterStore interface {
	CreateCluster(cluster *Cluster) error
	GetClusterByID(id string) (*Cluster, error)
	GetClusterByName(name string) (*Cluster, error)
	GetAllClusters() ([]Cluster, error)
	UpdateCluster(cluster *Cluster) error
	DeleteClusterByName(name string) error
	DeleteClusterByID(id string) error
}

// UserStore defines all methods required for interacting with user data persistent storage.
type UserStore interface {
	CreateUser(user *User) error
	GetUserByID(id uint) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
	ListUsers(offset, limit int) ([]*User, int64, error)
}

// RoleStore defines all methods required for interacting with role data persistent storage.
type RoleStore interface {
	CreateRole(role *Role) error
	GetRoleByID(id uint) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	UpdateRole(role *Role) error
	DeleteRole(id uint) error
	ListRoles() ([]*Role, error)
}

// UserRoleStore defines all methods required for managing user-role associations.
type UserRoleStore interface {
	AssignRole(userID, roleID uint) error
	RemoveRole(userID, roleID uint) error
	GetUserRoles(userID uint) ([]*Role, error)
	GetRoleUsers(roleID uint) ([]*User, error)
	HasRole(userID, roleID uint) (bool, error)
}

// OAuthStore defines all methods required for managing OAuth provider data.
type OAuthStore interface {
	CreateOAuthProvider(provider *OAuthProvider) error
	GetOAuthProvider(userID uint, provider string) (*OAuthProvider, error)
	GetOAuthProviderByProviderUserID(provider, providerUserID string) (*OAuthProvider, error)
	UpdateOAuthProvider(provider *OAuthProvider) error
	DeleteOAuthProvider(userID uint, provider string) error
	ListUserOAuthProviders(userID uint) ([]*OAuthProvider, error)
}

// AuditLogStore defines all methods required for managing audit logs.
type AuditLogStore interface {
	CreateAuditLog(log *AuditLog) error
	GetAuditLogsByUserID(userID uint, offset, limit int) ([]*AuditLog, int64, error)
	GetAuditLogsByAction(action string, offset, limit int) ([]*AuditLog, int64, error)
	ListAuditLogs(offset, limit int) ([]*AuditLog, int64, error)
}

// LoginAttemptStore defines all methods required for managing login attempts.
type LoginAttemptStore interface {
	CreateLoginAttempt(attempt *LoginAttempt) error
	GetLoginAttemptsByUserID(userID uint, since time.Time) ([]*LoginAttempt, error)
	GetLoginAttemptsByUsername(username string, since time.Time) ([]*LoginAttempt, error)
	GetLoginAttemptsByIP(ipAddress string, since time.Time) ([]*LoginAttempt, error)
	CleanupOldLoginAttempts(before time.Time) error
}

// UserSessionStore defines all methods required for managing user sessions.
type UserSessionStore interface {
	CreateUserSession(session *UserSession) error
	GetUserSession(sessionID string) (*UserSession, error)
	UpdateUserSession(session *UserSession) error
	DeleteUserSession(sessionID string) error
	GetUserSessions(userID uint) ([]*UserSession, error)
	DeleteUserSessions(userID uint) error
	CleanupExpiredSessions(before time.Time) error
}

// Store is the main interface that combines all storage interfaces
type Store interface {
	ClusterStore
	UserStore
	RoleStore
	UserRoleStore
	OAuthStore
	AuditLogStore
	LoginAttemptStore
	UserSessionStore

	// Initialize initializes the storage (creates tables, default data, etc.)
	Initialize() error
	// Close closes the storage connection
	Close() error
}
