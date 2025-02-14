package gluarunes

import (
	"slices"
	"unicode"
	"unicode/utf8"

	lua "github.com/yuin/gopher-lua"
)

// BytesToRune converts UTF-8 bytes to a rune.
// Takes a table of bytes and returns the corresponding rune value
// as a lua.LNumber, or nil if the bytes are not valid UTF-8.
func BytesToRune(L *lua.LState) int {
	table := L.CheckTable(1)

	bytes := make([]byte, 0, table.Len())
	table.ForEach(func(_, v lua.LValue) {
		if num, ok := v.(lua.LNumber); ok {
			bytes = append(bytes, byte(num))
		}
	})

	r, _ := utf8.DecodeRune(bytes)
	if r == utf8.RuneError {
		L.Push(lua.LNil)

		return 1
	}

	L.Push(lua.LNumber(r))

	return 1
}

// BytesToString converts UTF-8 bytes to a string.
// Takes a table of bytes and returns the corresponding string
// as a lua.LString.
func BytesToString(L *lua.LState) int {
	table := L.CheckTable(1)

	bytes := make([]byte, 0, table.Len())
	table.ForEach(func(_, v lua.LValue) {
		if num, ok := v.(lua.LNumber); ok {
			bytes = append(bytes, byte(num))
		}
	})

	L.Push(lua.LString(string(bytes)))

	return 1
}

// ContainsRune checks if a rune exists in a string.
// Parameters:
//   - string: The input string to search
//   - rune: The rune to search for
//
// Returns a boolean as lua.LBool indicating whether the rune exists in the string.
func ContainsRune(L *lua.LState) int {
	s := L.CheckString(1)
	search := rune(L.CheckInt(2))

	runes := []rune(s)
	L.Push(lua.LBool(slices.Contains(runes, search)))

	return 1
}

// IsControl checks if a rune is a Unicode control character.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsControl(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsControl(r)))

	return 1
}

// IsDigit checks if a rune is a Unicode decimal digit.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsDigit(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsDigit(r)))

	return 1
}

// IsGraphic checks if a rune is graphic.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// but not spaces or control characters.
func IsGraphic(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsGraphic(r)))

	return 1
}

// IsInRange checks if a rune is within a specified range.
// Parameters:
//   - rune: The rune to check
//   - lo: The lower bound of the range (inclusive)
//   - hi: The upper bound of the range (inclusive)
//
// Returns a boolean as lua.LBool indicating whether the rune is within the range.
func IsInRange(L *lua.LState) int {
	r := rune(L.CheckInt(1))
	lo := rune(L.CheckInt(2))
	hi := rune(L.CheckInt(3))

	L.Push(lua.LBool(r >= lo && r <= hi))

	return 1
}

// IsLetter checks if a rune is a Unicode letter.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsLetter(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsLetter(r)))

	return 1
}

// IsLower checks if a rune is a lowercase letter.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsLower(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsLower(r)))

	return 1
}

// IsMark checks if a rune is a Unicode mark character.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsMark(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsMark(r)))

	return 1
}

// IsNumber checks if a rune is a Unicode number (includes characters besides 0-9).
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsNumber(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsNumber(r)))

	return 1
}

// IsPrint checks if a rune is printable.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
// Printable characters include letters, marks, numbers, punctuation, symbols,
// and spaces, but not control characters.
func IsPrint(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsPrint(r)))

	return 1
}

// IsPunct checks if a rune is a Unicode punctuation character.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsPunct(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsPunct(r)))

	return 1
}

// IsSpace checks if a rune is a Unicode white space character.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsSpace(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsSpace(r)))

	return 1
}

// IsSymbol checks if a rune is a Unicode symbol character.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsSymbol(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsSymbol(r)))

	return 1
}

// IsTitle checks if a rune is a Unicode title case letter.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsTitle(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsTitle(r)))

	return 1
}

// IsUpper checks if a rune is an uppercase letter.
// Takes a rune value as an integer and returns a boolean as lua.LBool.
func IsUpper(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(unicode.IsUpper(r)))

	return 1
}

// IsValidUTF8 checks if a string contains valid UTF-8 encoding.
// Takes a string argument and returns a boolean indicating whether
// the string is valid UTF-8 as a lua.LBool.
func IsValidUTF8(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LBool(utf8.ValidString(s)))

	return 1
}

