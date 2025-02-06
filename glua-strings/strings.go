// Copyright 2017 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package strings

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule("strings", Loader)
}

func Loader(L *lua.LState) int {
	mod := L.NewTable()
	L.SetFuncs(mod, stringsFuncs)
	L.Push(mod)
	return 1
}

func RetBool(L *lua.LState, v bool) int {
	L.Push(lua.LBool(v))
	return 1
}

func RetInt(L *lua.LState, v int) int {
	L.Push(lua.LNumber(v))
	return 1
}

func RetString(L *lua.LState, v string) int {
	L.Push(lua.LString(v))
	return 1
}

func RetStringList(L *lua.LState, vs []string) int {
	tb := L.NewTable()
	for _, v := range vs {
		tb.Append(lua.LString(v))
	}
	L.Push(tb)
	return 1
}

var stringsFuncs = map[string]lua.LGFunction{
	"Compare": func(L *lua.LState) int {
		a := L.CheckString(1)
		b := L.CheckString(2)

		ret := strings.Compare(a, b)
		return RetInt(L, ret)
	},
	"Contains": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Contains(s, substr)
		return RetBool(L, ret)
	},
	"ContainsAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		chars := L.CheckString(2)

		ret := strings.ContainsAny(s, chars)
		return RetBool(L, ret)
	},
	"ContainsRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		r := L.CheckInt(2)

		ret := strings.ContainsRune(s, rune(r))
		return RetBool(L, ret)
	},
	"Count": func(L *lua.LState) int {
		s := L.CheckString(1)
		substr := L.CheckString(2)

		ret := strings.Count(s, substr)
		return RetInt(L, ret)
	},
	"EqualFold": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.EqualFold(s, t)
		return RetBool(L, ret)
	},
	"Fields": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Fields(s)
		return RetStringList(L, ret)
	},
	"FieldsFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.FieldsFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetStringList(L, ret)
	},
	"HasPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasPrefix(s, t)
		return RetBool(L, ret)
	},
	"HasSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.HasSuffix(s, t)
		return RetBool(L, ret)
	},
	"Index": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Index(s, t)
		return RetInt(L, ret)
	},
	"IndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.IndexAny(s, t)
		return RetInt(L, ret)
	},
	"IndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexByte(s, byte(t))
		return RetInt(L, ret)
	},
	"IndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.IndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetInt(L, ret)
	},
	"IndexRune": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.IndexRune(s, rune(t))
		return RetInt(L, ret)
	},
	"Join": func(L *lua.LState) int {
		tbl := L.CheckTable(1)
		sep := L.CheckString(2)

		strs := make([]string, 0, tbl.Len())

		tbl.ForEach(func(_, value lua.LValue) {
			if str, ok := value.(lua.LString); ok {
				strs = append(strs, string(str))
			}
		})

		ret := strings.Join(strs, sep)
		return RetString(L, ret)
	},
	"LastIndex": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndex(s, t)
		return RetInt(L, ret)
	},
	"LastIndexAny": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.LastIndexAny(s, t)
		return RetInt(L, ret)
	},
	"LastIndexByte": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.LastIndexByte(s, byte(t))
		return RetInt(L, ret)
	},
	"LastIndexFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.LastIndexFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetInt(L, ret)
	},
	"Map": func(L *lua.LState) int {
		fn := L.CheckFunction(1)
		s := L.CheckString(2)

		ret := strings.Map(
			func(r rune) rune {
				return callFunc_Rune_ret_Rune(
					L, fn, lua.LNumber(r),
				)
			},
			s,
		)
		return RetString(L, ret)
	},
	"Repeat": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckInt(2)

		ret := strings.Repeat(s, t)
		return RetString(L, ret)
	},
	"Replace": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		z := L.CheckString(3)
		n := L.CheckInt(4)

		ret := strings.Replace(s, t, z, n)
		return RetString(L, ret)
	},
	"Split": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.Split(s, t)
		return RetStringList(L, ret)
	},
	"SplitAfter": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)

		ret := strings.SplitAfter(s, t)
		return RetStringList(L, ret)
	},
	"SplitAfterN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		if n == 0 {
			L.Push(lua.LNil)
			return 1
		}

		ret := strings.SplitAfterN(s, t, n)
		return RetStringList(L, ret)
	},
	"SplitN": func(L *lua.LState) int {
		s := L.CheckString(1)
		t := L.CheckString(2)
		n := L.CheckInt(3)

		if n == 0 {
			L.Push(lua.LNil)
			return 1
		}

		ret := strings.SplitN(s, t, n)
		return RetStringList(L, ret)
	},
	"Title": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.Title(s)
		return RetString(L, ret)
	},
	"ToLower": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToLower(s)
		return RetString(L, ret)
	},
	"ToTitle": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToTitle(s)
		return RetString(L, ret)
	},
	"ToUpper": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.ToUpper(s)
		return RetString(L, ret)
	},
	"Trim": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.Trim(s, cutset)
		return RetString(L, ret)
	},
	"TrimFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetString(L, ret)
	},
	"TrimLeft": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimLeft(s, cutset)
		return RetString(L, ret)
	},
	"TrimLeftFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimLeftFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetString(L, ret)
	},
	"TrimPrefix": func(L *lua.LState) int {
		s := L.CheckString(1)
		prefix := L.CheckString(2)

		ret := strings.TrimPrefix(s, prefix)
		return RetString(L, ret)
	},
	"TrimRight": func(L *lua.LState) int {
		s := L.CheckString(1)
		cutset := L.CheckString(2)

		ret := strings.TrimRight(s, cutset)
		return RetString(L, ret)
	},
	"TrimRightFunc": func(L *lua.LState) int {
		s := L.CheckString(1)
		fn := L.CheckFunction(2)

		ret := strings.TrimRightFunc(s, func(r rune) bool {
			return callFunc_Rune_ret_Bool(
				L, fn, lua.LNumber(r),
			)
		})
		return RetString(L, ret)
	},
	"TrimSpace": func(L *lua.LState) int {
		s := L.CheckString(1)

		ret := strings.TrimSpace(s)
		return RetString(L, ret)
	},
	"TrimSuffix": func(L *lua.LState) int {
		s := L.CheckString(1)
		suffix := L.CheckString(2)

		ret := strings.TrimSuffix(s, suffix)
		return RetString(L, ret)
	},
}

// func(rune) bool
func callFunc_Rune_ret_Bool(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) bool {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckBool(-1)
	return ret
}

// func(rune) rune
func callFunc_Rune_ret_Rune(L *lua.LState, lf *lua.LFunction, args ...lua.LValue) rune {
	err := L.CallByParam(lua.P{Protect: true, Fn: lf, NRet: 1}, args...)
	if err != nil {
		panic(err)
	}
	defer L.Pop(1)

	ret := L.CheckInt(-1)
	return rune(ret)
}
