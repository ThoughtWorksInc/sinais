# Runas: passo 1

Vamos usar TDD para desenvolver esse projeto. A primeira coisa então é escrever um teste. Eis o estado inicial do arquivo-fonte `runefinder_test.go`:

```go
package runefinder // <1>

import "testing" // <2>

const linhaLetraA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;" // <3>

func TestAnalisarLinha(t *testing.T) { // <4>
	runa, _ := AnalisarLinha(linhaLetraA) // <5>
	if runa != 'A' {                      // <6>
		t.Errorf("Esperava 'A', veio %c", runa) // <7>
	}
}
```

Vejamos o que temos aqui:

<1> Todo arquivo-fonte em Go precisa declarar o pacote ao qual ele pertence.
<2> Importamos o pacote `testing` da biblioteca padrão.
<3> Definimos uma constante do tipo `string` (não é preciso declarar o tipo, explicitamente, o compilador sabe que o que tem à direita do `=` é uma `string`).
<4> Todas as funções de teste precisam começar com o prefixo `Test`, e recebem como argumento um ponteiro para o objeto `testing.T`, através do qual acessamos os métodos, como por exemplo `t.Errorf`. As declarações em Go tem a forma `x tipo`, onde `x` é o identificador sendo declarado, seguido de seu `tipo` (como em Pascal!)
<5> Nossa primeira função devolverá dois valores: um caractere e uma string - que por enquanto vamos ignorar (mais detalhes sobre essa linha a seguir).
<6> Note a ausência de parentesis ao redor da condição, e o uso de aspas simples para indicar que `'A'` é um caractere e não uma `string`.
<7> Um dos métodos para reportar erros em testes é `t.Errorf`.

A linha <5> traz algumas novidades peculiares da linguagem Go:

* Os caracteres Unicode em Go são chamados de "runas", e o tipo de dado usado para representar um caractere é `rune` (na verdade, o mesmo que `int32`). Assim como em C, um caractere é na verdade um número, que pode ser exibido como um caractere na saída se usarmos o código de formatação `"%c"`, como fizemos na linha <7>.

* Assim como Python, Go permite que uma função devolva mais de um valor, e esses valores são atribuídos de uma vez só a suas respectivas variáveis. O compilador recusa variáveis que não serão usadas, então se você precisa ignorar um valor devolvido por uma função, use o nome especial `_`, o chamado _identificador vazio_ ([blank identifier](https://golang.org/doc/effective_go.html#blank)).

* Em vez de declarar variáveis, em Go você pode usar uma _atribuição curta_ ([short assignment](https://tour.golang.org/basics/10)), com o sinal `:=`. Isso funciona somente na primeira vez que uma variável aparece dentro de um escopo. O compilador cria cada variável com o tipo compatível com a expressão à direita do `:=`. No caso, os tipos `rune` e `string` dos resultados que `AnalisarLinha` vai devolver.

## Rodando o teste

Após criar o arquivo `runefinder_test.go`, você pode usar o comando `go test` para executar o teste. O resultado será este:

```bash
$ go test
# github.com/labgo/runas-passo-a-passo
./runefinder_test.go:8: undefined: AnalisarLinha
FAIL	github.com/labgo/runas-passo-a-passo [build failed]
```

Obviamente, falta definir a função `AnalisarLinha` em algum lugar. Vamos lá.


## Primeira função: `AnalisarLinha`

Vamos criar outro arquivo-fonte, com nome `ucdlib.go` -- nossa biblioteca para lidar com a UCD (Unicode Character Database). O mínimo que precisamos para fazer passar o teste é isso:

```go
package runefinder // <1>

func AnalisarLinha(ucdLine string) (rune, string) { // <2>
	return 'A', "LETRA A" // <3>
}
```

Vejamos o que temos em `ucdlib.go`:

<1> Novamente declaramos o mesmo pacote, assim os identificadores públicos deste arquivo ficam acessíveis para outros arquivos do mesmo pacote.
<2> Esta função é pública porque seu nome começa com uma letra maiúscula. Isso não é apenas uma convenção, é algo que o compilador verifica. Esta linha também declara que a função recebe um argumento do tipo `string` e devolve um par de resultados, respectivamente `rune` e `string`.
<3> Apenas para fazer passar o teste, devolvemos uma `rune` e uma `string`.

Feito isso, rodamos o teste:

```bash
$ go test
PASS
ok  	github.com/labgo/runas-passo-a-passo	0.013s
```

Sucesso!

OK, completamos nosso primeiro _baby step_: fizemos a função mais simples possível que faz o teste passar. A função `AnalisarLinha` ignora totalmente o argumento passado, e o teste confere a runa mas ignora a string devolvida. Mas um teste passando é um bom começo!

Hora de mudar para o _branch_ `passo-02` e ler o arquivo `passo-02.md`.
