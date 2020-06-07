package ldap

// User : LDAP user struct
type User struct {
	Username  string
	FirstName string
	LastName  string
	Groups    []Group
}

// Group : LDAP group struct
type Group struct {
	Name    string
	Members []User
}

// LDAP : LDAP interface to be implemented
type LDAP interface {
	GetUser(username string) (User, error)
	GetGroup(name string) (Group, error)
}
