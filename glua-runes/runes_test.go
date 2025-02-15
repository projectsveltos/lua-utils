package gluarunes_test

import (
	"fmt"
	"testing"
	"unicode"
	"unicode/utf8"

	gluarunes "github.com/projectsveltos/lua-utils/glua-runes"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestBytesToRune(t *testing.T) {
	f := func(r rune) *rune {
		return &r
	}

	tests := []struct {
		input    []byte
		expected *rune
	}{
		{[]byte(fmt.Sprintf("%c", 'A')), f('A')},
		{[]byte(fmt.Sprintf("%c", '你')), f('你')},
		{[]byte(fmt.Sprintf("%c", '😀')), f('😀')},
		{[]byte(fmt.Sprintf("%c", 'é')), f('é')},
		{[]byte{}, nil},
		{[]byte{255}, nil},
		{[]byte{255, 254, 253}, nil},
		{[]byte{228}, nil},
		{[]byte{228, 189}, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			table := L.NewTable()
			for _, b := range tt.input {
				table.Append(lua.LNumber(b))
			}

			L.Push(table)

			gluarunes.BytesToRune(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			if tt.expected == nil {
				require.Equal(t, lua.LNil, result)
			} else {
				num, ok := result.(lua.LNumber)
				require.True(t, ok, "expected number return value")
				require.Equal(t, int64(*tt.expected), int64(num))
			}

			L.Pop(1)
		})
	}
}

func TestBytesToString(t *testing.T) {
	tests := []struct {
		input    []byte
		expected string
	}{
		{[]byte{}, ""},
		{[]byte{'A'}, "A"},
		{[]byte{'H', 'e', 'l', 'l', 'o'}, "Hello"},
		{[]byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}, "你好"},
		{[]byte{0xf0, 0x9f, 0x98, 0x80}, "😀"},
		{[]byte{0xc3, 0xa9}, "é"},
		{[]byte{255}, "\xff"},
		{[]byte{255, 254, 253}, "\xff\xfe\xfd"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			table := L.NewTable()
			for _, b := range tt.input {
				table.Append(lua.LNumber(b))
			}

			L.Push(table)

			gluarunes.BytesToString(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			str, ok := result.(lua.LString)
			require.True(t, ok, "expected string return value")
			require.Equal(t, tt.expected, string(str))

			L.Pop(1)
		})
	}
}

func TestContainsRune(t *testing.T) {
	tests := []struct {
		input    string
		search   rune
		expected bool
	}{
		{"Hello", 'H', true},
		{"Hello", 'l', true},
		{"Hello", 'x', false},
		{"你好", '你', true},
		{"你好", '他', false},
		{"Hello你好", '你', true},
		{"Hello你好", 'H', true},
		{"Hello你好", 'x', false},
		{"😀😃😄", '😃', true},
		{"😀😃😄", '😅', false},
		{"café", 'é', true},
		{"cafe", 'é', false},
		{"", 'a', false},
		{" ", ' ', true},
		{"∀x∈ℝ", '∈', true},
		{"∀x∈ℝ", '∉', false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s/search_%d", i, tt.input, tt.search), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))
			L.Push(lua.LNumber(tt.search))

			gluarunes.ContainsRune(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			contains, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")
			require.Equal(t, tt.expected, bool(contains))

			L.Pop(1)
		})
	}
}

func TestIsControl(t *testing.T) {
	tests := []rune{
		0x00,
		0x01,
		0x02,
		0x03,
		0x04,
		0x05,
		0x06,
		0x07,
		0x08,
		0x09,
		0x0A,
		0x0B,
		0x0C,
		0x0D,
		0x0E,
		0x0F,
		0x10,
		0x1F,
		0x7F,
		0x9F,
		'A',
		'1',
		' ',
		'你',
		'好',
		'😀',
		'\u0085',
		'\u009F',
		'\u2028',
		'\u2029',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsControl(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isControl, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsControl(tt), bool(isControl))

			L.Pop(1)
		})
	}
}

