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

Para começar, vamos criar a função que determina o caminho para salvar o `UnicodeData.txt`, começando pelo teste simulando o caso de existir a variável de ambiente `UCD_PATH`:

```go
func TestObterCaminhoUCD_setado(t *testing.T) {
	caminhoAntes := os.Getenv("UCD_PATH") // ➊
	defer func() { os.Setenv("UCD_PATH", caminhoAntes) }() // ➋
	UCDPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	os.Setenv("UCD_PATH", UCDPath) // ➍
	obtido := obterCaminhoUCD() // ➎
	if obtido != UCDPath {
		t.Errorf("obterUCDPath() [setado]\nesperado: %q; recebido: %q", UCDPath, obtido)
	}
}
```

➊ Obtemos o valor da variável de ambiente `UCD_PATH` e guardamos para restaurar depois.

➋ Usamos `defer` para restaurar no final do teste o valor inicial de `UCD_PATH`.

➌ Geramos um caminho contendo o momento atual em nanossegundos, assim a cada excetução do teste o caminho será diferente.

➍ Colocamos o caminho gerado na variável de ambiente.

➎ Invocamos a função que queremos testar: `obterCaminhoUCD` deve obter o caminho que acabamos de colocar na variável de ambiente.

Essa é a implementação mínima de `obterCaminhoUCD` que faria o teste acima passar:

```go
func obterCaminhoUCD() string {
	return os.Getenv("UCD_PATH")
}
```

Não tem muita graça esta função, e nem faria sentido o teste anterior: na prática estamos testando só função `os.Getenv`, e ao escrever testes automatizados devemos acreditar que as bibliotecas que são nossas dependências funcionam. Mas este teste faz sentido com o próximo, que verifica o caso contrário: quando não existe a variável de ambiente `UCD_PATH`, ou ela está vazia.

```go
func TestObterCaminhoUCD_default(t *testing.T) {
	caminhoAntes := os.Getenv("UCD_PATH")
	defer func() { os.Setenv("UCD_PATH", caminhoAntes) }()
	os.Unsetenv("UCD_PATH") // ➊
	sufixoUCDPath := "/UnicodeData.txt"  // ➋
	obtido := obterCaminhoUCD()
	if !strings.HasSuffix(obtido, sufixoUCDPath) { // ➌
		t.Errorf("obterUCDPath() [default]\nesperado (sufixo): %q; recebido: %q", sufixoUCDPath, obtido)
	}
}
```

➊ Depois de preservar seu valor, removemos a variável de ambiente `UCD_PATH`.

➋ Para não complicar demais o teste, vamos apenas checar que o caminho termina com o nome do arquivo que esperamos.

➌ `strings.HasSuffix` serve para testar se uma string termina com um dado sufixo.

Para fazer esse teste passar, precisamos de mais algumas linhas em `obterCaminhoUCD`:

```go
func obterCaminhoUCD() string {
	caminhoUCD := os.Getenv("UCD_PATH")
	if caminhoUCD == "" { // ➊
		usuário, err := user.Current() // ➋
		check(err) // ➌
		caminhoUCD = usuário.HomeDir + "/UnicodeData.txt" // ➍
	}
	return caminhoUCD
}
```

➊ Se a variável de ambiente `UCD_PATH` está vazia ou não existe (nos dois casos, `os.Getenv` devolve `""`)...

➋ ...invocamos `user.Current` para obter informações sobre o usuário logado.

➌ A função `check` é uma forma rápida e preguiçosa de lidar com erros. Em seguida falaremos sobre ela.

➍ Construímos o `caminnhoUCD` concatenando o nome do arquivo ao caminho do diretório _home_ do usuário, ex. `/home/luciano/UnicodeData.txt` no meu caso.
