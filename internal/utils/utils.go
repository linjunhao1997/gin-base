package gutils

import "strconv"

func Int2Strings(a []int) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return b
}

func Int2String(a int) string {
	return strconv.Itoa(a)
}
