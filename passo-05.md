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

> No entanto, para escrever testes vale a pena usar o `reflect.DeepEqual`. Desempenho não é uma prioridade nos testes, e as salvaguardas de tipo continuam valendo em nosso código principal (onde não usamos `reflect`), então podemos relaxá-las no código de teste.


```go
// AnalisarLinha devolve a runa e o nome de uma linha do UnicodeData.txt
func AnalisarLinha(linha string) (rune, string, []string) { // ➊
	campos := strings.Split(linha, ";")
	código, _ := strconv.ParseInt(campos[0], 16, 32)
	palavras := strings.Split(campos[1], " ") // ➋
	return rune(código), campos[1], palavras // ➌
}
```

➊

➋

➌
