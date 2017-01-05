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

Neste passo final do tutorial faremos o seguinte:

* Configuração do caminho local de `Unicode.txt`.
* Função que abre o `Unicode.txt`, depois de baixá-lo da Web se não for encontrado.
* Download concorrente com indicador de progresso.


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

A configuração do caminho local será feita com uma variável de ambiente, `UCD_PATH`. Se esta variável não estiver presente, o programa usará o caminho do diretório "home" do usuário, por exemplo, `/home/luciano/UnicodeData.txt` em uma máquina GNU/Linux.

Para começar, vamos criar a função que determina o caminho para salvar o `UnicodeData.txt`, começando pelo teste simulando o caso de existir a variável de ambiente `UCD_PATH`:

```go
func TestObterCaminhoUCD_setado(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH") // ➊
	defer restaurar("UCD_PATH", caminhoAntes, existia) // ➋
	UCDPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) //➌
	os.Setenv("UCD_PATH", UCDPath) // ➍
	obtido := obterCaminhoUCD() // ➎
	if obtido != UCDPath {
		t.Errorf("obterUCDPath() [setado]\nesperado: %q; recebido: %q", UCDPath, obtido)
	}
}
```

➊ Obtemos o estado da variável de ambiente `UCD_PATH` e guardamos para restaurar depois. `os.LookupEnv` devolve o valor da variável e `true` se ela existe, ou uma string vazia e `false` se ela não existe.

➋ Usamos `defer` para restaurar no final do teste o estado inicial de `UCD_PATH`. Veremos a seguir o código de `restaurar`.

➌ Geramos um caminho contendo o momento atual em nanossegundos, assim a cada execução do teste o caminho será diferente.

➍ Colocamos o caminho gerado na variável de ambiente.

➎ Invocamos a função que queremos testar: `obterCaminhoUCD` deve obter o caminho que acabamos de colocar na variável de ambiente.

A função `restaurar` é bem simples. Se a variável em questão existia, ela recebe o valor passado. Se ela não existia, ela é removida com `os.Unsetenv`.

```go
func restaurar(nomeVar, valor string, existia bool) {
	if existia {
		os.Setenv(nomeVar, valor)
	} else {
		os.Unsetenv(nomeVar)
	}
}
```

Essa é a implementação mínima de `obterCaminhoUCD` que faz o teste acima passar:

```go
func obterCaminhoUCD() string {
	return os.Getenv("UCD_PATH")
}
```

Não tem nenhuma graça esta função. Nem faria sentido o teste anterior: na prática estamos testando só a função `os.Getenv`, e ao escrever testes automatizados devemos acreditar que as bibliotecas que são nossas dependências funcionam. Mas este teste faz sentido junto com o próximo teste, que verifica o caso contrário: quando não existe a variável de ambiente `UCD_PATH`, ou ela está vazia.

```go
func TestObterCaminhoUCD_default(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH")
	defer restaurar("UCD_PATH", caminhoAntes, existia)
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

➊ Se a variável de ambiente `UCD_PATH` está vazia ou não existe (nos dois casos, `os.Getenv` devolve `""`), então...

➋ ...invocamos `user.Current` para obter informações sobre o usuário logado.

➌ A função `check` é uma forma rápida e preguiçosa de lidar com erros. Em seguida falaremos sobre ela.

➍ Construímos o `caminhoUCD` concatenando o nome do arquivo ao caminho do diretório _home_ do usuário, ex. `/home/luciano/UnicodeData.txt` no meu caso.

Nesta etapa faremos várias operações com o SO que podem gerar erros. Em vez de colocar testes `if err != nil` por toda parte, num exemplo simples como este vamos usar essa função `check` para verificar se houve erro e terminar o programa com `panic`:

```go
func check(e error) {
	if e != nil {
		panic(e)
	}
}
```

Se o programa fosse um serviço que precisa ficar no ar 24x7, `check` seria uma péssima maneira de tratar erros. Mas em uma ferramenta como `runas`, é um atalho razoável.

## O programa principal e a função que abre o UCD

Uma vez que temos `obterCaminhoUCD`, veja como fica nossa função `main`:

```go
func main() {
	ucd, err := abrirUCD(obterCaminhoUCD()) // ➊
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { ucd.Close() }()
	consulta := strings.Join(os.Args[1:], " ")
	Listar(ucd, strings.ToUpper(consulta))
}
```

➊ Em vez de `os.Open`, agora invocamos `abrirUCD` passando como argumento o caminho configurado pelo usuário, ou o default.

Não temos um teste unitário para `main`; ela é verificada pelos testes funcionais `Example`,  `Example_consultaDuasPalavras` e `Example_consultaComHífenECampo10` que fizemos nos passos 5 e 6.

Vamos escrever os testes para `abrirUCD`, primeiro um teste que assume a existência do arquivo `UnicodeData.txt` no diretório corrente:

```go
func TestAbrirUCD_local(t *testing.T) {
	caminhoUCD := "./UnicodeData.txt"
	ucd, err := abrirUCD(caminhoUCD)
	if err != nil { // ➊
		t.Errorf("AbrirUCD(%q):\n%v", caminhoUCD, err)
	}
	ucd.Close()
}
```

➌ Se `abrirUCD` não reportar erro, o teste falhou.

Supondo que existe o arquivo `UnicodeData.txt` no diretório corrente, esta versão super simples de `abrirUCD` satisfaz o teste anterior:

```go
func abrirUCD(caminhoUCD string) (*os.File, error) {
	ucd, err := os.Open(caminhoUCD)
	return ucd, err
}
```

Agora, um teste que gera um caminho novo a cada vez, forçando `abrirUCD` a baixar o arquivo `UnicodeData.txt` toda vez:

```go

