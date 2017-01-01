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

O trecho que precisa ser melhorado em `Listar` é este:

```go
		runa, nome, _ := AnalisarLinha(linha)
		if strings.Contains(nome, consulta) {
		 	fmt.Printf("U+%04X\t%[1]c\t%s\n", runa, nome)
		}
```

Em vez de procurar a string `consulta` dentro do `nome`, agora temos que procurar cada palavra da consulta na lista de palavras devolvida por `AnalisarLinha`. Em Python isso poderia ser feito facilmente em uma linha usando dois conjuntos (tipo `set`). Infelizmente, Go por enquanto não tem um tipo `set`, e nem mesmo uma função na biblioteca padrão que diga se uma string está presente em uma fatia de strings. Então o jeito é arregaçar a manga e fazer, guiados por testes.

Primeiro vamos implementar a função `contém`, que devolve `true` se uma fatia de strings contém uma determinada string. Para verificar três casos em uma função de teste, vamos usar um [teste em tabela](https://golang.org/doc/code.html#Testing).

Para decifrar a sintaxe marcada com ➊ e ➋ em `TestContém` (mais abaixo), vale a pena ver um caso mais simples da mesma sintaxe. Suponha que você quer declarar e inicializar uma variável com uma fatia de bytes. Essa seria uma forma de fazê-lo:

```go
var l = []byte{10, 20, 30}
```

Repare que temos a declaração `var`, seguida do identificador da variável `l`, um sinal `=`, e um valor literal do tipo `[]byte`. Valores literais de tipos compostos em Go são escritos assim: o identificador do tipo, seguido de zero ou mais itens ou campos entre chaves: `[]byte{10, 20, 30}`.

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
		obtido := contém(caso.fatia, caso.procurado) // ➏
		if obtido != caso.esperado {
			t.Errorf("contém(%#v, %#v) esperado: %v; recebido: %v",
				caso.fatia, caso.procurado, caso.esperado, obtido) // ➐
		}
	}
}
```

➊ Esta declaração `var` cria uma variável `casos` e atribui a ela uma fatia de `struct` anônima. A `struct` é declarada dentro do primeiro par de `{}` com três campos: uma fatia de strings, uma string e um booleano.

➋ Completando a declaração `var`, o segundo par de `{}` contém o valor literal da `[]struct`, que são três itens delimitados por `{}`, sendo que cada item é formado por uma fatia de strings, uma string e um booleano.

➌ É obrigatório incluir essa vírgula ao final do último item de um literal composto de várias linhas.

➍ Aqui termina a declaração `var` que começou em ➊.

➎ Usamos a sintaxe de laço `for/range` para percorrer os três itens de `casos`. A cada iteração, o `for/range` produz dois valores: um índice a partir de zero (que ignoramos atribuindo a `_`) e o valor do item correspondente, que atribuímos a `caso`.

➏ Invocamos `contém`, passando os valores de `caso.fatia` e `caso.procurado`. A função tem que devolver `true` se `caso.fatia` contém o item `caso.procurado`.

➐ Em caso de falha, mostramos os argumentos passados e os valores que recebemos de volta.



Pois bem,


mas analisando o `UnicodeData.txt` dá para ver dois requisitos adicionais que vamos implementar no _branch_ `passo-06`, texto em `passo-06.md`.
