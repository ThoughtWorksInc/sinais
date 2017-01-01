# Runas, passo 7 (bônus): download automático da UCD

Nosso programa `runas` depende da presença do arquivo `UnicodeData.txt` no diretório atual para funcionar. Neste passo, vamos criar uma função para baixar o arquivo direto do site `unicode.org`, caso ele não esteja presente em um caminho local configurado pelo usuário.

Antes de mais nada, vamos verificar que temos uma versão funcional de `runas`, após o exercício do `passo-06`.
