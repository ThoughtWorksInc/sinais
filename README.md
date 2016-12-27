# Runas: passo a passo

Neste repositório você pode ver o desenvolvimento passo a passo do exemplo `runefinder`: um utilitário em Go para localizar caracteres Unicode pelo nome.

## Nosso objetivo

Ao final desse tutorial, teremos um utilitário de linha de comando que faz isso:

```
$ runefinder face eyes
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


## Contexto

O projeto [Unicode](http://unicode.org) mantém um banco de dados chamado Unicode Character Database (UCD), com nomes descritivos e outros metadados sobre os mais de 128.000 caracteres que fazem parte da versão atual do padrão. A tabela mais interessante UCD é um arquivo ASCII de 1.6MB cuja versão mais atual pode ser obtida neste URL: `http://www.unicode.org/Public/UNIDATA/UnicodeData.txt`.

O `UnicodeData.txt` traz informações sobre os caracteres de praticamente todos os idiomas, incluindo também símbolos, ícones e emojis, somando 30.592 linhas na versão 9.0 do padrão Unicode. Isso corresponde a cerca de 24% do total de caracteres do UCD (a maior parte dos ideogramas CJK -- Chinês/Japonês/Coreano -- é documentada em outras tabelas).

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

Nosso programa `runefinder` vai usar o `UnicodeData.txt` para localizar caracteres pelo nome. Então, mãos à obra!

Para continuar, mude para o _branch_ `passo-01` e veja o arquivo `passo-01.md`.