func TestIsDigit(t *testing.T) {
	tests := []rune{
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9',
		'A',
		'z',
		'你',
		'好',
		'😀',
		' ',
		'-',
		'\n',
		'\t',
		'\u0000',
		'\u0660',
		'\u06F0',
		'\u0966',
		'\u09E6',
		'\u0CE6',
		'\u0E50',
		'\uFF10',
		'\u2070',
		'\u2080',
		'\u24EA',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsDigit(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isDigit, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsDigit(tt), bool(isDigit))

			L.Pop(1)
		})
	}
}

func TestIsGraphic(t *testing.T) {
	tests := []rune{
		'A',
		'1',
		'.',
		'你',
		'好',
		'😀',
		'é',
		'$',
		'@',
		'[',
		']',
		' ',
		'\t',
		'\n',
		'\r',
		'\u0000',
		'\u0002',
		'\u0010',
		'\u001F',
		'\u007F',
		'\u0080',
		'\u00A0',
		'\u2000',
		'\u2028',
		'\u2029',
		'\u202F',
		'\u205F',
		'\u2060',
		'\u3000',
		'\uFEFF',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsGraphic(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isGraphic, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsGraphic(tt), bool(isGraphic))

			L.Pop(1)
		})
	}
}

func TestIsInRange(t *testing.T) {
	tests := []struct {
		value    rune
		lo       rune
		hi       rune
		expected bool
	}{
		{'A', 'A', 'Z', true},
		{'Z', 'A', 'Z', true},
		{'M', 'A', 'Z', true},
		{'a', 'A', 'Z', false},
		{'1', '0', '9', true},
		{'5', '0', '9', true},
		{'9', '0', '9', true},
		{'A', '0', '9', false},
		{'你', '你', '好', true},
		{'您', '你', '好', false},
		{'好', '你', '好', true},
		{'A', '你', '好', false},
		{'😀', '😀', '😃', true},
		{'😂', '😀', '😃', true},
		{'😃', '😀', '😃', true},
		{'A', '😀', '😃', false},
		{0, 0, 10, true},
		{5, 0, 10, true},
		{10, 0, 10, true},
		{11, 0, 10, false},
		{-1, -10, 0, true},
		{-5, -10, 0, true},
		{-10, -10, 0, true},
		{1, -10, 0, false},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/value_%d/lo_%d/hi_%d", i, tt.value, tt.lo, tt.hi), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.value))
			L.Push(lua.LNumber(tt.lo))
			L.Push(lua.LNumber(tt.hi))

			gluarunes.IsInRange(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			inRange, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, tt.expected, bool(inRange))

			L.Pop(1)
		})
	}
}

func TestIsLetter(t *testing.T) {
	tests := []rune{
		'A',
		'z',
		'é',
		'你',
		'好',
		'😀',
		' ',
		'1',
		'-',
		'\n',
		'\t',
		'\u0000',
		'\u00B5',
		'\u00BA',
		'\u01C5',
		'\u0294',
		'\u0988',
		'\u0939',
		'\u0CA0',
		'\u0E01',
		'\u3042',
		'\u30A2',
		'\u3105',
		'\uAC00',
		'\uFB00',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsLetter(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isLetter, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsLetter(tt), bool(isLetter))

			L.Pop(1)
		})
	}
}

func TestIsLower(t *testing.T) {
	tests := []rune{
		'a',
		'b',
		'z',
		'A',
		'B',
		'Z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u00E0',
		'\u00E1',
		'\u00E2',
		'\u00E3',
		'\u00E4',
		'\u00E5',
		'\u00E6',
		'\u00E7',
		'\u00E8',
		'\u00E9',
		'\u00EA',
		'\u00EB',
		'\u00EC',
		'\u00ED',
		'\u00EE',
		'\u00EF',
		'\u0430',
		'\u0431',
		'\u0432',
		'\u0433',
		'\u0434',
		'\u0435',
		'\u0436',
		'\u0437',
		'\u0438',
		'\u0439',
		'\uFF41',
		'\uFF42',
		'\uFF43',
		'\uFF44',
		'\uFF45',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsLower(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isLower, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsLower(tt), bool(isLower))

			L.Pop(1)
		})
	}
}

