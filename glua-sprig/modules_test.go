package gluasprig_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash/adler32"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	gluasprig "github.com/projectsveltos/lua-utils/glua-sprig"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/crypto/bcrypt"
)

func TestAbbrevFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		maxWidth int
	}{
		{
			maxWidth: 10,
			input:    "hello",
			expected: "hello",
		},
		{
			maxWidth: 5,
			input:    "hello world",
			expected: "he...",
		},
		{
			maxWidth: 5,
			input:    "hello",
			expected: "hello",
		},
		{
			maxWidth: 4,
			input:    "hello",
			expected: "h...",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.maxWidth))
			L.Push(lua.LString(tt.input))

			gluasprig.AbbrevFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAbbrevbothFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		left     int
		right    int
	}{
		{
			left:     2,
			right:    7,
			input:    "hello world",
			expected: "hell...",
		},
		{
			left:     5,
			right:    3,
			input:    "hello world",
			expected: "hello world",
		},
		{
			left:     2,
			right:    6,
			input:    "hello world",
			expected: "hello world",
		},
		{
			left:     0,
			right:    10,
			input:    "hello world",
			expected: "hello w...",
		},
		{
			left:     3,
			right:    7,
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: "abcd...",
		},
		{
			left:     5,
			right:    10,
			input:    "hello",
			expected: "hello",
		},
		{
			left:     4,
			right:    10,
			input:    "hello world",
			expected: "hello w...",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.left))
			L.Push(lua.LNumber(tt.right))
			L.Push(lua.LString(tt.input))

			gluasprig.AbbrevbothFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestAdler32sumFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "hello",
		},
		{
			input: "",
		},
		{
			input: "test string",
		},
		{
			input: "This is a longer test string with various characters!@#$%^&*()",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			adlerChecksum := adler32.Checksum([]byte(tt.input))
			expected := strconv.FormatUint(uint64(adlerChecksum), 10)

			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.Adler32sumFunc(L)

			result := L.ToString(-1)
			require.Equal(t, expected, result)
		})
	}
}

func TestAgoFunc(t *testing.T) {
	tests := []struct {
		timestamp int64
	}{
		{
			timestamp: time.Now().Unix(),
		},
		{
			timestamp: time.Now().Add(-5 * time.Minute).Unix(),
		},
		{
			timestamp: time.Now().Add(-1 * time.Hour).Unix(),
		},
		{
			timestamp: time.Now().Add(-24 * time.Hour).Unix(),
		},
		{
			timestamp: time.Now().Add(-7 * 24 * time.Hour).Unix(),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.timestamp))
			gluasprig.AgoFunc(L)

			result := L.ToString(-1)
			require.NotEmpty(t, result, "Result should not be empty")

			matched, _ := regexp.MatchString(`\d+(\.\d+)?(ns|us|µs|ms|s|m|h)`, result)
			require.True(t, matched, "Expected duration string, got: %s", result)
		})
	}
}

func TestAllFunc(t *testing.T) {
	tests := []struct {
		inputs   []any
		expected bool
	}{
		{
			inputs:   []any{1, "hello", true, 42.5},
			expected: true,
		},
		{
			inputs:   []any{1, "", true, 42.5},
			expected: false,
		},
		{
			inputs:   []any{nil, 0, "", false},
			expected: false,
		},
		{
			inputs:   []any{nil, "", false, 0},
			expected: false,
		},
		{
			inputs:   []any{},
			expected: true,
		},
		{
			inputs:   []any{[]int{}, map[string]string{}},
			expected: false,
		},
		{
			inputs:   []any{[]int{1, 2}, map[string]string{"key": "value"}},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.CreateTable(len(tt.inputs), 0)

			for _, input := range tt.inputs {
				var luaValue lua.LValue

				switch v := input.(type) {
				case nil:
					luaValue = lua.LNil
				case int:
					luaValue = lua.LNumber(v)
				case int64:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case string:
					luaValue = lua.LString(v)
				case bool:
					luaValue = lua.LBool(v)
				case []int:
					tbl := L.CreateTable(len(v), 0)
					for _, val := range v {
						tbl.Append(lua.LNumber(val))
					}

					luaValue = tbl
				case map[string]string:
					tbl := L.CreateTable(0, len(v))
					for key, val := range v {
						tbl.RawSetString(key, lua.LString(val))
					}

					luaValue = tbl
				default:
					t.Fatalf("Unsupported test input type: %T", input)
				}

				tbl.Append(luaValue)
			}

			L.Push(tbl)

			gluasprig.AllFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")

			boolResult := bool(result.(lua.LBool))
			require.Equal(t, tt.expected, boolResult)
		})
	}
}

func TestAnyFunc(t *testing.T) {
	tests := []struct {
		inputs   []any
		expected bool
	}{
		{
			inputs:   []any{1, "hello", true, 42.5},
			expected: true,
		},
		{
			inputs:   []any{1, "", true, 42.5},
			expected: true,
		},
		{
			inputs:   []any{nil, 0, "", false, "something"},
			expected: true,
		},
		{
			inputs:   []any{nil, "", false, 0},
			expected: false,
		},
		{
			inputs:   []any{},
			expected: false,
		},
		{
			inputs:   []any{[]int{}, map[string]string{}},
			expected: false,
		},
		{
			inputs:   []any{[]int{1, 2}, map[string]string{}},
			expected: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.CreateTable(len(tt.inputs), 0)

			for _, input := range tt.inputs {
				var luaValue lua.LValue

				switch v := input.(type) {
				case nil:
					luaValue = lua.LNil
				case int:
					luaValue = lua.LNumber(v)
				case int64:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case string:
					luaValue = lua.LString(v)
				case bool:
					luaValue = lua.LBool(v)
				case []int:
					tbl := L.CreateTable(len(v), 0)
					for _, val := range v {
						tbl.Append(lua.LNumber(val))
					}

					luaValue = tbl
				case map[string]string:
					tbl := L.CreateTable(0, len(v))
					for key, val := range v {
						tbl.RawSetString(key, lua.LString(val))
					}

					luaValue = tbl
				default:
					t.Fatalf("Unsupported test input type: %T", input)
				}

				tbl.Append(luaValue)
			}

			L.Push(tbl)

			gluasprig.AnyFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")

			boolResult := bool(result.(lua.LBool))
			require.Equal(t, tt.expected, boolResult)
		})
	}
}

func TestB32decFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "ORSXG5A=",
			expected: "test",
		},
		{
			input:    "ORSXG5BRGIZQ====",
			expected: "test123",
		},
		{
			input:    "NBSWY3DPEB3W64TMMQ======",
			expected: "hello world",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.B32decFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestB32encFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "test",
			expected: "ORSXG5A=",
		},
		{
			input:    "test123",
			expected: "ORSXG5BRGIZQ====",
		},
		{
			input:    "hello world",
			expected: "NBSWY3DPEB3W64TMMQ======",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.B32encFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestB64decFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "dGVzdA==",
			expected: "test",
		},
		{
			input:    "dGVzdDEyMw==",
			expected: "test123",
		},
		{
			input:    "aGVsbG8gd29ybGQ=",
			expected: "hello world",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.B64decFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestB64encFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "test",
			expected: "dGVzdA==",
		},
		{
			input:    "test123",
			expected: "dGVzdDEyMw==",
		},
		{
			input:    "hello world",
			expected: "aGVsbG8gd29ybGQ=",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.B64encFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestBaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: "file.txt",
		},
		{
			input:    "/path/to/directory/",
			expected: "directory",
		},
		{
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			input:    "",
			expected: ".",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.BaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestBcryptFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "password123",
		},
		{
			input: "secure-password",
		},
		{
			input: "simple",
		},
		{
			input: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.BcryptFunc(L)

			hashedPassword := L.ToString(-1)
			if tt.input != "" {
				require.NotEmpty(t, hashedPassword)

				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.input))
				require.NoError(t, err, "Bcrypt hash verification failed")
			}
		})
	}
}

func TestCamelcaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "HelloWorld",
		},
		{
			input:    "hello-world",
			expected: "HelloWorld",
		},
		{
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			input:    "HELLO WORLD",
			expected: "HelloWorld",
		},
		{
			input:    "HelloWorld",
			expected: "HelloWorld",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.CamelcaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCatFunc(t *testing.T) {
	tests := []struct {
		inputs   []any
		expected string
	}{
		{
			inputs:   []any{"hello", "world"},
			expected: "hello world",
		},
		{
			inputs:   []any{"testing", 123, true},
			expected: "testing 123 true",
		},
		{
			inputs:   []any{"single"},
			expected: "single",
		},
		{
			inputs:   []any{1, 2, 3, 4, 5},
			expected: "1 2 3 4 5",
		},
		{
			inputs:   []any{"mixed", 42, true, 3.14},
			expected: "mixed 42 true 3.14",
		},
		{
			inputs:   []any{nil, "after nil", 42},
			expected: "after nil 42",
		},
		{
			inputs:   []any{},
			expected: "",
		},
		{
			inputs:   []any{42.0, 42.5},
			expected: "42 42.5",
		},
		{
			inputs:   []any{-123, -45.67},
			expected: "-123 -45.67",
		},
		{
			inputs:   []any{0, 0.0},
			expected: "0 0",
		},
		{
			inputs:   []any{1e6, 1e-6},
			expected: "1000000 1e-06",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			for _, input := range tt.inputs {
				var luaValue lua.LValue
				switch v := input.(type) {
				case nil:
					luaValue = lua.LNil
				case int:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case string:
					luaValue = lua.LString(v)
				case bool:
					luaValue = lua.LBool(v)
				default:
					luaValue = lua.LString(fmt.Sprint(v))
				}
				L.Push(luaValue)
			}

			gluasprig.CatFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: ".",
		},
		{
			input:    "/path/to//file/../dir/",
			expected: "/path/to/dir",
		},
		{
			input:    "./some/path/../file.txt",
			expected: "some/file.txt",
		},
		{
			input:    "a/./b/../../c/",
			expected: "c",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.CleanFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCoalesceFunc(t *testing.T) {
	tests := []struct {
		inputs        []any
		expectedType  lua.LValueType
		expectedValue any
	}{
		{
			inputs:        []any{"", nil, "hello", "world"},
			expectedType:  lua.LTString,
			expectedValue: "hello",
		},
		{
			inputs:        []any{nil, "", 0, false, 42, "test"},
			expectedType:  lua.LTNumber,
			expectedValue: float64(42),
		},
		{
			inputs:        []any{nil, "", 0, false, "test"},
			expectedType:  lua.LTString,
			expectedValue: "test",
		},
		{
			inputs:        []any{nil, "", 0},
			expectedType:  lua.LTNil,
			expectedValue: nil,
		},
		{
			inputs:        []any{nil, "", false},
			expectedType:  lua.LTNil,
			expectedValue: nil,
		},
		{
			inputs:        []any{nil, ""},
			expectedType:  lua.LTNil,
			expectedValue: nil,
		},
		{
			inputs:        []any{nil},
			expectedType:  lua.LTNil,
			expectedValue: nil,
		},
		{
			inputs:        []any{true, 42, "test"},
			expectedType:  lua.LTBool,
			expectedValue: true,
		},
		{
			inputs:        []any{3.14, "pi"},
			expectedType:  lua.LTNumber,
			expectedValue: float64(3.14),
		},
		{
			inputs:        []any{"hello", 42, true},
			expectedType:  lua.LTString,
			expectedValue: "hello",
		},
		{
			inputs:        []any{"", 0, false, "non-empty"},
			expectedType:  lua.LTString,
			expectedValue: "non-empty",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			for _, input := range tt.inputs {
				var luaValue lua.LValue

				switch v := input.(type) {
				case nil:
					luaValue = lua.LNil
				case int:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case string:
					luaValue = lua.LString(v)
				case bool:
					luaValue = lua.LBool(v)
				default:
					luaValue = lua.LString(fmt.Sprint(v))
				}

				L.Push(luaValue)
			}

			gluasprig.CoalesceFunc(L)

			result := L.Get(-1)
			require.Equal(t, tt.expectedType, result.Type(), "Expected type %v but got %v for inputs %v", tt.expectedType, result.Type(), tt.inputs)

			switch tt.expectedType {
			case lua.LTNil:
			case lua.LTString:
				require.Equal(t, tt.expectedValue, result.String(), "For inputs %v", tt.inputs)
			case lua.LTNumber:
				require.Equal(t, tt.expectedValue, float64(result.(lua.LNumber)), "For inputs %v", tt.inputs)
			case lua.LTBool:
				require.Equal(t, tt.expectedValue, bool(result.(lua.LBool)), "For inputs %v", tt.inputs)
			}
		})
	}
}

func TestCompactFunc(t *testing.T) {
	tests := []struct {
		input    []any
		expected []any
	}{
		{
			input:    []any{1, "a", "foo", ""},
			expected: []any{float64(1), "a", "foo"},
		},
		{
			input:    []any{nil, "", 0, false, "hello", 42, true},
			expected: []any{"hello", float64(42), true},
		},
		{
			input:    []any{"", "", ""},
			expected: []any{},
		},
		{
			input:    []any{},
			expected: []any{},
		},
		{
			input:    []any{1, 2, 3},
			expected: []any{float64(1), float64(2), float64(3)},
		},
		{
			input:    []any{"", nil, false, 0, "string", 42, true},
			expected: []any{"string", float64(42), true},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.CreateTable(len(tt.input), 0)
			for i, item := range tt.input {
				var luaValue lua.LValue

				switch v := item.(type) {
				case nil:
					luaValue = lua.LNil
				case string:
					luaValue = lua.LString(v)
				case int:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case bool:
					luaValue = lua.LBool(v)
				default:
					luaValue = lua.LString(fmt.Sprint(v))
				}

				tbl.RawSetInt(i+1, luaValue)
			}

			L.Push(tbl)

			gluasprig.CompactFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTTable, result.Type(), "Expected table return type")

			resultTbl := result.(*lua.LTable)
			resultArray := make([]any, 0)

			resultTbl.ForEach(func(idx, value lua.LValue) {
				require.Equal(t, lua.LTNumber, idx.Type(), "Expected numeric index")

				switch value.Type() {
				case lua.LTNil:
					resultArray = append(resultArray, nil)
				case lua.LTString:
					resultArray = append(resultArray, value.String())
				case lua.LTNumber:
					resultArray = append(resultArray, float64(value.(lua.LNumber)))
				case lua.LTBool:
					resultArray = append(resultArray, bool(value.(lua.LBool)))
				default:
					resultArray = append(resultArray, value.String())
				}
			})

			require.Equal(t, len(tt.expected), len(resultArray), "Result array length doesn't match expected length")
			for i, v := range tt.expected {
				require.Equal(t, v, resultArray[i], "Value at index %d doesn't match", i)
			}
		})
	}
}

func TestDecryptAESFunc(t *testing.T) {
	mustEncryptAES := func(password, text string) string {
		if text == "" {
			return ""
		}

		// create 32 byte key
		key := make([]byte, 32)
		copy(key, []byte(password))

		// create cipher block
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(fmt.Sprintf("failed to create cipher: %v", err))
		}

		// pad text to multiple of block size
		paddingSize := aes.BlockSize - (len(text) % aes.BlockSize)
		padded := make([]byte, len(text)+paddingSize)
		copy(padded, text)
		for i := len(text); i < len(padded); i++ {
			padded[i] = byte(paddingSize)
		}

		// create IV and ciphertext buffer
		ciphertext := make([]byte, aes.BlockSize+len(padded))
		iv := ciphertext[:aes.BlockSize]
		if _, err := rand.Read(iv); err != nil {
			panic(fmt.Sprintf("failed to read random IV: %v", err))
		}

		// encrypt
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

		// base64 encode
		return base64.StdEncoding.EncodeToString(ciphertext)
	}

	const key = "secret-key-123456"

	tests := []struct {
		ciphertext   string
		expectedText string
	}{
		{
			ciphertext:   mustEncryptAES(key, "hello world"),
			expectedText: "hello world",
		},
		{
			ciphertext:   mustEncryptAES(key, "test123"),
			expectedText: "test123",
		},
		{
			ciphertext:   mustEncryptAES(key, ""),
			expectedText: "",
		},
		{
			ciphertext:   "",
			expectedText: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(key))
			L.Push(lua.LString(tt.ciphertext))

			gluasprig.DecryptAESFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			require.Equal(t, lua.LTString, result.Type(), "Expected string return type")
			require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
			require.Equal(t, tt.expectedText, result.String())
		})
	}
}

func TestDerivePasswordFunc(t *testing.T) {
	tests := []struct {
		counter      uint32
		passType     string
		password     string
		user         string
		site         string
		expectedLen  int
		expectedType string
	}{
		{
			counter:      1,
			passType:     "medium",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  8,
			expectedType: "mixed",
		},
		{
			counter:      2,
			passType:     "medium",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  8,
			expectedType: "mixed",
		},
		{
			counter:      1,
			passType:     "long",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  14,
			expectedType: "mixed",
		},
		{
			counter:      1,
			passType:     "basic",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  8,
			expectedType: "mixed",
		},
		{
			counter:      1,
			passType:     "short",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  4,
			expectedType: "mixed",
		},
		{
			counter:      1,
			passType:     "pin",
			password:     "password123",
			user:         "user@example.com",
			site:         "example.com",
			expectedLen:  4,
			expectedType: "numeric",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			var derivedPassword1, derivedPassword2, derivedPassword3 string

			{
				L := lua.NewState()
				defer L.Close()

				L.Push(lua.LNumber(tt.counter))
				L.Push(lua.LString(tt.passType))
				L.Push(lua.LString(tt.password))
				L.Push(lua.LString(tt.user))
				L.Push(lua.LString(tt.site))

				gluasprig.DerivePasswordFunc(L)

				result := L.Get(-1)
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")

				derivedPassword1 = result.String()
			}

			{
				L := lua.NewState()
				defer L.Close()

				L.Push(lua.LNumber(tt.counter))
				L.Push(lua.LString(tt.passType))
				L.Push(lua.LString(tt.password))
				L.Push(lua.LString(tt.user))
				L.Push(lua.LString(tt.site))

				gluasprig.DerivePasswordFunc(L)

				result := L.Get(-1)
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")

				derivedPassword2 = result.String()
			}

			{
				L := lua.NewState()
				defer L.Close()

				L.Push(lua.LNumber(tt.counter + 1))
				L.Push(lua.LString(tt.passType))
				L.Push(lua.LString(tt.password))
				L.Push(lua.LString(tt.user))
				L.Push(lua.LString(tt.site))

				gluasprig.DerivePasswordFunc(L)

				result := L.Get(-1)
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")

				derivedPassword3 = result.String()
			}

			require.Equal(t, tt.expectedLen, len(derivedPassword1),
				"Derived password length doesn't match expected for type %s", tt.passType)

			switch tt.expectedType {
			case "numeric":
				require.Regexp(t, "^[0-9]+$", derivedPassword1,
					"PIN should only contain digits")
			case "mixed":
				hasLower := regexp.MustCompile(`[a-z]`).MatchString(derivedPassword1)
				hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(derivedPassword1)
				hasDigit := regexp.MustCompile(`[0-9]`).MatchString(derivedPassword1)

				if tt.passType == "long" || tt.passType == "medium" {
					require.True(t, hasLower && hasUpper && hasDigit,
						"Expected complex password with lower, upper and digits for type %s", tt.passType)
				} else {
					require.True(t, hasLower || hasUpper || hasDigit,
						"Expected at least some complexity for type %s", tt.passType)
				}
			}

			require.Equal(t, derivedPassword1, derivedPassword2,
				"Derived passwords should be deterministic for the same inputs")

			require.NotEqual(t, derivedPassword1, derivedPassword3,
				"Different counter should produce different password")
		})
	}
}

func TestDirFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: "/path/to",
		},
		{
			input:    "/path/to/directory/",
			expected: "/path/to/directory",
		},
		{
			input:    "file.txt",
			expected: ".",
		},
		{
			input:    "/",
			expected: "/",
		},
		{
			input:    "",
			expected: ".",
		},
		{
			input:    "./relative/path",
			expected: "relative",
		},
		{
			input:    "../parent/path",
			expected: "../parent",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.DirFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestDurationFunc(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{
			input:    int64(3600),
			expected: "1h0m0s",
		},
		{
			input:    int64(3661),
			expected: "1h1m1s",
		},
		{
			input:    "86400",
			expected: "24h0m0s",
		},
		{
			input:    "300",
			expected: "5m0s",
		},
		{
			input:    "invalid",
			expected: "0s",
		},
		{
			input:    nil,
			expected: "0s",
		},
		{
			input:    int64(0),
			expected: "0s",
		},
		{
			input:    "-120",
			expected: "-2m0s",
		},
		{
			input:    float64(60),
			expected: "1m0s",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			var luaValue lua.LValue

			switch v := tt.input.(type) {
			case nil:
				luaValue = lua.LNil
			case float64:
				luaValue = lua.LNumber(v)
			case int64:
				luaValue = lua.LNumber(v)
			case string:
				luaValue = lua.LString(v)
			default:
				luaValue = lua.LString(fmt.Sprint(v))
			}

			L.Push(luaValue)

			gluasprig.DurationFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestDurationRoundFunc(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{
			input:    "3600s",
			expected: "60m",
		},
		{
			input:    "1h1m1s",
			expected: "1h",
		},
		{
			input:    "24h",
			expected: "24h",
		},
		{
			input:    "5m",
			expected: "5m",
		},
		{
			input:    "invalid",
			expected: "0s",
		},
		{
			input:    int64(3600),
			expected: "60m",
		},
		{
			input:    float64(60),
			expected: "60s",
		},
		{
			input:    int64(86400),
			expected: "24h",
		},
		{
			input:    int64(0),
			expected: "0s",
		},
		{
			input:    nil,
			expected: "0s",
		},
		{
			input:    "-120s",
			expected: "2m",
		},
		{
			input:    int64(2592000),
			expected: "30d",
		},
		{
			input:    int64(31536000),
			expected: "12mo",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			var luaValue lua.LValue

			switch v := tt.input.(type) {
			case nil:
				luaValue = lua.LNil
			case float64:
				luaValue = lua.LNumber(v)
			case int64:
				luaValue = lua.LNumber(v)
			case string:
				luaValue = lua.LString(v)
			default:
				luaValue = lua.LString(fmt.Sprint(v))
			}

			L.Push(luaValue)

			gluasprig.DurationRoundFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestEmptyFunc(t *testing.T) {
	tests := []struct {
		input    lua.LValue
		expected bool
	}{
		{
			input:    lua.LNil,
			expected: true,
		},
		{
			input:    lua.LString(""),
			expected: true,
		},
		{
			input:    lua.LString("hello"),
			expected: false,
		},
		{
			input:    lua.LNumber(0),
			expected: true,
		},
		{
			input:    lua.LNumber(42),
			expected: false,
		},
		{
			input:    lua.LNumber(-1),
			expected: false,
		},
		{
			input:    lua.LBool(false),
			expected: true,
		},
		{
			input:    lua.LBool(true),
			expected: false,
		},
		{
			input: func() lua.LValue {
				L := lua.NewState()

				return L.NewTable()
			}(),
			expected: true,
		},
		{
			input: func() lua.LValue {
				L := lua.NewState()
				tbl := L.NewTable()

				tbl.RawSetString("key", lua.LString("value"))

				return tbl
			}(),
			expected: false,
		},
		{
			input: func() lua.LValue {
				L := lua.NewState()
				tbl := L.NewTable()

				tbl.Append(lua.LNumber(1))

				return tbl
			}(),
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(tt.input)

			gluasprig.EmptyFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")

			boolResult := bool(result.(lua.LBool))
			require.Equal(t, tt.expected, boolResult)
		})
	}
}

func TestEncryptAESFunc(t *testing.T) {
	mustDecryptAES := func(password, ciphertext string) string {
		if ciphertext == "" {
			return ""
		}

		// decode base64
		data, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			panic(fmt.Sprintf("failed to decode base64: %v", err))
		}

		// create cipher block
		key := make([]byte, 32)
		copy(key, []byte(password))
		block, err := aes.NewCipher(key)
		if err != nil {
			panic(fmt.Sprintf("failed to create cipher: %v", err))
		}

		// extract IV and ciphertext
		if len(data) < aes.BlockSize {
			panic("ciphertext too short")
		}
		iv := data[:aes.BlockSize]
		cipherData := data[aes.BlockSize:]

		// decrypt
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(cipherData, cipherData)

		// unpad
		padding := int(cipherData[len(cipherData)-1])
		if padding > aes.BlockSize || padding > len(cipherData) {
			panic("invalid padding size")
		}

		// ensure padding is correct
		for i := len(cipherData) - padding; i < len(cipherData); i++ {
			if int(cipherData[i]) != padding {
				panic("invalid padding")
			}
		}

		// return plaintext
		return string(cipherData[:len(cipherData)-padding])
	}

	tests := []struct {
		key  string
		text string
	}{
		{
			key:  "secret-key-12345678901234567890123456",
			text: "hello world",
		},
		{
			key:  "another-secret-key-1234567890123456",
			text: "test123",
		},
		{
			key:  "short-key",
			text: "test text",
		},
		{
			key:  "secret-key-12345678901234567890123456",
			text: "",
		},
		{
			key:  "",
			text: "some text",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.key))
			L.Push(lua.LString(tt.text))

			gluasprig.EncryptAESFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			require.Equal(t, lua.LTString, result.Type(), "Expected string return type")
			require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")

			if tt.text == "" {
				require.Equal(t, "", result.String())
			} else {
				require.NotEqual(t, "", result.String())

				decrypted := mustDecryptAES(tt.key, result.String())
				require.Equal(t, tt.text, decrypted, "Decrypted text should match original")
			}
		})
	}
}

func TestExtFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: ".txt",
		},
		{
			input:    "file.jpg",
			expected: ".jpg",
		},
		{
			input:    "/path/to/file",
			expected: "",
		},
		{
			input:    "file",
			expected: "",
		},
		{
			input:    ".htaccess",
			expected: ".htaccess",
		},
		{
			input:    "/path/with.dots/file.go",
			expected: ".go",
		},
		{
			input:    "/path/to/archive.tar.gz",
			expected: ".gz",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.ExtFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGenPrivateKeyFunc(t *testing.T) {
	tests := []struct {
		keyType      string
		expectedType string
		minLength    int
	}{
		{
			keyType:      "rsa",
			expectedType: "RSA PRIVATE KEY",
			minLength:    100,
		},
		{
			keyType:      "",
			expectedType: "RSA PRIVATE KEY",
			minLength:    100,
		},
		{
			keyType:      "dsa",
			expectedType: "DSA PRIVATE KEY",
			minLength:    100,
		},
		{
			keyType:      "ecdsa",
			expectedType: "EC PRIVATE KEY",
			minLength:    50,
		},
		{
			keyType:      "ed25519",
			expectedType: "PRIVATE KEY",
			minLength:    10,
		},
		{
			keyType:      "unknown",
			expectedType: "",
			minLength:    0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.keyType))

			gluasprig.GenPrivateKeyFunc(L)

			result := L.ToString(-1)

			if tt.keyType == "unknown" {
				require.Contains(t, result, "Unknown type", "Should return error for unknown key type")
			} else {
				block, _ := pem.Decode([]byte(result))
				require.NotNil(t, block, "Result should be decodable as PEM")
				require.Equal(t, tt.expectedType, block.Type, "PEM block should have correct type")

				require.True(t, len(block.Bytes) >= tt.minLength,
					"Key data should have reasonable length (at least %d bytes)", tt.minLength)
			}
		})
	}
}

func TestHtpasswdFunc(t *testing.T) {
	tests := []struct {
		username string
		password string
		expected string
	}{
		{
			username: "user1",
			password: "password123",
			expected: "user1:$2a$",
		},
		{
			username: "admin",
			password: "admin123",
			expected: "admin:$2a$",
		},
		{
			username: "",
			password: "emptyuser",
			expected: ":$2a$",
		},
		{
			username: "user:with:colons",
			password: "test",
			expected: "invalid username",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.username))
			L.Push(lua.LString(tt.password))

			gluasprig.HtpasswdFunc(L)

			result := L.ToString(-1)

			if tt.username == "user:with:colons" {
				require.Contains(t, result, tt.expected, "Should reject usernames with colons")
			} else {
				require.True(t, strings.HasPrefix(result, tt.expected),
					"Result should start with expected prefix. Got: %s, Expected prefix: %s",
					result, tt.expected)

				if strings.HasPrefix(tt.expected, tt.username+":$2a$") {
					parts := strings.SplitN(result, ":", 2)
					require.Len(t, parts, 2, "Result should have username and hash parts")
					require.Equal(t, tt.username, parts[0], "Username part should match")

					err := bcrypt.CompareHashAndPassword([]byte(parts[1]), []byte(tt.password))
					require.NoError(t, err, "Bcrypt hash should verify against original password")
				}
			}
		})
	}
}

