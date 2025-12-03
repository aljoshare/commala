package utils

func AllTrue(s map[string]bool) bool {
	for _, v := range s {
		if !v {
			return false
		}
	}
	return true
}

func AllFalse(s map[string]bool) bool {
	for _, v := range s {
		if v {
			return false
		}
	}
	return true
}
