# Runas: passo a passo com TDD

[Neste repositório](https://github.com/ThoughtWorksInc/runas) você pode ver o desenvolvimento passo a passo do exemplo `runas`: um utilitário em linguagem Go para localizar caracteres Unicode pelo nome.

Seguiremos o princípio do [TDD](http://tdd.caelum.com.br/) _(Test-Driven Development)_: primeiro escrevemos um teste, depois o código mais simples que faz o teste passar.

## Conteúdo

1. [Iniciando o TDD](passo-01)
2. [Primeira função completa](passo-02)
3. [Gerar a lista de caracteres](passo-03)
4. [MVP 1, o mínimo que é útil](passo-04)
5. [Busca por palavras inteiras](passo-05)
6. [Exercício: hífens e nomes antigos](passo-06)
7. [Bônus: download automático da UCD](passo-07)


## Nosso objetivo

Ao final desse tutorial, teremos um utilitário de linha de comando que faz isso:

```
$ runas face eyes
U+1F601	😁	GRINNING FACE WITH SMILING EYES
U+1F604	😄	SMILING FACE WITH OPEN MOUTH AND SMILING EYES
U+1F606	😆	SMILING FACE WITH OPEN MOUTH AND TIGHTLY-CLOSED EYES
U+1F60A	😊	SMILING FACE WITH SMILING EYES
U+1F60D	😍	SMILING FACE WITH HEART-SHAPED EYES
U+1F619	😙	KISSING FACE WITH SMILING EYES
U+1F61A	😚	KISSING FACE WITH CLOSED EYES
U+1F61D	😝	FACE WITH STUCK-OUT TONGUE AND TIGHTLY-CLOSED EYES
U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
U+1F63D	😽	KISSING CAT FACE WITH CLOSED EYES
U+1F644	🙄	FACE WITH ROLLING EYES
```

Você passa uma um mais palavras como argumento, e o programa devolve uma lista ordenada de caracteres Unicode cujas descrições contém todas as palavras que você passou.


## Sobre o banco de dados Unicode

O projeto [Unicode](http://unicode.org) mantém um banco de dados chamado Unicode Character Database (UCD), com nomes descritivos e outros metadados sobre os mais de 128.000 caracteres que fazem parte da versão atual do padrão. A tabela mais interessante do UCD é um arquivo ASCII de 1.6MB cuja versão mais atual pode ser obtida neste URL: [`http://www.unicode.org/Public/UNIDATA/UnicodeData.txt`](http://www.unicode.org/Public/UNIDATA/UnicodeData.txt).

O `UnicodeData.txt` traz informações sobre os caracteres de praticamente todos os idiomas, incluindo também símbolos, ícones e emojis. São 30.592 linhas na versão 9.0 do padrão Unicode. Isso corresponde a cerca de 24% do total de caracteres do UCD (a maior parte dos ideogramas CJK — Chinês/Japonês/Coreano — é documentada em outras tabelas).

Eis uma pequena amostra do `UnicodeData.txt`:

```
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
```

Os campos que nos interessam são o primeiro e o segundo: o código Unicode em hexadecimal e o nome oficial do caractere.

Nosso programa `runas` vai usar o `UnicodeData.txt` para localizar caracteres pelo nome. Então, mãos à obra!

Para iniciar o desenvolvimento, vá para o [Passo 1](passo-01). O código está no diretório `runas01` do repositório [https://github.com/ThoughtWorksInc/runas](https://github.com/ThoughtWorksInc/runas), mas eu recomendo que você copie o código em seu próprio espaço de trabalho, porque os códigos nos diretórios `runasNN` são o estado final de cada passo, porém existem passos intermediários que vale a pena você acompanhar escrevendo os testes, rodando os testes, e fazendo as mudanças necessárias para os testes passarem.