func TestIsMark(t *testing.T) {
	tests := []rune{
		'\u0300',
		'\u0301',
		'\u0302',
		'\u0303',
		'\u0304',
		'\u0305',
		'\u0306',
		'\u0307',
		'\u0308',
		'\u0309',
		'\u030A',
		'\u030B',
		'\u030C',
		'\u030D',
		'\u030E',
		'\u030F',
		'\u0310',
		'\u0311',
		'\u0312',
		'\u0313',
		'\u0314',
		'\u0315',
		'\u0316',
		'\u0317',
		'\u0318',
		'\u0319',
		'\u031A',
		'\u031B',
		'\u031C',
		'\u031D',
		'A',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u0020',
		'\u0041',
		'\u0061',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsMark(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isMark, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsMark(tt), bool(isMark))

			L.Pop(1)
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []rune{
		'0',
		'1',
		'9',
		'A',
		'z',
		'你',
		'好',
		'😀',
		' ',
		'-',
		'\n',
		'\u0000',
		'\u00B2',
		'\u00B3',
		'\u00B9',
		'\u0660',
		'\u0661',
		'\u0662',
		'\u0663',
		'\u06F0',
		'\u06F1',
		'\u06F2',
		'\u0966',
		'\u0967',
		'\u0968',
		'\u09E6',
		'\u09E7',
		'\u09E8',
		'\u0BE6',
		'\u0BE7',
		'\u0C66',
		'\u0C67',
		'\u0CE6',
		'\u0CE7',
		'\u0D66',
		'\u0D67',
		'\u0E50',
		'\u0E51',
		'\u0ED0',
		'\u0ED1',
		'\u0F20',
		'\u0F21',
		'\u2070',
		'\u2071',
		'\u2074',
		'\u2075',
		'\u2080',
		'\u2081',
		'\u2460',
		'\u2461',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsNumber(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isNumber, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsNumber(tt), bool(isNumber))

			L.Pop(1)
		})
	}
}

func TestIsPrint(t *testing.T) {
	tests := []rune{
		'A',
		'1',
		'.',
		'你',
		'好',
		'😀',
		'é',
		' ',
		'\t',
		'\n',
		'\r',
		'\u0000',
		'\u0001',
		'\u0002',
		'\u0003',
		'\u0004',
		'\u0005',
		'\u0006',
		'\u0007',
		'\u0008',
		'\u0009',
		'\u000A',
		'\u000B',
		'\u000C',
		'\u000D',
		'\u000E',
		'\u000F',
		'\u0010',
		'\u007F',
		'\u0080',
		'\u0081',
		'\u0082',
		'\u0083',
		'\u0084',
		'\u0085',
		'\u0086',
		'\u0087',
		'\u0088',
		'\u0089',
		'\u008A',
		'\u008B',
		'\u008C',
		'\u008D',
		'\u008E',
		'\u008F',
		'\u0090',
		'\u2028',
		'\u2029',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsPrint(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isPrint, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsPrint(tt), bool(isPrint))

			L.Pop(1)
		})
	}
}

func TestIsPunct(t *testing.T) {
	tests := []rune{
		'.',
		',',
		';',
		':',
		'!',
		'?',
		'"',
		'\'',
		'(',
		')',
		'[',
		']',
		'{',
		'}',
		'<',
		'>',
		'-',
		'_',
		'/',
		'\\',
		'@',
		'#',
		'$',
		'%',
		'&',
		'*',
		'+',
		'=',
		'~',
		'^',
		'`',
		'|',
		'A',
		'1',
		' ',
		'\n',
		'\u0000',
		'\u00A1',
		'\u00BF',
		'\u2010',
		'\u2011',
		'\u2012',
		'\u2013',
		'\u2014',
		'\u2018',
		'\u2019',
		'\u201C',
		'\u201D',
		'\u2026',
		'\u3001',
		'\u3002',
		'\uFF01',
		'\uFF0C',
		'\uFF0E',
		'\uFF1F',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsPunct(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isPunct, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsPunct(tt), bool(isPunct))

			L.Pop(1)
		})
	}
}

func TestIsSpace(t *testing.T) {
	tests := []rune{
		' ',
		'\t',
		'\n',
		'\r',
		'\v',
		'\f',
		'\u0085',
		'\u00A0',
		'\u2000',
		'\u3000',
		'A',
		'1',
		'你',
		'😀',
		'-',
		0,
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsSpace(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isSpace, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsSpace(tt), bool(isSpace))

			L.Pop(1)
		})
	}
}

