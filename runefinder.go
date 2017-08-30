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

func obterCaminhoUCD() string {
	caminhoUCD := os.Getenv("UCD_PATH")
	if caminhoUCD == "" { // ➊
		usuário, err := user.Current()                    // ➋
		terminarSe(err)                                   // ➌
		caminhoUCD = usuário.HomeDir + "/UnicodeData.txt" // ➍
	}
	return caminhoUCD
}

func terminarSe(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func baixarUCD(url, caminho string, feito chan<- bool) { // ➊
	resposta, err := http.Get(url)
	terminarSe(err)
	defer resposta.Body.Close()
	arquivo, err := os.Create(caminho)
	terminarSe(err)
	defer arquivo.Close()
	_, err = io.Copy(arquivo, resposta.Body)
	terminarSe(err)
	feito <- true // ➋
}

func progresso(feito <-chan bool) { // ➊
	for { // ➋
		select { // ➌
		case <-feito: // ➍
			fmt.Println()
			return
		default: // ➎
			fmt.Print(".")
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// URLUCD fica em http://www.unicode.org/Public/UNIDATA/UnicodeData.txt
// mas unicode.org não é confiável, então esta URL alternativa pode ser usada:
// http://turing.com.br/etc/UnicodeData.txt
const URLUCD = "http://turing.com.br/etc/UnicodeData.txt"

func abrirUCD(caminho string) (*os.File, error) {
	ucd, err := os.Open(caminho)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s não encontrado\nbaixando %s\n", caminho, URLUCD)
		feito := make(chan bool)             // ➊
		go baixarUCD(URLUCD, caminho, feito) // ➋
		progresso(feito)                     // ➌
		ucd, err = os.Open(caminho)          // ➌
	}
	return ucd, err // ➍
}

func main() {
	ucd, err := abrirUCD(obterCaminhoUCD()) // ➊
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
	consulta := strings.Join(os.Args[1:], " ")
	Listar(ucd, strings.ToUpper(consulta))
}
