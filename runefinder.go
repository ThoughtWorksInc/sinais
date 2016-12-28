package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const UCDFileName = "UnicodeData.txt" // ➊

func check(e error) { // ➊
	if e != nil {
		panic(e)
	}
}

// AnalisarLinha devolve a runa e o nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string) {
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	return rune(código), campos[1]
}

// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem o texto da consulta. Devolve o número de ocorrências.
func Listar(texto io.Reader, consulta string) int { // ➊
	varredor := bufio.NewScanner(texto)
	ocorrências := 0 // ➋
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
			ocorrências++ // ➌
		}
	}
	return ocorrências
}

func main() { // ➊
	ucd, err := os.Open(UCDFileName) // ➋
	check(err) // ➌
	consulta := strings.ToUpper(strings.Join(os.Args[1:], " ")) // ➍
	ocorrências := Listar(ucd, consulta) // ➎
	var plural string // ➏
	if ocorrências != 1 { // ➐
		plural = "s"
	}
	fmt.Println(ocorrências, "character"+plural, "found")
}
