package main

import (
  "fmt"
  "net/http"
)

type Pagina struct {
  Titulo  string
  Corpo   []byte
}

func CarregarPaginaPrincipal(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, corpoPaginaPrincipal)
}

func CarregarResultado(w http.ResponseWriter, r *http.Request) {
  if r.URL.Query().Encode() == "" {
    fmt.Fprintf(w, "Busca precisa de palavra")
    return
  }

  fmt.Fprintf(w, "Palavra não encontrada")
}

func main() {
  http.HandleFunc("/", CarregarPaginaPrincipal)
  http.HandleFunc("/buscar", CarregarResultado)
  http.ListenAndServe(":8080", nil)
}

const corpoPaginaPrincipal = "<html><head/>"+
                        "<body>"+
                        "  <form action=\"/buscar\" method=\"GET\">"+
                        "   Qual unicode está procurando?<br/>"+
                        "   <input type=\"text\" name=\"palavras\">"+
                        "   <input type=\"submit\" value=\"Buscar\">"+
                        "  </form>"+
                        "</body></html>"
