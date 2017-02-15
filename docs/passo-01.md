---
permalink: passo-01
---

# Runas, passo 1: iniciando o TDD

Vamos usar [TDD](http://tdd.caelum.com.br/) para desenvolver esse projeto. A primeira coisa então é escrever um teste, no arquivo-fonte `runefinder_test.go`. Nosso primeiro teste verifica a função `AnalisarLinha`, que deve extrair um caractere e um nome de uma linha do `UnicodeData.txt`:

```go
package main // ➊

import "testing" // ➋

const linhaLetraA = `0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;` // ➌

func TestAnalisarLinha(t *testing.T) { // ➍
	runa, _ := AnalisarLinha(linhaLetraA) // ➎
	if runa != 'A' {                      // ➏
		t.Errorf("Esperava 'A', veio %c", runa) // ➐
	}
}
```

Vejamos o que temos aqui:

➊ Todo arquivo-fonte em Go precisa declarar o pacote ao qual ele pertence. Para programas executáveis, o pacote deve ser `main`. Bibliotecas devem usar um nome igual à última parte do caminho até seu código-fonte, ex. `runas`.

➋ Importamos o pacote `testing` da biblioteca padrão.

➌ Definimos uma constante do tipo `string` (não é preciso declarar o tipo, porque o compilador identifica que o valor à direita do `=` é uma `string`).

➍ Todas as funções de teste precisam começar com o prefixo `Test`, e recebem como argumento um ponteiro para um objeto `testing.T`, através do qual acessamos métodos como `t.Errorf` neste exemplo. As declarações em Go tem a forma `x tipo`, onde `x` é o identificador sendo declarado, seguido de seu `tipo` (como em Pascal!)

➎ Nossa primeira função devolverá dois valores: um caractere e uma string — que por enquanto vamos ignorar (mais detalhes sobre essa linha a seguir).

➏ Note a ausência de parênteses ao redor da condição, e o uso de aspas simples para indicar que `'A'` é um caractere e não uma `string`.

➐ Um dos métodos para reportar erros em testes é `t.Errorf`. Em Go, funções que terminam com a letra `f` normalmente aceitam strings com códigos de formatação, semelhante à função `printf` em C.

A linha ➎ traz algumas peculiaridades da linguagem Go:

* Os caracteres Unicode em Go são chamados de "runas", e o tipo de dado usado para representar um caractere é `rune`. Assim como em C, um caractere é na verdade um número, que pode ser exibido como um caractere na saída se usarmos o código de formatação `"%c"`, como fizemos na linha ➐. Em Go, o tipo `rune` é o mesmo que `int32`, mas usamos `rune` para deixar claro quando estamos lidando com o código de um caractere, e não um número qualquer.

* Go permite que uma função devolva mais de um valor, e esses valores são atribuídos de uma vez só a suas respectivas variáveis. O compilador recusa variáveis que não serão usadas, então se você precisa ignorar um valor devolvido por uma função, use o nome especial `_`, o chamado _identificador vazio_ ([blank identifier](https://golang.org/doc/effective_go.html#blank)).

* Em Go nem sempre é necessário declarar o tipo das variáveis porque em muitas situações podemos usar o sinal `:=` para fazer uma _atribuição curta_ ([short assignment](https://tour.golang.org/basics/10)). Isso funciona somente na primeira vez que uma variável aparece dentro de um escopo. O compilador cria cada variável com o tipo compatível com a expressão à direita do `:=`. No caso, os tipos `rune` e `string` dos resultados que `AnalisarLinha` vai devolver.

## Rodando o teste

Após criar o arquivo `runefinder_test.go`, você pode usar o comando `go test` para executar o teste. O resultado será este:

```bash
$ go test
# github.com/ThoughtWorksInc/runas
./runefinder_test.go:8: undefined: AnalisarLinha
FAIL	github.com/ThoughtWorksInc/runas [build failed]
```

Obviamente, falta definir a função `AnalisarLinha` em algum lugar. Vamos lá.


## Primeira função: `AnalisarLinha`

Vamos criar outro arquivo-fonte, com nome `runefinder.go`. O mínimo que precisamos para fazer passar o teste é isso:

```go
package main // ➊

func AnalisarLinha(ucdLine string) (rune, string) { // ➋
	return 'A', "LETRA A" // ➌
}
```

Vejamos o que temos em `runefinder.go`:

➊ Novamente declaramos o mesmo pacote, assim os identificadores deste arquivo ficam acessíveis para outros arquivos do mesmo pacote.

➋ Esta função é pública porque seu nome começa com uma letra maiúscula. Isso não é apenas uma convenção, é algo que o compilador verifica. Esta linha também declara que a função recebe um argumento do tipo `string` e devolve um par de resultados, respectivamente `rune` e `string`.

➌ Apenas para fazer passar o teste, devolvemos uma `rune` e uma `string`.

Feito isso, rodamos o teste:

```bash
$ go test
PASS
ok  	github.com/ThoughtWorksInc/runas	0.013s
```

Sucesso!

OK, completamos nosso primeiro _baby step_: fizemos a função mais simples possível que faz o teste passar. A função `AnalisarLinha` ignora totalmente o argumento passado, e o teste confere a runa mas ignora a string devolvida. Mas um teste passando é um bom começo!

Hora de ler o [Passo 2](passo-02), código no diretório `runas02`.
