package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const lineLetterA = "0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;"

const lines3Dto43 = `
003D;EQUALS SIGN;Sm;0;ON;;;;;N;;;;;
003E;GREATER-THAN SIGN;Sm;0;ON;;;;;Y;;;;;
003F;QUESTION MARK;Po;0;ON;;;;;N;;;;;
0040;COMMERCIAL AT;Po;0;ON;;;;;N;;;;;
0041;LATIN CAPITAL LETTER A;Lu;0;L;;;;;N;;;;0061;
0042;LATIN CAPITAL LETTER B;Lu;0;L;;;;;N;;;;0062;
0043;LATIN CAPITAL LETTER C;Lu;0;L;;;;;N;;;;0063;
`

func TestParseLine(t *testing.T) {
	rune, name, words := ParseLine(lineLetterA) // ➊
	if rune != 'A' {
		t.Errorf("Esperado: 'A'; got: %q", rune)
	}
	const nameA = "LATIN CAPITAL LETTER A"
	if name != nameA {
		t.Errorf("Esperado: %q; got: %q", nameA, name)
	}
	wordsA := []string{"LATIN", "CAPITAL", "LETTER", "A"} // ➋
	if !reflect.DeepEqual(words, wordsA) {             // ➌
		t.Errorf("\n\tEsperado: %q\n\tgot: %q", wordsA, words) // ➍
	}
}

