---
permalink: passo-07
---

# Runas, passo 7 (bônus): download da UCD

Nosso programa `runas` depende da presença do arquivo `UnicodeData.txt` no diretório atual para funcionar. Neste passo, vamos criar uma função para baixar o arquivo direto do site `unicode.org`, caso ele não esteja presente em um caminho local configurado pelo usuário.

Antes de mais nada, vamos verificar que temos uma versão funcional de `runas`, após o exercício do `passo-06`.

```bash
$ go test
PASS
ok  	github.com/ThoughtWorksInc/runas	0.109s
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

A configuração do caminho local será feita com uma variável de ambiente, `UCD_PATH`. Se esta variável não existir, o programa usará o caminho do diretório "home" do usuário, por exemplo, `/home/luciano/UnicodeData.txt` em uma máquina GNU/Linux.

Para começar, vamos criar a função que determina o caminho para salvar o `UnicodeData.txt`, começando pelo teste simulando o caso de existir a variável de ambiente `UCD_PATH`:

```go
func TestObterCaminhoUCD_setado(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH") // ➊
	defer restaurar("UCD_PATH", caminhoAntes, existia) // ➋
	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano()) //➌
	os.Setenv("UCD_PATH", caminhoUCD) // ➍
	obtido := obterCaminhoUCD() // ➎
	if obtido != caminhoUCD {
		t.Errorf("obterCaminhoUCD() [setado]\nesperado: %q; recebido: %q", caminhoUCD, obtido)
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

Não tem nenhuma graça esta função. Nem faria sentido o teste anterior: na prática estamos testando só a função `os.GetEnv`, e ao escrever testes automatizados devemos acreditar que as bibliotecas que são nossas dependências funcionam, e não testá-las. Mas este teste faz sentido junto com o próximo teste, que verifica o caso contrário: quando não existe a variável de ambiente `UCD_PATH`, ou ela está vazia. Para satisfazer esses dois testes, `obterCaminhoUCD` vai ter que ficar mais esperta.

```go
func TestObterCaminhoUCD_default(t *testing.T) {
	caminhoAntes, existia := os.LookupEnv("UCD_PATH")
	defer restaurar("UCD_PATH", caminhoAntes, existia)
	os.Unsetenv("UCD_PATH") // ➊
	sufixoCaminhoUCD := "/UnicodeData.txt"  // ➋
	obtido := obterCaminhoUCD()
	if !strings.HasSuffix(obtido, sufixoCaminhoUCD) { // ➌
		t.Errorf("obterCaminhoUCD() [default]\nesperado (sufixo): %q; recebido: %q", sufixoCaminhoUCD, obtido)
	}
}
```

➊ Depois de copiar seu valor, removemos a variável de ambiente `UCD_PATH`.

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

Nesta etapa faremos várias operações com o SO que podem gerar erros. Em vez de colocar testes `if err != nil` por toda parte, num exemplo didático como este vamos usar essa função `check` para verificar se houve erro e terminar o programa com `panic`:

```go
func check(e error) {
	if e != nil {
		panic(e)
	}
}
```

Se o programa fosse um serviço que precisa ficar no ar 24x7, `check` seria uma péssima maneira de tratar erros. Mas em uma ferramenta como `runas`, é um atalho razoável.


## O programa principal e a função que abre aquivo UCD local

Uma vez que temos `obterCaminhoUCD`, veja como fica nossa função `main`:

```go
func main() {
	ucd, err := abrirUCD(obterCaminhoUCD()) // ➊
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { ucd.Close() }()
	consulta := strings.Join(os.Args[1:], " ")
	Listar(ucd, strings.ToUpper(consulta))
}
```

➊ Em vez de `os.Open`, agora invocamos `abrirUCD` passando como argumento o caminho configurado pelo usuário, ou o default.

Não temos um teste unitário para `main`; ela é verificada pelos testes funcionais `Example`, `Example_consultaDuasPalavras` e `Example_consultaComHífenECampo10` que fizemos nos passos 5 e 6.

Vamos escrever os testes para `abrirUCD`, primeiro um teste que assume a existência do arquivo `UnicodeData.txt` no diretório local configurado:

```go
func TestAbrirUCD_local(t *testing.T) {
	caminhoUCD := obterCaminhoUCD()
	ucd, err := abrirUCD(caminhoUCD)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", caminhoUCD, err)
	}
	ucd.Close()
}
```

Supondo que existe o arquivo `UnicodeData.txt` no diretório configurado, esta versão super simples de `abrirUCD` satisfaz o teste anterior:

```go
func abrirUCD(caminho string) (*os.File, error) {
	ucd, err := os.Open(caminho)
	return ucd, err
}
```

## Teste com mock: baixar UCD

Se não existir o `UnicodeData.txt` no diretório configurado, teremos que baixá-lo do site `unicode.org`. Para isso, vamos criar uma função `baixarUCD`, cujo teste é assim:

```go
func TestBaixarUCD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc( // ➊
		func(w http.ResponseWriter, r *http.Request) { // ➋
			w.Write([]byte(linhas3Da43)) // ➌
		}))
	defer srv.Close() // ➍

	caminhoUCD := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	baixarUCD(srv.URL, caminhoUCD) // ➎
	ucd, err := os.Open(caminhoUCD) // ➏
	if os.IsNotExist(err) { // ➐
		t.Errorf("baixarUCD não gerou:%v\n%v", caminhoUCD, err)
	}
	ucd.Close() // ➑
	os.Remove(caminhoUCD) // ➒
}
```

➊ Para testar um cliente HTTP sem acessar a Web, usamos o pacote `httptest`, cuja função `NewServer` devolve um objeto que imita um servidor HTTP. O servidor é criado com um objeto `http.Handler`, que pode ser construído passando uma função para `http.HandlerFunc`.

➋ Usamos uma função anônima para tratar as requisições ao nosso servidor fajuto.

➌ Essa função apenas escreve no objeto `w` (que é um `http.ResponseWriter`) a string da constante `linhas3Da43`, convertida para uma fatia de `byte`.

➍ Agendamos o fechamento do servidor para o final do teste.

➎ Invocamos a função que queremos testar, `baixarUCD`, passando a URL do servidor fajuto e um caminho que inclui momento atual em nanossegundos, como já fizemos antes.

➏ Tentamos abrir o arquivo baixado no tal caminho.

➐ Se o arquivo não existe, reportamos erro no teste.

➑ Fechamos o arquivo baixado...

➒ ...e o removemos, para não deixar sujeira no diretório de trabalho.

Esse teste exige um certo esforço para codar, mas o que ele faz é bem legal: "mocar" um servidor HTTP (se podemos "codar", então podemos "mocar": uma gíria derivada de "to mock", que significa "imitar").

Agora vejamos o código de `baixarUCD` que satisfaz aquele teste:

```go
func baixarUCD(url, caminho string) {
	resposta, err := http.Get(url) // ➊
	check(err) // ➋
	defer resposta.Body.Close() // ➌
	arquivo, err := os.Create(caminho) // ➍
	check(err)
	defer arquivo.Close() // ➎
	_, err = io.Copy(arquivo, resposta.Body) // ➏
	check(err)
}
```

➊ Invocamos `http.Get` para baixar a UCD.

➋ Verificamos qualquer erro com `check`, encerrando o programa se for o caso. Vamos invocar `check` mais duas vezes nesta função.

➌ Usamos `defer` para fechar o corpo da resposta HTTP no final dessa função.

➍ Criamos um arquvo local no `caminho`, para salvar os bytes baixados.

➎ Se o arquivo foi criado com sucesso, usamos `defer` para fechá-lo no final dessa função.

➏ Invocamos `io.Copy` para copiar os bytes do corpo da resposta HTTP para o arquivo local (estranhamente, a ordem dos parâmetros é destino, origem). `io.Copy` devolve o número de bytes copiados (que ignoramos atribuindo a `_`) e um possível erro, que verificaremos com `check`.

Agora que já temos como baixar um arquivo UCD, podemos melhorar a função `abrirUCD`.


## Abrindo o UCD remoto

Novamente, usamos um teste que gera um caminho novo a cada vez, forçando `abrirUCD` a baixar o arquivo `UnicodeData.txt` sempre:

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

➌ De novo, o truque de gerar um caminho com o momento atual em nanossegundos, garantindo que cada execução desse teste vai gerar um novo caminho, obrigando `abrirUCD` a detectar a falta do arquivo `UnicodeData.txt` e baixá-lo.

> __DICA__: Usei a IDE Atom com o plug-in `go-plus` para editar este tutorial, e notei que `go-plus` executa os testes com a opção `-test.short` cada vez que salvo um arquivo-fonte. Atom com `go-plus` é uma ótima IDE para Go!

Segue a implementação de `abrirUCD`. Note a chamada para `baixarUCD`, que implementamos antes. Esta função depende da constante `URLUCD`. Ela deve ser criada no topo do arquivo `runefinder.go`, mas colocamos aqui para facilitar a leitura do tutorial.

```go
// URLUCD é a URL canônica do arquivo UnicodeData.txt mais atual
const URLUCD = "http://www.unicode.org/Public/UNIDATA/UnicodeData.txt"

func abrirUCD(caminho string) (*os.File, error) {
	ucd, err := os.Open(caminho)
	if os.IsNotExist(err) { // ➊
		fmt.Printf("%s não encontrado\nbaixando %s\n", caminho, URLUCD)
		baixarUCD(URLUCD, caminho) // ➋
		ucd, err = os.Open(caminho) // ➌
	}
	return ucd, err // ➍
}
```

➊ Verificamos se `os.Open` devolveu especificamente um erro de arquivo não existente. Neste caso...

➋ ...depois de avisar o usuário, invocamos `baixarUCD`, passando o caminho onde será salvo o arquivo.

➌ Tentamos abrir de novo o arquivo.

➍ Seja qual for o caminho percorrido em `abrirUCD`, no final devolvemos o arquivo e o erro.

Neste ponto temos um programa bastante funcional: `runas` sabe procurar o arquivo `UnicodeData.txt` no local configurado, e sabe baixá-lo da Web se necessário.

O único incômdo é que, durante o download, nada acontece durante alguns segundos após o programa informar que está baixando o arquivo. Na seção final vamos resolver esse problema usando os recursoss mais empolgantes de Go: gorrotinas e canais.


## Download concorrente com indicador de progresso

Vamos gerar continuamente uma sequência de pontos `.....` durante o download, evitando que o usuário suspeite que o programa travou. Para fazê-lo, usaremos alguns recursos especiais da linguagem Go: um canal (_channel_) e uma gorrotina (_goroutine_).

Na `abrirUCD`, acrescentamos três linhas:

```go
func abrirUCD(caminho string) (*os.File, error) {
	ucd, err := os.Open(caminho)
	if os.IsNotExist(err) {
		fmt.Printf("%s não encontrado\nbaixando %s\n", caminho, URLUCD)
		feito := make(chan bool) // ➊
		go baixarUCD(URLUCD, caminho, feito) // ➋
		progresso(feito) // ➌
		ucd, err = os.Open(caminho)
	}
	return ucd, err
}
```

➊ Construímos um canal do tipo `chan bool`, ou seja, um canal por onde vão trafegar valores booleanos. Um canal permite a comunicação e a sincronização entre gorrotinas, que são como threads leves gerenciadas pelo ambiente de execução da linguagem Go.

➋ A instrução `go` dispara uma função em uma nova gorrotina, permitindo que ela execute de forma concorrente. A partir desse ponto, nosso programa opera com duas gorrotinas: a gorrotina principal e a gorrotina que executa `baixarUCD`. Note que, além de `URLUCD` e `caminho`, estamos passando o canal `feito`.

➌ Invocamos a função `progresso`. Ela vai ficar em _loop_ gerando `....` na saída, até que receba pelo canal `feito` um sinal de que `baixarUCD` terminou o download.

Temos apenas duas mudanças em `baixarUCD`:

```go
func baixarUCD(url, caminho string, feito chan<- bool) { // ➊
	resposta, err := http.Get(url)
	check(err)
	defer resposta.Body.Close()
	arquivo, err := os.Create(caminho)
	check(err)
	defer arquivo.Close()
	_, err = io.Copy(arquivo, resposta.Body)
	check(err)
	feito <- true // ➋
}
```
➊ O terceiro parâmetro é `feito chan<- bool`. A notação `chan<-` indica que o canal `feito` apenas consome e não produz valores, dentro dessa função. Ou seja, `baixarUCD` só pode enviar valores para `feito`.

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

➊ Aqui a notação `<-chan` indica que, dentro de `progresso`, o canal `feito` apenas produz valores, mas não consome. Portanto `progresso` só pode receber valores do canal `feito`.

➋ Inciamos um laço infinito com `for`.

➌ `select` é uma instrução de controle de fluxo especial para programar sistemas concorrentes. Funciona como uma `switch` com vários blocos `case`, mas a seleção é baseada no estado do canal em cada caso. O bloco `case` do primeiro canal que estiver pronto para consumir ou produzir um valor será executado. Se mais de um `case` estiver pronto, Go seleciona um deles aleatoriamente.

➍ O bloco `case <-feito` será executado quando o canal `feito` estiver pronto para produzir um valor. A expressão `<-feito` lê e descarta o primeiro valor no canal (se quiséssemos usar o valor, poderíamos atribuí-lo a uma variável). Nesse caso não interessa o valor, pois estamos usando o canal somente para sincronização. A gorrotina de `baixarUCD` informa a gorrotina principal que terminou seu processamento enviando `true` pelo canal `feito`. Quando isso acontecer, este bloco vai exibir uma quebra de linha com `fmt.Println` e encerrar a função `progresso` com `return`.

➎ Em um `select`, o bloco `default` é acionado quando nenhum `case` está pronto para executar. Neste caso, se o canal `feito` não produziu uma mensagem, então geramos um `"."` na saída, e congelamos esta gorrotina por 150 milissegundos (do contrário apareceriam milhares de `.....` por segundo na saída).

Como temos o laço `for`, após cada execução do `default`, o `select` vai novamente verificar se o `case <-feito` está pronto para produzir um valor.

Vale notar que, quando uma instrução `select` não tem um `default`, ela bloqueia até que algum `case` esteja pronto para produzir ou consumir um valor. Mas com um `default`, `select` é uma estrutura de controle não bloqueante.

Agora você pode compilar o programa com o comando `go build` e obter um executável `runas` (porque este é o nome do diretório onde está o código-fonte do passo 7, na raiz do repositório).


## os.Exit(0) // Fim!

Isso conclui a nossa degustação da linguagem Go. Uma deGostação! 😬

Você pode rodar o comando `go test -cover` para executar os testes com uma medida de cobertura. Aqui estou obtendo 97.1% de cobertura, um bom número. Se testar com `go test -cover -test.short`, a cobertura cai para 80.9%, porque pulamos `TestAbrirUCD_remoto`.

Nosso objetivo era mostrar elementos da linguagem através de um exemplo simples porém útil, e ao mesmo tempo ilustrar algumas técnicas básicas de testes automatizados para praticar TDD em Go.

Agradecemos se você [mandar feedback](https://github.com/ThoughtWorksInc/runas/issues) com sugestões para melhorias. Por exemplo: como melhorar a cobertura de testes neste passo final? Não deixe também de postar suas dúvidas. Cada dúvida é um possível _bug_ deste tutorial, pois sempre é possível explicar melhor.

_Happy hacking!_
