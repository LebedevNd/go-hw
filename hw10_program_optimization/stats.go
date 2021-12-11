package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (result []string, err error) {
	scanner := bufio.NewScanner(r)
	var sc fastjson.Scanner
	var email string

	for scanner.Scan() {
		bytes := scanner.Bytes()

		sc.InitBytes(bytes)
		for sc.Next() {
			if sc.Value().Get("Email") != nil {
				email = sc.Value().Get("Email").String()
				result = append(result, email[1:len(email)-1])
			}
		}
	}

	if len(result) == 0 {
		return []string{}, errors.New("no emails found")
	}
	return result, err
}

func countDomains(emails []string, domain string) (DomainStat, error) {
	result := make(DomainStat)

	var p []string
	var userDomain string
	for _, email := range emails {
		p = strings.Split(email, ".")
		if len(p) < 2 {
			continue
		}

		userDomain = strings.ToLower(p[len(p)-1])
		if userDomain == domain {
			address := strings.ToLower(strings.Split(email, "@")[1])
			result[address]++
		}
	}
	return result, nil
}