func TestIsSymbol(t *testing.T) {
	tests := []rune{
		'+',
		'=',
		'<',
		'>',
		'^',
		'$',
		'¢',
		'£',
		'¥',
		'€',
		'©',
		'®',
		'™',
		'°',
		'±',
		'×',
		'÷',
		'∀',
		'∂',
		'∃',
		'∅',
		'∇',
		'∈',
		'∉',
		'∋',
		'∏',
		'∑',
		'√',
		'∝',
		'∞',
		'∠',
		'∧',
		'∨',
		'∩',
		'∪',
		'∫',
		'∴',
		'∼',
		'≅',
		'≈',
		'≠',
		'≡',
		'≤',
		'≥',
		'⊂',
		'⊃',
		'⊄',
		'⊆',
		'⊇',
		'⊕',
		'⊗',
		'⊥',
		'⋅',
		'⌈',
		'⌉',
		'⌊',
		'⌋',
		'A',
		'1',
		'.',
		' ',
		'\n',
		'\u0000',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsSymbol(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isSymbol, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsSymbol(tt), bool(isSymbol))

			L.Pop(1)
		})
	}
}

func TestIsTitle(t *testing.T) {
	tests := []rune{
		'A',
		'a',
		'Z',
		'z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u01C5',
		'\u01C8',
		'\u01CB',
		'\u01F2',
		'\u1F88',
		'\u1F89',
		'\u1F8A',
		'\u1F8B',
		'\u1F8C',
		'\u1F8D',
		'\u1F8E',
		'\u1F8F',
		'\u1F98',
		'\u1F99',
		'\u1F9A',
		'\u1F9B',
		'\u1F9C',
		'\u1F9D',
		'\u1F9E',
		'\u1F9F',
		'\u1FA8',
		'\u1FA9',
		'\u1FAA',
		'\u1FAB',
		'\u1FAC',
		'\u1FAD',
		'\u1FAE',
		'\u1FAF',
		'\u0000',
		'\u0041',
		'\u0061',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsTitle(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isTitle, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsTitle(tt), bool(isTitle))

			L.Pop(1)
		})
	}
}

func TestIsUpper(t *testing.T) {
	tests := []rune{
		'A',
		'B',
		'Z',
		'a',
		'b',
		'z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u0041',
		'\u0042',
		'\u0043',
		'\u0044',
		'\u0045',
		'\u0046',
		'\u0047',
		'\u0048',
		'\u0049',
		'\u004A',
		'\u004B',
		'\u004C',
		'\u004D',
		'\u004E',
		'\u004F',
		'\u0391',
		'\u0392',
		'\u0393',
		'\u0394',
		'\u0395',
		'\u0396',
		'\u0397',
		'\u0398',
		'\u0399',
		'\u039A',
		'\u039B',
		'\uFF21',
		'\uFF22',
		'\uFF23',
		'\uFF24',
		'\uFF25',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.IsUpper(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			isUpper, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, unicode.IsUpper(tt), bool(isUpper))

			L.Pop(1)
		})
	}
}

func TestIsValidUTF8(t *testing.T) {
	tests := []string{
		"",
		"Hello",
		"你好",
		"café",
		"😀",
		"\xed\xa0\x80",
		"\xff",
		string([]byte{0xff, 0xfe, 0xfd}),
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt))

			gluarunes.IsValidUTF8(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			valid, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, utf8.ValidString(tt), bool(valid))

			L.Pop(1)
		})
	}
}

func TestReverseRunes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"a", "a"},
		{"ab", "ba"},
		{"abc", "cba"},
		{"Hello", "olleH"},
		{"你好", "好你"},
		{"Hello你好", "好你olleH"},
		{"😀😃😄", "😄😃😀"},
		{"café", "éfac"},
		{"Hello, 世界！", "！界世 ,olleH"},
		{"🌟star✨", "✨rats🌟"},
		{"Go语言", "言语oG"},
		{" ", " "},
		{"    ", "    "},
		{"a b c", "c b a"},
		{"汉字漢字", "字漢字汉"},
		{"12345", "54321"},
		{"!@#$%", "%$#@!"},
		{"Hello\nWorld", "dlroW\nolleH"},
		{"société", "étéicos"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt.input), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluarunes.ReverseRunes(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			str, ok := result.(lua.LString)
			require.True(t, ok, "expected string return value")
			require.Equal(t, tt.expected, string(str))

			L.Pop(1)
		})
	}
}

