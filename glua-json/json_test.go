package json_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	luajson "github.com/projectsveltos/lua-utils/glua-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func luaValuesEqual(v1, v2 lua.LValue) bool {
	if v1.Type() != v2.Type() {
		return false
	}

	switch v1.Type() {
	case lua.LTNil:
		return true
	case lua.LTBool:
		l, ok := v1.(lua.LBool)
		if !ok {
			panic("unexpected type, expecting LBool")
		}

		r, ok := v2.(lua.LBool)
		if !ok {
			panic("unexpected type, expecting LBool")
		}

		return bool(l) == bool(r)
	case lua.LTNumber:
		l, ok := v1.(lua.LNumber)
		if !ok {
			panic("unexpected type, expecting LNumber")
		}

		r, ok := v2.(lua.LNumber)
		if !ok {
			panic("unexpected type, expecting LNumber")
		}

		return float64(l) == float64(r)
	case lua.LTString:
		l, ok := v1.(lua.LString)
		if !ok {
			panic("unexpected type, expecting LString")
		}

		r, ok := v2.(lua.LString)
		if !ok {
			panic("unexpected type, expecting LString")
		}

		return string(l) == string(r)
	case lua.LTTable:
		t1, ok := v1.(*lua.LTable)
		if !ok {
			panic("unexpected type, expecting LTable")
		}

		t2, ok := v2.(*lua.LTable)
		if !ok {
			panic("unexpected type, expecting LTable")
		}

		if t1.Len() != t2.Len() {
			return false
		}

		equal := true

		t1.ForEach(func(k, v lua.LValue) {
			if !luaValuesEqual(v, t2.RawGet(k)) {
				equal = false
			}
		})

		return equal
	default:
		return false
	}
}

func TestSimple(t *testing.T) {
	const str = `
	local json = require("json")
	assert(type(json) == "table")
	assert(type(json.decode) == "function")
	assert(type(json.encode) == "function")

	assert(json.encode(true) == "true")
	assert(json.encode(1) == "1")
	assert(json.encode(-10) == "-10")
	assert(json.encode(nil) == "null")
	assert(json.encode({}) == "[]")
	assert(json.encode({1, 2, 3}) == "[1,2,3]")

	local _, err = json.encode({1, 2, [10] = 3})
	assert(string.find(err, "sparse array"))

	local _, err = json.encode({1, 2, 3, name = "Tim"})
	assert(string.find(err, "mixed or invalid key types"))

	local _, err = json.encode({name = "Tim", [false] = 123})
	assert(string.find(err, "mixed or invalid key types"))

	local obj = {"a",1,"b",2,"c",3}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	for i = 1, #obj do
		assert(obj[i] == jsonObj[i])
	end

	local obj = {name="Tim",number=12345}
	local jsonStr = json.encode(obj)
	local jsonObj = json.decode(jsonStr)
	assert(obj.name == jsonObj.name)
	assert(obj.number == jsonObj.number)

	assert(json.decode("null") == nil)

	assert(json.decode(json.encode({person={name = "tim",}})).person.name == "tim")

	local obj = {
		abc = 123,
		def = nil,
	}
	local obj2 = {
		obj = obj,
	}
	obj.obj2 = obj2
	assert(json.encode(obj) == nil)

	local a = {}
	for i=1, 5 do
		a[i] = i
	end
	assert(json.encode(a) == "[1,2,3,4,5]")
	`

	s := lua.NewState()
	defer s.Close()

	luajson.Preload(s)

	if err := s.DoString(str); err != nil {
		t.Error(err)
	}
}

func TestCustomRequire(t *testing.T) {
	const str = `
	local j = require("JSON")
	assert(type(j) == "table")
	assert(type(j.decode) == "function")
	assert(type(j.encode) == "function")
	`

	s := lua.NewState()
	defer s.Close()

	s.PreloadModule("JSON", luajson.Loader)

	if err := s.DoString(str); err != nil {
		t.Error(err)
	}
}

