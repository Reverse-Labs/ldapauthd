package main

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap"
)

// Identity type
type Identity struct {
	FullName  string
	UserName  string
	Groups    []string
	UIDNumber string
	GIDNumber string
}

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

	return conn.Bind(uid+","+baseDN, password)
}

func ldapQueryUser(username string) (Identity, error) {
	conn, err := ldap.Dial("tcp", ldapHost+":"+ldapPort)

	if err != nil {
		return Identity{}, err
	}

	defer conn.Close()

	filter := fmt.Sprintf(ldapFilter, ldap.EscapeFilter(username))

	searchReq := ldap.NewSearchRequest(
		baseDN, ldap.ScopeWholeSubtree, 0, 0, 0, false, filter, []string{}, []ldap.Control{},
	)

	results, err := conn.Search(searchReq)
	if len(results.Entries) > 0 && err == nil {
		return Identity{
			FullName:  results.Entries[0].GetAttributeValue("cn"),
			UserName:  results.Entries[0].GetAttributeValue("uid"),
			Groups:    results.Entries[0].GetAttributeValues("objectClass"),
			UIDNumber: results.Entries[0].GetAttributeValue("uidNumber"),
			GIDNumber: results.Entries[0].GetAttributeValue("gidNumber"),
		}, nil
	}

	return Identity{}, err
}