func TestIndentFunc(t *testing.T) {
	tests := []struct {
		spaces   int
		input    string
		expected string
	}{
		{
			spaces:   4,
			input:    "hello\nworld",
			expected: "    hello\n    world",
		},
		{
			spaces:   2,
			input:    "line1\nline2\nline3",
			expected: "  line1\n  line2\n  line3",
		},
		{
			spaces:   0,
			input:    "no indentation",
			expected: "no indentation",
		},
		{
			spaces:   3,
			input:    "",
			expected: "   ",
		},
		{
			spaces:   5,
			input:    "single line",
			expected: "     single line",
		},
		{
			spaces:   1,
			input:    "\n\n\n",
			expected: " \n \n \n ",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.spaces))
			L.Push(lua.LString(tt.input))

			gluasprig.IndentFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestInitialsFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "John Doe",
			expected: "JD",
		},
		{
			input:    "jane smith johnson",
			expected: "jsj",
		},
		{
			input:    "JAMES BOND",
			expected: "JB",
		},
		{
			input:    "J.R.R. Tolkien",
			expected: "JT",
		},
		{
			input:    "single",
			expected: "s",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "multiple   spaces  between words",
			expected: "msbw",
		},
		{
			input:    " leading-and-trailing-spaces ",
			expected: "l",
		},
		{
			input:    "hyphenated-name",
			expected: "h",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.InitialsFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestIsAbsFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "/absolute/path",
			expected: true,
		},
		{
			input:    "/",
			expected: true,
		},
		{
			input:    "relative/path",
			expected: false,
		},
		{
			input:    "./relative",
			expected: false,
		},
		{
			input:    "../parent",
			expected: false,
		},
		{
			input:    "",
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.IsAbsFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")

			boolResult := bool(result.(lua.LBool))
			require.Equal(t, tt.expected, boolResult)
		})
	}
}

func TestKebabcaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "HelloWorld",
			expected: "hello-world",
		},
		{
			input:    "helloWorld",
			expected: "hello-world",
		},
		{
			input:    "hello world",
			expected: "hello-world",
		},
		{
			input:    "HELLO WORLD",
			expected: "hello-world",
		},
		{
			input:    "hello_world",
			expected: "hello-world",
		},
		{
			input:    "hello-world",
			expected: "hello-world",
		},
		{
			input:    "hello123world",
			expected: "hello-123world",
		},
		{
			input:    "Hello-World_Example",
			expected: "hello-world-example",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.KebabcaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestMustRegexFindAllFunc(t *testing.T) {
	tests := []struct {
		regex    string
		s        string
		n        int
		expected []string
		hasError bool
	}{
		{
			regex:    `\d+`,
			s:        "abc123def456ghi",
			n:        -1,
			expected: []string{"123", "456"},
			hasError: false,
		},
		{
			regex:    `\w{3}`,
			s:        "abc123def456ghi",
			n:        2,
			expected: []string{"abc", "123"},
			hasError: false,
		},
		{
			regex:    `[a-z]+`,
			s:        "abc123def456ghi",
			n:        -1,
			expected: []string{"abc", "def", "ghi"},
			hasError: false,
		},
		{
			regex:    `\d`,
			s:        "no digits here",
			n:        -1,
			expected: []string{},
			hasError: false,
		},
		{
			regex:    `[`,
			s:        "invalid regex pattern",
			n:        -1,
			expected: nil,
			hasError: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.s))
			L.Push(lua.LNumber(tt.n))

			gluasprig.MustRegexFindAllFunc(L)

			if tt.hasError {
				result := L.Get(-2)
				errMsg := L.Get(-1)

				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result for invalid regex")
				require.Equal(t, lua.LTString, errMsg.Type(), "Expected error message")
			} else {
				table := L.Get(-2)
				errMsg := L.Get(-1)

				require.Equal(t, lua.LTTable, table.Type(), "Expected table result")
				require.Equal(t, lua.LTNil, errMsg.Type(), "Expected nil error")

				tbl := table.(*lua.LTable)
				results := make([]string, 0, tbl.Len())

				tbl.ForEach(func(idx lua.LValue, value lua.LValue) {
					require.Equal(t, lua.LTNumber, idx.Type(), "Expected numeric index")
					require.Equal(t, lua.LTString, value.Type(), "Expected string value")

					results = append(results, value.String())
				})

				require.Equal(t, len(tt.expected), len(results), "Table length doesn't match expected length")

				if len(tt.expected) == 0 {
					require.Empty(t, results, "Expected empty results")
				} else {
					require.Equal(t, tt.expected, results, "Table results don't match expected values")
				}

				require.Equal(t, len(tt.expected), tbl.Len(), "Table length doesn't match expected length")
			}
		})
	}
}

func TestMustRegexFindFunc(t *testing.T) {
	tests := []struct {
		regex    string
		input    string
		expected string
		wantErr  bool
	}{
		{
			regex:    "\\d+",
			input:    "abc123def",
			expected: "123",
			wantErr:  false,
		},
		{
			regex:    "[a-z]+",
			input:    "123abc456",
			expected: "abc",
			wantErr:  false,
		},
		{
			regex:    "foo",
			input:    "bar",
			expected: "",
			wantErr:  false,
		},
		{
			regex:    "^hello",
			input:    "hello world",
			expected: "hello",
			wantErr:  false,
		},
		{
			regex:    "(\\w+)@(\\w+).com",
			input:    "contact me at user@example.com",
			expected: "user@example.com",
			wantErr:  false,
		},
		{
			regex:    "(",
			input:    "test",
			expected: "",
			wantErr:  true,
		},
		{
			regex:    "",
			input:    "test",
			expected: "",
			wantErr:  false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.input))

			gluasprig.MustRegexFindFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			if tt.wantErr {
				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result when error")
				require.NotEqual(t, lua.LTNil, errValue.Type(), "Expected non-nil error")
			} else {
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")
				require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
				require.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestMustRegexMatchFunc(t *testing.T) {
	tests := []struct {
		regex    string
		input    string
		expected bool
		wantErr  bool
	}{
		{
			regex:    "^\\d+$",
			input:    "12345",
			expected: true,
			wantErr:  false,
		},
		{
			regex:    "^\\d+$",
			input:    "12345abc",
			expected: false,
			wantErr:  false,
		},
		{
			regex:    "^[a-z]+$",
			input:    "abcdef",
			expected: true,
			wantErr:  false,
		},
		{
			regex:    "^hello.*world$",
			input:    "hello world",
			expected: true,
			wantErr:  false,
		},
		{
			regex:    "^hello.*world$",
			input:    "hello wonderful world",
			expected: true,
			wantErr:  false,
		},
		{
			regex:    "^hello.*world$",
			input:    "hello but no world here",
			expected: false,
			wantErr:  false,
		},
		{
			regex:    "(",
			input:    "test",
			expected: false,
			wantErr:  true,
		},
		{
			regex:    "",
			input:    "",
			expected: true,
			wantErr:  false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.input))

			gluasprig.MustRegexMatchFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			if tt.wantErr {
				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result when error")
				require.NotEqual(t, lua.LTNil, errValue.Type(), "Expected non-nil error")
			} else {
				require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")
				require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
				require.Equal(t, tt.expected, bool(result.(lua.LBool)))
			}
		})
	}
}

func TestMustRegexReplaceAllFunc(t *testing.T) {
	tests := []struct {
		regex       string
		replacement string
		input       string
		expected    string
		wantErr     bool
	}{
		{
			regex:       "\\d+",
			replacement: "NUM",
			input:       "abc123def456",
			expected:    "abcNUMdefNUM",
			wantErr:     false,
		},
		{
			regex:       "[a-z]+",
			replacement: "WORD",
			input:       "123abc456def",
			expected:    "123WORD456WORD",
			wantErr:     false,
		},
		{
			regex:       "foo",
			replacement: "bar",
			input:       "footest foobar",
			expected:    "bartest barbar",
			wantErr:     false,
		},
		{
			regex:       "(\\w+)@(\\w+)",
			replacement: "$1@example",
			input:       "contact user@domain.com",
			expected:    "contact user@example.com",
			wantErr:     false,
		},
		{
			regex:       "\\s+",
			replacement: " ",
			input:       "multiple    spaces\t\tbetween\nwords",
			expected:    "multiple spaces between words",
			wantErr:     false,
		},
		{
			regex:       "(",
			replacement: "invalid",
			input:       "test with invalid regex",
			expected:    "",
			wantErr:     true,
		},
		{
			regex:       "",
			replacement: "x",
			input:       "empty regex",
			expected:    "xexmxpxtxyx xrxexgxexxx",
			wantErr:     false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.input))
			L.Push(lua.LString(tt.replacement))

			gluasprig.MustRegexReplaceAllFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			if tt.wantErr {
				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result when error")
				require.NotEqual(t, lua.LTNil, errValue.Type(), "Expected non-nil error")
			} else {
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")
				require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
				require.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestMustRegexReplaceAllLiteralFunc(t *testing.T) {
	tests := []struct {
		regex       string
		replacement string
		input       string
		expected    string
		wantErr     bool
	}{
		{
			regex:       "\\d+",
			replacement: "NUM",
			input:       "abc123def456",
			expected:    "abcNUMdefNUM",
			wantErr:     false,
		},
		{
			regex:       "[a-z]+",
			replacement: "WORD",
			input:       "123abc456def",
			expected:    "123WORD456WORD",
			wantErr:     false,
		},
		{
			regex:       "foo",
			replacement: "bar",
			input:       "footest foobar",
			expected:    "bartest barbar",
			wantErr:     false,
		},
		{
			regex:       "(\\w+)@(\\w+)",
			replacement: "$1@example",
			input:       "contact user@domain.com",
			expected:    "contact $1@example.com",
			wantErr:     false,
		},
		{
			regex:       "\\$\\d",
			replacement: "DOLLAR",
			input:       "Price: $5 and $9",
			expected:    "Price: DOLLAR and DOLLAR",
			wantErr:     false,
		},
		{
			regex:       "(",
			replacement: "invalid",
			input:       "test with invalid regex",
			expected:    "",
			wantErr:     true,
		},
		{
			regex:       "",
			replacement: "x",
			input:       "empty regex",
			expected:    "xexmxpxtxyx xrxexgxexxx",
			wantErr:     false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.input))
			L.Push(lua.LString(tt.replacement))

			gluasprig.MustRegexReplaceAllLiteralFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			if tt.wantErr {
				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result when error")
				require.NotEqual(t, lua.LTNil, errValue.Type(), "Expected non-nil error")
			} else {
				require.Equal(t, lua.LTString, result.Type(), "Expected string return type")
				require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
				require.Equal(t, tt.expected, result.String())
			}
		})
	}
}

