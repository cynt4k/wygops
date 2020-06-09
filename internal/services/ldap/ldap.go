package ldap

// User : LDAP user struct
type User struct {
	Username  string
	FirstName string
	LastName  string
	Mail      string
	Groups    []string
}

// Group : LDAP group struct
type Group struct {
	Name    string
	Members []User
}

// LDAP : LDAP interface to be implemented
type LDAP interface {
	GetUser(username string, filterLdap bool) (*User, error)
	GetGroup(name string, recursive bool) (*Group, error)
}