func TestParseLineWithHyphenAndField10(t *testing.T) {
	var tests = []struct { // ➊
		line    string
		rune     rune
		name     string
		words []string
	}{ // ➋
		{"0021;EXCLAMATION MARK;Po;0;ON;;;;;N;;;;;",
			'!', "EXCLAMATION MARK", []string{"EXCLAMATION", "MARK"}},
		{"002D;HYPHEN-MINUS;Pd;0;ES;;;;;N;;;;;",
			'-', "HYPHEN-MINUS", []string{"HYPHEN", "MINUS"}},
		{"0027;APOSTROPHE;Po;0;ON;;;;;N;APOSTROPHE-QUOTE;;;",
			'\'', "APOSTROPHE (APOSTROPHE-QUOTE)", []string{"APOSTROPHE", "QUOTE"}},
	}
	for _, tt := range tests { // ➌
		rune, name, words := ParseLine(tt.line) // ➍
		if rune != tt.rune || name != tt.name ||
			!reflect.DeepEqual(words, tt.words) {
			t.Errorf("\nParseLine(%q)\n-> (%q, %q, %q)", // ➎
				tt.line, rune, name, words)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct { // ➊
		haystack []string
		needle string
		want  bool
	}{ // ➋
		{[]string{"A", "B"}, "B", true},
		{[]string{}, "A", false},
		{[]string{"A", "B"}, "Z", false}, // ➌
	} // ➍
	for _, tt := range tests { // ➎
		got := contains(tt.haystack, tt.needle) // ➏
		if got != tt.want {                 // ➐
			t.Errorf("contains(%#v, %#v) want: %v; got: %v",
				tt.haystack, tt.needle, tt.want, got) // ➑
		}
	}
}

func TestContainsAll(t *testing.T) {
	tests := []struct { // ➊
		slice      []string
		needles []string
		want   bool
	}{ // ➋
		{[]string{"A", "B"}, []string{"B"}, true},
		{[]string{}, []string{"A"}, false},
		{[]string{"A"}, []string{}, true}, // ➌
		{[]string{"A", "B"}, []string{"Z"}, false},
		{[]string{"A", "B", "C"}, []string{"A", "C"}, true},
		{[]string{"A", "B", "C"}, []string{"A", "Z"}, false},
		{[]string{"A", "B"}, []string{"A", "B", "C"}, false},
	}
	for _, tt := range tests {
		got := containsAll(tt.slice, tt.needles) // ➍
		if got != tt.want {
			t.Errorf("containsAll(%#v, %#v)\nwant: %v; got: %v",
				tt.slice, tt.needles, tt.want, got) // ➎
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		text    string
		want []string
	}{
		{"A", []string{"A"}},
		{"A B", []string{"A", "B"}},
		{"A B-C", []string{"A", "B", "C"}},
	}
	for _, tt := range tests {
		got := split(tt.text)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("split(%q)\nwant: %#v; got: %#v",
				tt.text, tt.want, got)
		}
	}
}

func ExampleList() {
	text := strings.NewReader(lines3Dto43)
	List(text, "MARK")
	// Output: U+003F	?	QUESTION MARK
}

func ExampleList_2Results() {
	text := strings.NewReader(lines3Dto43)
	List(text, "SIGN")
	// Output:
	// U+003D	=	EQUALS SIGN
	// U+003E	>	GREATER-THAN SIGN
}

func ExampleList_2Words() {
	text := strings.NewReader(lines3Dto43)
	List(text, "CAPITAL LATIN")
	// Output:
	// U+0041	A	LATIN CAPITAL LETTER A
	// U+0042	B	LATIN CAPITAL LETTER B
	// U+0043	C	LATIN CAPITAL LETTER C
}

func Example() {
	argsBefore := os.Args
	defer func() { os.Args = argsBefore }()
	os.Args = []string{"", "cruzeiro"}
	main()
	// Output:
	// U+20A2	₢	CRUZEIRO SIGN
}

func Example_2WordQuery() { // ➊
	argsBefore := os.Args // ➋
	defer func() { os.Args = argsBefore }()
	os.Args = []string{"", "cat", "smiling"}
	main() // ➌
	// Output:
	// U+1F638	😸	GRINNING CAT FACE WITH SMILING EYES
	// U+1F63A	😺	SMILING CAT FACE WITH OPEN MOUTH
	// U+1F63B	😻	SMILING CAT FACE WITH HEART-SHAPED EYES
}

func Example_queryWithHiphenAndField10() {
	argsBefore := os.Args
	defer func() { os.Args = argsBefore }()
	os.Args = []string{"", "quote"}
	main()
	// Output:
	// U+0027	'	APOSTROPHE (APOSTROPHE-QUOTE)
	// U+2358	⍘	APL FUNCTIONAL SYMBOL QUOTE UNDERBAR
	// U+235E	⍞	APL FUNCTIONAL SYMBOL QUOTE QUAD
}

func restore(nameVar, value string, existed bool) {
	if existed {
		os.Setenv(nameVar, value)
	} else {
		os.Unsetenv(nameVar)
	}
}

func TestGetUCDPath_isSet(t *testing.T) {
	pathBefore, existed := os.LookupEnv("UCD_PATH")
	defer restore("UCD_PATH", pathBefore, existed)
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	os.Setenv("UCD_PATH", ucdPath)
	got := getUCDPath()
	if got != ucdPath {
		t.Errorf("getUCDPath() [setado]\nwant: %q; got: %q", ucdPath, got)
	}
}

func TestGetUCDPath_default(t *testing.T) {
	pathBefore, existed := os.LookupEnv("UCD_PATH")
	defer restore("UCD_PATH", pathBefore, existed)
	os.Unsetenv("UCD_PATH")
	ucdPathSuffix := "/UnicodeData.txt"
	got := getUCDPath()
	if !strings.HasSuffix(got, ucdPathSuffix) {
		t.Errorf("getUCDPath() [default]\nwant (sufixo): %q; got: %q", ucdPathSuffix, got)
	}
}

func TestFetchUCD(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(lines3Dto43))
		}))
	defer srv.Close()

	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	done := make(chan bool)
	go fetchUCD(srv.URL, ucdPath, done)
	_ = <-done
	ucd, err := os.Open(ucdPath)
	if os.IsNotExist(err) {
		t.Errorf("fetchUCD não gerou:%v\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}

func TestOpenUCD_local(t *testing.T) {
	ucdPath := getUCDPath()
	ucd, err := openUCD(ucdPath)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", ucdPath, err)
	}
	ucd.Close()
}

func TestOpenUCD_remote(t *testing.T) {
	if testing.Short() {
		t.Skip("teste ignorado [opção -test.short]")
	}
	ucdPath := fmt.Sprintf("./TEST%d-UnicodeData.txt", time.Now().UnixNano())
	ucd, err := openUCD(ucdPath)
	if err != nil {
		t.Errorf("AbrirUCD(%q):\n%v", ucdPath, err)
	}
	ucd.Close()
	os.Remove(ucdPath)
}
