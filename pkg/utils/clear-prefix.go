package utils

func CleanPrefix(s string) string {
	prefixes := []string{"test_", "acc_", "prod_"}
	for _, prefix := range prefixes {
		if len(s) > len(prefix) && s[:len(prefix)] == prefix {
			return s[len(prefix):]
		}
	}
	return s
}
