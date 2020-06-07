package ldap

import "github.com/cynt4k/wygops/internal/repository"

// Service : Ldap Service struct
type Service struct {
	repo repository.Repository
}

// NewService : Create an new LDAP Service
func NewService(repo repository.Repository) (LDAP, error) {
	ldap := &Service{
		repo: repo,
	}
	return ldap, nil
}

// GetUser : Get an LDAP User by its unique name
func (s *Service) GetUser(username string) (User, error) {
	return User{}, nil
}

// GetGroup : Get an LDAP Group by its unique name
func (s *Service) GetGroup(name string) (Group, error) {
	return Group{}, nil
}