func TestRuneAt(t *testing.T) {
	f := func(r rune) *rune {
		return &r
	}

	tests := []struct {
		input    string
		pos      int
		expected *rune
	}{
		{"Hello", 1, f('H')},
		{"Hello", 2, f('e')},
		{"你好", 1, f('你')},
		{"Hi你", 3, f('你')},
		{"😀", 1, f('😀')},
		{"", 1, nil},
		{"Hello", 6, nil},
		{"Hello", -1, nil},
		{"Hello", -100500, nil},
		{"Hello", 10, nil},
		{"café", 4, f('é')},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s/pos_%d", i, tt.input, tt.pos), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))
			L.Push(lua.LNumber(tt.pos))

			gluarunes.RuneAt(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)

			if tt.expected == nil {
				require.Equal(t, lua.LNil, result)
			} else {
				num, ok := result.(lua.LNumber)
				require.True(t, ok, "expected number return value")
				require.Equal(t, int64(*tt.expected), int64(num))
			}

			L.Pop(1)
		})
	}
}

func TestRuneCount(t *testing.T) {
	tests := []string{
		``,
		`A`,
		`Hello`,
		`你好`,
		`Hi你`,
		`😀`,
		`Hello 你好 😀`,
		`\u0041`,
		`café`,
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt))

			gluarunes.RuneCount(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			count, ok := result.(lua.LNumber)
			require.True(t, ok, "expected number return value")

			require.Equal(t, utf8.RuneCountInString(tt), int(count))

			L.Pop(1)
		})
	}
}

func TestRuneIndex(t *testing.T) {
	f := func(i int) *int {
		return &i
	}

	tests := []struct {
		input    string
		search   rune
		start    int
		expected *int
	}{
		{"Hello", 'H', 1, f(1)},
		{"Hello", 'l', 1, f(3)},
		{"Hello", 'o', 1, f(5)},
		{"Hello", 'x', 1, nil},
		{"你好", '你', 1, f(1)},
		{"你好", '好', 1, f(2)},
		{"你好", '们', 1, nil},
		{"Hello你好", '你', 1, f(6)},
		{"Hello你好", 'l', 4, f(4)},
		{"Hello", 'l', 4, f(4)},
		{"Hello", 'l', 5, nil},
		{"", 'a', 1, nil},
		{"Hello", 'H', 2, nil},
		{"Hello", 'H', 0, f(1)},
		{"Hello", 'H', -1, f(1)},
		{"Hello", 'H', -100500, f(1)},
		{"😀😃😄", '😃', 1, f(2)},
		{"café", 'é', 1, f(4)},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s/search_%d/start_%d", i, tt.input, tt.search, tt.start), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))
			L.Push(lua.LNumber(tt.search))
			L.Push(lua.LNumber(tt.start))

			gluarunes.RuneIndex(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			if tt.expected == nil {
				require.Equal(t, lua.LNil, result)
			} else {
				num, ok := result.(lua.LNumber)
				require.True(t, ok, "expected number return value")
				require.Equal(t, *tt.expected, int(num))
			}

			L.Pop(1)
		})
	}
}

func TestRuneRange(t *testing.T) {
	tests := []struct {
		input    string
		start    int
		end      int
		expected string
	}{
		{"Hello", 1, 3, "He"},
		{"Hello", 0, 3, "He"},
		{"Hello", 2, 10, "ello"},
		{"Hello", 3, 5, "ll"},
		{"Hello", 1, -1, "Hello"},
		{"Hello", 1, -100500, "Hello"},
		{"你好世界", 1, 3, "你好"},
		{"", 1, 1, ""},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s/start_%d/end_%d", i, tt.input, tt.start, tt.end), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))
			L.Push(lua.LNumber(tt.start))
			L.Push(lua.LNumber(tt.end))

			gluarunes.RuneRange(L)

			result := L.ToString(-1)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}

			L.Pop(1)
		})
	}
}

