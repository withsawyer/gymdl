package utils

import (
	"bufio"
	"os"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// normalizeDomain 去掉开头的点并统一小写
func normalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimPrefix(domain, "."))
}

// isRootDomain 判断是否为注册域名（根域名）
func isRootDomain(domain string) bool {
	domain = normalizeDomain(domain)
	_, icann := publicsuffix.PublicSuffix(domain)
	eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return false
	}
	// 如果 eTLD+1 等于 domain 本身，则是根域名
	return domain == eTLDPlusOne && icann
}

// GetCookieValue 从 cookies.txt 文件中获取指定 domain 和 name 的 cookie value
func GetCookieValue(filePath, targetDomain, targetName string) string {
	targetDomain = normalizeDomain(targetDomain)

	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			parts = strings.Fields(line)
			if len(parts) < 7 {
				continue
			}
		}

		domain := normalizeDomain(parts[0])
		name := parts[5]
		value := parts[6]

		if domain == targetDomain && name == targetName {
			return value
		}
	}

	return ""
}

// GetCookiesByDomain 获取指定域名相关的所有 cookie（map[name]value）
func GetCookiesByDomain(filePath, inputDomain string) map[string]string {
	result := make(map[string]string)
	domainNormalized := normalizeDomain(inputDomain)
	matchRoot := isRootDomain(domainNormalized)

	f, err := os.Open(filePath)
	if err != nil {
		return result
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			parts = strings.Fields(line)
			if len(parts) < 7 {
				continue
			}
		}

		cookieDomain := normalizeDomain(parts[0])
		name := parts[5]
		value := parts[6]

		if matchRoot {
			if strings.HasSuffix(cookieDomain, domainNormalized) {
				result[name] = value
			}
		} else {
			if cookieDomain == domainNormalized {
				result[name] = value
			}
		}
	}

	return result
}
