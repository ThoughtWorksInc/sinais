# Runas, passo 3: gerar a listagem

Nosso objetivo agora é produzir a saída do programa, uma listagem com este formato:

```
U+0041→ A→  LATIN CAPITAL LETTER A
U+0042→ B→  LATIN CAPITAL LETTER B
U+0043→ C→  LATIN CAPITAL LETTER C
```

(acrescentei setas `→` para indicar as tabulações, que serão feitas com `"\t"`)

A maneira mais simples de conferir a saída padrão de um programa em Go é usar um exemplo: um tipo especial de teste, feito com uma função nomeada com o prefixo `Example`. Vamos criá-la no arquivo de testes `runefinder_test.go` assim:

```go
func ExampleListar() { // <1>
	texto := strings.NewReader(linhas3Da43) // <2>
	Listar(texto, "MARK")                   // <3>
	// Output: U+003F	?	QUESTION MARK
}
```

Observe:

<1> O nome da função tem que começar com `Example`; em seguida obrigatoriamente vem o nome da função a ser testada, no caso `Listar`.
<2> Nossa função listar receberá um argumento `io.Reader` (ao final, será o arquivo `UnicodeData.txt`). Para testar, construímos um buffer de leitura `strings.Reader` a partir da constante `linhas3Da43`, cujo conteúdo veremos abaixo.
<3> Aqui invocamos a função a testar, passando o buffer e a consulta, `"MARK"`. O comentário na linha final da função `ExampleListar` define o resultado esperado. O sistema de testes vai comparar o texto gerado pela função `Listar` na saída padrão com o que vier após a string `"Output: "` no comentário.

Neste exemplo, a função `Listar` produz uma listagem delimitada por tabs, portanto o comentário `// Output:` precisa ser escrito com tabs entre os campos `U+003F`, `?` e `QUESTION MARK`. Se você colocar espaços em vez de tabs entre esses campos, o teste não passará.

A constante `linhas3Da43` que usamos nesse teste é definida assim, no topo do arquivo de testes:

```go
const linhas3Da43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`
```

Note que usamos o sinal de crase (_grave accent_, Unicode U+0060) para delimitar uma string que tem múltiplas linhas. As quebras de linha farão parte do valor da constante.  

Seguindo a filosofia do TDD, vamos rodar os testes:

```bash
$ go test
# github.com/labgo/runas-passo-a-passo
./runefinder_test.go:33: undefined: Listar
FAIL	github.com/labgo/runas-passo-a-passo [build failed]
```

Claro que falhou porque ainda não escrevemos a função `Listar`. Vamos começar implementando essa função da forma mais simples possível, apenas para fazer o teste passar:

```go
func Listar(texto io.Reader, consulta string) { // <1>
	runa, nome := '?', "QUESTION MARK"            // <2>
	fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome) // <3>
}
```

Passo a passo:
<1> `Listar` recebe um `texto` do tipo `io.Reader` e uma `consulta` do tipo `string`.
<2> Por enquanto vamos chumbar aqui os valores esperados por nosso teste.
<3> Aqui usamos `fmt.Printf` com três _verbos de formatação_ indicados por `%`, explicados a seguir. Note os tabs `\t` e a quebra de linha `\n` para formatar a saída.

Verbos de formatação que usamos:

* `%04X` para exibir um inteiro (o valor de `runa`) em formato hexadeximal com 4 casas, preenchendo com zeros à esquerda.
* `%[1]c` para exibir a runa como caractere; o modificador `[1]` é para indicar que queremos usar o argumento 1 (`runa`) novamente nesta posição. Assim geramos três campos na saída usando apenas dois argumentos (`runa`, `nome`), porque usamos `runa` duas vezes.
* `%s` para exibir um valor `string`.

Na realidade, `io.Reader` é uma _interface_, o que significa que nossa função `Listar` aceita como primeiro argumento qualquer objeto que implemente o método `Read` conforme a [documentação](https://golang.org/pkg/io/#Reader). Isso facilita os testes.

Para podermos usar `fmt.Printf` e `io.Reader` na função `Listar`, temos que acrescentar esses dois pacotes à declaração `import` no arquivo `ucdlib.go`, mantendo a ordem alfabética como pede a boa educação na comunidade Go:

```go
import (
	"fmt"
	"io"
	"strconv"
	"strings"
)
```

Essa implementação marota de `Listar` é suficiente para fazer passar o teste `ExampleListar`. Veja o resultado de `go test` com a opção `-v` (saída verbosa):

```bash
$ go test -v
=== RUN   TestAnalisarLinha
--- PASS: TestAnalisarLinha (0.00s)
=== RUN   ExampleListar
--- PASS: ExampleListar (0.00s)
PASS
ok  	github.com/labgo/runas-passo-a-passo	0.001s
```

Agora vamos codar a lógica da função `Listar`.


## Implementando `Listar` de verdade

Antes de mais nada, criar outro teste para expor o problema da nossa função `Listar`, que simplesmente ignora os argumentos passados. Usaremos outra função exemplo:

```go
func ExampleListar_doisResultados() { // <1>
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "SIGN") // <2>
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}
```

<1> Como já vimos, a função exemplo sempre tem o prefixo `Example` seguido do nome da função a ser testada. Para fazer mais de um exemplo testando a mesma função `Listar`, coloque um `_` seguido de uma descrição iniciando com caixa baixa (obrigatoriamente). Por isso temos: `ExampleListar_doisResultados`.

<2> Agora vamos testar a palavra `"SIGN"`, que ocorre duas vezes em nosso texto de testes. Para verificar saídas de várias linhas, escreva `Output:` na primeira linha do comentário, e coloque as linhas esperadas nos comentários seguintes, sempre colocando um espaço após o sinal `//`. Esse espaço será ignorado no teste.