func TestRuneSlice(t *testing.T) {
	tests := []string{
		``,
		`A`,
		`Hello`,
		`你好`,
		`Hi你`,
		`😀`,
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt))

			gluarunes.RuneSlice(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			tbl, ok := result.(*lua.LTable)
			require.True(t, ok, "expected table return value")

			expected := []rune(tt)
			require.Equal(t, len(expected), tbl.Len())

			for j, r := range expected {
				val := tbl.RawGetInt(j + 1)

				num, ok := val.(lua.LNumber)
				require.True(t, ok, "expected number in table")
				require.Equal(t, int64(r), int64(num))
			}

			L.Pop(1)
		})
	}
}

func TestRuneSplit(t *testing.T) {
	tests := []struct {
		input    string
		sep      rune
		expected []string
	}{
		{"a,b,c", ',', []string{"a", "b", "c"}},
		{"hello world", ' ', []string{"hello", "world"}},
		{"one", ',', []string{"one"}},
		{"", ',', []string{""}},
		{"你,好,世,界", ',', []string{"你", "好", "世", "界"}},
		{"hello你好world", '你', []string{"hello", "好world"}},
		{"a😀b😀c", '😀', []string{"a", "b", "c"}},
		{"cafété", 'é', []string{"caf", "t", ""}},
		{",,a,,b,,", ',', []string{"", "", "a", "", "b", "", ""}},
		{"🌟star🌟light🌟", '🌟', []string{"", "star", "light", ""}},
		{"no-split-char", 'x', []string{"no-split-char"}},
		{" ", ' ', []string{"", ""}},
		{"世界世界世", '世', []string{"", "界", "界", ""}},
		{"  a  b  c  ", ' ', []string{"", "", "a", "", "b", "", "c", "", ""}},
		{"hello世界goodbye世界", '世', []string{"hello", "界goodbye", "界"}},
		{"🎈party🎈time🎈end", '🎈', []string{"", "party", "time", "end"}},
		{"e\u0301", '\u0301', []string{"e", ""}},
		{"\u200Ba\u200Bb\u200B", '\u200B', []string{"", "a", "b", ""}},
		{"\na\nb\n", '\n', []string{"", "a", "b", ""}},
		{"\ta\tb\t", '\t', []string{"", "a", "b", ""}},
		{"∀x∈ℝ", '∈', []string{"∀x", "ℝ"}},
		{"吉🈲な🈲の", '🈲', []string{"吉", "な", "の"}},
		{"aaaaaaaaaaaaaaa,bbbbbbbbbbbbbbb", ',', []string{"aaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbb"}},
		{"a,b,c,d,e,f,g", ',', []string{"a", "b", "c", "d", "e", "f", "g"}},
		{",,,,", ',', []string{"", "", "", "", ""}},
		{"a\nb\nc", '\n', []string{"a", "b", "c"}},
		{"a\rb\rc", '\r', []string{"a", "b", "c"}},
		{"a\r\nb\r\nc", '\n', []string{"a\r", "b\r", "c"}},
		{"a\nb\rc\n", '\n', []string{"a", "b\rc", ""}},
		{"a\r\nb\rc\n\r\n", '\n', []string{"a\r", "b\rc", "\r", ""}},
		{"café\u0301", 'é', []string{"caf", "\u0301"}},
		{"𝄞music𝄞notes𝄞", '𝄞', []string{"", "music", "notes", ""}},
		{"١،٢،٣", '،', []string{"١", "٢", "٣"}},
		{"一,二,三", ',', []string{"一", "二", "三"}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s/sep_%d", i, tt.input, tt.sep), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))
			L.Push(lua.LNumber(tt.sep))

			gluarunes.RuneSplit(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			tbl, ok := result.(*lua.LTable)
			require.True(t, ok, "expected table return value")

			require.Equal(t, len(tt.expected), tbl.Len())

			for j, expected := range tt.expected {
				val := tbl.RawGetInt(j + 1)

				str, ok := val.(lua.LString)
				require.True(t, ok, "expected string in table")
				require.Equal(t, expected, string(str))
			}

			L.Pop(1)
		})
	}
}

func TestRuneString(t *testing.T) {
	f := func(s string) []int {
		r := []rune(s)
		n := make([]int, len(r))

		for i, v := range r {
			n[i] = int(v)
		}

		return n
	}

	tests := []string{
		``,
		`A`,
		`Hello`,
		`你好`,
		`Hi你`,
		`😀`,
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			in := f(tt)

			for _, v := range in {
				L.Push(lua.LNumber(v))
			}

			gluarunes.RuneString(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.ToString(-1)
			require.Equal(t, tt, result)

			L.Pop(1)
		})
	}
}

