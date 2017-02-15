---
permalink: passo-04
---

# Runas, passo 4: MVP 1, o mínimo que é útil

Estamos prontos para fazer a interface de linha de comando e ter um produto mínimo viável (MVP) que permite encontrar caracteres pelo nome.

Seguindo com o TDD, vamos fazer um teste para nosso programa principal. Como precisamos testar a saída padrão, usaremos outra função `Example`:

```go
func Example() { // ➊
	oldArgs := os.Args  // ➋
	defer func() { os.Args = oldArgs }()  // ➌
	os.Args = []string{"", "cruzeiro"}  // ➍
	main() // ➎
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}
```

Esse teste traz várias novidades:

➊ A função chamada simplesmente `Example` é o exemplo do pacote: um teste funcional para o pacote como um todo ([documentação sobre funções exemplo](https://golang.org/pkg/testing/#hdr-Examples)).

➋ Vamos simular a passagem de argumentos pela linha de comando. O primeiro passo é copiar os argumentos de `os.Args` para `oldArgs`, assim poderemos restaurá-los depois. Para acessar `os.Args`, não esqueça de incluir o pacote `os` na declaração `import` de `rundefinder_test.go`.

➌ Criamos uma função anônima que vai restaurar o valor de `os.Args` no final da nossa função `Example`. Leia mais sobre a instrução `defer` logo adiante.

➍ Mudamos os valor de `os.Args` para fazer o teste. Observe a sintaxe de uma fatia literal: primeiro o tipo `[]string`, depois os itens entre `{}`. O primeiro item de `os.Args` é o nome do programa (irrelevante para o nosso teste). O segundo item é a palavra que vamos buscar, `"cruzeiro"`, cuidadosamente escolhida porque só existe um caractere Unicode que contém essa palavra em seu nome.

➎ Invocamos a função `main`, a mesma que será chamada quando nosso programa for acionado na linha de comando. A saída que aparece aqui é o que nosso programa vai gerar quando alguém buscar um caractere com a palavra "cruzeiro".

A instrução `defer` é uma inovação simples porém genial da linguagem Go. Ela serve para invocar uma função no final da execução da função atual (`Example`). `defer` é útil para fechar arquivos, encerrar conexões, liberar mutexes, etc. É como se o corpo da função `Example` estivesse dentro de um `try/finally` de Java ou Python, e as funções chamadas em `defer` seriam executadas no bloco `finally`, ou seja, após o `return` e mesmo que ocorram exceções. No exemplo, o uso de `defer` garante que o valor de `os.Args` será restaurado ao valor original, independente do sucesso ou fracasso do teste.

> __Nota__: Alterar uma variável global como `os.Args` pode produzir resultados
> inesperados em um sistema concorrente, mas Go só executa testes em paralelo se
> usamos o método [`T.Parallel`](https://golang.org/pkg/testing/#T.Parallel).


## A função `main`

Afinal vamos implementar a função `main`, que permite executar o `runefinder` como um programa direto da linha de comando.

```go
func main() { // ➊
	ucd, err := os.Open("UnicodeData.txt") // ➋
	if err != nil {                        // ➌
		log.Fatal(err.Error()) // ➍
	}
	defer func() { ucd.Close() }()             // ➎
	consulta := strings.Join(os.Args[1:], " ") // ➏
	Listar(ucd, strings.ToUpper(consulta))     // ➐
}
```

➊ Em Go, a função `main` não recebe argumentos.

➋ Abrimos o arquivo "UnicodeData.txt", assumindo que ele está no diretório atual. A maioria das funções de E/S em Go devolve dois resultados, e o segundo é do tipo `error`, uma interface usada para reportar erros. No caso de `os.Open`, o primeiro resultado é um `*File`, ponteiro para um objeto arquivo.

➌ Se `err` é diferente `nil`, houve erro em `os.Open`. Nesse caso vamos exibir a mensagem de erro e terminar o programa.

➍ A função `log.Fatal` faz duas coisas: exibe a mensagem passada como argumento e invoca `os.Exit(1)`, encerrando o programa. O tipo `error` tem o método `Error()` que devolve uma string com a mensagem de erro. Chamando `log.Fatal` (ou `os.Exit`), as funções em `defer` não são executadas.

➎ Usamos `defer` para fechar o arquivo que abrimos em ➋.

➏ Montamos a string de consulta concatenando os argumentos. A notação `os.Args[1:]` lembra Python ou Ruby: ela devolve uma nova fatia formada pelos itens de índice 1 em diante. Assim omitimos o nome do programa invocado, que fica em `os.Args[0]`. A função `strings.Join` monta uma string intercalando os itens da fatia com o segundo argumento, `" "` neste caso.

➐ Invocamos a função `Listar` com o arquivo `ucd` e a `consulta` convertida em caixa alta (porque na UCD os nomes aparecem em maiúsculas).

Agora você precisa baixar o arquivo `UnicodeData.txt` ([URL oficial](http://www.unicode.org/Public/UNIDATA/UnicodeData.txt)). Depois faremos o `runefinder` baixar este arquivo, se necessário, mas agora você precisa buscar e colocá-lo no diretório atual (onde está o `runefinder.go`). Feito isso, você pode rodar os testes:

```bash
$ go test -v
=== RUN   TestAnalisarLinha
--- PASS: TestAnalisarLinha (0.00s)
=== RUN   ExampleListar
--- PASS: ExampleListar (0.00s)
=== RUN   ExampleListar_doisResultados
--- PASS: ExampleListar_doisResultados (0.00s)
=== RUN   Example
--- PASS: Example (0.02s)
PASS
ok  	github.com/ThoughtWorksInc/runas	0.033s
```

## Experimentando o `runefinder`

Agora já é possível brincar com o programa na linha de comando. A forma mais simples de experimentar um programa Go em desenvolvimento é o comando `go run`. Veja como funciona:

```bash
$ go run runefinder.go cat face
U+1F431	🐱	CAT FACE
U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
U+1F639	😹	CAT FACE WITH TEARS OF JOY
U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
U+1F63C	😼	CAT FACE WITH WRY SMILE
U+1F63D	😽	KISSING CAT FACE WITH CLOSED EYES
U+1F63E	😾	POUTING CAT FACE
U+1F63F	😿	CRYING CAT FACE
U+1F640	🙀	WEARY CAT FACE
```

Experimente fazer várias consultas, tem muita coisa interessante no banco de dados Unicode. Alguns exemplos para você experimentar:

```bash
$ go run runefinder.go chess
...
$ go run runefinder.go runic
...
$ go run runefinder.go roman
...
$ go run runefinder.go clock face
...
$ go run runefinder.go alchemical
...
$ go run runefinder.go hexagram  # I Ching!
```

Outra forma de usar o programa é gerar um executável, com o comando `go build`, assim:

```bash
$ go build
$ ls -lah runas
-rwxr-xr-x  1 lramalho  staff   1.9M Dec 28 20:10 runas04
```

Se der tudo certo, o comando `go build` não exibe nenhuma mensagem. Mas ele produz um binário executável com o nome do projeto, no caso `runas04` (que é o nome do diretório onde está o projeto). Por convenção, o nome do projeto é o nome do repositório, mas neste tutorial temos na verdade vários projetos, um em cada diretório `runasNN`. Note o executável de 1.9MB no `ls` acima.

Para rodar o binário, é só rodar!

```bash
$ ./runas04 flag
U+2690	⚐	WHITE FLAG
U+2691	⚑	BLACK FLAG
U+26F3	⛳	FLAG IN HOLE
U+26FF	⛿	WHITE FLAG WITH HORIZONbTAL MIDDLE BLACK STRIPE
U+1D16E	𝅮	MUSICAL SYMBOL COMBINING FLAG-1
U+1D16F	𝅯	MUSICAL SYMBOL COMBINING FLAG-2
U+1D170	𝅰	MUSICAL SYMBOL COMBINING FLAG-3
U+1D171	𝅱	MUSICAL SYMBOL COMBINING FLAG-4
U+1D172	𝅲	MUSICAL SYMBOL COMBINING FLAG-5
U+1F38C	🎌	CROSSED FLAGS
U+1F3C1	🏁	CHEQUERED FLAG
U+1F3F3	🏳	WAVING WHITE FLAG
U+1F3F4	🏴	WAVING BLACK FLAG
U+1F4EA	📪	CLOSED MAILBOX WITH LOWERED FLAG
U+1F4EB	📫	CLOSED MAILBOX WITH RAISED FLAG
U+1F4EC	📬	OPEN MAILBOX WITH RAISED FLAG
U+1F4ED	📭	OPEN MAILBOX WITH LOWERED FLAG
U+1F6A9	🚩	TRIANGULAR FLAG ON POST
```

## Próximos passos

Esse foi o nosso MVP1, a primeira versão usável do programa.

Ele tem algumas limitações que resolveremos nos próximos passos:

* Só funciona na presença do arquivo `UnicodeData.txt`. O ideal é que, se o arquivo não está presente, o programa baixe-o direto do site `unicode.org`.
* Nossa busca por substring é bem tosca. Por exemplo, se você busca "cat", todos os caracteres que têm essa sequência de letras em qualquer parte do nome serão exibidos, e a maioria deles não tem nada a ver com gatinhos. Seria mais legal fazer a busca por palavras inteiras.
* Também seria bom ignorar a ordem das palavras, assim as pesquisas "chess black" e "black chess" devolveriam os mesmos resultados.
* Seria legal exibir no final o número de caracteres encontrados. Isso é útil principalmente quando não vem nenhum ou quando vem centenas.
* Se você não passar nenhum argumento, todos os caracteres do UCD serão exibidos, veja só:

```bash
$ ./runas | wc
   30593  182344 1181756
```

Vamos deixar a questão do download do UCD para o final deste tutorial. Primeiro vamos melhorar a busca no [Passo 5](passo-05).
