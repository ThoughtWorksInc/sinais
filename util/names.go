package main

import (
  "io/ioutil"
  "go/scanner"
  "go/token"
  "log"
  "os"
  "fmt"
)

func main() {
  src, err := ioutil.ReadFile(os.Args[1] + ".go")
  if err != nil {
    log.Fatal(err.Error())
  }

  // Initialize the scanner.
  var s scanner.Scanner
  fset := token.NewFileSet()                      // positions are relative to fset
  file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
  s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

  // Repeated calls to Scan yield the token sequence found in the input.
  for {
      _, tok, lit := s.Scan()
      if tok == token.EOF {
          break
      }
      if tok == token.IDENT {
        fmt.Println(lit)
      }
  }
}