Note que a saída esperada neste teste é formatada em três colunas separadas por tab.

Se rodarmos os testes agora, veremos isto:

```bash
$ go test
--- FAIL: ExampleListar_doisResultados (0.00s)
got:
U+003F	?	QUESTION MARK
want:
U+003D	=	EQUALS SIGN
U+003E	>	GREATER-THAN SIGN
FAIL
exit status 1
FAIL	github.com/labgo/runas-passo-a-passo	0.001s
```

A palavra `got:` (recebido) indica a saída que foi produzida, e `want:` (desejado), a saída que era esperada. Naturalmente o teste não passou porque nossa função `Listar` mostra sempre o mesmo resultado. Vamos consertar isso, com este código que traz várias novidades:

```go
// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem o texto da consulta // <1>
func Listar(texto io.Reader, consulta string) {
	varredor := bufio.NewScanner(texto) // <2>
	for varredor.Scan() {               // <3>
		linha := varredor.Text()            // <4>
		if strings.TrimSpace(linha) == "" { // <5>
			continue
		}
		runa, nome := AnalisarLinha(linha)    // <6>
		if strings.Contains(nome, consulta) { // <7>
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
	}
}
```

O que temos de novo:

<1> Por convenção, funções exportadas (públicas) devem ser documentadas com um comentário logo acima. Mais detalhes sobre essa convenção na seção __Documentando funções__, no final deste passo.
<2> Para percorrer um `io.Reader` linha-a-linha, usamos a função `bufio.NewScanner`, que devolve um objeto `Scanner` ([documentação]https://golang.org/pkg/bufio/#NewScanner).
<3> Um dos métodos do tipo `Scanner` é `Scan`: ele avança o `Scanner` até a próxima quebra de linha, e devolve `true` enquanto não encontrar o final do texto, e enquanto não ocorrer um erro. Aqui usamos o resultado de `Scan` como condição para um laço `for` (não existe `while` na linguagem Go; o comando `for` pode ser usado como um `while` dessa forma).
<4> Cada vez que invocamos `Scan` podemos usar o método `Text()` para obter a linha que acabou de ser lida.
<5> Aqui retiramos os caracteres brancos (_whitespace_) à esquerda e à direita da `linha`; se o resultado for uma string vazia, usamos `continue` para iniciar a próxima volta do laço porque não há o que fazer.
<6> Passamos a linha para a função `AnalisarLinha`, que devolve a `runa` e seu `nome`.
<7> Se o `nome` contém a string `consulta`, então geramos uma linha na saída, no formato que já vimos anteriormente.

Para que essa função funcione, precisamos importar o pacote `bufio`. A declaração `import` ficará assim:

```go
import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)
```

Note que, por convenção, os pacotes importados devem aparecer em ordem alfabética. O utilitário `go fmt` ordena os pacotes automaticamente, e o `goimports` insere/remove pacotes da declaração `import` automaticamente para satisfazer o compilador.

Com isso, todos os testes passam novamente:

```bash
$ go test -v
=== RUN   TestAnalisarLinha
--- PASS: TestAnalisarLinha (0.00s)
=== RUN   ExampleListar
--- PASS: ExampleListar (0.00s)
=== RUN   ExampleListar_doisResultados
--- PASS: ExampleListar_doisResultados (0.00s)
PASS
ok  	github.com/labgo/runas-passo-a-passo	0.001s
```

## Documentando funções

Por convenção, funções exportadas (públicas) dever ser precedidas de um comentário, como vimos no exemplo da função `Listar`. Eis o cabeçalho desta função:

```go
// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem o texto da consulta // <1>
func Listar(texto io.Reader, consulta string) {
```

O comentário deve começar com uma frase completa iniciando com o nome da função e terminando com ponto final. A ferramenta `golint` aponta a falta de tais comentários.

O comando `go doc` exibe as assinaturas das funções de um pacote:

```bash
$ go doc
package runefinder // import "github.com/labgo/runas-passo-a-passo"

func AnalisarLinha(linha string) (rune, string)
func Listar(texto io.Reader, consulta string)
```

E se você informar o nome de uma função, `go doc` mostra sua assinatura e documentação, assim:

```bash
$ go doc listar
func Listar(texto io.Reader, consulta string)
    Listar exibe na saída padrão o código, a runa e o nome dos caracteres
    Unicode cujo nome contem o texto da consulta

```

Com isso terminamos o passo 3. Estamos próximos do MVP: o produto mínimo viável é implementado no _branch_ `passo-04`, com explicações em `passo-04.md`.
