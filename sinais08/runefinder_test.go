package main

import (
  "testing"
  "io"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
)

func TestBuscarRuna(t *testing.T) {
  testes := []struct {
    descrição     string
    metodo        string
    caminho       string
    corpo         io.Reader
    corpoEsperado string
  } {{
    descrição: "Carrega primeira vez",
    caminho:   "/",
    corpoEsperado : `<html><head/>
                    <body>
                       <form action="/" method="GET">
                       Qual unicode está procurando?<br/>
                       <input type="text" name="palavras">
                       <input type="submit" value="Buscar">
                      </form>
                      <div></div>
                    </body></html>`,
  },
  {
    descrição: "Busca sem palavra digitada",
    caminho:   "/?palavras=",
    corpoEsperado : `<html><head/>
                    <body>
                       <form action="/" method="GET">
                       Qual unicode está procurando?<br/>
                       <input type="text" name="palavras">
                       <input type="submit" value="Buscar">
                      </form>
                      <div>Palavra não encontrada</div>
                    </body></html>`,
  },
  {
    descrição: "Primeira Busca de verdade",
    caminho:   "/?palavras=SIGN",
    corpoEsperado: `<html><head/>
                    <body>
                       <form action="/" method="GET">
                       Qual unicode está procurando?<br/>
                       <input type="text" name="palavras">
                       <input type="submit" value="Buscar">
                      </form>
                      <div>U+003D&#09;=&#09;EQUALS SIGN<br/>U+003E&#09;>&#09;GREATER-THAN SIGN<br/></div>
                    </body></html>`,
  },
  {
    descrição: "Que acontece se busco 2 vezes?",
    caminho:   "/?palavras=SIGN",
    corpoEsperado: `<html><head/>
                    <body>
                       <form action="/" method="GET">
                       Qual unicode está procurando?<br/>
                       <input type="text" name="palavras">
                       <input type="submit" value="Buscar">
                      </form>
                      <div>U+003D&#09;=&#09;EQUALS SIGN<br/>U+003E&#09;>&#09;GREATER-THAN SIGN<br/></div>
                    </body></html>`,
  }}

  servidor := httptest.NewServer(&meuManipulador{ucd: linhas3Da43})
  defer servidor.Close()

  for _, teste := range testes {
    requisição, _ := http.NewRequest("Get", servidor.URL+teste.caminho, teste.corpo)

    resposta, _ := http.DefaultClient.Do(requisição)
    corpo, _ := ioutil.ReadAll(resposta.Body)
    if string(corpo) != teste.corpoEsperado {
      t.Errorf("Teste: %s\n", teste.descrição)
      t.Errorf("Corpo esperado: %q; Corpo recebido: %q", teste.corpoEsperado, corpo)
    }
  }
}