func TestMustRegexSplitFunc(t *testing.T) {
	tests := []struct {
		regex    string
		s        string
		n        int
		expected []string
		hasError bool
	}{
		{
			regex:    `,`,
			s:        "a,b,c,d",
			n:        -1,
			expected: []string{"a", "b", "c", "d"},
			hasError: false,
		},
		{
			regex:    `\s+`,
			s:        "hello  world   test",
			n:        -1,
			expected: []string{"hello", "world", "test"},
			hasError: false,
		},
		{
			regex:    `;`,
			s:        "one;two;three",
			n:        2,
			expected: []string{"one", "two;three"},
			hasError: false,
		},
		{
			regex:    `\d+`,
			s:        "abc123def456",
			n:        -1,
			expected: []string{"abc", "def", ""},
			hasError: false,
		},
		{
			regex:    `,`,
			s:        "no commas here",
			n:        -1,
			expected: []string{"no commas here"},
			hasError: false,
		},
		{
			regex:    `[`,
			s:        "invalid regex pattern",
			n:        -1,
			expected: nil,
			hasError: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.regex))
			L.Push(lua.LString(tt.s))
			L.Push(lua.LNumber(tt.n))

			gluasprig.MustRegexSplitFunc(L)

			if tt.hasError {
				result := L.Get(-2)
				errMsg := L.Get(-1)

				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result for invalid regex")
				require.Equal(t, lua.LTString, errMsg.Type(), "Expected error message")
			} else {
				table := L.Get(-2)
				errMsg := L.Get(-1)

				require.Equal(t, lua.LTTable, table.Type(), "Expected table result")
				require.Equal(t, lua.LTNil, errMsg.Type(), "Expected nil error")

				tbl := table.(*lua.LTable)
				results := make([]string, 0, tbl.Len())

				tbl.ForEach(func(idx lua.LValue, value lua.LValue) {
					require.Equal(t, lua.LTNumber, idx.Type(), "Expected numeric index")
					require.Equal(t, lua.LTString, value.Type(), "Expected string value")

					results = append(results, value.String())
				})

				require.Equal(t, len(tt.expected), len(results), "Table length doesn't match expected length")

				if len(tt.expected) == 0 {
					require.Empty(t, results, "Expected empty results")
				} else {
					require.Equal(t, tt.expected, results, "Table results don't match expected values")
				}

				require.Equal(t, len(tt.expected), tbl.Len(), "Table length doesn't match expected length")
			}
		})
	}
}

func TestNindentFunc(t *testing.T) {
	tests := []struct {
		spaces   int
		input    string
		expected string
	}{
		{
			spaces:   4,
			input:    "hello\nworld",
			expected: "\n    hello\n    world",
		},
		{
			spaces:   2,
			input:    "line1\nline2\nline3",
			expected: "\n  line1\n  line2\n  line3",
		},
		{
			spaces:   0,
			input:    "no indentation",
			expected: "\nno indentation",
		},
		{
			spaces:   3,
			input:    "",
			expected: "\n   ",
		},
		{
			spaces:   5,
			input:    "single line",
			expected: "\n     single line",
		},
		{
			spaces:   1,
			input:    "\n\n\n",
			expected: "\n \n \n \n ",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.spaces))
			L.Push(lua.LString(tt.input))

			gluasprig.NindentFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNospaceFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "helloworld",
		},
		{
			input:    "  spaces   at   ends  ",
			expected: "spacesatends",
		},
		{
			input:    "tabs\tand\tspaces",
			expected: "tabsandspaces",
		},
		{
			input:    "new\nlines\r\nand\rcarriage returns",
			expected: "newlinesandcarriagereturns",
		},
		{
			input:    "multiple    consecutive     spaces",
			expected: "multipleconsecutivespaces",
		},
		{
			input:    "   ",
			expected: "",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "special!@#$%^&*() characters",
			expected: "special!@#$%^&*()characters",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.NospaceFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestOsBaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: "file.txt",
		},
		{
			input:    "/path/to/directory/",
			expected: "directory",
		},
		{
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			input:    "",
			expected: ".",
		},
		{
			input:    "/",
			expected: "/",
		},
		{
			input:    "path/with/../relative/parts",
			expected: "parts",
		},
		{
			input:    "./relative/path",
			expected: "path",
		},
		{
			input:    "../parent/file",
			expected: "file",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.OsBaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestOsCleanFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to//file/../dir/",
			expected: "/path/to/dir",
		},
		{
			input:    "//multiple/slashes///path",
			expected: "/multiple/slashes/path",
		},
		{
			input:    "./some/path/../file.txt",
			expected: "some/file.txt",
		},
		{
			input:    "a/./b/../../c/",
			expected: "c",
		},
		{
			input:    "",
			expected: ".",
		},
		{
			input:    ".",
			expected: ".",
		},
		{
			input:    "..",
			expected: "..",
		},
		{
			input:    "/../path",
			expected: "/path",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.OsCleanFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestOsDirFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: "/path/to",
		},
		{
			input:    "/path/to/directory/",
			expected: "/path/to/directory",
		},
		{
			input:    "file.txt",
			expected: ".",
		},
		{
			input:    "/",
			expected: "/",
		},
		{
			input:    "",
			expected: ".",
		},
		{
			input:    "./relative/path",
			expected: "relative",
		},
		{
			input:    "../parent/path",
			expected: "../parent",
		},
		{
			input:    "path/with/../parts",
			expected: "path",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.OsDirFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestOsExtFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "/path/to/file.txt",
			expected: ".txt",
		},
		{
			input:    "file.jpg",
			expected: ".jpg",
		},
		{
			input:    "/path/to/file",
			expected: "",
		},
		{
			input:    "file",
			expected: "",
		},
		{
			input:    ".htaccess",
			expected: ".htaccess",
		},
		{
			input:    "/path/with.dots/file.go",
			expected: ".go",
		},
		{
			input:    "/path/to/archive.tar.gz",
			expected: ".gz",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.OsExtFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestOsIsAbsFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "/absolute/path",
			expected: true,
		},
		{
			input:    "/",
			expected: true,
		},
		{
			input:    "relative/path",
			expected: false,
		},
		{
			input:    "./relative",
			expected: false,
		},
		{
			input:    "../parent",
			expected: false,
		},
		{
			input:    "",
			expected: false,
		},
		{
			input:    "just\\backslashes",
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.OsIsAbsFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")

			boolResult := bool(result.(lua.LBool))
			require.Equal(t, tt.expected, boolResult)
		})
	}
}

func TestPluralFunc(t *testing.T) {
	tests := []struct {
		singular string
		plural   string
		count    int
		expected string
	}{
		{
			singular: "one anchovy",
			plural:   "many anchovies",
			count:    1,
			expected: "one anchovy",
		},
		{
			singular: "one anchovy",
			plural:   "many anchovies",
			count:    2,
			expected: "many anchovies",
		},
		{
			singular: "one anchovy",
			plural:   "many anchovies",
			count:    0,
			expected: "many anchovies",
		},
		{
			singular: "one fish",
			plural:   "many fish",
			count:    1,
			expected: "one fish",
		},
		{
			singular: "one fish",
			plural:   "many fish",
			count:    10,
			expected: "many fish",
		},
		{
			singular: "1 item",
			plural:   "$COUNT items",
			count:    1,
			expected: "1 item",
		},
		{
			singular: "1 item",
			plural:   "$COUNT items",
			count:    5,
			expected: "$COUNT items",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.singular))
			L.Push(lua.LString(tt.plural))
			L.Push(lua.LNumber(tt.count))

			gluasprig.PluralFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result, "Plural result incorrect")
		})
	}
}

func TestQuoteFunc(t *testing.T) {
	tests := []struct {
		input    []any
		expected string
	}{
		{
			input:    []any{"hello"},
			expected: `"hello"`,
		},
		{
			input:    []any{"hello", "world"},
			expected: `"hello" "world"`,
		},
		{
			input:    []any{"hello", 123, true},
			expected: `"hello" "123" "true"`,
		},
		{
			input:    []any{"string with \"quotes\""},
			expected: `"string with \"quotes\""`,
		},
		{
			input:    []any{"string", nil, "after nil"},
			expected: `"string" "after nil"`,
		},
		{
			input:    []any{},
			expected: ``,
		},
		{
			input:    []any{"multi\nline\nstring"},
			expected: `"multi\nline\nstring"`,
		},
		{
			input:    []any{`string\with\backslashes`},
			expected: `"string\\with\\backslashes"`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.NewTable()
			for _, item := range tt.input {
				var luaValue lua.LValue

				switch v := item.(type) {
				case nil:
					luaValue = lua.LNil
				case string:
					luaValue = lua.LString(v)
				case int:
					luaValue = lua.LNumber(v)
				case bool:
					luaValue = lua.LBool(v)
				default:
					luaValue = lua.LString(fmt.Sprint(v))
				}

				tbl.Append(luaValue)
			}

			L.Push(tbl)

			gluasprig.QuoteFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRandIntFunc(t *testing.T) {
	tests := []struct {
		min int
		max int
	}{
		{
			min: 1,
			max: 10,
		},
		{
			min: 0,
			max: 100,
		},
		{
			min: -10,
			max: 10,
		},
		{
			min: 1000,
			max: 2000,
		},
		{
			min: 5,
			max: 6,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.min))
			L.Push(lua.LNumber(tt.max))

			gluasprig.RandIntFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTNumber, result.Type(), "Expected number return type")

			numResult, ok := result.(lua.LNumber)
			require.True(t, ok)

			require.GreaterOrEqual(t, int(numResult), tt.min, "Result should be greater than or equal to min")
			require.LessOrEqual(t, int(numResult), tt.max, "Result should be less than or equal to max")
		})
	}
}

