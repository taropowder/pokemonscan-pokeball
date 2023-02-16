package utils

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DelFromSlice(str string, s []string) (ns []string) {
	for _, t := range s {
		if t != str {
			ns = append(ns, t)
		}
	}
	return
}
