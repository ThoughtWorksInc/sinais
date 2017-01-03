---
permalink: passo-07
---

# Runas, passo 7 (bônus): download da UCD

Nosso programa `runas` depende da presença do arquivo `UnicodeData.txt` no diretório atual para funcionar. Neste passo, vamos criar uma função para baixar o arquivo direto do site `unicode.org`, caso ele não esteja presente em um caminho local configurado pelo usuário.

Antes de mais nada, vamos verificar que temos uma versão funcional de `runas`, após o exercício do `passo-06`.

 ```bash
$ go test
PASS
ok  	github.com/labgo/runas	0.109s
$ go run runefinder.go minus hyphen
U+002D	-	HYPHEN-MINUS
U+207B	⁻	SUPERSCRIPT MINUS (SUPERSCRIPT HYPHEN-MINUS)
U+208B	₋	SUBSCRIPT MINUS (SUBSCRIPT HYPHEN-MINUS)
U+FE63	﹣	SMALL HYPHEN-MINUS
U+FF0D	－	FULLWIDTH HYPHEN-MINUS
U+E002D		TAG HYPHEN-MINUS
```

## Configuração do caminho local de `Unicode.txt`

A função `main` que fizemos no `passo-04` ficou assim:

```go
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
```

Vamos trocar a chamada `os.Open` por uma função nossa, `abrirUCD`, que vai tentar abrir o arquivo em um caminho local configurado e, caso não encontre, vai baixar o arquivo do site `unicode.org`.

A configuração do caminho local será feita com uma variável de ambiente, `UCD_PATH`. Se esta variável não estiver presente, o programa usará o caminho do diretório "home" do usuário, por exemplo, `/home/luciano` em um GNU/Linux.

Para começar, vamos criar a função que determina o caminho para salvar o `UnicodeData.txt`, começando pelo teste:
