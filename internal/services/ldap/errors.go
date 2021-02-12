package ldap

import (
	"fmt"
)

var (
	ErrConnectionFailed = newLdapError("connection to ldap failed")
	ErrBindFailed       = newLdapError("user binding failed - possible wrong passsword")
)

type Error struct {
	Context string
	Err     error
}

func (l *Error) Error() string {
	return fmt.Sprintf("%s: %v", l.Context, l.Err)
}

func (l *Error) Wrap(err error) *Error {
	l.Err = err
	return l
}

func newLdapError(info string) *Error {
	return &Error{
		Context: info,
	}
}
