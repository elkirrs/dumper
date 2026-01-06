package mask

import (
	"regexp"
	"strings"
)

var mask = "********"

func Mask(s string) string {
	out := s

	reLongFlag := regexp.MustCompile(`(?i)(--(password|passwd|pwd|user|username|login)|PWD|MYSQL_PWD|PGPASSWORD|MYSQL_USER|PGUSER|--uri|uri|mongo_uri)(=|\s+)(["']?)([^"' \t\r\n;|&>]+)(["']?)`)
	out = reLongFlag.ReplaceAllString(out, `${1}${3}${4}`+mask+`${6}`)

	reShortFlag := regexp.MustCompile(`(^|\s)-([pu])([^ \t\r\n;|&>]*)`)
	out = reShortFlag.ReplaceAllString(out, `${1}-$2`+mask)

	reURI := regexp.MustCompile(`([A-Za-z][A-Za-z0-9+.-]*://)([^:@/\s]+):([^@/\s]+)@`)
	out = reURI.ReplaceAllString(out, `${1}`+mask+`:`+mask+`@`)

	reUserPassTCP := regexp.MustCompile(`([A-Za-z0-9._%+\-]{1,64}):([^@/\s]+)@tcp\(`)
	out = reUserPassTCP.ReplaceAllString(out, mask+`:`+mask+`@tcp(`)

	reUserPassSimple := regexp.MustCompile(`([A-Za-z0-9._%+\-]{1,64}):([^@/\s]+)@`)
	out = reUserPassSimple.ReplaceAllString(out, mask+`:`+mask+`@`)

	reEnv := regexp.MustCompile(`(?i)\b([A-Z_]*(PWD|PASSWORD|PGPASSWORD|MYSQL_PWD|MONGO_URI|URI|USER|USERNAME|LOGIN|AWS_SECRET_ACCESS_KEY|AWS_SECRET))\s*=\s*(["']?)([^"' \t\r\n;|&>]+)(["']?)`)
	out = reEnv.ReplaceAllStringFunc(out, func(m string) string {
		sub := reEnv.FindStringSubmatch(m)
		if len(sub) < 6 {
			return m
		}
		key := sub[1]
		leadQuote := sub[3]
		endQuote := sub[5]
		return key + "=" + leadQuote + mask + endQuote
	})

	out = strings.ReplaceAll(out, `\"`, `"`)
	out = strings.ReplaceAll(out, `\'`, `'`)

	return out
}
