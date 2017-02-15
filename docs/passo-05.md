---
permalink: passo-05
---

# Runas, passo 5: busca por palavras inteiras

A versão MVP1 do programa `runas` busca caracteres comparando uma substring do nome. Isso gera dois problemas:

* Resultados demais: pesquisando "cat" vêm 82 caracteres, sendo que a maioria não tem nada a ver com gatos, por exemplo "MULTIPLICATION SIGN".
* Resultados de menos: a ordem das palavras na consulta deveria ser ignorada: "chess black" e "black chess" deveriam devolver os mesmos resultados, e "cat smiling" deveria encontrar todos estes caracteres:

```
U+1F638 😸 	GRINNING CAT FACE WITH SMILING EYES
U+1F63A 😺 	SMILING CAT FACE WITH OPEN MOUTH
U+1F63B 😻 	SMILING CAT FACE WITH HEART-SHAPED EYES
```

> __TEORIA__: na área de recuperação de informação (_information retrieval_) esses problemas são caracterizados por duas métricas: [precisão e revocação](https://pt.wikipedia.org/wiki/Precis%C3%A3o_e_revoca%C3%A7%C3%A3o) (_precision_, _recall_). Resultados demais é falta de precisão: o sistema está recuperando resultados irrelevantes, ou encontrando falsos positivos. Resultados de menos é falta de revocação: o sistema está deixando de recuperar resultados relevantes, ou seja, falsos negativos.

Vamos melhorar a precisão e a revocação pesquisando por palavras inteiras. Poderíamos resolver o problema todo mexendo apenas na função `Listar`, mas isso deixaria ela muito grande e difícil de testar. Então vamos colocar um pouco das novas funcionalidades na função `AnalisarLinha` e em outras funções que criaremos aos poucos.

## Melhorias em `AnalisarLinha`

Em vez de devolver apenas o código e o nome do caractere, vamos fazer a função `AnalisarLinha` devolver também as palavras do nome, na forma de uma lista de strings. Em Go, uma lista de strings é representada pela notação `[]string`, que lê-se como uma fatia de strings (_slice of strings_).

Para começar, mudamos o teste `TestAnalisarLinha`:

```go
func TestAnalisarLinha(t *testing.T) {
	runa, nome, palavras := AnalisarLinha(linhaLetraA) // ➊
	if runa != 'A' {
		t.Errorf("Esperado: 'A'; recebido: %q", runa)
	}
	const nomeA = "LATIN CAPITAL LETTER A"
	if nome != nomeA {
		t.Errorf("Esperado: %q; recebido: %q", nomeA, nome)
	}
  palavrasA := []string{"LATIN", "CAPITAL", "LETTER", "A"} // ➋
	if ! reflect.DeepEqual(palavras, palavrasA) { // ➌
		t.Errorf("\n\tEsperado: %q\n\trecebido: %q", palavrasA, palavras) // ➍
	}
}
```

➊ Incluímos a variável `palavras`, que vai receber a `[]string`.

➋ Criamos a variável `palavrasA`, com o valor esperado.

➌ Em Go, fatias não são comparáveis diretamente, ou seja, os operadores `==` e `!=` não funcionam com elas. Porém o pacote `reflect` oferece a função `DeepEqual`, que compara estruturas de dados em profundidade. `reflect.DeepEqual` é particularmente útil em testes, mas em geral deve ser evitada no código do programa em si, por razões apresentadas logo abaixo.

➍ Usamos `"\n\t"` para exibir este erro em linhas separadas e indentadas no mesmo nível, para facilitar a comparação visual do esperado com o recebido. Coloque `"X"` no lugar de `"A"` na variável `palavrasA` para forçar o erro e ver o formato da mensagem. Também algeramos as outras mensagens de erro para usar as palavras "esperado/recebido", por consistência.

> __NOTA__: O pacote `reflect` traz para a Go recursos poderosos e convenientes de linguagens dinâmicas, como por exemplo comparar ou copiar facilmente estruturas de dados complexas. Para quem tem experiência com linguagens de nível mais alto como Python, Ruby, JavaScript e PHP, é tentador sair usando `reflect` em seus programas Go. No entanto, a comunidade Go recomenda evitar abusar de `reflect`, por dois motivos principais: desempenho e salvaguardas de tipo (_type safety_).

> O desempenho de uma função como `DeepEqual` pode ser uma ordem de grandeza inferior ao código equivalente otimizado para os tipos de dados envolvidos. E a natureza dinâmica das funções de `reflect` possibilita a ocorrência de erros em tempo de execução que seriam capturados pelo compilador, se o seu código fosse escrito declarando os tipos específicos.

> No entanto, para escrever testes vale a pena usar o `reflect.DeepEqual`. Desempenho não é uma prioridade nos testes, e podemos abrir mão de algumas salvaguardas de tipo nos testes, porque elas continuam valendo em nosso código principal (onde não usamos `reflect`).

As mudanças necessárias para satisfazer este teste são simples:

```go
// AnalisarLinha devolve a runa, o nome e uma fatia de palavras que
// ocorrem no campo nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) { // ➊
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	palavras := strings.Fields(campos[1]) // ➋
	return rune(código), campos[1], palavras // ➌
}
```

➊ Na declaração de `AnalisarLinha`, acrescentamos o tipo de mais um valor a ser devolvido: `[]string`.

➋ Produzimos a fatia de palavras do nome, usando `strings.Fields` que é como `strings.Split`, mas usa como separador qualquer caractere Unicode considerado _whitespace_.

➌ Devolvemos a fatia de palavras, além da runa e seu nome.

Além disso, para poder compilar o programa e rodar o teste, precisamos mexer na função `Listar` onde invocamos `AnalisarLinha`, para aceitar a fatia de palavras devolvida como terceiro resultado, mesmo ignorando esse valor por enquanto:

```go
    runa, nome, _ := AnalisarLinha(linha)
```

Isso satisfaz o teste de `AnalisarLinha`. Mas para fazer `Listar` trabalhar com a fatia de palavras, várias mudanças serão necessárias.

## Consultas com várias palavras em `Listar`

O trecho que precisa ser melhorado em `Listar` é este:

```go
		runa, nome, _ := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
		 	fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
```

Em vez de procurar a string `consulta` dentro do `nome`, agora vamos procurar cada palavra da consulta na lista de palavras devolvida por `AnalisarLinha`. Em Python isso poderia ser feito facilmente em uma linha de código, pela subtração de conjuntos (tipo `set`). Infelizmente, Go por enquanto não tem um tipo `set`. Go não tem sequer uma função na biblioteca padrão que diga se uma string está presente em uma fatia de strings. Então o jeito é arregaçar a manga e codar, guiados por testes.

Primeiro vamos implementar a função `contém`, que devolve `true` se uma fatia de strings contém uma determinada string. Para verificar três casos em uma função de teste, vamos usar um [teste em tabela](https://golang.org/doc/code.html#Testing).

Para decifrar a elaborada sintaxe marcada com ➊, ➋, ➌ e ➍ em `TestContém` (mais abaixo), vale a pena ver um caso mais simples da mesma sintaxe. Suponha que você quer declarar e inicializar uma variável com uma fatia de bytes. Essa seria uma forma de fazê-lo:

```go
var octetos = []byte{10, 20, 30}
```

Repare que temos a palavra reservada `var`, seguida do identificador da variável `octetos`, um sinal `=`, e um valor literal do tipo `[]byte`. Valores literais de tipos compostos em Go são escritos assim: o identificador do tipo, seguido de zero ou mais itens ou campos entre chaves: `[]byte{10, 20, 30}`.

Agora vamos analisar `TestContém`, que usa uma declaração `var` semelhante, apenas mais extensa:

```go
func TestContém(t *testing.T) {
	var casos = []struct { // ➊
		fatia     []string
		procurado string
		esperado  bool
	}{ // ➋
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false}, // ➌
	} // ➍
	for _, caso := range casos { // ➎
		recebido := contém(caso.fatia, caso.procurado) // ➏
		if obtido != caso.esperado {                 // ➐
			t.Errorf("contém(%#v, %#v) esperado: %v; recebido: %v",
				caso.fatia, caso.procurado, caso.esperado, recebido) // ➑
		}
	}
}
```

➊ Esta declaração `var` cria a variável `casos` e atribui a ela uma fatia de `struct` anônima. A `struct` é declarada dentro do primeiro par de `{}` com três campos: uma fatia de strings, uma string e um booleano.

➋ Completando a declaração `var`, o segundo par de `{}` contém o valor literal da `[]struct`, que são três itens delimitados por `{}`, sendo que cada item é formado por uma fatia de strings, uma string e um booleano.

➌ É obrigatório incluir essa vírgula ao final do último item de um literal composto de várias linhas, se você quiser fechar a chave do literal na próxima linha como fizemos aqui.

➍ Aqui termina a declaração `var` que começou em ➊.

➎ Usamos a sintaxe de laço `for/range` para percorrer os três itens de `casos`. A cada iteração, o `for/range` produz dois valores: um índice a partir de zero (que descartamos atribuindo a `_`) e o valor do item correspondente, que atribuímos a `caso`.

➏ Invocamos `contém`, passando os valores de `caso.fatia` e `caso.procurado`. A função tem que devolver `true` se `caso.fatia` contém o item `caso.procurado`.

➐ Comparamos o resultado `recebido` com `caso.esperado`. Se forem diferentes...

➑ ...mostramos os argumentos passados e os valor obtido.

A implementação de `contém` é bem mais simples que o `TestContém`:

```go
func contém(fatia []string, procurado string) bool { // ➊
	for _, item := range fatia {
		if item == procurado {
			return true // ➋
		}
	}
	return false // ➌
}
```

➊ `contém` aceita uma fatia de strings e uma string, devolvendo `true` se a string é igual a um dos itens da fatia.

➋ Devolvemos `true` imediatamente assim que um `item` da fatia é igual ao texto `procurado`.

➌ Se chegamos até aqui, é porque o `procurado` não foi encontrado; devolvemos `false`.

A função `contém` é o primeiro tijolo da solução de busca por várias palavras. Agora precisamos de outra função auxiliar, `contémTodos` que devolve `true` se uma fatia contém todos os itens de uma segunda fatia. Ou seja, se a segunda fatia é um sub-conjunto da primeira (isso já estaria pronto se Go tivesse o conceito de conjuntos em sua biblioteca padrão).

Usamos outro teste de tabela:

```go
func TestContémTodos(t *testing.T) {
	casos := []struct { // ➊
		fatia      []string
		procurados []string
		esperado   bool
	}{ // ➋
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true}, // ➌
		{[]string{}, []string{}, true},
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}
	for _, caso := range casos {
		obtido := contémTodos(caso.fatia, caso.procurados) // ➍
		if obtido != caso.esperado {
			t.Errorf("contémTodos(%#v, %#v)\nesperado: %v; recebido: %v",
				caso.fatia, caso.procurados, caso.esperado, obtido) // ➎
		}
	}
}
```

➊ Agora usamos uma declaração curta (_short declaration_), com o sinal `:=` em vez de var. O efeito é o mesmo, assim como o resto da sintaxe.

➋ Aqui temos 7 casos de teste.

➌ Caso a fatia `caso.procurados` seja vazia, o resultado será sempre `true`.

➍ Para cada `caso`, invocamos `contémTodos` com os campos `.fatia` e `.procurados`.

➎ Caso o `obtido` não seja igual ao `caso.esperado`, mostramos os argumentos passados, o resultado obtido e o esperado. O verbo de formatação `%#v` mostra o valor usando a sintaxe literal de Go.

Veja a diferença na formatação. Aqui a mensagem usando apenas `%v`:

```
--- FAIL: TestContémTodos (0.00s)
	runefinder_test.go:73: contémTodos([A B C], [A B])
		esperado: true; recebido: false
```

E aqui, usando `%#v` para formatar os argumentos de `contémTodos`

```
$ go test
--- FAIL: TestContémTodos (0.00s)
	runefinder_test.go:73: contémTodos([]string{"A", "B", "C"}, []string{"A", "B"})
		esperado: true; recebido: false
```

Eisa a implementação de `contémTodos`, bem simples porque já temos `contém`:

```go
func contémTodos(fatia []string, procurados []string) bool {
	for _, procurado := range procurados {
		if !contém(fatia, procurado) {
			return false
		}
	}
	return true
}
```

Aqui não há nenhuma novidade de sintaxe.

Vamos incluir outra função exemplo nos testes de `Listar` para cobrir uma consulta com mais de uma palavra:

```go
func ExampleListar_duasPalavras() {
	texto := strings.NewReader(linhas3Da43)
	Listar(texto, "CAPITAL LATIN")
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}
```

Finalmente, faremos o ajuste em `Listar` para satisfazer o teste `ExampleListar_duasPalavras`. As mudanças são simples, porque toda a lógica interessante está em `contém` e `contémTodos`.

```go
// Listar exibe na saída padrão o código, a runa e o nome dos caracteres Unicode
// cujo nome contem as palavras da consulta.
func Listar(texto io.Reader, consulta string) {
	termos := strings.Fields(consulta) // ➊
	varredor := bufio.NewScanner(texto)
	for varredor.Scan() {
		linha := varredor.Text()
		if strings.TrimSpace(linha) == "" {
			continue
		}
		runa, nome, palavrasNome := AnalisarLinha(linha) // ➋
		if contémTodos(palavrasNome, termos) {           // ➌
			fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
	}
}
```

➊ Criamos uma fatia `termos` com as palavras da string `consulta`.

➋ O terceiro resultado de `AnalisarLinha` é a lista de palavras do nome.

➌ Usamos `contémTodos` para checar se `palavrasNome` contém cada um dos `termos`.

Podemos criar um teste funcional do pacote para demonstrar o funcionamento de uma consulta com duas palavras, exibindo resultados onde tais palavras não aparecem em sequência no nome do caractere:

```go
func Example_consultaDuasPalavras() { // ➊
	oldArgs := os.Args // ➋
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"", "cat", "smiling"}
	main() // ➌
	// Output:
	// U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
	// U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
	// U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
}
```

Agora você pode experimentar o programa com `go run` ou criar outro executável com `go build` para ver a nova funcionalidade em ação. Por exemplo, pesquisar peças pretas do Xadrez:

```bash
$ go build
$ ./runas chess black
U+265A	♚	BLACK CHESS KING
U+265B	♛	BLACK CHESS QUEEN
U+265C	♜	BLACK CHESS ROOK
U+265D	♝	BLACK CHESS BISHOP
U+265E	♞	BLACK CHESS KNIGHT
U+265F	♟	BLACK CHESS PAWN
```

Ou ainda, o trem-bala japonês:

```bash
$ ./runas bullet train
U+1F685	🚅	HIGH-SPEED TRAIN WITH BULLET NOSE
```

E mesmo com apenas uma palavra, os resultados são melhores. A busca por "cat" traz principalmente emojis com gatos, e não mais caracteres com as letras "CAT" em qualquer parte do nome.

```bash
$ ./runas cat
U+A2B6	ꊶ	YI SYLLABLE CAT
U+101EC	𐇬	PHAISTOS DISC SIGN CAT
U+1F408	🐈	CAT
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

Agora é um bom momento para exercícios. Veja instruções no [Passo 6](passo-06).
