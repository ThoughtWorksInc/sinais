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
	// 1 character found
}
```

Esse teste traz várias novidades:

➊ A função chamada simplesmente `Example` é o exemplo do pacote: um teste funcional para o pacote como um todo ([documentação sobre funções exemplo](https://golang.org/pkg/testing/#hdr-Examples)).

➋ Vamos simular a passagem de argumentos pela linha de comando. O primeiro passo é copiar os argumentos de `os.Args` para `oldArgs`, assim poderemos restaurá-los depois.

➌ Aqui criamos uma função anônima que vai restaurar o valor de `os.Args` no final da nossa função `Example`. Leia mais sobre o comando `defer` logo adiante.

➍ Mudamos os valor de `os.Args` para fazer o teste. Observe a sintaxe de uma fatia literal: primeiro o tipo `[]string`, depois os itens entre `{}`. O primeiro item de `os.Args` é o nome do programa (irrelevante para o nosso teste). O segundo item é a palavra que vamos buscar, `"cruzeiro"`, cuidadosamente escolhida porque só existe um caractere Unicode que contém essa palavra em seu nome.

➎ Invocamos a função `main`, a mesma que será chamada quando nosso programa for acionado na linha de comando.

O comando `defer` é uma inovação simples porém genial da linguagem Go. Ele serve para invocar uma função no final da função atual (`Example`). `defer` é útil para fechar arquivos, encerrar conexões, liberar travas, etc. É como se o corpo da função `Example` estivesse dentro de um `try/finally`, e as funções chamadas em `defer` seriam executadas no bloco `finally`, ou seja, após o `return` e mesmo que ocorram exceções. No exemplo, o uso de `defer` garante que o valor de `os.Args` será restaurado ao valor original, independente do sucesso ou fracasso do teste.

> _Nota_: Alterar uma variável global como `os.Args` pode ser perigoso em um
> sistema concorrente, mas Go só executa testes em paralelo se usamos o método
> [`T.Parallel`](https://golang.org/pkg/testing/#T.Parallel).

## A função `main`