func TestRuneToBytes(t *testing.T) {
	tests := []struct {
		input    []rune
		expected []byte
	}{
		{[]rune{'A'}, []byte{'A'}},
		{[]rune{'H', 'i'}, []byte{'H', 'i'}},
		{[]rune{'你'}, []byte{0xe4, 0xbd, 0xa0}},
		{[]rune{'好'}, []byte{0xe5, 0xa5, 0xbd}},
		{[]rune{'你', '好'}, []byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}},
		{[]rune{'😀'}, []byte{0xf0, 0x9f, 0x98, 0x80}},
		{[]rune{'é'}, []byte{0xc3, 0xa9}},
		{[]rune{}, []byte{}},
		{[]rune{'H', '你', '😀'}, []byte{0x48, 0xe4, 0xbd, 0xa0, 0xf0, 0x9f, 0x98, 0x80}},
		{[]rune{0x20AC}, []byte{0xe2, 0x82, 0xac}},
		{[]rune{0x0000}, []byte{0x00}},
		{[]rune{0x007F}, []byte{0x7F}},
		{[]rune{0x80}, []byte{0xc2, 0x80}},
		{[]rune{0x7FF}, []byte{0xdf, 0xbf}},
		{[]rune{0x800}, []byte{0xe0, 0xa0, 0x80}},
		{[]rune{0xFFFF}, []byte{0xef, 0xbf, 0xbf}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			table := L.NewTable()
			for _, r := range tt.input {
				table.Append(lua.LNumber(r))
			}

			L.Push(table)

			gluarunes.RuneToBytes(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			resultTable, ok := result.(*lua.LTable)
			require.True(t, ok, "expected table return value")

			require.Equal(t, len(tt.expected), resultTable.Len())

			bytes := make([]byte, 0, resultTable.Len())
			resultTable.ForEach(func(_, v lua.LValue) {
				num, ok := v.(lua.LNumber)
				require.True(t, ok, "expected number in table")

				bytes = append(bytes, byte(num))
			})

			require.Equal(t, tt.expected, bytes)

			L.Pop(1)
		})
	}
}

