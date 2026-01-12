package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"sync"

	"github.com/mailru/easyjson/jlexer"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stringChannel := make(chan []byte, 25)
	domainsChannel := make(chan string, 20)
	domainRegexp := regexp.MustCompile("@(.*\\." + domain + ".*)$")

	go func() {
		readFile(r, stringChannel)
	}()

	linesWg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		linesWg.Add(1)
		go func() {
			defer linesWg.Done()

			parseUsers(stringChannel, domainRegexp, domainsChannel)
		}()
	}
	go func() {
		linesWg.Wait()
		close(domainsChannel)
	}()

	result := DomainStat{}
	for domainName := range domainsChannel {
		result[domainName]++
	}

	return result, nil
}

func readFile(r io.Reader, output chan<- []byte) {
	defer close(output)

	reader := bufio.NewReaderSize(r, 1024*1024)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) > 0 {
				output <- line
			}
			break
		}
		if err != nil {
			panic(err)
		}
		output <- line
	}
}

func parseUsers(input <-chan []byte, domain *regexp.Regexp, output chan<- string) {
	var user User
	for line := range input {
		user.UnmarshalEasyJSON(&jlexer.Lexer{Data: line})
		email := user.Email

		matched := domain.Find([]byte(email))
		if len(matched) == 0 {
			continue
		}

		output <- string(bytes.ToLower(matched[1:]))
	}
}