func TestDecodeValueJSONNumber(t *testing.T) {
	s := lua.NewState()
	defer s.Close()

	v := luajson.DecodeValue(s, json.Number("124.11"))
	if v.Type() != lua.LTString || v.String() != "124.11" {
		t.Fatalf("expecting LString, got %T", v)
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    lua.LValue
		expected string
		wantErr  string
	}{
		{
			name:     "encode nil",
			input:    lua.LNil,
			expected: "null",
		},
		{
			name:     "encode boolean true",
			input:    lua.LBool(true),
			expected: "true",
		},
		{
			name:     "encode boolean false",
			input:    lua.LBool(false),
			expected: "false",
		},
		{
			name:     "encode integer",
			input:    lua.LNumber(42),
			expected: "42",
		},
		{
			name:     "encode float",
			input:    lua.LNumber(3.14),
			expected: "3.14",
		},
		{
			name:     "encode string",
			input:    lua.LString("hello"),
			expected: `"hello"`,
		},
		{
			name:     "encode empty table as array",
			input:    &lua.LTable{},
			expected: "[]",
		},
		{
			name: "encode array table",
			input: func() *lua.LTable {
				tbl := &lua.LTable{}
				tbl.Append(lua.LNumber(1))
				tbl.Append(lua.LNumber(2))

				return tbl
			}(),
			expected: "[1,2]",
		},
		{
			name: "encode object table",
			input: func() *lua.LTable {
				tbl := &lua.LTable{}
				tbl.RawSetString("name", lua.LString("test"))
				tbl.RawSetString("value", lua.LNumber(42))

				return tbl
			}(),
			expected: `{"name":"test","value":42}`,
		},
		{
			name: "error on sparse array",
			input: func() *lua.LTable {
				tbl := &lua.LTable{}
				tbl.RawSetInt(1, lua.LNumber(1))
				tbl.RawSetInt(3, lua.LNumber(3))

				return tbl
			}(),
			wantErr: "cannot encode sparse array",
		},
		{
			name: "error on mixed keys",
			input: func() *lua.LTable {
				tbl := &lua.LTable{}
				tbl.RawSetString("name", lua.LString("test"))
				tbl.RawSetInt(1, lua.LNumber(42))

				return tbl
			}(),
			wantErr: "cannot encode mixed or invalid key types",
		},
		{
			name: "error on nested tables",
			input: func() *lua.LTable {
				tbl1 := &lua.LTable{}
				tbl2 := &lua.LTable{}
				tbl1.RawSetString("nested", tbl2)
				tbl2.RawSetString("parent", tbl1)

				return tbl1
			}(),
			wantErr: "cannot encode recursively nested tables",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := luajson.Encode(tt.input)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestYAMLIntegration(t *testing.T) {
	yamlTestData := `
name: test
value: 42
items:
  - one
  - two`

	luaSetup := `local json = require("json")`

	yamlToLuaTest := fmt.Sprintf(`
		-- Test YAML to Lua conversion
		local yamlStr = %q
		local luaObj = json.fromYAML(yamlStr)
		assert(luaObj.name == "test")
		assert(luaObj.value == 42)
		assert(luaObj.items[1] == "one")
		assert(luaObj.items[2] == "two")`, strings.TrimSpace(yamlTestData))

	luaToYamlTest := `
		-- Test Lua to YAML conversion
		local obj = {name = "example", count = 123, tags = {"a", "b"}}
		local yamlResult = json.toYAML(obj)
		assert(type(yamlResult) == "string")
		assert(string.find(yamlResult, "name: example"))
		assert(string.find(yamlResult, "count: 123"))`

	roundTripTest := `
		-- Test round-trip conversion
		local original = {greeting = "hello", numbers = {1, 2, 3}}
		local yamlString = json.toYAML(original)
		local restored = json.fromYAML(yamlString)
		assert(restored.greeting == original.greeting)
		assert(restored.numbers[1] == original.numbers[1])
		assert(restored.numbers[2] == original.numbers[2])
		assert(restored.numbers[3] == original.numbers[3])`

	luaScript := luaSetup + yamlToLuaTest + luaToYamlTest + roundTripTest

	s := lua.NewState()
	defer s.Close()

	luajson.Preload(s)

	if err := s.DoString(luaScript); err != nil {
		t.Error(err)
	}
}

func TestFromYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func(*lua.LState) lua.LValue
		wantErr  string
	}{
		{
			name: "simple YAML object",
			input: `name: test
value: 42`,
			expected: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(0, 2)
				tbl.RawSetString("name", lua.LString("test"))
				tbl.RawSetString("value", lua.LNumber(42))
				return tbl
			},
		},
		{
			name: "YAML array",
			input: `- one
- two
- three`,
			expected: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(3, 0)
				tbl.Append(lua.LString("one"))
				tbl.Append(lua.LString("two"))
				tbl.Append(lua.LString("three"))
				return tbl
			},
		},
		{
			name:    "invalid YAML",
			input:   "invalid: yaml: content:",
			wantErr: "yaml:",
		},
	}

	L := lua.NewState()
	defer L.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := luajson.FromYAML(L, []byte(tt.input))
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			expected := tt.expected(L)
			assert.True(t, luaValuesEqual(expected, result),
				"expected %v but got %v", expected, result)
		})
	}
}

