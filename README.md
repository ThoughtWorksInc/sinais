# Runas: passo a passo

Neste repositório você pode ver o desenvolvimento passo a passo do exemplo `runefinder`: um utilitário em Go para localizar caracteres Unicode pelo nome.

Cada etapa do desenvolvimento é documentada explicando os recursos da linguagem Go usados no código do exemplo.

Você não precisa saber nada de Go para acompanhar. Os requisitos são conhecer alguma linguagem de programação moderna, e saber o básico de `git` se quiser acessar os _branches_ contendo o código em cada passo da implementação.


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

Leia mais nas [páginas do projeto](https://labgo.github.io/runas/).
