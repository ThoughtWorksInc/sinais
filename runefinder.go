package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// AnalisarLinha devolve a runa, o nome e uma fatia de palavras que
// ocorrem no campo nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) { // ➊
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	palavras := strings.Split(campos[1], " ") // ➋
	return rune(código), campos[1], palavras // ➌
}

// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem o texto da consulta.
func Listar(texto io.Reader, consulta string) {
	varredor := bufio.NewScanner(texto)
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome, _ := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
	}
}

func main() {
	ucd, err := os.Open("UnicodeData.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { ucd.Close() }()
	consulta := strings.Join(os.Args[1:], " ")
	Listar(ucd, strings.ToUpper(consulta))
}
