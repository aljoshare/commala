package utils

func AllTrueSlice(s []bool) bool {
	for _, v := range s {
		if !v {
			return false
		}
	}
	return true
}

func AllFalseSlice(s []bool) bool {
	for _, v := range s {
		if v {
			return false
		}
	}
	return true
}

func AllTrueMap(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}