func TestRoundFunc(t *testing.T) {
	tests := []struct {
		value     any
		precision int
		roundOn   *float64
		expected  float64
	}{
		{
			value:     123.555555,
			precision: 3,
			expected:  123.556,
		},
		{
			value:     123.5555,
			precision: 2,
			expected:  123.56,
		},
		{
			value:     123.554,
			precision: 2,
			expected:  123.55,
		},
		{
			value:     "123.555",
			precision: 2,
			expected:  123.56,
		},
		{
			value:     123.555,
			precision: 0,
			expected:  124.0,
		},
		{
			value:     123.45,
			precision: 1,
			expected:  123.5,
		},
		{
			value:     123.45,
			precision: 1,
			roundOn:   func() *float64 { v := 0.4; return &v }(),
			expected:  123.5,
		},
		{
			value:     123.41,
			precision: 1,
			roundOn:   func() *float64 { v := 0.4; return &v }(),
			expected:  123.4,
		},
		{
			value:     -123.555,
			precision: 2,
			expected:  -123.56,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			switch v := tt.value.(type) {
			case string:
				L.Push(lua.LString(v))
			case float64:
				L.Push(lua.LNumber(v))
			default:
				L.Push(lua.LString(fmt.Sprint(v)))
			}

			L.Push(lua.LNumber(tt.precision))

			if tt.roundOn != nil {
				L.Push(lua.LNumber(*tt.roundOn))
			}

			gluasprig.RoundFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTNumber, result.Type(), "Expected number result")

			value := float64(result.(lua.LNumber))
			require.InDelta(t, tt.expected, value, 0.0001, "Result doesn't match expected value")
		})
	}
}

func TestSemverCompareFunc(t *testing.T) {
	tests := []struct {
		constraint string
		version    string
		expected   bool
		wantErr    bool
	}{
		{
			constraint: ">=1.0.0",
			version:    "1.0.0",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: ">=1.0.0",
			version:    "0.9.9",
			expected:   false,
			wantErr:    false,
		},
		{
			constraint: "<2.0.0",
			version:    "1.9.9",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: "=1.2.3",
			version:    "1.2.3",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: "=1.2.3",
			version:    "1.2.4",
			expected:   false,
			wantErr:    false,
		},
		{
			constraint: "~1.2.3",
			version:    "1.2.9",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: "~1.2.3",
			version:    "1.3.0",
			expected:   false,
			wantErr:    false,
		},
		{
			constraint: "^1.2.3",
			version:    "1.9.9",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: "^1.2.3",
			version:    "2.0.0",
			expected:   false,
			wantErr:    false,
		},
		{
			constraint: "1.2.x",
			version:    "1.2.3",
			expected:   true,
			wantErr:    false,
		},
		{
			constraint: "1.2.x",
			version:    "1.3.0",
			expected:   false,
			wantErr:    false,
		},
		{
			constraint: "invalid constraint",
			version:    "1.0.0",
			expected:   false,
			wantErr:    true,
		},
		{
			constraint: ">=1.0.0",
			version:    "invalid version",
			expected:   false,
			wantErr:    true,
		},
		{
			constraint: "",
			version:    "",
			expected:   false,
			wantErr:    true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.constraint))
			L.Push(lua.LString(tt.version))

			gluasprig.SemverCompareFunc(L)

			result := L.Get(-2)
			errValue := L.Get(-1)

			if tt.wantErr {
				require.Equal(t, lua.LTNil, result.Type(), "Expected nil result when error")
				require.NotEqual(t, lua.LTNil, errValue.Type(), "Expected non-nil error")
			} else {
				require.Equal(t, lua.LTBool, result.Type(), "Expected boolean return type")
				require.Equal(t, lua.LTNil, errValue.Type(), "Expected nil error")
				require.Equal(t, tt.expected, bool(result.(lua.LBool)))
			}
		})
	}
}

func TestSeqFunc(t *testing.T) {
	tests := []struct {
		input    []int
		expected string
	}{
		{
			input:    []int{1, 5},
			expected: "1 2 3 4 5",
		},
		{
			input:    []int{5, 1},
			expected: "5 4 3 2 1",
		},
		{
			input:    []int{1, 1},
			expected: "1",
		},
		{
			input:    []int{-3, 3},
			expected: "-3 -2 -1 0 1 2 3",
		},
		{
			input:    []int{3, -3},
			expected: "3 2 1 0 -1 -2 -3",
		},
		{
			input:    []int{5},
			expected: "1 2 3 4 5",
		},
		{
			input:    []int{},
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.NewTable()
			for _, num := range tt.input {
				tbl.Append(lua.LNumber(num))
			}

			L.Push(tbl)

			gluasprig.SeqFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSha1sumFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "hello",
		},
		{
			input: "test",
		},
		{
			input: "sha1sum",
		},
		{
			input: "",
		},
		{
			input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			input: "123456789",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			h := sha1.New()
			h.Write([]byte(tt.input))
			expected := hex.EncodeToString(h.Sum(nil))

			L.Push(lua.LString(tt.input))

			gluasprig.Sha1sumFunc(L)

			result := L.ToString(-1)
			require.Equal(t, expected, result, "SHA1 hash of %q should be %q, got %q", tt.input, expected, result)
		})
	}
}

func TestSha256sumFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "hello",
		},
		{
			input: "test",
		},
		{
			input: "sha256sum",
		},
		{
			input: "",
		},
		{
			input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			input: "123456789",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			h := sha256.New()
			h.Write([]byte(tt.input))
			expected := hex.EncodeToString(h.Sum(nil))

			L.Push(lua.LString(tt.input))

			gluasprig.Sha256sumFunc(L)

			result := L.ToString(-1)
			require.Equal(t, expected, result, "SHA256 hash of %q should be %q, got %q", tt.input, expected, result)
		})
	}
}

func TestSha512sumFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "hello",
		},
		{
			input: "test",
		},
		{
			input: "sha512sum",
		},
		{
			input: "",
		},
		{
			input: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			input: "123456789",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			h := sha512.New()
			h.Write([]byte(tt.input))
			expected := hex.EncodeToString(h.Sum(nil))

			L.Push(lua.LString(tt.input))

			gluasprig.Sha512sumFunc(L)

			result := L.ToString(-1)
			require.Equal(t, expected, result, "SHA512 hash of %q should be %q, got %q", tt.input, expected, result)
		})
	}
}

func TestShuffleFunc(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "abcdefg",
		},
		{
			input: "hello world",
		},
		{
			input: "123456789",
		},
		{
			input: "a",
		},
		{
			input: "",
		},
		{
			input: "aaa",
		},
	}

	sortString := func(s string) string {
		r := []rune(s)
		slices.Sort(r)

		return string(r)
	}

	canBeShuffled := func(s string) bool {
		// cannot shuffle empty string or single character
		if len(s) <= 1 {
			return false
		}

		// cannot meaningfully shuffle string with all identical characters
		first := s[0]

		for i := 1; i < len(s); i++ {
			if s[i] != first {
				return true
			}
		}

		return false
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.ShuffleFunc(L)

			result := L.ToString(-1)
			require.Equal(t, len(tt.input), len(result), "Shuffled result should have same length as input")

			if canBeShuffled(tt.input) {
				require.NotEqual(t, tt.input, result, "Shuffled result should not be equal to input")
			}

			sortedInput := sortString(tt.input)
			sortedResult := sortString(result)
			require.Equal(t, sortedInput, sortedResult, "Shuffled result should contain the same characters as input")
		})
	}
}

func TestSnakecaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "HelloWorld",
			expected: "hello_world",
		},
		{
			input:    "helloWorld",
			expected: "hello_world",
		},
		{
			input:    "hello world",
			expected: "hello_world",
		},
		{
			input:    "HELLO WORLD",
			expected: "hello_world",
		},
		{
			input:    "hello_world",
			expected: "hello_world",
		},
		{
			input:    "hello-world",
			expected: "hello_world",
		},
		{
			input:    "hello123world",
			expected: "hello_123world",
		},
		{
			input:    "Hello-World_Example",
			expected: "hello_world_example",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "HTTPRequest",
			expected: "http_request",
		},
		{
			input:    "XMLHttpRequest",
			expected: "xml_http_request",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.SnakecaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSortAlphaFunc(t *testing.T) {
	tests := []struct {
		input    any
		expected []string
	}{
		{
			input:    []string{"banana", "apple", "cherry", "date"},
			expected: []string{"apple", "banana", "cherry", "date"},
		},
		{
			input:    []any{5, "banana", "apple", 10, "cherry"},
			expected: []string{"10", "5", "apple", "banana", "cherry"},
		},
		{
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			input:    []string{},
			expected: []string{},
		},
		{
			input:    []int{5, 3, 1, 4, 2},
			expected: []string{"1", "2", "3", "4", "5"},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			var luaValue lua.LValue

			switch v := tt.input.(type) {
			case string:
				luaValue = lua.LString(v)
			case []string:
				tbl := L.CreateTable(len(v), 0)
				for i, str := range v {
					tbl.RawSetInt(i+1, lua.LString(str))
				}

				luaValue = tbl
			case []any:
				tbl := L.CreateTable(len(v), 0)
				for i, val := range v {
					switch val := val.(type) {
					case string:
						tbl.RawSetInt(i+1, lua.LString(val))
					case int:
						tbl.RawSetInt(i+1, lua.LNumber(val))
					default:
						tbl.RawSetInt(i+1, lua.LString(fmt.Sprint(val)))
					}
				}

				luaValue = tbl
			case []int:
				tbl := L.CreateTable(len(v), 0)
				for i, num := range v {
					tbl.RawSetInt(i+1, lua.LNumber(num))
				}

				luaValue = tbl
			default:
				luaValue = lua.LString(fmt.Sprint(v))
			}

			L.Push(luaValue)

			gluasprig.SortAlphaFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTTable, result.Type(), "Expected table result")

			tbl := result.(*lua.LTable)
			results := make([]string, 0, tbl.Len())

			tbl.ForEach(func(idx lua.LValue, value lua.LValue) {
				require.Equal(t, lua.LTNumber, idx.Type(), "Expected numeric index")
				require.Equal(t, lua.LTString, value.Type(), "Expected string value")

				results = append(results, value.String())
			})

			require.Equal(t, len(tt.expected), len(results), "Table length doesn't match expected length")

			if len(tt.expected) == 0 {
				require.Empty(t, results, "Expected empty results")
			} else {
				require.Equal(t, tt.expected, results, "Table results don't match expected values")
			}

			require.Equal(t, len(tt.expected), tbl.Len(), "Table length doesn't match expected length")
		})
	}
}

