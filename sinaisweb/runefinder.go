package main

import (
	"bufio"
	"bytes"
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

// URLUCD fica em http://www.unicode.org/Public/UNIDATA/UnicodeData.txt
// mas unicode.org não é confiável, então esta URL alternativa pode ser usada:
// http://turing.com.br/etc/UnicodeData.txt
const URLUCD = "http://turing.com.br/etc/UnicodeData.txt"

const ENDEREÇO = ":8080"

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

func carregar(texto io.Reader) []string {
	linhas := []string{}
	varredor := bufio.NewScanner(texto)
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		linhas = append(linhas, linha)
	}
	return linhas
}

// Listar produz texto com listagem com código, runa e nome dos
// caracteres Unicode cujo nome contem as palavras da consulta.
func Listar(linhas []string, consulta string) string {
	termos := separar(consulta)

	var buffer bytes.Buffer

	for _, linha := range linhas {
		runa, nome, palavrasNome := AnalisarLinha(linha) // ➊
		if contémTodos(palavrasNome, termos) {           // ➋
			buffer.WriteString(fmt.Sprintf("U+%04X\t%[1]c\t%s\n", runa, nome))
		}
	}

	return buffer.String()
}

// Exibir exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem as palavras da consulta.
func Exibir(linhas []string, consulta string) {
	fmt.Print(Listar(linhas, consulta))
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

func extrairOpções(args []string) (opções []string, resto []string) {
	opções = []string{}
	resto = []string{}
	for _, item := range args {
		if item[0] == '-' {
			opções = append(opções, item)
		} else {
			resto = append(resto, item)
		}
	}
	return opções, resto
}

const html = `<html><head/>
<body>
   <form action="/" method="GET">
   <input type="text" name="consulta">
   <input type="submit" value="Buscar">
  </form>
  <pre>%s</pre>
</body></html>`

func fazRespondedor(linhas []string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		saida := ""
		if r.URL.Query().Encode() != "" {
			consulta := r.URL.Query().Get("consulta")
			if consulta != "" {
				consulta = strings.ToUpper(consulta)
				saida = Listar(linhas, consulta)
			}
		}
		fmt.Fprintf(w, html, saida)
	}
}

// IniciarServidor sobe um servidor HTTP para receber consultas
func IniciarServidor(linhas []string, consulta string) {
	http.HandleFunc("/", fazRespondedor(linhas))
	fmt.Println("Servindo HTTP em", ENDEREÇO)
	http.ListenAndServe(ENDEREÇO, nil)
}

func main() {
	opções, palavras := extrairOpções(os.Args[1:])
	consulta := strings.Join(palavras, " ")
	consulta = strings.ToUpper(consulta)
	ucd, err := abrirUCD(obterCaminhoUCD()) // ➊
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ucd.Close()
	if contém(opções, "-w") {
		IniciarServidor(carregar(ucd), consulta)
	} else {
		Exibir(carregar(ucd), consulta)
	}
}
