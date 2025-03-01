package main

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
	"unicode"

	sprig "github.com/Masterminds/sprig/v3"
)

// go run tools/*.go

type FunctionInfo struct {
	Name       string
	SafeName   string
	ReturnType string
	ParamTypes []string
	HasError   bool
}

func GetSafeFunctionName(name string) string {
	if name == "" {
		return ""
	}

	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_'
	})

	for i, part := range parts {
		if part == "" {
			continue
		}

		runes := []rune(part)

		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}

	return strings.Join(parts, "") + "Func"
}

func GetTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.ArrayType:
		return "[]" + GetTypeString(t.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", GetTypeString(t.Key), GetTypeString(t.Value))
	case *ast.InterfaceType:
		return "any"
	case *ast.StarExpr:
		return "*" + GetTypeString(t.X)
	case *ast.SelectorExpr:
		return GetTypeString(t.X) + "." + t.Sel.Name
	default:
		return "any"
	}
}

func NormalizeTypeName(typeName string) string {
	sr := strings.NewReplacer(
		"interface {}", "any",
		"interface{}", "any",
	)

	return sr.Replace(typeName)
}

func AnalyzeFunctionSignature(name string, x any) (FunctionInfo, error) {
	info := FunctionInfo{
		Name:     name,
		SafeName: GetSafeFunctionName(name),
	}

	fn := reflect.TypeOf(x)
	if fn.Kind() != reflect.Func {
		return FunctionInfo{}, fmt.Errorf("not a function: %s", name)
	}

	{ // get function parameter types
		numOfParams := fn.NumIn()

		info.ParamTypes = make([]string, numOfParams)

		for i := range numOfParams {
			paramType := fn.In(i)

			if fn.IsVariadic() && i == numOfParams-1 {
				elemType := paramType.Elem().String()
				info.ParamTypes[i] = "..." + NormalizeTypeName(elemType)
			} else {
				info.ParamTypes[i] = NormalizeTypeName(paramType.String())
			}
		}
	}

	{ // get function return types
		numOfResults := fn.NumOut()
		if numOfResults > 0 {
			info.ReturnType = NormalizeTypeName(fn.Out(0).String())
		}

		if numOfResults > 1 {
			if fn.Out(1).String() == "error" {
				info.HasError = true
			}
		}
	}

	return info, nil
}

func main() {
	sprigFuncs := sprig.HermeticTxtFuncMap()
	functions := make([]FunctionInfo, 0, len(sprigFuncs))

	for name, fn := range sprigFuncs {
		if reflect.ValueOf(fn).Kind() != reflect.Func {
			// skip non-function values in the map
			continue
		}

		info, err := AnalyzeFunctionSignature(name, fn)
		if err != nil {
			log.Fatalf("Analyze function %q signature error: %v", name, err)
		}

		functions = append(functions, info)
	}

	slices.SortFunc(functions, func(a, b FunctionInfo) int {
		return strings.Compare(a.SafeName, b.SafeName)
	})

	for _, fn := range functions {
		fmt.Fprintf(os.Stdout, "Function: %s\n", fn.Name)
		fmt.Fprintf(os.Stdout, "  SafeName: %s\n", fn.SafeName)
		fmt.Fprintf(os.Stdout, "  ParamTypes: %s\n", strings.Join(fn.ParamTypes, ","))
		fmt.Fprintf(os.Stdout, "  ReturnType: %s\n", fn.ReturnType)
		fmt.Fprintf(os.Stdout, "  HasError: %t\n\n", fn.HasError)
	}
}