func TestAbrirUCD_remoto(t *testing.T) {
	if testing.Short() {  // ➊
		t.Skip("teste ignorado [opção -test.short]") // ➋
	}
	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) // ➌
	ucd, err := abrirUCD(caminhoUCD)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", caminhoUCD, err)
	}
	ucd.Close()
	os.Remove(caminhoUCD)
}
```

➊ Como este teste envolve um download que pode levar alguns segundos, este `if` utiliza a função `testing.Short` para ver se o teste foi acionado com a opção `-test.short`, assim: `go test -test.short`.

➋ Se `-test.short` foi informado, então o método `t.Skip` reporta que esse teste foi pulado, mas somente se for usada a opção `-v`; do contrário, o teste é pulado silenciosamente.

➌ Novamente usamos a técnica de gerar um caminho com o momento atual em nanossegundos, garantindo que cada execução desse teste vai gerar um novo caminho, obrigando `abrirUCD` a detectar a falta do arquivo `UnicodeData.txt` e baixá-lo.

Segue a implementação de `abrirUCD`. Note a chamada para `baixarUCD`, que implementaremos em seguida.

```go
func abrirUCD(caminhoUCD string) (*os.File, error) {
	ucd, err := os.Open(caminhoUCD)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s não encontrado\nbaixando %s\n", caminhoUCD, URLUCD)
		baixarUCD(caminhoUCD) // ➋
		ucd, err = os.Open(caminhoUCD) // ➌
	}
	return ucd, err // ➍
}
```

➊ Verificamos se `os.Open` devolveu especificamente um erro de arquivo não existente. Neste caso...

➋ ...depois de informar o usuário, invocamos `baixarUCD`, passando o caminho onde será salvo o arquivo.

➌ Tentamos abrir de novo o arquivo.

➍ Seja qual for o caminho percorrido em `abrirUCD`, no final devolvemos o arquivo e o erro.

Agora vejamos o código de `baixarUCD`:

```go
// URLUCD é a URL canônica do arquivo UnicodeData.txt mais atual
const URLUCD = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

