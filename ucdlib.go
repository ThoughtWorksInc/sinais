package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
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
// cujo nome contem o texto da consulta
func Listar(texto io.Reader, consulta string) int { // <1>
	varredor := bufio.NewScanner(texto)
	count := 0 // <2>
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
			count++ // <3>
		}
	}
	return count // <4>
}

const UCDFileName = "UnicodeData.txt"

func main() { // <1>
	ucd, err := os.Open(UCDFileName)
	check(err)
	consulta := strings.ToUpper(strings.Join(os.Args[1:], " "))
	count := Listar(ucd, consulta)
	var plural string
	if count != 1 {
		plural = "s"
	}
	fmt.Println(count, "character"+plural, "found")
}