func TestToYAML(t *testing.T) {
	tests := []struct {
		name     string
		input    func(*lua.LState) lua.LValue
		contains []string
		wantErr  string
	}{
		{
			name: "simple object",
			input: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(0, 2)
				tbl.RawSetString("name", lua.LString("test"))
				tbl.RawSetString("value", lua.LNumber(42))
				return tbl
			},
			contains: []string{"name: test", "value: 42"},
		},
		{
			name: "array",
			input: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(2, 0)
				tbl.Append(lua.LString("one"))
				tbl.Append(lua.LString("two"))
				return tbl
			},
			contains: []string{"- one", "- two"},
		},
		{
			name: "nested object",
			input: func(L *lua.LState) lua.LValue {
				inner := L.CreateTable(0, 1)
				inner.RawSetString("count", lua.LNumber(5))

				outer := L.CreateTable(0, 1)
				outer.RawSetString("inner", inner)
				return outer
			},
			contains: []string{"inner:", "count: 5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			input := tt.input(L)
			result, err := luajson.ToYAML(input)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			resultStr := string(result)
			for _, expected := range tt.contains {
				assert.Contains(t, resultStr, expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func(*lua.LState) lua.LValue
		wantErr  bool
	}{
		{
			name:  "decode null",
			input: "null",
			expected: func(_ *lua.LState) lua.LValue {
				return lua.LNil
			},
		},
		{
			name:  "decode boolean",
			input: "true",
			expected: func(_ *lua.LState) lua.LValue {
				return lua.LBool(true)
			},
		},
		{
			name:  "decode number",
			input: "42",
			expected: func(_ *lua.LState) lua.LValue {
				return lua.LNumber(42)
			},
		},
		{
			name:  "decode string",
			input: `"hello"`,
			expected: func(_ *lua.LState) lua.LValue {
				return lua.LString("hello")
			},
		},
		{
			name:  "decode empty array",
			input: "[]",
			expected: func(L *lua.LState) lua.LValue {
				return L.CreateTable(0, 0)
			},
		},
		{
			name:  "decode array",
			input: "[1,2,3]",
			expected: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(3, 0)
				tbl.Append(lua.LNumber(1))
				tbl.Append(lua.LNumber(2))
				tbl.Append(lua.LNumber(3))

				return tbl
			},
		},
		{
			name:  "decode object",
			input: `{"name":"test","value":42}`,
			expected: func(L *lua.LState) lua.LValue {
				tbl := L.CreateTable(0, 2)
				tbl.RawSetString("name", lua.LString("test"))
				tbl.RawSetString("value", lua.LNumber(42))

				return tbl
			},
		},
		{
			name:    "error on invalid JSON",
			input:   "{invalid}",
			wantErr: true,
		},
	}

	L := lua.NewState()
	defer L.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := luajson.Decode(L, []byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)

			expected := tt.expected(L)

			assert.True(t, luaValuesEqual(expected, result),
				"expected %v but got %v", expected, result)
		})
	}
}
