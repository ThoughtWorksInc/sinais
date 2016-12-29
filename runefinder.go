package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// AnalisarLinha devolve a runa e o nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) {
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	palavras := strings.Split(campos[1], " ")
	return rune(código), campos[1], palavras
}

// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem o texto da consulta.
func Listar(texto io.Reader, consulta string) {
	varredor := bufio.NewScanner(texto)
	ocorrências := 0 // ➋
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome, _ := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
			ocorrências++ // ➌
		}
	}
}

func main() { // ➊
	ucd, err := os.Open("UnicodeData.txt") // ➋
	if err != nil {                        // ➌
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { ucd.Close() }()             // ➍
	consulta := strings.Join(os.Args[1:], " ") // ➎
	Listar(ucd, strings.ToUpper(consulta))     // ➏
}
