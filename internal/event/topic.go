package event

const (
	// UserCreated : User was created
	UserCreated = "user.created"

	// UserUpdated : User was updated
	UserUpdated = "user.updated"

	// UserDeleted : User was deleted
	UserDeleted = "user.deleted"

	// GroupCreated : Group was created
	GroupCreated = "group.created"

	// GroupDeleted : Group was deleted
	GroupDeleted = "group.removed"

	// UserAddedToGroup : User was added to a group
	UserAddedToGroup = "group.user.added"

	// DeviceCreated : Device was created
	DeviceCreated = "device.created"

	// DeviceDeleted : Device was deleted
	DeviceDeleted = "device.deleted"
)

// UserCreatedEvent : Event when user was created
type UserCreatedEvent struct {
	UserID   uint
	Username string
}

// UserDeletedEvent : Event when user was deleted
type UserDeletedEvent struct {
	UserID uint
}

// GroupCreatedEvent : Event when group was created
type GroupCreatedEvent struct {
	GroupID uint
	Name    string
}

// GroupDeletedEvent : Event when group was removed
type GroupDeletedEvent struct {
	GroupID uint
}

// UserAddedToGroupEvent : Event when user was added to group
type UserAddedToGroupEvent struct {
	UserID  uint
	GroupID uint
}

// DeviceCreatedEvent : Event when device was created
type DeviceCreatedEvent struct {
	DeviceID uint
	UserID   uint
}

// DeviceDeletedEvent : Event when device was deleted
type DeviceDeletedEvent struct {
	DeviceID uint
	UserID   uint
}
