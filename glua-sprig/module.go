package gluasprig

import (
	"fmt"

	sprig "github.com/Masterminds/sprig/v3"
	lua "github.com/yuin/gopher-lua"
)

func isEmptyLuaValue(value lua.LValue) bool {
	switch value.Type() {
	case lua.LTNil:
		return true
	case lua.LTBool:
		return value == lua.LFalse
	case lua.LTNumber:
		return float64(value.(lua.LNumber)) == 0
	case lua.LTString:
		return len(value.String()) == 0
	case lua.LTTable:
		tbl := value.(*lua.LTable)
		isEmpty := true

		tbl.ForEach(func(_, _ lua.LValue) {
			isEmpty = false
		})

		return isEmpty
	default:
		return false
	}
}

// AbbrevFunc wraps the sprig.abbrev function.
func AbbrevFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("abbrev: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "abbrev requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["abbrev"].(func(int, string) string)
	if !ok {
		L.RaiseError("abbrev: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// AbbrevbothFunc wraps the sprig.abbrevboth function.
func AbbrevbothFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("abbrevboth: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "abbrevboth requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["abbrevboth"].(func(int, int, string) string)
	if !ok {
		L.RaiseError("abbrevboth: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := int(L.CheckNumber(2))
	param2 := L.CheckString(3)
	result := fn(param0, param1, param2)

	L.Push(lua.LString(result))

	return 1
}

// Adler32sumFunc wraps the sprig.adler32sum function.
func Adler32sumFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("adler32sum: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "adler32sum requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["adler32sum"].(func(string) string)
	if !ok {
		L.RaiseError("adler32sum: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// AgoFunc wraps the sprig.ago function.
func AgoFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("ago: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "ago requires 1 argument")

		return 0
	}

	if L.Get(1).Type() != lua.LTNumber {
		L.ArgError(1, "ago requires a number (Unix timestamp)")

		return 0
	}

	fn, ok := sprig.FuncMap()["ago"].(func(any) string)
	if !ok {
		L.RaiseError("ago: invalid function assertion")

		return 0
	}

	timestamp := int64(L.CheckNumber(1))
	result := fn(timestamp)

	L.Push(lua.LString(result))

	return 1
}

// AllFunc implements the sprig.all function.
func AllFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("all: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "all requires 1 arguments")

		return 0
	}

	tbl := L.CheckTable(1)
	result := true

	tbl.ForEach(func(_, v lua.LValue) {
		if isEmptyLuaValue(v) {
			result = false
		}
	})

	L.Push(lua.LBool(result))

	return 1
}

// AnyFunc implements the sprig.any function.
func AnyFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("any: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "any requires 1 arguments")

		return 0
	}

	tbl := L.CheckTable(1)
	result := false

	tbl.ForEach(func(_, v lua.LValue) {
		if !isEmptyLuaValue(v) {
			result = true
		}
	})

	L.Push(lua.LBool(result))

	return 1
}

// B32decFunc wraps the sprig.b32dec function.
func B32decFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("b32dec: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "b32dec requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["b32dec"].(func(string) string)
	if !ok {
		L.RaiseError("b32dec: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// B32encFunc wraps the sprig.b32enc function.
func B32encFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("b32enc: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "b32enc requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["b32enc"].(func(string) string)
	if !ok {
		L.RaiseError("b32enc: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// B64decFunc wraps the sprig.b64dec function.
func B64decFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("b64dec: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "b64dec requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["b64dec"].(func(string) string)
	if !ok {
		L.RaiseError("b64dec: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// B64encFunc wraps the sprig.b64enc function.
func B64encFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("b64enc: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "b64enc requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["b64enc"].(func(string) string)
	if !ok {
		L.RaiseError("b64enc: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// BaseFunc wraps the sprig.base function.
func BaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("base: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "base requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["base"].(func(string) string)
	if !ok {
		L.RaiseError("base: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// BcryptFunc wraps the sprig.bcrypt function.
func BcryptFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("bcrypt: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "bcrypt requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["bcrypt"].(func(string) string)
	if !ok {
		L.RaiseError("bcrypt: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// CamelcaseFunc wraps the sprig.camelcase function.
func CamelcaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("camelcase: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "camelcase requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["camelcase"].(func(string) string)
	if !ok {
		L.RaiseError("camelcase: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// CatFunc wraps the sprig.cat function.
func CatFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("cat: %v", r)
		}
	}()

	top := L.GetTop()
	args := make([]any, 0, top)

	for i := 1; i <= top; i++ {
		v := L.Get(i)
		if v == lua.LNil {
			continue
		}

		args = append(args, v.String())
	}

	fn, ok := sprig.FuncMap()["cat"].(func(...any) string)
	if !ok {
		L.RaiseError("cat: invalid function assertion")

		return 0
	}

	result := fn(args...)

	L.Push(lua.LString(result))

	return 1
}

// CleanFunc wraps the sprig.clean function.
func CleanFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("clean: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "clean requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["clean"].(func(string) string)
	if !ok {
		L.RaiseError("clean: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// CoalesceFunc implements the sprig.coalesce function.
func CoalesceFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("coalesce: %v", r)
		}
	}()

	top := L.GetTop()
	if top < 1 {
		L.ArgError(1, "coalesce requires at least 1 argument")

		return 0
	}

	for i := 1; i <= top; i++ {
		value := L.Get(i)

		if !isEmptyLuaValue(value) {
			L.Push(value)

			return 1
		}
	}

	L.Push(lua.LNil)

	return 1
}

// CompactFunc implements the sprig.compact function.
func CompactFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("compact: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "compact requires 1 argument")

		return 0
	}

	inputTable := L.CheckTable(1)
	resultTable := L.CreateTable(0, 0)
	newIndex := 1

	inputTable.ForEach(func(_, value lua.LValue) {
		if !isEmptyLuaValue(value) {
			resultTable.RawSetInt(newIndex, value)

			newIndex++
		}
	})

	L.Push(resultTable)

	return 1
}

// DecryptAESFunc wraps the sprig.decryptAES function.
func DecryptAESFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("decryptAES: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "decryptAES requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["decryptAES"].(func(string, string) (string, error))
	if !ok {
		L.RaiseError("decryptAES: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)

	result, err := fn(param0, param1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)

	return 2
}

// DerivePasswordFunc wraps the sprig.derivePassword function.
func DerivePasswordFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("derivePassword: %v", r)
		}
	}()

	if L.GetTop() < 5 {
		L.ArgError(1, "derivePassword requires 5 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["derivePassword"].(func(uint32, string, string, string, string) string)
	if !ok {
		L.RaiseError("derivePassword: invalid function assertion")

		return 0
	}

	param0 := uint32(L.CheckNumber(1)) // counter value
	param1 := L.CheckString(2)         // passwordType - the type like "medium", "short", etc...
	param2 := L.CheckString(3)         // password
	param3 := L.CheckString(4)         // username
	param4 := L.CheckString(5)         // site name

	result := fn(param0, param1, param2, param3, param4)

	L.Push(lua.LString(result))

	return 1
}

// DirFunc wraps the sprig.dir function.
func DirFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("dir: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "dir requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["dir"].(func(string) string)
	if !ok {
		L.RaiseError("dir: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// DurationFunc wraps the sprig.duration function.
func DurationFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("duration: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "duration requires 1 argument")

		return 0
	}

	fn, ok := sprig.FuncMap()["duration"].(func(any) string)
	if !ok {
		L.RaiseError("duration: invalid function assertion")

		return 0
	}

	var param any
	switch L.Get(1).Type() {
	case lua.LTNumber:
		param = int64(L.CheckNumber(1))
	case lua.LTString:
		param = L.CheckString(1)
	case lua.LTNil:
		param = int64(0)
	default:
		param = L.Get(1).String()
	}

	result := fn(param)

	L.Push(lua.LString(result))

	return 1
}

// DurationRoundFunc wraps the sprig.durationRound function.
func DurationRoundFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("durationRound: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "durationRound requires 1 argument")

		return 0
	}

	fn, ok := sprig.FuncMap()["durationRound"].(func(any) string)
	if !ok {
		L.RaiseError("durationRound: invalid function assertion")

		return 0
	}

	var param string
	switch L.Get(1).Type() {
	case lua.LTNumber:
		param = fmt.Sprintf("%d%s", int64(L.CheckNumber(1)), "s")
	case lua.LTString:
		param = L.CheckString(1)
	case lua.LTNil:
		param = "0s"
	default:
		param = L.Get(1).String() + "s"
	}

	result := fn(param)

	L.Push(lua.LString(result))

	return 1
}

// EmptyFunc implements the sprig.empty function.
func EmptyFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("empty: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "empty requires 1 argument")

		return 0
	}

	value := L.Get(1)
	isEmpty := isEmptyLuaValue(value)

	L.Push(lua.LBool(isEmpty))

	return 1
}

// EncryptAESFunc wraps the sprig.encryptAES function.
func EncryptAESFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("encryptAES: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "encryptAES requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["encryptAES"].(func(string, string) (string, error))
	if !ok {
		L.RaiseError("encryptAES: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)

	result, err := fn(param0, param1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)

	return 2
}

// ExtFunc wraps the sprig.ext function.
func ExtFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("ext: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "ext requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["ext"].(func(string) string)
	if !ok {
		L.RaiseError("ext: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// GenPrivateKeyFunc wraps the sprig.genPrivateKey function.
func GenPrivateKeyFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("genPrivateKey: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "genPrivateKey requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["genPrivateKey"].(func(string) string)
	if !ok {
		L.RaiseError("genPrivateKey: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// HtpasswdFunc wraps the sprig.htpasswd function.
func HtpasswdFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("htpasswd: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "htpasswd requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["htpasswd"].(func(string, string) string)
	if !ok {
		L.RaiseError("htpasswd: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// IndentFunc wraps the sprig.indent function.
func IndentFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("indent: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "indent requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["indent"].(func(int, string) string)
	if !ok {
		L.RaiseError("indent: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// InitialsFunc wraps the sprig.initials function.
func InitialsFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("initials: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "initials requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["initials"].(func(string) string)
	if !ok {
		L.RaiseError("initials: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// IsAbsFunc wraps the sprig.isAbs function.
func IsAbsFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("isAbs: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "isAbs requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["isAbs"].(func(string) bool)
	if !ok {
		L.RaiseError("isAbs: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LBool(result))

	return 1
}

// KebabcaseFunc wraps the sprig.kebabcase function.
func KebabcaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("kebabcase: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "kebabcase requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["kebabcase"].(func(string) string)
	if !ok {
		L.RaiseError("kebabcase: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// NindentFunc wraps the sprig.nindent function.
func NindentFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("nindent: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "nindent requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["nindent"].(func(int, string) string)
	if !ok {
		L.RaiseError("nindent: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// NospaceFunc wraps the sprig.nospace function.
func NospaceFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("nospace: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "nospace requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["nospace"].(func(string) string)
	if !ok {
		L.RaiseError("nospace: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// OsBaseFunc wraps the sprig.osBase function.
func OsBaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("osBase: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "osBase requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["osBase"].(func(string) string)
	if !ok {
		L.RaiseError("osBase: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// OsCleanFunc wraps the sprig.osClean function.
func OsCleanFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("osClean: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "osClean requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["osClean"].(func(string) string)
	if !ok {
		L.RaiseError("osClean: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// OsDirFunc wraps the sprig.osDir function.
func OsDirFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("osDir: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "osDir requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["osDir"].(func(string) string)
	if !ok {
		L.RaiseError("osDir: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// OsExtFunc wraps the sprig.osExt function.
func OsExtFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("osExt: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "osExt requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["osExt"].(func(string) string)
	if !ok {
		L.RaiseError("osExt: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// OsIsAbsFunc wraps the sprig.osIsAbs function.
func OsIsAbsFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("osIsAbs: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "osIsAbs requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["osIsAbs"].(func(string) bool)
	if !ok {
		L.RaiseError("osIsAbs: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LBool(result))

	return 1
}

// PluralFunc implements the sprig.plural function.
func PluralFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("plural: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "plural requires 3 arguments: singular, plural, count")

		return 0
	}

	singular := L.CheckString(1)
	plural := L.CheckString(2)
	count := int(L.CheckNumber(3))

	result := ""
	if count == 1 {
		result = singular
	} else {
		result = plural
	}

	L.Push(lua.LString(result))

	return 1
}

// QuoteFunc wraps the sprig.quote function.
func QuoteFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("quote: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "quote requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["quote"].(func(...any) string)
	if !ok {
		L.RaiseError("quote: invalid function assertion")

		return 0
	}

	tbl := L.CheckTable(1)
	args := make([]any, 0, tbl.Len())

	tbl.ForEach(func(_, v lua.LValue) {
		if v == lua.LNil {
			return
		}

		var val any
		switch v.Type() {
		case lua.LTString:
			val = string(v.(lua.LString))
		case lua.LTNumber:
			num := float64(v.(lua.LNumber))
			if num == float64(int(num)) {
				val = int(num)
			} else {
				val = num
			}
		case lua.LTBool:
			val = bool(v.(lua.LBool))
		default:
			val = v.String()
		}

		args = append(args, val)
	})

	result := fn(args...)

	L.Push(lua.LString(result))

	return 1
}

// RandIntFunc wraps the sprig.randInt function.
func RandIntFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("randInt: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "randInt requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["randInt"].(func(int, int) int)
	if !ok {
		L.RaiseError("randInt: invalid function assertion")

		return 0
	}

	min := int(L.CheckNumber(1))
	max := int(L.CheckNumber(2))

	if min == max {
		L.Push(lua.LNumber(min))

		return 1
	}

	if min > max {
		min, max = max, min
	}

	result := fn(min, max)

	L.Push(lua.LNumber(result))

	return 1
}

// RegexFindAllFunc wraps the sprig.mustRegexFindAll function.
func RegexFindAllFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexFindAll: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "mustRegexFindAll requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexFindAll"].(func(string, string, int) ([]string, error))
	if !ok {
		L.RaiseError("mustRegexFindAll: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)
	param2 := int(L.CheckNumber(3))

	result, err := fn(param0, param1, param2)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	resultTable := L.CreateTable(len(result), 0)
	for i, v := range result {
		resultTable.RawSetInt(i+1, lua.LString(v))
	}

	L.Push(resultTable)
	L.Push(lua.LNil)

	return 2
}

// RegexFindFunc wraps the sprig.mustRegexFind function.
func RegexFindFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexFind: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "mustRegexFind requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexFind"].(func(string, string) (string, error))
	if !ok {
		L.RaiseError("mustRegexFind: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)

	result, err := fn(param0, param1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)

	return 2
}

// RegexMatchFunc wraps the sprig.mustRegexMatch function.
func RegexMatchFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexMatch: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "mustRegexMatch requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexMatch"].(func(string, string) (bool, error))
	if !ok {
		L.RaiseError("mustRegexMatch: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)

	result, err := fn(param0, param1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LBool(result))
	L.Push(lua.LNil)

	return 2
}

// RegexReplaceAllFunc wraps the sprig.mustRegexReplaceAll function.
func RegexReplaceAllFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexReplaceAll: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "mustRegexReplaceAll requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexReplaceAll"].(func(string, string, string) (string, error))
	if !ok {
		L.RaiseError("mustRegexReplaceAll: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)
	param2 := L.CheckString(3)

	result, err := fn(param0, param1, param2)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)

	return 2
}

// RegexReplaceAllLiteralFunc wraps the sprig.mustRegexReplaceAllLiteral function.
func RegexReplaceAllLiteralFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexReplaceAllLiteral: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "mustRegexReplaceAllLiteral requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexReplaceAllLiteral"].(func(string, string, string) (string, error))
	if !ok {
		L.RaiseError("mustRegexReplaceAllLiteral: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)
	param2 := L.CheckString(3)

	result, err := fn(param0, param1, param2)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(result))
	L.Push(lua.LNil)

	return 2
}

// RegexSplitFunc wraps the sprig.mustRegexSplit function.
func RegexSplitFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("mustRegexSplit: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "mustRegexSplit requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["mustRegexSplit"].(func(string, string, int) ([]string, error))
	if !ok {
		L.RaiseError("mustRegexSplit: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)
	param2 := int(L.CheckNumber(3))

	result, err := fn(param0, param1, param2)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	resultTable := L.CreateTable(len(result), 0)
	for i, v := range result {
		resultTable.RawSetInt(i+1, lua.LString(v))
	}

	L.Push(resultTable)
	L.Push(lua.LNil)

	return 2
}

// RoundFunc wraps the sprig.round function.
func RoundFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("round: %v", r)
		}
	}()

	top := L.GetTop()
	if top < 2 {
		L.ArgError(1, "round requires at least 2 arguments: value and precision")

		return 0
	}

	roundFn := sprig.FuncMap()["round"]

	var value any
	switch L.Get(1).Type() {
	case lua.LTNumber:
		value = float64(L.CheckNumber(1))
	case lua.LTString:
		value = L.CheckString(1)
	default:
		value = L.Get(1).String()
	}

	precision := int(L.CheckNumber(2))

	var result float64
	if top >= 3 {
		result = roundFn.(func(any, int, ...float64) float64)(value, precision, float64(L.CheckNumber(3)))
	} else {
		result = roundFn.(func(any, int, ...float64) float64)(value, precision)
	}

	L.Push(lua.LNumber(result))

	return 1
}

// SemverCompareFunc wraps the sprig.semverCompare function.
func SemverCompareFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("semverCompare: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "semverCompare requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["semverCompare"].(func(string, string) (bool, error))
	if !ok {
		L.RaiseError("semverCompare: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	param1 := L.CheckString(2)

	result, err := fn(param0, param1)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LBool(result))
	L.Push(lua.LNil)

	return 2
}

// SeqFunc wraps the sprig.seq function.
func SeqFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("seq: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "seq requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["seq"].(func(...int) string)
	if !ok {
		L.RaiseError("seq: invalid function assertion")

		return 0
	}

	tbl := L.CheckTable(1)
	params := make([]int, 0, tbl.Len())

	tbl.ForEach(func(_, v lua.LValue) {
		if num, ok := v.(lua.LNumber); ok {
			params = append(params, int(num))
		}
	})

	result := fn(params...)

	L.Push(lua.LString(result))

	return 1
}

// Sha1sumFunc wraps the sprig.sha1sum function.
func Sha1sumFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("sha1sum: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "sha1sum requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["sha1sum"].(func(string) string)
	if !ok {
		L.RaiseError("sha1sum: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// Sha256sumFunc wraps the sprig.sha256sum function.
func Sha256sumFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("sha256sum: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "sha256sum requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["sha256sum"].(func(string) string)
	if !ok {
		L.RaiseError("sha256sum: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// Sha512sumFunc wraps the sprig.sha512sum function.
func Sha512sumFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("sha512sum: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "sha512sum requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["sha512sum"].(func(string) string)
	if !ok {
		L.RaiseError("sha512sum: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// ShuffleFunc wraps the sprig.shuffle function.
func ShuffleFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("shuffle: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "shuffle requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["shuffle"].(func(string) string)
	if !ok {
		L.RaiseError("shuffle: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// SnakecaseFunc wraps the sprig.snakecase function.
func SnakecaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("snakecase: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "snakecase requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["snakecase"].(func(string) string)
	if !ok {
		L.RaiseError("snakecase: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// SortAlphaFunc wraps the sprig.sortAlpha function.
func SortAlphaFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("sortAlpha: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "sortAlpha requires 1 argument")

		return 0
	}

	fn, ok := sprig.FuncMap()["sortAlpha"].(func(any) []string)
	if !ok {
		L.RaiseError("sortAlpha: invalid function assertion")

		return 0
	}

	var param any

	switch L.Get(1).Type() {
	case lua.LTTable:
		tbl := L.CheckTable(1)
		strSlice := make([]string, 0, tbl.Len())

		tbl.ForEach(func(_, v lua.LValue) {
			strSlice = append(strSlice, v.String())
		})

		param = strSlice
	case lua.LTString:
		param = L.CheckString(1)
	default:
		param = L.Get(1).String()
	}

	result := fn(param)

	resultTable := L.CreateTable(len(result), 0)
	for i, v := range result {
		resultTable.RawSetInt(i+1, lua.LString(v))
	}

	L.Push(resultTable)

	return 1
}

// SquoteFunc wraps the sprig.squote function.
func SquoteFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("squote: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "squote requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["squote"].(func(...any) string)
	if !ok {
		L.RaiseError("squote: invalid function assertion")

		return 0
	}

	tbl := L.CheckTable(1)
	params := make([]any, 0, tbl.Len())

	tbl.ForEach(func(_, v lua.LValue) {
		if v == lua.LNil {
			return
		}

		switch v.Type() {
		case lua.LTString:
			params = append(params, v.String())
		case lua.LTNumber:
			num := float64(v.(lua.LNumber))

			if num == float64(int(num)) {
				params = append(params, fmt.Sprintf("%d", int(num)))
			} else {
				params = append(params, fmt.Sprintf("%v", num))
			}
		case lua.LTBool:
			params = append(params, fmt.Sprintf("%v", bool(v.(lua.LBool))))
		default:
			params = append(params, v.String())
		}
	})

	result := fn(params...)

	L.Push(lua.LString(result))

	return 1
}

// SubstrFunc wraps the sprig.substr function.
func SubstrFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("substr: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "substr requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["substr"].(func(int, int, string) string)
	if !ok {
		L.RaiseError("substr: invalid function assertion")

		return 0
	}

	start := int(L.CheckNumber(1))
	end := int(L.CheckNumber(2))
	str := L.CheckString(3)

	runes := []rune(str)
	length := len(runes)

	if start < 0 {
		start = length + start
	}

	if start < 0 {
		start = 0
	}

	if end > length {
		end = length
	}

	if start >= length || start >= end {
		L.Push(lua.LString(""))

		return 1
	}

	result := fn(start, end, str)

	L.Push(lua.LString(result))

	return 1
}

// SwapcaseFunc wraps the sprig.swapcase function.
func SwapcaseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("swapcase: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "swapcase requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["swapcase"].(func(string) string)
	if !ok {
		L.RaiseError("swapcase: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// TernaryFunc implements the sprig.ternary function.
func TernaryFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("ternary: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "ternary requires 3 arguments")

		return 0
	}

	condition := L.CheckBool(3)

	trueValue := L.Get(1)
	falseValue := L.Get(2)

	if condition {
		L.Push(trueValue)
	} else {
		L.Push(falseValue)
	}

	return 1
}

// ToDecimalFunc wraps the sprig.toDecimal function.
func ToDecimalFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("toDecimal: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "toDecimal requires 1 argument")

		return 0
	}

	fn, ok := sprig.FuncMap()["toDecimal"].(func(any) int64)
	if !ok {
		L.RaiseError("toDecimal: invalid function assertion")

		return 0
	}

	var param string

	switch L.Get(1).Type() {
	case lua.LTString:
		param = L.CheckString(1)
	case lua.LTNumber:
		param = fmt.Sprintf("%v", L.CheckNumber(1))
	case lua.LTBool:
		if L.CheckBool(1) {
			param = "1"
		} else {
			param = "0"
		}
	case lua.LTNil:
		param = "0"
	default:
		param = L.Get(1).String()
	}

	result := fn(param)

	L.Push(lua.LNumber(result))

	return 1
}

// TruncFunc wraps the sprig.trunc function.
func TruncFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("trunc: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "trunc requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["trunc"].(func(int, string) string)
	if !ok {
		L.RaiseError("trunc: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// UniqFunc implements the sprig.uniq function.
func UniqFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("uniq: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "uniq requires 1 argument")

		return 0
	}

	inputTable := L.CheckTable(1)
	uniqueTable := L.CreateTable(0, 0)
	seen := make(map[string]bool)
	i := 1

	inputTable.ForEach(func(_, v lua.LValue) {
		var key string

		switch v.Type() {
		case lua.LTString:
			key = v.String()
		case lua.LTNumber:
			key = fmt.Sprintf("%v", v)
		case lua.LTBool:
			key = fmt.Sprintf("%v", v)
		case lua.LTNil:
			key = "nil"
		case lua.LTTable:
			key = fmt.Sprintf("%p", v)
		default:
			key = v.String()
		}

		if !seen[key] {
			seen[key] = true

			uniqueTable.RawSetInt(i, v)
			i++
		}
	})

	L.Push(uniqueTable)

	return 1
}

// UntitleFunc wraps the sprig.untitle function.
func UntitleFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("untitle: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "untitle requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["untitle"].(func(string) string)
	if !ok {
		L.RaiseError("untitle: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	L.Push(lua.LString(result))

	return 1
}

// UrlJoinFunc wraps the sprig.urlJoin function.
func UrlJoinFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("urlJoin: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "urlJoin requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["urlJoin"].(func(map[string]any) string)
	if !ok {
		L.RaiseError("urlJoin: invalid function assertion")

		return 0
	}

	tbl := L.CheckTable(1)

	param := make(map[string]any)
	tbl.ForEach(func(k, v lua.LValue) {
		if ks, ok := k.(lua.LString); ok {
			param[string(ks)] = v.String()
		}
	})

	result := fn(param)

	L.Push(lua.LString(result))

	return 1
}

// UrlParseFunc wraps the sprig.urlParse function.
func UrlParseFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("urlParse: %v", r)
		}
	}()

	if L.GetTop() < 1 {
		L.ArgError(1, "urlParse requires 1 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["urlParse"].(func(string) map[string]any)
	if !ok {
		L.RaiseError("urlParse: invalid function assertion")

		return 0
	}

	param0 := L.CheckString(1)
	result := fn(param0)

	table := L.NewTable()

	for k, v := range result {
		if v == nil {
			table.RawSetString(k, lua.LString(""))
		} else {
			table.RawSetString(k, lua.LString(fmt.Sprintf("%v", v)))
		}
	}

	L.Push(table)

	return 1
}

// WrapFunc wraps the sprig.wrap function.
func WrapFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("wrap: %v", r)
		}
	}()

	if L.GetTop() < 2 {
		L.ArgError(1, "wrap requires 2 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["wrap"].(func(int, string) string)
	if !ok {
		L.RaiseError("wrap: invalid function assertion")

		return 0
	}

	param0 := int(L.CheckNumber(1))
	param1 := L.CheckString(2)
	result := fn(param0, param1)

	L.Push(lua.LString(result))

	return 1
}

// WrapWithFunc wraps the sprig.wrapWith function.
func WrapWithFunc(L *lua.LState) int {
	defer func() {
		if r := recover(); r != nil {
			L.RaiseError("wrapWith: %v", r)
		}
	}()

	if L.GetTop() < 3 {
		L.ArgError(1, "wrapWith requires 3 arguments")

		return 0
	}

	fn, ok := sprig.FuncMap()["wrapWith"].(func(int, string, string) string)
	if !ok {
		L.RaiseError("wrapWith: invalid function assertion")

		return 0
	}

	width := int(L.CheckNumber(1))
	separator := L.CheckString(3)
	text := L.CheckString(2)

	result := fn(width, separator, text)

	L.Push(lua.LString(result))

	return 1
}

// Loader is the entrypoint to load the sprig library into a LState.
func Loader(L *lua.LState) int {
	mod := L.RegisterModule("sprig", map[string]lua.LGFunction{
		"abbrev":                 AbbrevFunc,
		"abbrevboth":             AbbrevbothFunc,
		"adler32sum":             Adler32sumFunc,
		"ago":                    AgoFunc,
		"all":                    AllFunc,
		"any":                    AnyFunc,
		"b32dec":                 B32decFunc,
		"b32enc":                 B32encFunc,
		"b64dec":                 B64decFunc,
		"b64enc":                 B64encFunc,
		"base":                   BaseFunc,
		"bcrypt":                 BcryptFunc,
		"camelcase":              CamelcaseFunc,
		"cat":                    CatFunc,
		"clean":                  CleanFunc,
		"coalesce":               CoalesceFunc,
		"compact":                CompactFunc,
		"decryptAES":             DecryptAESFunc,
		"derivePassword":         DerivePasswordFunc,
		"dir":                    DirFunc,
		"duration":               DurationFunc,
		"durationRound":          DurationRoundFunc,
		"empty":                  EmptyFunc,
		"encryptAES":             EncryptAESFunc,
		"ext":                    ExtFunc,
		"genPrivateKey":          GenPrivateKeyFunc,
		"htpasswd":               HtpasswdFunc,
		"indent":                 IndentFunc,
		"initials":               InitialsFunc,
		"isAbs":                  IsAbsFunc,
		"kebabcase":              KebabcaseFunc,
		"nindent":                NindentFunc,
		"nospace":                NospaceFunc,
		"osBase":                 OsBaseFunc,
		"osClean":                OsCleanFunc,
		"osDir":                  OsDirFunc,
		"osExt":                  OsExtFunc,
		"osIsAbs":                OsIsAbsFunc,
		"plural":                 PluralFunc,
		"quote":                  QuoteFunc,
		"randInt":                RandIntFunc,
		"regexFind":              RegexFindFunc,
		"regexFindAll":           RegexFindAllFunc,
		"regexMatch":             RegexMatchFunc,
		"regexReplaceAll":        RegexReplaceAllFunc,
		"regexReplaceAllLiteral": RegexReplaceAllLiteralFunc,
		"regexSplit":             RegexSplitFunc,
		"round":                  RoundFunc,
		"semverCompare":          SemverCompareFunc,
		"seq":                    SeqFunc,
		"sha1sum":                Sha1sumFunc,
		"sha256sum":              Sha256sumFunc,
		"sha512sum":              Sha512sumFunc,
		"shuffle":                ShuffleFunc,
		"snakecase":              SnakecaseFunc,
		"sortAlpha":              SortAlphaFunc,
		"squote":                 SquoteFunc,
		"substr":                 SubstrFunc,
		"swapcase":               SwapcaseFunc,
		"ternary":                TernaryFunc,
		"toDecimal":              ToDecimalFunc,
		"trunc":                  TruncFunc,
		"uniq":                   UniqFunc,
		"untitle":                UntitleFunc,
		"urlJoin":                UrlJoinFunc,
		"urlParse":               UrlParseFunc,
		"wrap":                   WrapFunc,
		"wrapWith":               WrapWithFunc,
	})

	L.Push(mod)

	return 1
}
