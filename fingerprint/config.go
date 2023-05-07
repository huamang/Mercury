package fingerprint

import "regexp"

// 指纹识别规则
var (
	bannerRules = map[string]*regexp.Regexp{
		"SSH":   regexp.MustCompile(`(?i)ssh`),
		"FTP":   regexp.MustCompile(`(?i)ftp`),
		"MYSQL": regexp.MustCompile(`(?i)mysql`),
	}
)
