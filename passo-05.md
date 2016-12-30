# Runas, passo 5: busca por palavras inteiras

A versão MVP1 do programa `runas` busca caracteres comparando uma substring do nome. Isso gera dois problemas: baixa precisão e baixa re

* Resultados demais: pesquisando "cat" vêm 82 caracteres, sendo que a maioria não tem nada a ver com gatos, por exemplo "MULTIPLICATION SIGN".
* Resultados de menos: a ordem das palavras na consulta deveria ser ignorada: "chess black" e "black chess" deveriam devolver os mesmos resultados, e "cat smiling" deveria encontrar todos estes caracteres:

```
U+1F638 😸 	GRINNING CAT FACE WITH SMILING EYES
U+1F63A 😺 	SMILING CAT FACE WITH OPEN MOUTH
U+1F63B 😻 	SMILING CAT FACE WITH HEART-SHAPED EYES
```

> __TEORIA__: na área de recuperação de informação (_information retrieval_) esses problemas são caracterizados por duas métricas: [precisão e revocação](https://pt.wikipedia.org/wiki/Precis%C3%A3o_e_revoca%C3%A7%C3%A3o) (_precision_, _recall_). Resultados demais é falta de precisão: o sistema está recuperando resultados irrelevantes, ou encontrando falsos positivos. Resultados de menos é falta de revocação: o sistema está deixando de recuperar resultados relevantes, ou seja, falsos negativos.

Vamos melhorar a precisão e a revocação pesquisando sempre por palavras inteiras. Poderíamos resolver o problea todo mexendo apenas na função `Listar`, mas isso deixaria ela muito grande e difícil de testar. Então vamos colocar um pouco das novas funcionalidades na função `AnalisarLinha` e em outras funções que criaremos aos poucos.

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
	palavras := strings.Split(campos[1], " ") // ➋
	return rune(código), campos[1], palavras // ➌
}
```

➊ Na declaração de `AnalisarLinha`, acrescentamos o tipo de mais um valor a ser devolvido: `[]string`.

➋ Produzimos a fatia de palavras do nome, usando `strings.Split`.

➌ Devolvemos a fatia de palavras, além da runa e seu nome.

Isso faz passar o teste, mas analisando o `UnicodeData.txt` dá para ver dois requisitos adicionais que vamos implementar em seguida.


## Tratando nomes com hífen e nomes antigos

Veja esta parte da tabela `UnicodeData.txt`:

```
0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;;
0028;LEFT PARENTHESIS;Ps;0;ON;;;;;Y;OPENING PARENTHESIS;;;;
0029;RIGHT PARENTHESIS;Pe;0;ON;;;;;Y;CLOSING PARENTHESIS;;;;
002A;ASTERISK;Po;0;ON;;;;;N;;;;;
002B;PLUS SIGN;Sm;0;ES;;;;;N;;;;;
002C;COMMA;Po;0;CS;;;;;N;;;;;
002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;
002E;FULL STOP;Po;0;CS;;;;;N;PERIOD;;;;
```

Duas coisas me chamaram atenção aqui:

* Alguns nomes têm palavras hifenadas, como "HYPHEN-MINUS" (por coincidência)! Seria interessante que o usuário pudesse encontrar esses caracteres digitando apenas uma das palavras, "HYPHEN" ou "MINUS".
* Algumas linhas tem no campo índice 10 um nome diferente, que era o nome adotado no Unicode 1.0 (veja documentação do [UCD 9.0](http://www.unicode.org/reports/tr44/tr44-18.html#UnicodeData.txt)). Por exemplo o caractere U+002E, "FULL STOP", era "PERIOD". Incluir esses nomes também pode facilitar a vida dos usuários.

Então para atender esses requisitos a função `AnalisarLinha` precisa devolver uma fatia de palavras que inclua as partes de cada termo com hífen, e também as palavras do campo índice 10. Em vez de um simples caso de teste, agora teremos pelo menos três:

* Caso mais simples: campo 10 vazio e nenhum hífen.
* Caso simples: campo 10 utilizado e nenhum hífen.
* Caso mais complexo: campo 10 utilizado e hífens presentes.

Para testar isso sem duplicar muito código em `TestAnalisarLinha`, vamos usar um [teste em tabela](https://golang.org/doc/code.html#Testing). A nova versão dessa função de teste vai ficar assim:

```go
func TestAnalisarLinha(t *testing.T) {
	var casos = []struct { // ➊
		linha    string
		runa     rune
		nome     string
		palavras []string
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK", []string{"EXCLAMATION", "MARK"}},
		{"002E;FULL STOP;Po;0;CS;;;;;N;PERIOD;;;;",
			'.', "FULL STOP (PERIOD)", []string{"FULL", "STOP", "PERIOD"}},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)", []string{"APOSTROPHE", "QUOTE"}},
	}
	for _, caso := range casos { // ➌
		runa, nome, palavras := AnalisarLinha(caso.linha) // ➍
		if runa != caso.runa || nome != caso.nome ||
			!reflect.DeepEqual(palavras, caso.palavras) {
			t.Errorf("\nAnalisarLinha(%q)\n-> (%q, %q, %q)", // ➎
				caso.linha, runa, nome, palavras)
		}
	}
}
```

Várias novidades neste teste. Vejamos:

➊ Aqui usamos a declaração `var` para definir o tipo e inicializar a variável `casos`, tipo `[]struct` -- uma fatia de `struct` (pense em uma lista de registros). A `struct` anônima é definida em seguida, com quatro campos: `linha`, `runa`, `nome` e `palavras`.

➋ Ainda continuando a declaração `var`, o segundo bloco contém os valores da fatia de structs com três itens, ou seja, os valores de cada um dos quatro campos, para cada um dos três itens. Resumindo: criamos uma série de registros na forma de uma fatia onde cada item é uma `struct`.

➌ Usamos a sintaxe de laço `for/range` para percorrer os três itens de `casos`. A cada iteração, o `for/range` produz dois valores: um índice a partir de zero (que ignoramos atribuindo a `_`) e o valor do item correspondente, que atribuímos a `caso`.

➍ Invocamos `AnalisarLinha`, passando o valor do campo `caso.linha`.

➎ Em caso de falha, mostramos o argumento que foi passado e os valores que recebemos de volta.

Veja o resultado de executar o teste agora:

```bash
$ go test
--- FAIL: TestAnalisarLinha (0.00s)
	runefinder_test.go:41:
		AnalisarLinha("002E;FULL STOP;Po;0;CS;;;;;N;PERIOD;;;;")
		-> ('.', "FULL STOP", ["FULL" "STOP"])
	runefinder_test.go:41:
		AnalisarLinha("0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;")
		-> ('\'', "APOSTROPHE", ["APOSTROPHE"])
FAIL
exit status 1
FAIL	github.com/labgo/runas	0.064s
```

Repare que duas falhas foram reportadas, porque o primeiro caso de teste (o mais simples), passou. Isso demonstra que a chamada para `t.Errorf` não aborta o teste, mas apenas reporta o erro, e o teste continua rodando.
