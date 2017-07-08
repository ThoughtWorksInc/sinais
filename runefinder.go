package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// UCD_URL é a URL canônica do file UnicodeData.txt mais atual
const UCD_URL = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

// ParseLine devolve a rune, o name e uma slice de words que
// ocorrem no campo name de uma line do UnicodeData.txt
func ParseLine(line string) (rune, string, []string) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	name := fields[1]
	words := split(fields[1])
	if fields[10] != "" { // ➊
		name += fmt.Sprintf(" (%s)", fields[10])
		for _, word := range split(fields[10]) { // ➋
			if !contains(words, word) { // ➌
				words = append(words, word) // ➍
			}
		}
	}
	return rune(code), name, words
}

func contains(slice []string, needle string) bool {
	for _, item := range slice {
		if item == needle {
			return true // ➋
		}
	}
	return false // ➌
}

func containsAll(slice []string, needles []string) bool {
	for _, needle := range needles {
		if !contains(slice, needle) {
			return false
		}
	}
	return true
}

func split(s string) []string { // ➊
	separator := func(c rune) bool { // ➋
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(s, separator) // ➌
}

// List exibe na saída padrão o code, a rune e o name dos caracteres Unicode
// cujo name contem as words da query.
func List(text io.Reader, query string) {
	terms := split(query)
	scanner := bufio.NewScanner(text)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		rune, name, wordsName := ParseLine(line) // ➊
		if containsAll(wordsName, terms) {           // ➋
			fmt.Printf("U+%04X\t%[1]c\t%s\n", rune, name)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getUCDPath() string {
	ucdPath := os.Getenv("UCD_PATH")
	if ucdPath == "" {
		user, err := user.Current()
		check(err)
		ucdPath = user.HomeDir + "/UnicodeData.txt"
	}
	return ucdPath
}

func progress(done <-chan bool) {
	for {
		select {
		case <-done:
			fmt.Println()
			return
		default:
			fmt.Print(".")
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func fetchUCD(url, path string, done chan<- bool) {
	response, err := http.Get(url)
	check(err)
	defer response.Body.Close()
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	check(err)
	done <- true
}

func openUCD(path string) (*os.File, error) {
	ucd, err := os.Open(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s não encontrado\nbaixando %s\n", path, UCD_URL)
		done := make(chan bool)
		go fetchUCD(UCD_URL, path, done)
		progress(done)
		ucd, err = os.Open(path)
	}
	return ucd, err
}

func main() {
	ucd, err := openUCD(getUCDPath())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { ucd.Close() }()
	query := strings.Join(os.Args[1:], " ")
	List(ucd, strings.ToUpper(query))
}
