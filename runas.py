#!/usr/bin/env python3

import os, sys
import urllib.request
import threading

# URL_UCD é a URL canônica do arquivo UnicodeData.txt mais atual
URL_UCD = 'http://www.unicode.org/Public/UNIDATA/UnicodeData.txt'


def analisar_linha(linha):
    campos = linha.split(';')
    código = int(campos[0], 16)
    nome = campos[1]
    palavras = set(nome.replace('-', ' ').split())
    if campos[10]:
        nome = '{} ({})'.format(nome, campos[10])
        palavras.update(campos[10].replace('-', ' ').split())
    return chr(código), nome, palavras


def listar(texto, consulta):
    consulta = set(consulta.replace('-', ' ').split())
    for linha in texto:
        runa, nome, palavras_nome = analisar_linha(linha)
        if consulta <= palavras_nome:
            try:
                print('U+{:04X}\t{}\t{}'.format(ord(runa), runa, nome))
            except UnicodeEncodeError:
                print('U+{:04X}\t\uFFFD\t{}'.format(ord(runa), nome))


def obter_caminho_UCD():
    caminho_UCD = os.environ.get('UCD_PATH')
    if caminho_UCD is None:
        caminho_UCD = os.path.join(os.environ['HOME'], "UnicodeData.txt")
    return caminho_UCD


def abrir_UCD(caminho):
    try:
        ucd = open(caminho)
    except FileNotFoundError:
        print('%s não encontrado\nbaixando %s' % (caminho, URL_UCD))
        feito = threading.Event()
        threading.Thread(target=baixar_UCD, args=(URL_UCD, caminho, feito)).start()
        progresso(feito)
        ucd = open(caminho)
    return ucd


def progresso(feito):
    while not feito.wait(.150):
        print('.', end='', flush=True)
    print()


def baixar_UCD(url, caminho, feito):
    with urllib.request.urlopen(url) as resposta:
        octetos = resposta.read()
    with open(caminho, 'wb') as arquivo:
        arquivo.write(octetos)
    feito.set()


def main():
    with abrir_UCD(obter_caminho_UCD()) as ucd:
        consulta = ' '.join(sys.argv[1:])
        listar(ucd, consulta.upper())


if __name__ == '__main__':
    main()
