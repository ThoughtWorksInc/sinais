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

➌ Criamos uma função anônima que vai restaurar o valor de `os.Args` no final da nossa função `Example`. Leia mais sobre o comando `defer` logo adiante.

➍ Mudamos os valor de `os.Args` para fazer o teste. Observe a sintaxe de uma fatia literal: primeiro o tipo `[]string`, depois os itens entre `{}`. O primeiro item de `os.Args` é o nome do programa (irrelevante para o nosso teste). O segundo item é a palavra que vamos buscar, `"cruzeiro"`, cuidadosamente escolhida porque só existe um caractere Unicode que contém essa palavra em seu nome.

➎ Invocamos a função `main`, a mesma que será chamada quando nosso programa for acionado na linha de comando.

O comando `defer` é uma inovação simples porém genial da linguagem Go. Ele serve para invocar uma função no final da função atual (`Example`). `defer` é útil para fechar arquivos, encerrar conexões, liberar travas, etc. É como se o corpo da função `Example` estivesse dentro de um `try/finally`, e as funções chamadas em `defer` seriam executadas no bloco `finally`, ou seja, após o `return` e mesmo que ocorram exceções. No exemplo, o uso de `defer` garante que o valor de `os.Args` será restaurado ao valor original, independente do sucesso ou fracasso do teste.

> _Nota_: Alterar uma variável global como `os.Args` pode ser perigoso em um
> sistema concorrente, mas Go só executa testes em paralelo se usamos o método
> [`T.Parallel`](https://golang.org/pkg/testing/#T.Parallel).

## A função `main`

Afinal vamos implementar a função `main`, que permite executar o `runefinder` como um programa direto da linha de comando.

```go
func main() { // ➊
	ucd, err := os.Open("UnicodeData.txt") // ➋
	if err != nil {                        // ➌
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer func() { ucd.Close() }()             // ➍
	consulta := strings.Join(os.Args[1:], " ") // ➎
	Listar(ucd, strings.ToUpper(consulta))     // ➏
}
```

➊ Em Go, a função `main` não recebe argumentos.

➋ Abrimos o arquivo "UnicodeData.txt", assumindo que ele está no diretório atual. A maioria das funções de E/S em go devolve dois resultados, e o segundo é do tipo `error`, uma interface usada para reportar erros. No caso de `os.Open`, o primeiro resultado é um `*File`, ponteiro para um objeto arquivo.

➌ Se `err` é diferente `nil`, houve erro em `os.Open`. Nesse caso vamos exibir a mensagem de erro e terminar o programa. Chamando `os.Exit`, as funções em `defer` não são executadas.

➍ Usamos `defer` para fechar o arquivo que abrimos em ➋.

➎ Montamos a string de consulta concatenando os argumentos. A notação `os.Args[1:]` devolve uma nova fatia formada pelos itens de índice 1 em diante, assim omitimos o nome do programa invocado, que fica em `os.Args[0]`. A função `strings.Join` monta uma string intercalando os itens da fatia com o segundo argumento, `" "` neste caso.

➏ Invocamos a função `Listar` com o arquivo `ucd` e a `consulta` convertida em caixa alta (porque os nomes na UCD aparecem assim).

Agora precisamos do arquivo `"UnicodeData.txt"` ([URL oficial](http://www.unicode.org/Public/UNIDATA/UnicodeData.txt)). Depois faremos o `runefinder` baixar este arquivo, se necessário, mas agora você precisa buscar e colocá-lo no diretório atual (onde está o `runefinder.go`). Feito isso, você pode rodar os testes:

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
ok  	github.com/labgo/runas-passo-a-passo	0.033s
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
$ go run runefinder.go hexagram  # I Ching!
...
$ go run runefinder.go roman
...
$ go run runefinder.go clock face
```
