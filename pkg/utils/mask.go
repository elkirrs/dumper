package utils

import (
	"regexp"
	"strings"
)

var mask = "********"

func Mask(s string) string {
	out := s

	reURI := regexp.MustCompile(`([A-Za-z][A-Za-z0-9+.-]*://)([^:@/\s]+):([^@/\s]+)@`)
	out = reURI.ReplaceAllString(out, `${1}`+mask+`:`+mask+`@`)

	reUserPass := regexp.MustCompile(`(?i)(\b[A-Za-z0-9._%+\-]{1,64}):([^@\s,;()<>]+)@`)
	out = reUserPass.ReplaceAllString(out, mask+`:`+mask+`@`)

	reKeyEq := regexp.MustCompile(`(?i)(--(password|passwd|pwd|user|username|login)|password|user|username|login|PWD|MYSQL_PWD|PGPASSWORD|MYSQL_USER|PGUSER|--uri|uri|mongo_uri)(=)(["']?)([^"' \t\r\n;|&>]+)(["']?)`)
	out = reKeyEq.ReplaceAllString(out, `${1}${3}${4}`+mask+`${6}`)

	reKeySpace := regexp.MustCompile(`(?i)(--(password|passwd|pwd|user|username|login)|password|user|username|login|PWD|MYSQL_PWD|PGPASSWORD|MYSQL_USER|PGUSER|--uri|uri|mongo_uri)(\s+)(["']?)([^"' \t\r\n;|&>]+)(["']?)`)
	out = reKeySpace.ReplaceAllString(out, `${1}${3}${4}`+mask+`${6}`)

	reShortP := regexp.MustCompile(`(?i)(-[pu])([^ \t\r\n;|&>]*)`)
	out = reShortP.ReplaceAllString(out, `${1}`+mask)

	reEnvCb := regexp.MustCompile(`(?i)\b([A-Z_]*(PWD|PASSWORD|PGPASSWORD|MYSQL_PWD|MONGO_URI|URI|USER|USERNAME|LOGIN|AWS_SECRET_ACCESS_KEY|AWS_SECRET))\s*=\s*(["']?)([^"' \t\r\n;|&>]+)(["']?)`)
	out = reEnvCb.ReplaceAllStringFunc(out, func(m string) string {
		sub := reEnvCb.FindStringSubmatch(m)
		if len(sub) < 6 {
			return m
		}
		key := sub[1]
		leadQuote := sub[3]
		endQuote := sub[5]
		return key + "=" + leadQuote + mask + endQuote
	})

	reUserPassTCP := regexp.MustCompile(`([A-Za-z0-9._%+\-]{1,64}):([^@/\s]+)@tcp\(`)
	out = reUserPassTCP.ReplaceAllString(out, mask+`:`+mask+`@tcp(`)

	out = strings.ReplaceAll(out, `\"`, `"`)
	out = strings.ReplaceAll(out, `\'`, `'`)

	return out
}
