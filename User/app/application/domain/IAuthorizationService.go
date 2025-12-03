package domain

// Permission represents a specific permission that can be granted
type Permission string

const (
	// User permissions
	PermissionCreateUser Permission = "user:create"
	PermissionReadUser   Permission = "user:read"
	PermissionUpdateUser Permission = "user:update"
	PermissionDeleteUser Permission = "user:delete"

	// Event permissions (for future use)
	PermissionCreateEvent Permission = "event:create"
	PermissionReadEvent   Permission = "event:read"
	PermissionUpdateEvent Permission = "event:update"
	PermissionDeleteEvent Permission = "event:delete"

	// Ticket permissions (for future use)
	PermissionCreateTicket Permission = "ticket:create"
	PermissionReadTicket   Permission = "ticket:read"
	PermissionUpdateTicket Permission = "ticket:update"
	PermissionDeleteTicket Permission = "ticket:delete"
)

// Role represents a user role with associated permissions
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleOrganizer Role = "organizer"
	RoleGuest    Role = "guest"
)

// Resource represents a resource that can be accessed
type Resource struct {
	Type string // e.g., "user", "event", "ticket"
	ID   int    // Resource ID
}

// Action represents an action that can be performed on a resource
type Action string

const (
	ActionCreate Action = "create"
	ActionRead   Action = "read"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// IAuthorizationService defines the contract for authorization operations
type IAuthorizationService interface {
	// Authorize checks if a user can perform a specific action on a resource
	Authorize(userID int, resource *Resource, action Action) (bool, error)

	// CheckPermission checks if a user has a specific permission
	CheckPermission(userID int, permission Permission) (bool, error)

	// CheckRole checks if a user has a specific role
	CheckRole(userID int, role Role) (bool, error)

	// GrantPermission grants a permission to a user
	GrantPermission(userID int, permission Permission) error

	// RevokePermission revokes a permission from a user
	RevokePermission(userID int, permission Permission) error

	// AssignRole assigns a role to a user
	AssignRole(userID int, role Role) error

	// RemoveRole removes a role from a user
	RemoveRole(userID int, role Role) error

	// GetUserPermissions retrieves all permissions for a user
	GetUserPermissions(userID int) ([]Permission, error)

	// GetUserRoles retrieves all roles for a user
	GetUserRoles(userID int) ([]Role, error)
}
