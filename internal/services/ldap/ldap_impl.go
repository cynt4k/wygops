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
func NewService(repo repository.Repository, config *config.Config) (LDAP, error) {
	ldap := &Service{
		repo:   repo,
		config: &config.Provider.Ldap,
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

// FindUser : Find an LDAP user by login attributes
func (s *Service) FindUser(name string, recursiveGroup bool) (*User, error) {
	return s.getUser(name, true, recursiveGroup)
}

// GetUser : Get an LDAP User by its unique name
func (s *Service) GetUser(name string, recursiveGroup bool) (*User, error) {
	return s.getUser(name, false, recursiveGroup)
}

func (s *Service) getUser(username string, filterLdap bool, recursiveGroup bool) (*User, error) {
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
		[]string{"cn", "sn", "givenName", "memberOf", "sAMAccountName", "mail"},
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

	if recursiveGroup {
		searchRequest = ldapgo.NewSearchRequest(
			s.config.BaseDn,
			ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
			s.config.GroupFilter,
			nil,
			nil,
		)

		sr, err = s.connection.Search(searchRequest)

		if err != nil {
			return nil, err
		}

		memberGroups, err := s.findGroupsNestedGroup(foundUser, sr.Entries)

		if err != nil {
			return nil, err
		}

		user.Groups = memberGroups
	}
	return user, nil

}

func (s *Service) getUserForGroup(group *ldapgo.Entry) ([]*User, error) {
	var users []*User
	for _, attr := range group.Attributes {
		switch attr.Name {
		case "uniqueMember":
			fallthrough
		case "member":
			for _, entry := range attr.Values {
				switch entry {
				case "person":
					fallthrough
				case "inetOrgPerson":
					user, err := s.getUser(entry, false, false)

					if err != nil {
						return nil, err
					}

					users = append(users, user)
				}
			}
		}
	}
	return users, nil
}

func (s *Service) findGroupsNestedGroup(user *ldapgo.Entry, groups []*ldapgo.Entry) ([]string, error) {
	return s.findMembersNestedGroup(user, groups, false)
}

func (s *Service) findUsersNestedGroup(group *ldapgo.Entry, groups []*ldapgo.Entry) ([]*User, error) {
	var users []*User

	var findUsers func(group *ldapgo.Entry) error

	findUsers = func(group *ldapgo.Entry) error {
		for _, attr := range group.Attributes {
			switch attr.Name {
			case "uniqueMember":
				fallthrough
			case "member":
				for _, member := range attr.Values {
					user, err := s.getUser(member, false, false)

					if err != nil {
						return err
					}

					if user == nil {
						var groupInfo *ldapgo.Entry
						for _, groupEntry := range groups {
							if groupEntry.DN == member {
								groupInfo = groupEntry
							}
						}
						if groupInfo == nil {
							continue
						}
						return findUsers(groupInfo)
					}
					var isDuplicated bool
					for _, entry := range users {
						if entry.Username == user.Username {
							isDuplicated = true
						}
					}
					if isDuplicated {
						continue
					}
					users = append(users, user)
				}
			}
		}
		return nil
	}
	if err := findUsers(group); err != nil {
		return nil, err
	}
	return users, nil

	// memberGroups, err := s.findMembersNestedGroup(group, groups, false)

	// if err != nil {
	// 	return nil, err
	// }

	// var foundRdn bool
	// for _, entry := range group.Attributes {
	// 	switch entry.Name {
	// 	case s.config.GroupRDN:
	// 		foundRdn = true
	// 		memberGroups = append(memberGroups, entry.Values[0])
	// 	default:
	// 		continue
	// 	}
	// }

	// if !foundRdn {
	// 	return nil, fmt.Errorf("group rdn not found")
	// }

	// for _, memberGroup := range memberGroups {
	// 	filter := fmt.Sprintf("(&(%s=%s)%s)", s.config.GroupRDN, memberGroup, s.config.GroupFilter)
	// 	searchRequest := ldapgo.NewSearchRequest(
	// 		s.config.BaseDn,
	// 		ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
	// 		filter,
	// 		nil,
	// 		nil,
	// 	)

	// 	sr, err := s.connection.Search(searchRequest)

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	if len(sr.Entries) == 0 {
	// 		return nil, fmt.Errorf("no group found")
	// 	}

	// 	if len(sr.Entries) > 1 {
	// 		return nil, fmt.Errorf("too much groups found - not unique rdn")
	// 	}

	// 	foundUsers, err := s.getUserForGroup(sr.Entries[0])

	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	for _, user := range foundUsers {
	// 		users = append(users, user.Username)
	// 	}
	// }
}

func (s *Service) findMembersNestedGroup(user *ldapgo.Entry, groups []*ldapgo.Entry, findUser bool) ([]string, error) {
	if len(groups) == 0 {
		return nil, nil
	}

	var findMembership func(group *ldapgo.Entry) chan bool

	findMembership = func(group *ldapgo.Entry) chan bool {
		out := make(chan bool)
		go func() {
			for _, entry := range group.Attributes {
				if entry.Name == "member" || entry.Name == "uniqueMember" {
					for _, member := range entry.Values {
						if member == user.DN {
							out <- true
							return
						}
						var groupInfos *ldapgo.Entry
						for _, groupDetail := range groups {
							if groupDetail.DN == member {
								groupInfos = groupDetail
							}
						}
						if groupInfos == nil {
							continue
						}
						if resultNestedFind := <-findMembership(groupInfos); resultNestedFind {
							out <- true
							return
						}
					}
				}
			}
			out <- false
		}()
		return out
	}

	members := make(map[string]chan bool)

	var findStr string

	if findUser {
		findStr = s.config.UserRDN
	} else {
		findStr = s.config.GroupRDN
	}

	for _, entry := range groups {
		var foundRdn bool
		for _, groupAttr := range entry.Attributes {
			if groupAttr.Name == findStr {
				foundRdn = true
				members[groupAttr.Values[0]] = findMembership(entry)
			}
		}
		if !foundRdn {
			return nil, fmt.Errorf("group rdn not found")
		}
	}

	var resultMembers []string

	for name, entry := range members {
		isMember := <-entry

		if isMember {
			resultMembers = append(resultMembers, name)
		}
	}

	return resultMembers, nil
}

// GetGroup : Get an LDAP Group by its unique name
func (s *Service) GetGroup(name string, recursive bool) (*Group, error) {
	return s.getGroup(name, true, false)
}

// GetGroupAndUsers : Get an LDAP Group with users by its unique name
func (s *Service) GetGroupAndUsers(name string, recursive bool) (*Group, error) {
	return s.getGroup(name, false, recursive)
}

func (s *Service) getGroup(name string, noUser bool, recursive bool) (*Group, error) {
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
			if !noUser {
				for _, userDn := range attr.Values {
					user, err := s.getUser(userDn, false, false)
					if err != nil {
						return nil, err
					}
					if user == nil {
						continue
					}
					user.Groups = nil
					group.Members = append(group.Members, user)
				}
			}
		default:
			continue
		}
	}

	if recursive {
		for _, attr := range foundGroup.Attributes {
			switch attr.Name {
			case "uniqueMember":
				fallthrough
			case "member":
				searchRequest = ldapgo.NewSearchRequest(
					s.config.BaseDn,
					ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
					s.config.GroupFilter,
					nil,
					nil,
				)

				sr, err = s.connection.Search(searchRequest)

				if err != nil {
					return nil, err
				}

				nestedGroups, err := s.findUsersNestedGroup(foundGroup, sr.Entries)

				if err != nil {
					return nil, err
				}

				for _, entry := range nestedGroups {
					entry.Groups = nil
					group.Members = append(group.Members, entry)
				}

				// for _, entry := range nestedGroups {
				// 	user, err := s.getUser(entry, true, false)

				// 	if err != nil {
				// 		return nil, err
				// 	}
				// 	user.Groups = nil

				// 	group.Members = append(group.Members, *user)
				// }
			}
		}
	}

	return &group, nil
}
