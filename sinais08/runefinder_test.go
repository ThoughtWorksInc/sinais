package main

import (
  "testing"
  "bytes"
  "net/http"
  "net/http/httptest"
)

func TestCarregarPaginaPrincipal(t *testing.T) {
  servidor := httptest.NewServer(http.HandlerFunc(CarregarPaginaPrincipal))
  defer servidor.Close()
  resposta, err := http.Get(servidor.URL)
  if err != nil {
    t.Errorf("Esperado: Página principal; Recebido %q", err)
  }

  corpo_esperado := "<html><head/>"+
                    "<body>"+
                    "  <form action=\"/buscar\" method=\"GET\">"+
                    "   Qual unicode está procurando?<br/>"+
                    "   <input type=\"text\" name=\"palavras\">"+
                    "   <input type=\"submit\" value=\"Buscar\">"+
                    "  </form>"+
                    "</body></html>"

  buffer := new(bytes.Buffer)
  buffer.ReadFrom(resposta.Body)
  resposta.Body.Close()

  if buffer.String() != corpo_esperado {
    t.Errorf("Esperado: %v; recebido: %v", corpo_esperado, buffer.String())
  }
}

func TestCarregarResultado(t *testing.T) {
  servidor := httptest.NewServer(http.HandlerFunc(CarregarResultado))
  defer servidor.Close()
  resposta, _ := http.Get(servidor.URL)

  buffer := new(bytes.Buffer)
  buffer.ReadFrom(resposta.Body)
  resposta.Body.Close()

  if buffer.String() != "Palavra não encontrada"{
    t.Errorf("Esperado: Palavra não encontrada; recebido: %v", buffer.String())
  }
}
