package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func contem(fatia []string, s string) bool { // ➊
	for _, item := range fatia {
		if s == item {
			return true
		}
	}
	return false
}

func separar(s string) []string { // ➊
	separador := func(c rune) bool { // ➋
		return c == ' ' || c == '-'
	}
	return strings.FieldsFunc(s, separador)
}

// AnalisarLinha devolve a runa, o nome e uma fatia de palavras que
// ocorrem nos campo 1 e 10 de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) {
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	nome := campos[1] // ➊
	palavras := separar(campos[1])
	if campos[10] != "" { // ➋
		nome += fmt.Sprintf(" (%s)", campos[10])
		for _, palavra := range separar(campos[10]) { // ➌
			if !contem(palavras, palavra) { // ➍
				palavras = append(palavras, palavra) // ➎
			}
		}
	}
	return rune(código), nome, palavras
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