// ReverseRunes reverses the runes in a string.
// Takes a string argument and returns a new string with the runes in reverse order
// as a lua.LString.
func ReverseRunes(L *lua.LState) int {
	s := L.CheckString(1)
	runes := []rune(s)

	slices.Reverse(runes)

	L.Push(lua.LString(string(runes)))

	return 1
}

// RuneAt returns the rune at a specific byte position in a string.
// Parameters:
//   - string: The input string
//   - position: The 1-based index of the desired rune
//
// Returns nil if the position is invalid or the rune is not valid UTF-8,
// otherwise returns the rune value as a lua.LNumber.
func RuneAt(L *lua.LState) int {
	s := L.CheckString(1)
	pos := L.CheckInt(2) - 1

	if pos < 0 || pos >= len(s) {
		L.Push(lua.LNil)

		return 1
	}

	r, _ := utf8.DecodeRuneInString(s[pos:])
	if r == utf8.RuneError {
		L.Push(lua.LNil)

		return 1
	}

	L.Push(lua.LNumber(r))

	return 1
}

// RuneCount returns the number of runes in a string.
// Takes a string argument and returns the count of Unicode code points
// in that string as a lua.LNumber.
func RuneCount(L *lua.LState) int {
	s := L.CheckString(1)
	count := utf8.RuneCountInString(s)

	L.Push(lua.LNumber(count))

	return 1
}

// RuneIndex finds the first occurrence of a rune in a string.
// Parameters:
//   - string: The input string to search
//   - rune: The rune to search for
//   - start: Optional 1-based starting position (defaults to 1)
//
// Returns the 1-based index of the first occurrence as lua.LNumber,
// or nil if the rune is not found.
func RuneIndex(L *lua.LState) int {
	s := L.CheckString(1)
	search := rune(L.CheckInt(2))
	pos := L.OptInt(3, 1) - 1

	runes := []rune(s)

	if pos < 0 {
		pos = 0
	}

	for i := pos; i < len(runes); i++ {
		if runes[i] == search {
			L.Push(lua.LNumber(i + 1))

			return 1
		}
	}

	L.Push(lua.LNil)

	return 1
}

// RuneRange extracts a substring by rune indices.
// Parameters:
//   - string: The input string
//   - start: Optional 1-based start index (defaults to 1)
//   - end: Optional 1-based end index (defaults to -1, meaning end of string)
//
// Returns the substring as lua.LString.
func RuneRange(L *lua.LState) int {
	s := L.CheckString(1)
	start := L.OptInt(2, 1) - 1
	end := L.OptInt(3, -1)

	runes := []rune(s)
	runeLen := len(runes)

	if start < 0 {
		start = 0
	}

	if start > runeLen {
		start = runeLen
	}

	if end < 0 {
		end = runeLen
	} else {
		end--

		if end > runeLen {
			end = runeLen
		}

		if end < 0 {
			end = 0
		}
	}

	result := string(runes[start:end])

	L.Push(lua.LString(result))

	return 1
}

// RuneSlice converts a string to a slice of rune values.
// Takes a string argument and returns a Lua table containing the numeric
// values of each rune in the string.
func RuneSlice(L *lua.LState) int {
	s := L.CheckString(1)
	result := L.NewTable()

	for _, r := range s {
		result.Append(lua.LNumber(r))
	}

	L.Push(result)

	return 1
}

// RuneSplit splits a string on a specified rune delimiter.
// Parameters:
//   - string: The input string to split
//   - separator: The rune to use as the delimiter
//
// Returns a Lua table containing the resulting substrings.
func RuneSplit(L *lua.LState) int {
	s := L.CheckString(1)
	sep := rune(L.CheckInt(2))

	result := L.NewTable()
	runes := []rune(s)

	if len(runes) == 0 {
		result.Append(lua.LString(""))
		L.Push(result)

		return 1
	}

	lastIdx := 0

	for i, r := range runes {
		if r == sep {
			result.Append(lua.LString(string(runes[lastIdx:i])))
			lastIdx = i + 1
		}
	}

	if lastIdx <= len(runes) {
		result.Append(lua.LString(string(runes[lastIdx:])))
	}

	L.Push(result)

	return 1
}

// RuneString converts a slice of integers to a string of runes.
// Each integer argument is converted to a rune and concatenated into a string.
// Returns the resulting string as a lua.LString.
func RuneString(L *lua.LState) int {
	top := L.GetTop()
	runes := make([]rune, top)

	for i := 1; i <= top; i++ {
		runes[i-1] = rune(L.CheckInt(i))
	}

	L.Push(lua.LString(string(runes)))

	return 1
}