func baixarUCD(caminhoUCD string) {
	response, err := http.Get(URLUCD) // ➊
	check(err) // ➋
	defer response.Body.Close() // ➌
	file, err := os.Create(caminhoUCD) // ➍
	check(err)
	defer file.Close() // ➎
	_, err = io.Copy(file, response.Body) // ➏
	check(err)
}
```

➊ Invocamos `http.Get` para baixar a UCD, cuja URL está na constante `URLUCD` (que deve ser criada no topo do arquivo `runefinder.go`)

➋ Verificamos qualquer erro com `check`, encerrando o programa se for o caso. Vamos invocar `check` mais duas vezes nesta função.

➌ Usamos `defer` para fechar o corpo da resposta HTTP no final dessa função.

➍ Criamos um arquvo local no `caminhoUCD`, para salvar os bytes baixados.

➎ Se o arquivo foi criado com sucesso, usamos `defer` para fechá-lo no final dessa função.

➏ Invocamos `io.Copy` para copiar os bytes do corpo da resposta HTTP para o arquivo local (estranhamente, a ordem dos parâmetros é destino, origem). `io.Copy` devolve o número de bytes copiados (que ignoramos atribuindo a `_`) e um possível erro, que verificaremos com `check`.

Aqui deixamos de lado o TDD. Não é fácil estar adequadamente esta função, e talvez nem valha o esforço, porque ela não tem nenhum lógica elaborada: é simplesmente uma sequência de passos realizados um depois do outro.

Neste ponto temos um programa bastante funcional: `runas` sabe procurar o arquivo `UnicodeData.txt` no local configurado, e sabe baixá-lo da Web se necessário.

O único incômdo é que, durante o download, nada acontece durante alguns segundos após o programa informar que está baixando o arquivo. Na seção final vamos resolver esse problema usando os recursoss mais empolgantes de Go: gorrotinas e canais.


## Download concorrente com indicador de progresso

Vamos gerar continuamente uma sequência de pontos `.....` durante o download, evitando que o usuário suspeite que o programa travou. Para fazê-lo, usaremos alguns recursos especiais da linguagem Go.

Veja como fica a função `abrirUCD`:

```go
func abrirUCD(caminhoUCD string) (*os.File, error) {
	ucd, err := os.Open(caminhoUCD)
	if os.IsNotExist(err) {
		fmt.Printf("%s não encontrado\nbaixando %s\n", caminhoUCD, URLUCD)
		feito := make(chan bool) // ➊
		go baixarUCD(caminhoUCD, feito) // ➋
		progresso(feito) // ➌
		ucd, err = os.Open(caminhoUCD)
	}
	return ucd, err
}
```

➊ Construímos um _channel_ ou canal, do tipo `chan bool`, ou seja, um canal por onde vão trafegar valores booleanos. Um canal permite a comunicação e a sincronização entre gorrotinas, que são como threads leves gerenciadas pelo ambiente de execução da linguagem Go.

➋ O comando `go` dispara uma função em uma nova gorrotina, permitindo que ela execute de forma concorrente. A partir desse ponto, nosso programa opera com duas gorrotinas: a gorrotina principal e a gorrotina que executa `baixarUCD`. Note que, além do `caminhoUCD`, estamos passando o canal `feito`.

➌ Invocamos a função `progresso`. Ela vai ficar em _loop_ gerando `....` na saída, até que receba pelo canal `feito` um sinal de que `baixarUCD` terminou o download.

Temos apenas duas mudanças em `baixarUCD`:

```go
func baixarUCD(caminhoUCD string, feito chan<- bool) { // ➊
	response, err := http.Get(URLUCD)
	check(err)
	defer response.Body.Close()
	file, err := os.Create(caminhoUCD)
	check(err)
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	check(err)
	feito <- true // ➋
}
```
➊ O segundo parâmetro é `feito chan<- bool`. A notação `chan<-` indica que, dentro dessa função, o canal `feito` apenas consome e não produz valores.

➋ Uma vez terminado o download, enviamos para o canal `feito` o sinal `true`. Isso terminará a função `progresso`, como veremos a seguir.

```go
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
```

➊ Aqui a notação `<-chan` indica que, dentro de `progresso`, o canal `feito` apenas produz valores, mas não consome.

➋ Inciamos um laço infinito com `for`.

➌ `select` é uma comando de controle de fluxo especial para suportar com sistemas concorrentes. Funciona como um `switch` com vários blocos `case`, mas a seleção é baseada no estado do canal em cada caso. O `select` executa o bloco `case` do primeiro canal que estiver pronto para consumir ou produzir um valor.

➍ O bloco `case <-feito` será executado quando o canal `feito` estiver pronto para produzir um valor; isso só vai acontecer quando `feito` receber o valor `true` na última linha de `baixarUCD`. Dessa maneira a gorrotina auxiliar informa a gorrotina principal que terminou seu processamento. Neste caso, este bloco vai exibir uma quebra de linha com `fmt.Println` e encerrar a função `progresso` com `return`.

➎ Em um `select`, o bloco `default` é acionado quando nenhum `case` está pronto para executar. Neste caso, se o canal `feito` não produziu uma mensagem, então geramos um `"."` na saída, e congelamos esta gorrotina por 150 milissegundos (do contrário milhares de `.....` por segundo apareceriam na saída).

Como temos o laço `for`, após cada execução do `default`, o `select` vai novamente verificar se o `case <-feito` está pronto para produzir um valor.

Vale notar que, quando um `select` não tem um `default`, ele bloqueia até que algum `case` esteja pronto para produzir ou consumir um valor. Mas com um `default`, o `select` é uma estrutura de controle não bloqueante.

## os.Exit(0)

Isso conclui a nossa degustação da linguagem Go. Uma deGostação!

Nosso objetivo era mostrar elementos da linguagem através de um exemplo simples porém útil, e ao mesmo tempo ilustrar algumas técnicas básicas de testes automatizados em Go.

Agradecemos se você mandar feedback com sugestões para melhorias. Por exemplo: como melhorar a cobertura de testes neste passo final? Não deixe também de postar suas dúvidas, pois sempre é possível explicar melhor.

Happy hacking!