func TestRuneWidth(t *testing.T) {
	tests := []rune{
		'A',
		'你',
		'😀',
		'é',
		'\u0000',
		'\uffff',
		-1,
		0x110000,
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.RuneWidth(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			expected := utf8.RuneLen(tt)

			num, ok := result.(lua.LNumber)
			if expected >= 0 {
				require.True(t, ok, "expected number return value")
				require.Equal(t, utf8.RuneLen(tt), int(num))
			} else {
				require.Equal(t, 0, int(num))
			}

			L.Pop(1)
		})
	}
}

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
	}{
		{"", []byte{}},
		{"A", []byte{'A'}},
		{"Hello", []byte{'H', 'e', 'l', 'l', 'o'}},
		{"你", []byte{0xe4, 0xbd, 0xa0}},
		{"好", []byte{0xe5, 0xa5, 0xbd}},
		{"你好", []byte{0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}},
		{"😀", []byte{0xf0, 0x9f, 0x98, 0x80}},
		{"café", []byte{0x63, 0x61, 0x66, 0xc3, 0xa9}},
		{"Hello你好😀", []byte{
			0x48, 0x65, 0x6c, 0x6c, 0x6f,
			0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd,
			0xf0, 0x9f, 0x98, 0x80,
		}},
		{"\x00", []byte{0x00}},
		{"\x7F", []byte{0x7F}},
		{"\u0080", []byte{0xc2, 0x80}},
		{"\u07FF", []byte{0xdf, 0xbf}},
		{"\u0800", []byte{0xe0, 0xa0, 0x80}},
		{"\uffff", []byte{0xef, 0xbf, 0xbf}},
		{"€", []byte{0xe2, 0x82, 0xac}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/string_%s", i, tt.input), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluarunes.StringToBytes(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			resultTable, ok := result.(*lua.LTable)
			require.True(t, ok, "expected table return value")

			require.Equal(t, len(tt.expected), resultTable.Len())

			bytes := make([]byte, 0, resultTable.Len())
			resultTable.ForEach(func(_, v lua.LValue) {
				num, ok := v.(lua.LNumber)
				require.True(t, ok, "expected number in table")

				bytes = append(bytes, byte(num))
			})

			require.Equal(t, tt.expected, bytes)

			L.Pop(1)
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []rune{
		'A',
		'B',
		'Z',
		'a',
		'b',
		'z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u0041',
		'\u0042',
		'\u0043',
		'\u0391',
		'\u0392',
		'\u0393',
		'\u0394',
		'\u0395',
		'\u0396',
		'\u0397',
		'\u0398',
		'\u0399',
		'\u039A',
		'\u039B',
		'\u039C',
		'\u039D',
		'\u039E',
		'\u039F',
		'\u0400',
		'\u0401',
		'\u0402',
		'\u0403',
		'\u0404',
		'\u0405',
		'\u0406',
		'\u0407',
		'\u0408',
		'\u0409',
		'\u040A',
		'\uFF21',
		'\uFF22',
		'\uFF23',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.ToLower(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			lower, ok := result.(lua.LNumber)
			require.True(t, ok, "expected number return value")

			require.Equal(t, unicode.ToLower(tt), rune(lower))

			L.Pop(1)
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []rune{
		'a',
		'b',
		'z',
		'A',
		'B',
		'Z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u0061',
		'\u0062',
		'\u0063',
		'\u03B1',
		'\u03B2',
		'\u03B3',
		'\u03B4',
		'\u03B5',
		'\u03B6',
		'\u03B7',
		'\u03B8',
		'\u03B9',
		'\u03BA',
		'\u03BB',
		'\u03BC',
		'\u03BD',
		'\u03BE',
		'\u03BF',
		'\u0430',
		'\u0431',
		'\u0432',
		'\u0433',
		'\u0434',
		'\u0435',
		'\u0436',
		'\u0437',
		'\u0438',
		'\u0439',
		'\u043A',
		'\uFF41',
		'\uFF42',
		'\uFF43',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.ToUpper(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			upper, ok := result.(lua.LNumber)
			require.True(t, ok, "expected number return value")

			require.Equal(t, unicode.ToUpper(tt), rune(upper))

			L.Pop(1)
		})
	}
}

func TestToTitle(t *testing.T) {
	tests := []rune{
		'a',
		'b',
		'z',
		'A',
		'B',
		'Z',
		'1',
		'.',
		' ',
		'你',
		'好',
		'😀',
		'\u0000',
		'\u01C5',
		'\u01C8',
		'\u01CB',
		'\u01F2',
		'\u03B1',
		'\u03B2',
		'\u03B3',
		'\u03B4',
		'\u03B5',
		'\u03B6',
		'\u03B7',
		'\u03B8',
		'\u03B9',
		'\u03BA',
		'\u03BB',
		'\u03BC',
		'\u03BD',
		'\u03BE',
		'\u03BF',
		'\u0430',
		'\u0431',
		'\u0432',
		'\u0433',
		'\u0434',
		'\u0435',
		'\u0436',
		'\u0437',
		'\u0438',
		'\u0439',
		'\u043A',
		'\uFF41',
		'\uFF42',
		'\uFF43',
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.ToTitle(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			title, ok := result.(lua.LNumber)
			require.True(t, ok, "expected number return value")

			require.Equal(t, unicode.ToTitle(tt), rune(title))

			L.Pop(1)
		})
	}
}

func TestValidRune(t *testing.T) {
	tests := []rune{
		'A',
		'1',
		'你',
		'好',
		'😀',
		'\u0000',
		'\uFFFF',
		0x10FFFF,
		-1,
		0x110000,
		0x200000,
		0xFFFFFFF,
		rune(1<<31 - 1),
		rune(-1 << 31),
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d/rune_%d", i, tt), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt))

			gluarunes.ValidRune(L)

			if L.GetTop() == 0 {
				t.Fatal("expected a return value, got none")
			}

			result := L.Get(-1)
			valid, ok := result.(lua.LBool)
			require.True(t, ok, "expected boolean return value")

			require.Equal(t, utf8.ValidRune(tt), bool(valid))

			L.Pop(1)
		})
	}
}