// RuneToBytes converts a table of runes to their UTF-8 byte representation.
// Takes a table of rune values as integers and returns a Lua table containing
// the UTF-8 bytes of all runes concatenated together.
func RuneToBytes(L *lua.LState) int {
	table := L.CheckTable(1)

	maxSize := table.Len() * 4
	buf := make([]byte, 0, maxSize)

	tmpBuf := make([]byte, 4)

	table.ForEach(func(_, v lua.LValue) {
		if num, ok := v.(lua.LNumber); ok {
			n := utf8.EncodeRune(tmpBuf, rune(num))
			buf = append(buf, tmpBuf[:n]...)
		}
	})

	result := L.NewTable()
	for _, b := range buf {
		result.Append(lua.LNumber(b))
	}

	L.Push(result)

	return 1
}

// RuneWidth returns the number of bytes needed to encode a rune.
// Takes a rune value as an integer and returns its UTF-8 encoding width
// as a lua.LNumber, or nil if the rune is invalid.
func RuneWidth(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	width := utf8.RuneLen(r)
	if width == -1 {
		L.Push(lua.LNil)

		return 1
	}

	L.Push(lua.LNumber(width))

	return 1
}

// StringToBytes converts a string to its UTF-8 byte representation.
// Takes a string and returns a Lua table containing the UTF-8 bytes.
func StringToBytes(L *lua.LState) int {
	s := L.CheckString(1)

	result := L.NewTable()
	for _, b := range []byte(s) {
		result.Append(lua.LNumber(b))
	}

	L.Push(result)

	return 1
}

// ToLower converts a rune to lowercase.
// Takes a rune value as an integer and returns the lowercase version as lua.LNumber.
func ToLower(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LNumber(unicode.ToLower(r)))

	return 1
}

// ToUpper converts a rune to uppercase.
// Takes a rune value as an integer and returns the uppercase version as lua.LNumber.
func ToUpper(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LNumber(unicode.ToUpper(r)))

	return 1
}

// ToTitle converts a rune to title case.
// Takes a rune value as an integer and returns the title case version as lua.LNumber.
func ToTitle(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LNumber(unicode.ToTitle(r)))

	return 1
}

// ValidRune checks if an integer is a valid Unicode code point.
// Takes an integer value and returns a boolean as lua.LBool indicating
// whether it represents a valid Unicode code point.
func ValidRune(L *lua.LState) int {
	r := rune(L.CheckInt(1))

	L.Push(lua.LBool(utf8.ValidRune(r)))

	return 1
}

// Loader is the module loader function for the runes package.
// It creates a new table and populates it with the package's functions.
func Loader(L *lua.LState) int {
	mod := L.NewTable()

	funcs := map[string]lua.LGFunction{
		"bytestorune":   BytesToRune,
		"bytetostring":  BytesToString,
		"containsrune":  ContainsRune,
		"iscontrol":     IsControl,
		"isdigit":       IsDigit,
		"isgraphic":     IsGraphic,
		"isinrange":     IsInRange,
		"isletter":      IsLetter,
		"islower":       IsLower,
		"ismark":        IsMark,
		"isnumber":      IsNumber,
		"isprint":       IsPrint,
		"ispunct":       IsPunct,
		"isspace":       IsSpace,
		"issymbol":      IsSymbol,
		"istitle":       IsTitle,
		"isupper":       IsUpper,
		"isvalidutf8":   IsValidUTF8,
		"reverserunes":  ReverseRunes,
		"runeat":        RuneAt,
		"runecount":     RuneCount,
		"runeindex":     RuneIndex,
		"runerange":     RuneRange,
		"runeslice":     RuneSlice,
		"runesplit":     RuneSplit,
		"runestring":    RuneString,
		"runetobytes":   RuneToBytes,
		"runewidth":     RuneWidth,
		"stringtobytes": StringToBytes,
		"tolower":       ToLower,
		"totitle":       ToTitle,
		"toupper":       ToUpper,
		"validrune":     ValidRune,
	}

	L.SetFuncs(mod, funcs)
	L.Push(mod)

	return 1
}

// Preload registers the runes package loader function.
// It should be called during Lua state initialization to make the package available.
func Preload(L *lua.LState) {
	L.PreloadModule("runes", Loader)
}
