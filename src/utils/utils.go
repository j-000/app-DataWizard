package utils

func Substring(data string, n int) string {
	if len(data) < n {
		return data
	}
	return data[:n]
}

func ContainsString(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}
