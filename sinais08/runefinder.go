package main

import (
  "fmt"
  "io"
  "strings"
  "strconv"
  "bufio"
  "net/http"
)

func Listar(texto io.Reader, consulta string) (saida string) {
  termos := strings.Fields(consulta)
	varredor := bufio.NewScanner(texto)
  base := "U+%04X&#09;%[1]c&#09;%s<br/>"
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome, palavrasNome := AnalisarLinha(linha) // ➊
		if contémTodos(palavrasNome, termos) {           // ➋
      saida += base
      saida = fmt.Sprintf(saida, runa, nome)
		}
	}
  return saida
}

func contém(fatia []string, procurado string) bool { // ➊
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

func AnalisarLinha(linha string) (rune, string, []string) { // ➊
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	palavras := strings.Fields(campos[1])    // ➋
	return rune(código), campos[1], palavras // ➌
}

type meuManipulador struct{
  ucd string
}

func (m *meuManipulador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  saida := ""

  if r.URL.Query().Encode() != "" {
    palavra := r.URL.Query().Get("palavras")
    if palavra == "" {
      saida = "Palavra não encontrada"
    } else {
      saida = Listar(strings.NewReader(m.ucd), palavra)
    }
  }

  fmt.Fprintf(w, pagina, saida)
}

func main() {
  manipulador := &meuManipulador{ucd: linhas3Da43}
  http.Handle("/", manipulador)
  http.ListenAndServe(":8080", nil)
}

const linhas3Da43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

const pagina = `<html><head/>
                    <body>
                       <form action="/" method="GET">
                       Qual unicode está procurando?<br/>
                       <input type="text" name="palavras">
                       <input type="submit" value="Buscar">
                      </form>
                      <div>%s</div>
                    </body></html>`
