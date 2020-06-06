package event

const (
	// UserCreated : User was created
	UserCreated = "user.created"

	// UserUpdated : User was updated
	UserUpdated = "user.updated"

	// GroupCreated : Group was created
	GroupCreated = "group.created"

	// UserAddedToGroup : User was added to a group
	UserAddedToGroup = "group.user.added"
)

// UserCreatedEvent : Event when user was created
type UserCreatedEvent struct {
	UserID   uint
	Username string
}

// GroupCreatedEvent : Event when group was created
type GroupCreatedEvent struct {
	GroupID uint
	Name    string
}

// UserAddedToGroupEvent : Event when user was added to group
type UserAddedToGroupEvent struct {
	UserID  uint
	GroupID uint
}
