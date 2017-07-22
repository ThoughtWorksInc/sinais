import unittest
import unittest.mock

import io, sys

import sinais

LINHAS_3D_A_43 = '''
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
'''.strip()


class TestAnalise(unittest.TestCase):

    def test_analisar_linha(self):
        linha_letra_a = '0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;'
        runa, nome, palavras = sinais.analisar_linha(linha_letra_a)
        self.assertEqual(runa, 'A')
        self.assertEqual(nome, 'LATIN CAPITAL LETTER A')
        self.assertEqual(palavras, {'LATIN', 'CAPITAL', 'LETTER', 'A'})

    def test_analisar_linha_com_hifen_e_campo_10(self):
        casos = [
            ('0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;',
             '!', 'EXCLAMATION MARK', {'EXCLAMATION', 'MARK'}),
            ('002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;',
             '-', 'HYPHEN-MINUS', {'HYPHEN', 'MINUS'}),
            ('0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;',
             "'", 'APOSTROPHE (APOSTROPHE-QUOTE)', {'APOSTROPHE', 'QUOTE'}),
        ]
        for linha, runa_ok, nome_ok, palavras_ok in casos:
            runa, nome, palavras = sinais.analisar_linha(linha)
            self.assertEqual(runa, runa_ok)
            self.assertEqual(nome, nome_ok)
            self.assertEqual(palavras, palavras_ok)


class TestListagem(unittest.TestCase):

    def setUp(self):
        self.linhas_3D_a_43 = io.StringIO(LINHAS_3D_A_43)
        if not hasattr(sys.stdout, 'getvalue'):
            self.fail('run tests in buffered mode (use -b flag in command line)')

    def test_listar(self):
        sinais.listar(self.linhas_3D_a_43, 'MARK')
        saida = sys.stdout.getvalue()
        self.assertEqual(saida, 'U+003F\t?\tQUESTION MARK\n')

    def test_listar_dois_resultados(self):
        sinais.listar(self.linhas_3D_a_43, 'SIGN')
        saida = sys.stdout.getvalue()
        self.assertEqual(saida, 'U+003D\t=\tEQUALS SIGN\n'
                                'U+003E\t>\tGREATER-THAN SIGN\n')

    def test_listar_duas_palavras(self):
        sinais.listar(self.linhas_3D_a_43, 'CAPITAL LATIN')
        saida = sys.stdout.getvalue()
        self.assertEqual(saida, 'U+0041\tA\tLATIN CAPITAL LETTER A\n'
                                'U+0042\tB\tLATIN CAPITAL LETTER B\n'
                                'U+0043\tC\tLATIN CAPITAL LETTER C\n')

class TestPrograma(unittest.TestCase):

    def test_um_arg(self):
        args =  ['', 'cruzeiro']
        with unittest.mock.patch.object(sys, 'argv', args):
            sinais.main()
            saida = sys.stdout.getvalue()
            self.assertEqual(saida, 'U+20A2\t₢\tCRUZEIRO SIGN\n')

    def test_dois_args(self):
        args =  ['', 'cat', 'smiling']
        with unittest.mock.patch.object(sys, 'argv', args):
            sinais.main()
            saida = sys.stdout.getvalue()
            self.assertEqual(saida,
                'U+1F638\t😸\tGRINNING CAT FACE WITH SMILING EYES\n'
                'U+1F63A\t😺\tSMILING CAT FACE WITH OPEN MOUTH\n'
                'U+1F63B\t😻\tSMILING CAT FACE WITH HEART-SHAPED EYES\n'
            )

    def test_hifen_e_campo_10(self):
        args =  ['', 'quote']
        with unittest.mock.patch.object(sys, 'argv', args):
            sinais.main()
            saida = sys.stdout.getvalue()
            self.assertEqual(saida,
                "U+0027\t'\tAPOSTROPHE (APOSTROPHE-QUOTE)\n"
                'U+2358\t⍘\tAPL FUNCTIONAL SYMBOL QUOTE UNDERBAR\n'
                'U+235E\t⍞\tAPL FUNCTIONAL SYMBOL QUOTE QUAD\n'
            )


if __name__ == '__main__':
    unittest.main(buffer=True)
