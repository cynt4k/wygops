package ldap

import (
	ldapgo "github.com/go-ldap/ldap/v3"
)

func (s *Service) CheckLogin(username string, password string) (bool, error) {
	con, err := s.connect()

	if err != nil {
		return false, ErrConnectionFailed.Wrap(err)
	}

	defer con.Close()

	filter := createFilter(username, s.config.UserAttr, s.config.UserFilter)

	searchRequest := ldapgo.NewSearchRequest(
		s.config.BaseDn,
		ldapgo.ScopeWholeSubtree, ldapgo.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn", "userAccountControl"},
		nil,
	)

	sr, err := con.Search(searchRequest)

	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, nil
	}

	userDN := sr.Entries[0].DN

	uac := sr.Entries[0].GetAttributeValue("userAccountControl")

	if uac != "" && IsLdapUserDisabled(uac) {
		return false, nil
	}

	err = con.Bind(userDN, password)

	if err != nil {
		return false, ErrBindFailed.Wrap(err)
	}

	return true, nil
}
