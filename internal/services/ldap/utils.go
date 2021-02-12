package ldap

import "strconv"

func IsLdapUserDisabled(userAccountControl string) bool {
	uacInt, err := strconv.ParseInt(userAccountControl, 0, 10)
	if err != nil {
		return true
	}
	if int32(uacInt)&0x2 != 0 {
		return true // bit 2 set means account is disabled
	}

	return false
}
