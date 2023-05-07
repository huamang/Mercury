package fingerprint

// ParseBanner tcpBanner指纹识别
func ParseBanner(banner []byte) string {
	for name, rule := range bannerRules {
		if rule.Match(banner) {
			return name
		}
	}
	return ""
}
