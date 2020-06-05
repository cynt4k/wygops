package event

const (
	// UserCreated : User was created
	UserCreated = "user.created"

	// UserUpdated : User was updated
	UserUpdated = "user.updated"
)

// UserCreatedEvent : Event when user was created
type UserCreatedEvent struct {
	UserID   uint
	Username string
}