func TestSquoteFunc(t *testing.T) {
	tests := []struct {
		input    []any
		expected string
	}{
		{
			input:    []any{"hello"},
			expected: `'hello'`,
		},
		{
			input:    []any{"hello", "world"},
			expected: `'hello' 'world'`,
		},
		{
			input:    []any{"hello", 123, true},
			expected: `'hello' '123' 'true'`,
		},
		{
			input:    []any{"string", nil, "after nil"},
			expected: `'string' 'after nil'`,
		},
		{
			input:    []any{},
			expected: ``,
		},
		{
			input:    []any{42},
			expected: `'42'`,
		},
		{
			input:    []any{3.14159},
			expected: `'3.14159'`,
		},
		{
			input:    []any{-100},
			expected: `'-100'`,
		},
		{
			input:    []any{-2.718},
			expected: `'-2.718'`,
		},
		{
			input:    []any{42.0},
			expected: `'42'`,
		},
		{
			input:    []any{10, 20.5, 30},
			expected: `'10' '20.5' '30'`,
		},
		{
			input:    []any{1234567890},
			expected: `'1234567890'`,
		},
		{
			input:    []any{0.00001},
			expected: `'1e-05'`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.NewTable()
			for _, item := range tt.input {
				var luaValue lua.LValue
				switch v := item.(type) {
				case nil:
					luaValue = lua.LNil
				case string:
					luaValue = lua.LString(v)
				case int:
					luaValue = lua.LNumber(v)
				case int64:
					luaValue = lua.LNumber(v)
				case float64:
					luaValue = lua.LNumber(v)
				case bool:
					luaValue = lua.LBool(v)
				default:
					luaValue = lua.LString(fmt.Sprint(v))
				}
				tbl.Append(luaValue)
			}

			L.Push(tbl)

			gluasprig.SquoteFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSubstrFunc(t *testing.T) {
	tests := []struct {
		start    int
		end      int
		input    string
		expected string
	}{
		{
			start:    0,
			end:      5,
			input:    "hello world",
			expected: "hello",
		},
		{
			start:    6,
			end:      11,
			input:    "hello world",
			expected: "world",
		},
		{
			start:    0,
			end:      11,
			input:    "hello world",
			expected: "hello world",
		},
		{
			start:    -5,
			end:      11,
			input:    "hello world",
			expected: "world",
		},
		{
			start:    0,
			end:      20,
			input:    "hello world",
			expected: "hello world",
		},
		{
			start:    7,
			end:      8,
			input:    "hello world",
			expected: "o",
		},
		{
			start:    0,
			end:      0,
			input:    "hello world",
			expected: "",
		},
		{
			start:    100,
			end:      105,
			input:    "hello world",
			expected: "",
		},
		{
			start:    0,
			end:      5,
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.start))
			L.Push(lua.LNumber(tt.end))
			L.Push(lua.LString(tt.input))

			gluasprig.SubstrFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result, "Substring of %q from %d to %d", tt.input, tt.start, tt.end)
		})
	}
}

func TestSwapcaseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "HELLO WORLD",
		},
		{
			input:    "HELLO WORLD",
			expected: "hello world",
		},
		{
			input:    "Hello World",
			expected: "hELLO wORLD",
		},
		{
			input:    "HeLlO wOrLd",
			expected: "hElLo WoRlD",
		},
		{
			input:    "hello 123 WORLD",
			expected: "HELLO 123 world",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "123456789",
			expected: "123456789",
		},
		{
			input:    "Hello, World! 123",
			expected: "hELLO, wORLD! 123",
		},
		{
			input:    "lowercase",
			expected: "LOWERCASE",
		},
		{
			input:    "UPPERCASE",
			expected: "uppercase",
		},
		{
			input:    "Günther Über",
			expected: "gÜNTHER üBER",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.SwapcaseFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestTernaryFunc(t *testing.T) {
	tests := []struct {
		trueValue      lua.LValue
		falseValue     lua.LValue
		condition      bool
		expectedType   lua.LValueType
		expectedNumber float64
		expectedString string
		expectedBool   bool
	}{
		{
			trueValue:      lua.LString("true value"),
			falseValue:     lua.LString("false value"),
			condition:      true,
			expectedType:   lua.LTString,
			expectedString: "true value",
		},
		{
			trueValue:      lua.LString("true value"),
			falseValue:     lua.LString("false value"),
			condition:      false,
			expectedType:   lua.LTString,
			expectedString: "false value",
		},
		{
			trueValue:      lua.LNumber(42),
			falseValue:     lua.LNumber(17),
			condition:      true,
			expectedType:   lua.LTNumber,
			expectedNumber: 42,
		},
		{
			trueValue:      lua.LNumber(42),
			falseValue:     lua.LNumber(17),
			condition:      false,
			expectedType:   lua.LTNumber,
			expectedNumber: 17,
		},
		{
			trueValue:    lua.LBool(true),
			falseValue:   lua.LBool(false),
			condition:    true,
			expectedType: lua.LTBool,
			expectedBool: true,
		},
		{
			trueValue:    lua.LBool(true),
			falseValue:   lua.LBool(false),
			condition:    false,
			expectedType: lua.LTBool,
			expectedBool: false,
		},
		{
			trueValue:    lua.LNil,
			falseValue:   lua.LString("not nil"),
			condition:    true,
			expectedType: lua.LTNil,
		},
		{
			trueValue:    lua.LString("not nil"),
			falseValue:   lua.LNil,
			condition:    false,
			expectedType: lua.LTNil,
		},
		{
			trueValue:      lua.LString("string value"),
			falseValue:     lua.LNumber(42),
			condition:      true,
			expectedType:   lua.LTString,
			expectedString: "string value",
		},
		{
			trueValue:      lua.LString("string value"),
			falseValue:     lua.LNumber(42),
			condition:      false,
			expectedType:   lua.LTNumber,
			expectedNumber: 42,
		},
		{
			trueValue: func() lua.LValue {
				L := lua.NewState()
				t := L.NewTable()

				t.RawSetString("key", lua.LString("value"))

				return t
			}(),
			falseValue:   lua.LString("not table"),
			condition:    true,
			expectedType: lua.LTTable,
		},
		{
			trueValue: lua.LString("not table"),
			falseValue: func() lua.LValue {
				L := lua.NewState()
				t := L.NewTable()

				t.RawSetString("key", lua.LString("value"))

				return t
			}(),
			condition:    false,
			expectedType: lua.LTTable,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(tt.trueValue)
			L.Push(tt.falseValue)
			L.Push(lua.LBool(tt.condition))

			gluasprig.TernaryFunc(L)

			result := L.Get(-1)
			require.Equal(t, tt.expectedType, result.Type(), "Result type does not match expected")

			switch tt.expectedType {
			case lua.LTString:
				require.Equal(t, tt.expectedString, result.String())
			case lua.LTNumber:
				require.Equal(t, tt.expectedNumber, float64(result.(lua.LNumber)))
			case lua.LTBool:
				require.Equal(t, tt.expectedBool, bool(result.(lua.LBool)))
			case lua.LTTable:
				require.IsType(t, &lua.LTable{}, result)
			case lua.LTNil:
				require.Equal(t, lua.LNil, result)
			}
		})
	}
}

func TestToDecimalFunc(t *testing.T) {
	tests := []struct {
		octal   any
		decimal int64
	}{
		{
			octal:   "10",
			decimal: 8,
		},
		{
			octal:   "177",
			decimal: 127,
		},
		{
			octal:   "0",
			decimal: 0,
		},
		{
			octal:   "7",
			decimal: 7,
		},
		{
			octal:   "8",
			decimal: 0,
		},
		{
			octal:   "hello",
			decimal: 0,
		},
		{
			octal:   nil,
			decimal: 0,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			var luaValue lua.LValue

			switch v := tt.octal.(type) {
			case nil:
				luaValue = lua.LNil
			case string:
				luaValue = lua.LString(v)
			case int:
				luaValue = lua.LNumber(v)
			case float64:
				luaValue = lua.LNumber(v)
			case bool:
				luaValue = lua.LBool(v)
			default:
				luaValue = lua.LString(fmt.Sprint(v))
			}

			L.Push(luaValue)

			gluasprig.ToDecimalFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTNumber, result.Type(), "Expected number result")

			value := int64(result.(lua.LNumber))
			require.Equal(t, tt.decimal, value, "Result doesn't match expected value")
		})
	}
}

func TestTruncFunc(t *testing.T) {
	tests := []struct {
		length   int
		input    string
		expected string
	}{
		{
			length:   5,
			input:    "hello world",
			expected: "hello",
		},
		{
			length:   20,
			input:    "hello world",
			expected: "hello world",
		},
		{
			length:   0,
			input:    "hello world",
			expected: "",
		},
		{
			length:   1,
			input:    "hello world",
			expected: "h",
		},
		{
			length:   11,
			input:    "hello world",
			expected: "hello world",
		},
		{
			length:   10,
			input:    "hello world",
			expected: "hello worl",
		},
		{
			length:   -5,
			input:    "hello world",
			expected: "world",
		},
		{
			length:   5,
			input:    "",
			expected: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.length))
			L.Push(lua.LString(tt.input))

			gluasprig.TruncFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result, "Truncating %q to length %d", tt.input, tt.length)
		})
	}
}

