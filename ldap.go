package main

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap"
)

func authenticateLdap(username, password string) error {
	conn, err := ldap.Dial("tcp", ldapHost+":"+ldapPort)

	if err != nil {
		return err
	}

	defer conn.Close()

	if err := conn.Bind(bindDN, bindPassword); err != nil {
		return err
	}
	uid := fmt.Sprintf(strings.Trim(ldapFilter, "()"), username)

	return conn.Bind(uid+",cn=users,cn=accounts,dc=revlabs,dc=xyz", password)
}
