package main

import "strings"

type Env struct {
	TelApi        string
	TelChan       int64
	FilterInclude []string
	FilterExclude []string
}

func TitleContains(s string, include []string, exclude []string) bool {
	s2 := strings.ToLower(s)
	for _, v := range exclude {
		v := strings.ToLower(v)
		if strings.Contains(s2, v) {
			return false
		}
	}

	for _, v := range include {
		v := strings.ToLower(v)
		if strings.Contains(s2, v) {
			return true
		}
	}

	return false
}