func TestUniqFunc(t *testing.T) {
	tests := []struct {
		input    []lua.LValue
		expected []lua.LValue
	}{
		{
			input:    []lua.LValue{lua.LString("a"), lua.LString("b"), lua.LString("a"), lua.LString("c")},
			expected: []lua.LValue{lua.LString("a"), lua.LString("b"), lua.LString("c")},
		},
		{
			input:    []lua.LValue{lua.LNumber(1), lua.LNumber(3), lua.LNumber(1), lua.LNumber(2)},
			expected: []lua.LValue{lua.LNumber(1), lua.LNumber(3), lua.LNumber(2)},
		},
		{
			input:    []lua.LValue{lua.LBool(true), lua.LBool(false), lua.LBool(true)},
			expected: []lua.LValue{lua.LBool(true), lua.LBool(false)},
		},
		{
			input:    []lua.LValue{lua.LNil, lua.LNumber(1), lua.LNil},
			expected: []lua.LValue{lua.LNumber(1)},
		},
		{
			input:    []lua.LValue{},
			expected: []lua.LValue{},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			inputTable := L.CreateTable(len(tt.input), 0)
			for j, v := range tt.input {
				inputTable.RawSetInt(j+1, v)
			}

			L.Push(inputTable)
			gluasprig.UniqFunc(L)

			result := L.CheckTable(-1)

			expectedTable := L.CreateTable(len(tt.expected), 0)
			for j, v := range tt.expected {
				expectedTable.RawSetInt(j+1, v)
			}
			require.Equal(t, len(tt.expected), result.Len(), "Table length mismatch")

			for j, expectedVal := range tt.expected {
				actualVal := result.RawGetInt(j + 1)
				require.Equal(t, expectedVal, actualVal, "Value mismatch at index %d", j+1)
			}
		})
	}
}

func TestUntitleFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello World",
			expected: "hello world",
		},
		{
			input:    "HELLO WORLD",
			expected: "hELLO wORLD",
		},
		{
			input:    "hello world",
			expected: "hello world",
		},
		{
			input:    "Hello",
			expected: "hello",
		},
		{
			input:    "Title With Multiple Words",
			expected: "title with multiple words",
		},
		{
			input:    "Title-With-Hyphens",
			expected: "title-With-Hyphens",
		},
		{
			input:    "Title_With_Underscores",
			expected: "title_With_Underscores",
		},
		{
			input:    "Title With 123 Numbers",
			expected: "title with 123 numbers",
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "123",
			expected: "123",
		},
		{
			input:    "Günther Über",
			expected: "günther über",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.UntitleFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUrlJoinFunc(t *testing.T) {
	tests := []struct {
		input    map[string]string
		expected string
	}{
		{
			input: map[string]string{
				"scheme": "https",
				"host":   "example.com",
			},
			expected: "https://example.com",
		},
		{
			input: map[string]string{
				"scheme": "https",
				"host":   "example.com",
				"path":   "/api/v1",
			},
			expected: "https://example.com/api/v1",
		},
		{
			input: map[string]string{
				"scheme": "https",
				"host":   "example.com",
				"path":   "/api/v1",
				"query":  "page=1&q=search",
			},
			expected: "https://example.com/api/v1?page=1&q=search",
		},
		{
			input: map[string]string{
				"scheme":   "https",
				"host":     "example.com:8443",
				"path":     "/secure",
				"fragment": "section-1",
			},
			expected: "https://example.com:8443/secure#section-1",
		},
		{
			input: map[string]string{
				"scheme":   "http",
				"userinfo": "user:pass",
				"host":     "example.com",
			},
			expected: "http://user:pass@example.com",
		},
		{
			input: map[string]string{
				"scheme":   "http",
				"userinfo": "user",
				"host":     "example.com",
			},
			expected: "http://user@example.com",
		},
		{
			input:    map[string]string{},
			expected: "",
		},
		{
			input: map[string]string{
				"host": "example.com",
			},
			expected: "//example.com",
		},
		{
			input: map[string]string{
				"path": "/just/path",
			},
			expected: "/just/path",
		},
		{
			input: map[string]string{
				"scheme": "https",
				"opaque": "example.com:443",
				"path":   "/ignored",
			},
			expected: "https:example.com:443",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			tbl := L.NewTable()
			for k, v := range tt.input {
				tbl.RawSetString(k, lua.LString(v))
			}

			L.Push(tbl)

			gluasprig.UrlJoinFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestUrlParseFunc(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input: "https://example.com",
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "https://example.com/api/v1",
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "/api/v1",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "https://example.com/api/v1?page=1&q=search",
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "/api/v1",
				"query":    "page=1&q=search",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "https://example.com:8443/secure#section-1",
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com:8443",
				"hostname": "example.com",
				"path":     "/secure",
				"query":    "",
				"opaque":   "",
				"fragment": "section-1",
				"userinfo": "",
			},
		},
		{
			input: "http://user:pass@example.com",
			expected: map[string]string{
				"scheme":   "http",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "user:pass",
			},
		},
		{
			input: "http://user@example.com",
			expected: map[string]string{
				"scheme":   "http",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "user",
			},
		},
		{
			input: "//example.com",
			expected: map[string]string{
				"scheme":   "",
				"host":     "example.com",
				"hostname": "example.com",
				"path":     "",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "/just/path",
			expected: map[string]string{
				"scheme":   "",
				"host":     "",
				"hostname": "",
				"path":     "/just/path",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "mailto:user@example.com",
			expected: map[string]string{
				"scheme":   "mailto",
				"host":     "",
				"hostname": "",
				"path":     "",
				"query":    "",
				"opaque":   "user@example.com",
				"fragment": "",
				"userinfo": "",
			},
		},
		{
			input: "https://user:pass@example.com:8443/path?query=value#fragment",
			expected: map[string]string{
				"scheme":   "https",
				"host":     "example.com:8443",
				"hostname": "example.com",
				"path":     "/path",
				"query":    "query=value",
				"opaque":   "",
				"fragment": "fragment",
				"userinfo": "user:pass",
			},
		},
		{
			input: "",
			expected: map[string]string{
				"scheme":   "",
				"host":     "",
				"hostname": "",
				"path":     "",
				"query":    "",
				"opaque":   "",
				"fragment": "",
				"userinfo": "",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LString(tt.input))

			gluasprig.UrlParseFunc(L)

			result := L.Get(-1)
			require.Equal(t, lua.LTTable, result.Type(), "Expected table return type")

			tbl := result.(*lua.LTable)
			resultMap := make(map[string]string)

			tbl.ForEach(func(k, v lua.LValue) {
				if k.Type() == lua.LTString {
					key := string(k.(lua.LString))

					resultMap[key] = v.String()
				}
			})

			for key, expectedValue := range tt.expected {
				value, exists := resultMap[key]
				require.True(t, exists, "Expected key %q is missing", key)
				require.Equal(t, expectedValue, value, "Value mismatch for key %q", key)
			}
		})
	}
}

func TestWrapFunc(t *testing.T) {
	tests := []struct {
		width    int
		input    string
		expected string
	}{
		{
			width:    5,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    20,
			input:    "hello world",
			expected: "hello world",
		},
		{
			width:    7,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    2,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    1,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    0,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    -1,
			input:    "hello world",
			expected: "hello\nworld",
		},
		{
			width:    10,
			input:    "This is a longer sentence that needs to be wrapped at multiple points.",
			expected: "This is a\nlonger\nsentence\nthat needs\nto be\nwrapped at\nmultiple\npoints.",
		},
		{
			width:    5,
			input:    "",
			expected: "",
		},
		{
			width:    5,
			input:    "supercalifragilisticexpialidocious",
			expected: "supercalifragilisticexpialidocious",
		},
		{
			width:    5,
			input:    "hello\nworld",
			expected: "hello\nworld",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.width))
			L.Push(lua.LString(tt.input))

			gluasprig.WrapFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result, "Wrapping %q at width %d", tt.input, tt.width)
		})
	}
}

func TestWrapWithFunc(t *testing.T) {
	tests := []struct {
		width      int
		input      string
		wrapString string
		expected   string
	}{
		{
			width:      5,
			input:      "hello world",
			wrapString: "<br>",
			expected:   "hello<br>world",
		},
		{
			width:      20,
			input:      "hello world",
			wrapString: "<br>",
			expected:   "hello world",
		},
		{
			width:      7,
			input:      "hello world",
			wrapString: "\n\t",
			expected:   "hello\n\tworld",
		},
		{
			width:      2,
			input:      "hello world",
			wrapString: " ... ",
			expected:   "he ... ll ... o ... wo ... rl ... d",
		},
		{
			width:      1,
			input:      "hello world",
			wrapString: " ",
			expected:   "h e l l o w o r l d",
		},
		{
			width:      0,
			input:      "hello world",
			wrapString: "-",
			expected:   "h-e-l-l-o-w-o-r-l-d",
		},
		{
			width:      -1,
			input:      "hello world",
			wrapString: "|",
			expected:   "h|e|l|l|o|w|o|r|l|d",
		},
		{
			width:      10,
			input:      "This is a longer sentence that needs to be wrapped at multiple points.",
			wrapString: " ... ",
			expected:   "This is a ... longer ... sentence ... that needs ... to be ... wrapped at ... multiple ... points.",
		},
		{
			width:      5,
			input:      "",
			wrapString: "<br>",
			expected:   "",
		},
		{
			width:      5,
			input:      "supercalifragilisticexpialidocious",
			wrapString: "-",
			expected:   "super-calif-ragil-istic-expia-lidoc-ious",
		},
		{
			width:      5,
			input:      "hello\nworld",
			wrapString: "<br>",
			expected:   "hello<br>\nworl<br>d",
		},
		{
			width:      5,
			input:      "hello world",
			wrapString: "",
			expected:   "hello\nworld",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.Push(lua.LNumber(tt.width))
			L.Push(lua.LString(tt.input))
			L.Push(lua.LString(tt.wrapString))

			gluasprig.WrapWithFunc(L)

			result := L.ToString(-1)
			require.Equal(t, tt.expected, result, "Wrapping %q at width %d with separator %q",
				tt.input, tt.width, tt.wrapString)
		})
	}
}
