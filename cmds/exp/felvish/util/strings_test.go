package util

import "testing"

var findContextTests = []struct {
	text          string
	pos           int
	lineno, colno int
	line          string
}{
	{"a\nb", 2, 1, 0, "b"},
}

func TestFindContext(t *testing.T) {
	for _, tt := range findContextTests {
		lineno, colno, line := FindContext(tt.text, tt.pos)
		if lineno != tt.lineno || colno != tt.colno || line != tt.line {
			t.Errorf("FindContext(%v, %v) => (%v, %v, %v), want (%v, %v, %v)",
				tt.text, tt.pos, lineno, colno, line, tt.lineno, tt.colno, tt.line)
		}
	}
}

var SubstringByRuneTests = []struct {
	s         string
	low, high int
	wantedStr string
	wantedErr error
}{
	{"Hello world", 1, 4, "ell", nil},
	{"你好世界", 0, 0, "", nil},
	{"你好世界", 1, 1, "", nil},
	{"你好世界", 1, 2, "好", nil},
	{"你好世界", 1, 4, "好世界", nil},
	{"你好世界", -1, -1, "", ErrIndexOutOfRange},
	{"你好世界", 0, 5, "", ErrIndexOutOfRange},
	{"你好世界", 5, 5, "", ErrIndexOutOfRange},
}

func TestSubstringByRune(t *testing.T) {
	for _, tt := range SubstringByRuneTests {
		s, e := SubstringByRune(tt.s, tt.low, tt.high)
		if s != tt.wantedStr || e != tt.wantedErr {
			t.Errorf("SubstringByRune(%q, %v, %d) => (%q, %v), want (%q, %v)",
				tt.s, tt.low, tt.high, s, e, tt.wantedStr, tt.wantedErr)
		}
	}
}

var NthRuneTests = []struct {
	s          string
	n          int
	wantedRune rune
	wantedErr  error
}{
	{"你好世界", -1, 0, ErrIndexOutOfRange},
	{"你好世界", 0, '你', nil},
	{"你好世界", 4, 0, ErrIndexOutOfRange},
}

func TestNthRune(t *testing.T) {
	for _, tt := range NthRuneTests {
		r, e := NthRune(tt.s, tt.n)
		if r != tt.wantedRune || e != tt.wantedErr {
			t.Errorf("NthRune(%q, %v) => (%q, %v), want (%q, %v)",
				tt.s, tt.n, r, e, tt.wantedRune, tt.wantedErr)
		}
	}
}

var EOLSOLTests = []struct {
	s                         string
	wantFirstEOL, wantLastSOL int
}{
	{"0", 1, 0},
	{"\n12", 0, 1},
	{"01\n", 2, 3},
	{"01\n34", 2, 3},
}

func TestEOLSOL(t *testing.T) {
	for _, tc := range EOLSOLTests {
		eol := FindFirstEOL(tc.s)
		if eol != tc.wantFirstEOL {
			t.Errorf("FindFirstEOL(%q) => %d, want %d", tc.s, eol, tc.wantFirstEOL)
		}
		sol := FindLastSOL(tc.s)
		if sol != tc.wantLastSOL {
			t.Errorf("FindLastSOL(%q) => %d, want %d", tc.s, sol, tc.wantLastSOL)
		}
	}
}

var MatchSubseqTests = []struct {
	s, p string
	want bool
}{
	{"elvish", "e", true},
	{"elvish", "elh", true},
	{"elvish", "sh", true},
	{"elves/elvish", "l/e", true},
	{"elves/elvish", "e/e", true},
	{"elvish", "le", false},
	{"elvish", "evii", false},
}

func TestMatchSubseq(t *testing.T) {
	for _, tc := range MatchSubseqTests {
		b := MatchSubseq(tc.s, tc.p)
		if b != tc.want {
			t.Errorf("MatchSubseq(%q, %q) -> %v, want %v", tc.s, tc.p, b, tc.want)
		}
	}
}
