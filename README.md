# Runas: passo a passo com TDD

Neste repositório você pode ver o desenvolvimento passo a passo do exemplo `runas`: um utilitário em Go para localizar caracteres Unicode pelo nome.

Cada etapa do desenvolvimento é documentada explicando os recursos da linguagem Go usados no código do exemplo.

Você não precisa saber nada de Go para acompanhar. Os requisitos são conhecer alguma linguagem de programação moderna.


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

Leia mais nas [páginas do projeto](https://ThoughtWorksInc.github.io/runas/).


## Créditos

Este tutorial é baseado no exemplo `charfinder` do capítulo 18 de [Python Fluente](http://novatec.com.br/livros/pythonfluente/), de Luciano Ramalho. A versão Go, chamada `runefinder`, foi iniciada no grupo de estudos [Garoa Gophers](https://garoa.net.br/wiki/Garoa_Gophers), com a participação de Afonso Coutinho (@afonso), Alexandre Souza (@alexandre), Andrews Medina (@andrewsmedina), João "JC" Martins (@jcmartins), Luciano Ramalho (@ramalho), Marcio Ribeiro (@mmr) e Michael Howard.
