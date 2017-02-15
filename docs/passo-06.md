---
permalink: passo-06
---

# Runas, passo 06: hífens e nomes antigos

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

* Alguns nomes têm palavras hifenadas, como "HYPHEN-MINUS" (por coincidência). É desejável que o usuário encontre esses caracteres digitando apenas uma das palavras, "HYPHEN" ou "MINUS".
* Algumas linhas tem no campo índice 10 um nome diferente, que era o nome adotado no Unicode 1.0 (veja documentação do [UCD 9.0](http://www.unicode.org/reports/tr44/tr44-18.html#UnicodeData.txt)). Por exemplo o caractere U+002E, "FULL STOP", era "PERIOD". Incluir esses nomes também pode facilitar a vida dos usuários.

Nesta parte do tutorial a proposta é que você implemente as mudadanças para exercitar os conceitos vistos até agora.

No final do exercício, o programa deverá se comportar assim:

```
$ ./runas quote
U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
U+2358	⍘	APL FUNCTIONAL SYMBOL QUOTE UNDERBAR
U+235E	⍞	APL FUNCTIONAL SYMBOL QUOTE QUAD
$ ./runas minus hyphen
U+002D	-	HYPHEN-MINUS
U+207B	⁻	SUPERSCRIPT MINUS (SUPERSCRIPT HYPHEN-MINUS)
U+208B	₋	SUBSCRIPT MINUS (SUBSCRIPT HYPHEN-MINUS)
U+FE63	﹣	SMALL HYPHEN-MINUS
U+FF0D	－	FULLWIDTH HYPHEN-MINUS
U+E002D		TAG HYPHEN-MINUS

```

Comportamentos a observar:

* A busca pelas palavra "quote" encontra qualquer caractere onde essa palavra apareça, mesmo como parte de uma palavra hifenada como "APOSTROPHE-QUOTE"
* Buscar as palavras "minus hyphen" encontra qualquer caractere onde ambas palavras apareçam, em qualquer ordem, mesmo como partes parte de uma palavra hifenada como "HYPHEN-MINUS"
* A busca inclui o campo 10 dos registros em "UnicodeData.txt", onde aparecem os nomes do Unicode 1.0 como "APOSTROPHE-QUOTE" e "SUPERSCRIPT HYPHEN-MINUS".
* Quando há informações no campo 10, elas são exibidas nos resultados entre parêntesis depois do nome do caractere. Por exemplo:

```
U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
```

Mãos à obra! Se precisar de ajuda, veja dicas a seguir.

Independente de usar dicas ou não, depois me conte quanto tempo você levou para fazer este exercício.

## Dicas

### 1. Mudanças em `AnalisarLinha`

Para atender os requisitos do exercício, a função `AnalisarLinha` precisa devolver uma fatia de palavras que inclua as partes de cada termo com hífen, e também as palavras do campo índice 10. Além disso, havendo conteúdo no campo 10, esse texto deverá ser concatenado ao nome, entre parêntesis.

> Tente resolver o exercício com a dica acima, antes de ler a próxima dica.


### 2. Casos de teste para `AnalisarLinha`

Crie outra função de teste, em vez de apagar ou mudar o `TestAnalisarLinha` que já existe. O comportamento que já verificamos em `TestAnalisarLinha` continua valendo.

Agora teremos pelo menos três casos de teste:

* Campo 10 vazio e nenhum hífen.
* Campo 10 vazio e hífen no campo 1.
* Campo 10 utilizado e hífens presentes.

Para testar isso sem duplicar muito código de `TestAnalisarLinha`, sugiro fazer um [teste em tabela](https://golang.org/doc/code.html#Testing). Use a estrutura de `TestContém` como inspiração.

> Tente resolver o exercício com a dica acima, antes de ler a próxima dica.


### 3. Testando com o "baby steps"

A metodologia TDD recomenda "baby steps" - passos bem simples. Ao criar um teste em tabela, coloque inicialmente apenas um caso na tabela. Faça este caso passar antes de colocar outro caso. No final, sua tabela pode ter vários casos, mas você só deve incluir e fazer passar um caso de cada vez.

> Tente resolver o exercício com a dica acima, antes de ler a próxima dica.


### 4. Resolva primeiro o tratamento dos hífens

Há várias formas de transformar `"SMALL HYPHEN-MINUS"` em uma lista de três palavras: `[]string {"SMALL", "HYPHEN", "MINUS"}`. Você pode passar o texto original por [`strings.Replace`](https://golang.org/pkg/strings/#Replace) para substituir `"-"` por `" "` antes de usar `strings.Fields` para separar as palavras. Ou então você pode usar [`strings.FieldsFunc`](`https://golang.org/pkg/strings/#FieldsFunc`) para fazer as duas operações de uma vez só.

Seja como for, recomendo criar uma função auxiliar para fazer essa separação por espaços ou hífens. E não esqueça de fazer TDD: escreva o teste antes de implementar a funcionalidade!

> Tente resolver o exercício com a dica acima, antes de ler a próxima dica.


### 5. Lembre-se de `reflect.DeepEqual`

Para testar a função que transforma `"SMALL HYPHEN-MINUS"` em `[]string {"SMALL", "HYPHEN", "MINUS"}`, você vai precisar comparar a fatia produzida com a fatia esperada, mas Go só permite comparar uma fatia com `nil`. Para comparar uma fatia com outra, lembre-se de usar a função `reflect.DeepEqual` como fizemos em `TestAnalisarLinha`.

> Tente resolver o exercício com a dica acima, antes de ler a próxima dica.


### 6. Atenção para palavras duplicadas nos campos 1 e 10

Observe este caso:

```
0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;
```
Queremos que `AnalisarLinha` devolva como lista de palavras apenas isto: `[]string{"APOSTROPHE", "QUOTE"}` e não `[]string{"APOSTROPHE", "APOSTROPHE", "QUOTE"}`. Você pode usar função auxilar `contém` que criamos no `passo-05` para resolver este problema.


## Fim do exercício

Quando tiver terminado, anote o tempo que levou para fazer o exercício e conte para o instrutor. Isso ajuda a melhorar o tutorial.

O [Passo 7](passo-07) é uma seção bônus, onde faremos o download automático do arquivo `UnicodeData.txt`. Siga em frente se tiver feito o exercício. A solução deste exercício e o código que faz download estão no diretório raiz do repositório `runas` (não há um diretório `runas07`).
