package gutils

import "strconv"

func Int2String(a []int) []string {
	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return b
}
