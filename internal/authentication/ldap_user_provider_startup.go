package authentication

import (
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func (p *LDAPUserProvider) checkServer() (err error) {
	conn, err := p.connect(p.configuration.User, p.configuration.Password)
	if err != nil {
		return err
	}

	defer conn.Close()

	searchRequest := ldap.NewSearchRequest("", ldap.ScopeBaseObject, ldap.NeverDerefAliases,
		1, 0, false, "(objectClass=*)", []string{ldapSupportedExtensionAttribute}, nil)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(sr.Entries) != 1 {
		return nil
	}

	// Iterate the attribute values to see what the server supports.
	for _, attr := range sr.Entries[0].Attributes {
		if attr.Name == ldapSupportedExtensionAttribute {
			p.logger.Tracef("LDAP Supported Extension OIDs: %s", strings.Join(attr.Values, ", "))

			for _, oid := range attr.Values {
				if oid == ldapOIDPasswdModifyExtension {
					p.supportExtensionPasswdModify = true
					break
				}
			}

			break
		}
	}

	return nil
}

func (p *LDAPUserProvider) parseDynamicUsersConfiguration() {
	p.configuration.UsersFilter = strings.ReplaceAll(p.configuration.UsersFilter, "{username_attribute}", p.configuration.UsernameAttribute)
	p.configuration.UsersFilter = strings.ReplaceAll(p.configuration.UsersFilter, "{mail_attribute}", p.configuration.MailAttribute)
	p.configuration.UsersFilter = strings.ReplaceAll(p.configuration.UsersFilter, "{display_name_attribute}", p.configuration.DisplayNameAttribute)

	p.logger.Tracef("Dynamically generated users filter is %s", p.configuration.UsersFilter)

	p.usersAttributes = []string{
		p.configuration.DisplayNameAttribute,
		p.configuration.MailAttribute,
		p.configuration.UsernameAttribute,
	}

	if p.configuration.AdditionalUsersDN != "" {
		p.usersBaseDN = p.configuration.AdditionalUsersDN + "," + p.configuration.BaseDN
	} else {
		p.usersBaseDN = p.configuration.BaseDN
	}

	p.logger.Tracef("Dynamically generated users BaseDN is %s", p.usersBaseDN)

	if strings.Contains(p.configuration.UsersFilter, "{input}") {
		p.usersFilterReplacementInput = true
	}

	if strings.Contains(p.configuration.UsersFilter, "{epoch:win32}") {
		p.usersFilterReplacementEpochWin32 = true
	}

	p.logger.Tracef("Detected user filter replacements that need to be resolved per lookup are: input=%v, epoch:win32=%v", p.usersFilterReplacementInput, p.usersFilterReplacementEpochWin32)
}

func (p *LDAPUserProvider) parseDynamicGroupsConfiguration() {
	p.groupsAttributes = []string{
		p.configuration.GroupNameAttribute,
	}

	if p.configuration.AdditionalGroupsDN != "" {
		p.groupsBaseDN = ldap.EscapeFilter(p.configuration.AdditionalGroupsDN + "," + p.configuration.BaseDN)
	} else {
		p.groupsBaseDN = p.configuration.BaseDN
	}

	p.logger.Tracef("Dynamically generated groups BaseDN is %s", p.groupsBaseDN)

	if strings.Contains(p.configuration.GroupsFilter, "{input}") {
		p.groupsFilterReplacementInput = true
	}

	if strings.Contains(p.configuration.GroupsFilter, "{username}") {
		p.groupsFilterReplacementUsername = true
	}

	if strings.Contains(p.configuration.GroupsFilter, "{dn}") {
		p.groupsFilterReplacementDN = true
	}

	p.logger.Tracef("Detected group filter replacements that need to be resolved per lookup are: input=%v, username=%v, dn=%v", p.groupsFilterReplacementInput, p.groupsFilterReplacementUsername, p.groupsFilterReplacementDN)
}
