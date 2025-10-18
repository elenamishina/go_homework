package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	//nolint:depguard
	json "github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	lineScanner := bufio.NewScanner(r)
	domainSuffix := "." + domain

	for lineScanner.Scan() {
		var user User
		if err := json.Unmarshal(lineScanner.Bytes(), &user); err != nil {
			return nil, err
		}
		if !strings.HasSuffix(user.Email, domainSuffix) {
			continue
		}

		atPos := strings.LastIndex(user.Email, "@")
		if atPos < 0 {
			continue
		}
		domainName := strings.ToLower(user.Email[atPos+1:])
		result[domainName]++
	}

	return result, nil
}
