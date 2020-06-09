package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/cynt4k/wygops/cmd/config"
	"github.com/cynt4k/wygops/internal/repository"
	ldapgo "github.com/go-ldap/ldap/v3"
)

// Service : Ldap Service struct
type Service struct {
	repo       repository.Repository
	connection *ldapgo.Conn
	config     *config.ProviderLdap
}

// NewService : Create an new LDAP Service
func NewService(repo repository.Repository, config *config.ProviderLdap) (LDAP, error) {
	ldap := &Service{
		repo:   repo,
		config: config,
	}
	return ldap, ldap.connect()
}

func (s *Service) connect() error {
	connection, err := ldapgo.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))

	if err != nil {
		return err
	}

	if s.config.Type == "tls" {
		err = connection.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return err
		}
	}

	err = connection.Bind(s.config.BindDn, s.config.BindPassword)

	if err != nil {
		return err
	}

	s.connection = connection
	return nil
}

func (s *Service) disconnect() {
	s.connection.Close()
}

// createFilter : Create a ldap search filter for all searchAttr
func createFilter(searchAttr string, attr []string, baseFilter string) string {
	filterAttr := ""

	for _, attr := range attr {
		filterAttr += fmt.Sprintf("(%s=%s)", attr, searchAttr)
	}

	filterAllArgs := fmt.Sprintf("(|%s)", filterAttr)

	return fmt.Sprintf("(&(%s)%s)", baseFilter, filterAllArgs)
}

// GetUser : Get an LDAP User by its unique name
func (s *Service) GetUser(username string, filterLdap bool) (*User, error) {
	var filter string
	if filterLdap {
		filter = createFilter(username, s.config.UserAttr, s.config.UserFilter)
	} else {
		filter = s.config.UserFilter
	}

	searchRequest := ldapgo.NewSearchRequest(
		s.config.BaseDn,
		ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"cn", "sn", "givenName", "memberOf"},
		nil,
	)

	sr, err := s.connection.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	if !filterLdap {
		found := false
		for _, attr := range sr.Entries {
			if attr.DN == username {
				found = true
				sr.Entries = []*ldapgo.Entry{attr}
				break
			}
		}

		if !found {
			return nil, nil
		}
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("to much entries for searchString %s", username)
	}

	foundUser := sr.Entries[0]
	user := &User{}

	for _, attr := range foundUser.Attributes {

		switch attr.Name {
		case s.config.UserRDN:
			user.Username = attr.Values[0]
		case "givenName":
			user.FirstName = attr.Values[0]
		case "sn":
			user.LastName = attr.Values[0]
		case "mail":
			user.Mail = attr.Values[0]
		case "memberOf":
			for _, group := range attr.Values {
				user.Groups = append(user.Groups, group)
			}
		default:
			continue
		}
		// switch attr.Name {
		// case s.config.UserRDN:
		// 	user.Username = attr.Values[0]
		// 	default:

		// }
	}

	// user := &User{
	// 	Username: foundUser.Attributes["asdf"],
	// }

	return user, nil
}

// GetGroup : Get an LDAP Group by its unique name
func (s *Service) GetGroup(name string, recursive bool) (*Group, error) {
	filter := createFilter(name, s.config.GroupAttr, s.config.GroupFilter)

	searchRequest := ldapgo.NewSearchRequest(
		s.config.BaseDn,
		ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
		filter,
		nil,
		nil,
	)

	sr, err := s.connection.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, nil
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("to much entries for searchString %s", name)
	}

	foundGroup := sr.Entries[0]
	group := Group{}

	for _, attr := range foundGroup.Attributes {
		switch attr.Name {
		case s.config.GroupRDN:
			group.Name = attr.Values[0]
		case "uniqueMember":
			fallthrough
		case "member":
			for _, userDn := range attr.Values {
				user, err := s.GetUser(userDn, false)
				if err != nil {
					return nil, err
				}
				group.Members = append(group.Members, *user)
			}
		default:
			continue
		}
	}

	return &group, nil
}
