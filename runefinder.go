package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// UCDURL is the canonical URL of the current UnicodeData.txt file
const UCDURL = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

// AnalisarLinha devolve a runa, o nome e uma fatia de palavras que
// ocorrem no campo nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) {
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	nome := campos[1]
	palavras := separar(campos[1])
	if campos[10] != "" { // ➊
		nome += fmt.Sprintf(" (%s)", campos[10])
		for _, palavra := range separar(campos[10]) { // ➋
			if !contém(palavras, palavra) { // ➌
				palavras = append(palavras, palavra) // ➍
			}
		}
	}
	return rune(código), nome, palavras
}

func contém(fatia []string, procurado string) bool {
	for _, item := range fatia {
		if item == procurado {
			return true // ➋
		}
	}
	return false // ➌
}

func contémTodos(fatia []string, procurados []string) bool {
	for _, procurado := range procurados {
		if !contém(fatia, procurado) {
			return false
		}
	}
	return true
}

func separar(s string) []string { // ➊
	separador := func(c rune) bool { // ➋
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(s, separador) // ➌
}

// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem as palavras da consulta.
func Listar(texto io.Reader, consulta string) {
	termos := separar(consulta)
	varredor := bufio.NewScanner(texto)
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome, palavrasNome := AnalisarLinha(linha) // ➊
		if contémTodos(palavrasNome, termos) {           // ➋
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func obterCaminhoUCD() string {
	UCDPath := os.Getenv("UCD_PATH")
	if UCDPath == "" {
		usuário, err := user.Current()
		check(err)
		UCDPath = usuário.HomeDir + "/UnicodeData.txt"
	}
	return UCDPath
}

func downloadUCD(ucdPath string) {
	fmt.Printf("%s not found\ndownloading %s\n", ucdPath, UCDURL)
	feito := make(chan bool)
	progresso := func(feito <-chan bool) {
		for {
			select {
			case <-feito:
				fmt.Println("concluído!")
			case <-time.After(200 * time.Millisecond):
				fmt.Print(".")
			}
		}
	}
	go progresso(feito)
	defer func() {
		feito <- false
	}()
	response, err := http.Get(UCDURL)
	check(err)
	defer response.Body.Close()
	file, err := os.Create(ucdPath)
	check(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	check(err)
}

func abrirUCD(ucdPath string) (*os.File, bool, error) {
	var remoto bool // acesso remoto ao arquivo?
	ucd, err := os.Open(ucdPath)
	if os.IsNotExist(err) {
		remoto = true
		downloadUCD(ucdPath)
		ucd, err = os.Open(ucdPath)
	}
	return ucd, remoto, err
}

func main() {
	ucd, _, err := abrirUCD(obterCaminhoUCD())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { ucd.Close() }()
	consulta := strings.Join(os.Args[1:], " ")
	Listar(ucd, strings.ToUpper(consulta))
}
