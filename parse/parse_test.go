package parse_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tehbilly/genny/parse"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

var tests = []struct {
	// input
	filename string
	pkgName  string
	in       string
	tag      string
	imports  []string
	types    []map[string]parse.TypeRef

	// expectations
	expectedOut string
	expectedErr error

	suppressForAstImpl    bool
	suppressForLegacyImpl bool
}{
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/queue/int_queue.go`,
	},
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{Alias: "CrazyNumber", Type: "int"}}},
		expectedOut: `test/queue/int_queue_aliased.go`,
	},
	{
		filename:    "generic_queue.go",
		pkgName:     "changed",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/queue/changed/int_queue.go`,
	},
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{Alias: "float32", Type: "float32"}}},
		expectedOut: `test/queue/float32_queue.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{Alias: "string", Type: "string"},
			"ValueType": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/multipletypes/string_int_simplemap.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{Alias: "interface{}", Type: "interface{}"},
			"ValueType": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/multipletypes/interface_int_simplemap.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{Alias: "*MyType1", Type: "*MyType1"},
			"ValueType": parse.TypeRef{Alias: "*MyOtherType", Type: "*MyOtherType"}}},
		expectedOut: `test/multipletypes/custom_types_simplemap.go`,
	},
	{
		filename:    "generic_internal.go",
		in:          `test/unexported/generic_internal.go`,
		types:       []map[string]parse.TypeRef{{"secret": parse.TypeRef{Alias: "*myType", Type: "*myType"}}},
		expectedOut: `test/unexported/mytype_internal.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypesets/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{
			{"KeyType": parse.TypeRef{Alias: "int", Type: "int"}, "ValueType": parse.TypeRef{Alias: "string", Type: "string"}},
			{"KeyType": parse.TypeRef{Alias: "float64", Type: "float64"}, "ValueType": parse.TypeRef{Alias: "bool", Type: "bool"}},
		},
		expectedOut: `test/multipletypesets/many_simplemaps.go`,
	},
	{
		filename:    "generic_number.go",
		in:          `test/numbers/generic_number.go`,
		types:       []map[string]parse.TypeRef{{"NumberType": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/numbers/int_number.go`,
	},
	{
		filename:    "generic_digraph.go",
		in:          `test/bugreports/generic_digraph.go`,
		types:       []map[string]parse.TypeRef{{"Node": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/bugreports/int_digraph.go`,
	},
	{
		filename:    "renamed_pkg.go",
		in:          `test/renamed/renamed_pkg.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/renamed/renamed_pkg_int.go`,
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/buildtags/buildtags_expected.go`,
		tag:         "genny",
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{Alias: "string", Type: "string"}}},
		expectedOut: `test/buildtags/buildtags_expected_nostrip.go`,
		tag:         "",
	},
	{
		filename: "buildtags.go",
		in:       `test/buildtags/buildtags.go`,
		types: []map[string]parse.TypeRef{
			{"_t_": parse.TypeRef{Alias: "string", Type: "string"}},
			{"_t_": parse.TypeRef{Alias: "int", Type: "int"}},
		},
		expectedOut: `test/buildtags/buildtags_expected_multiple.go`,
		tag:         "genny",
	},
	{
		filename:    "join.go",
		in:          `test/interfaces/join.go`,
		types:       []map[string]parse.TypeRef{{"Stringer": parse.TypeRef{Alias: "MyStr", Type: "MyStr"}}},
		expectedOut: `test/interfaces/join_expected.go`,
		tag:         "",
	},
	{
		filename: "syntax.go",
		in:       `test/syntax/syntax.go`,
		types: []map[string]parse.TypeRef{
			{"myType": parse.TypeRef{Alias: "timeSpan", Type: "time.Duration"}},
			{"myType": parse.TypeRef{Alias: "Fractional", Type: "float64"}},
		},
		expectedOut:           `test/syntax/syntax_expected.go`,
		tag:                   "",
		suppressForLegacyImpl: true,
	},
	{
		filename: "generic.go",
		in:       `test/interface-template/generic.go`,
		types: []map[string]parse.TypeRef{
			{"TypeParam": parse.TypeRef{Alias: "string", Type: "string"}},
		},
		expectedOut: `test/interface-template/generic_expected.go`,
		tag:         "",
	},
	{
		filename:    "generic_new_and_make_slice.go",
		in:          `test/bugreports/generic_new_and_make_slice.go`,
		types:       []map[string]parse.TypeRef{{"NumberType": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/bugreports/int_new_and_make_slice.go`,
	},
	{
		filename:    "cell_x.go",
		in:          `test/bugreports/cell_x.go`,
		types:       []map[string]parse.TypeRef{{"X": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/bugreports/cell_int.go`,
	},
	{
		filename:    "interface_generic_type.go",
		in:          `test/bugreports/interface_generic_type.go`,
		types:       []map[string]parse.TypeRef{{"GenericType": parse.TypeRef{Alias: "uint8", Type: "uint8"}}},
		expectedOut: `test/bugreports/interface_uint8.go`,
	},
	{
		filename:    "negation_generic.go",
		in:          `test/bugreports/negation_generic.go`,
		types:       []map[string]parse.TypeRef{{"SomeThing": parse.TypeRef{Alias: "string", Type: "string"}}},
		expectedOut: `test/bugreports/negation_string.go`,
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{Alias: "int", Type: "int"}}},
		expectedOut: `test/buildtags/buildtags_expected.go`,
		tag:         "genny",
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{Alias: "string", Type: "string"}}},
		expectedOut: `test/buildtags/buildtags_expected_nostrip.go`,
		tag:         "",
	},
}

func TestParse(t *testing.T) {
	for testNo, test := range tests {

		for _, useAst := range []bool{true, false} {
			if (useAst && test.suppressForAstImpl) || (!useAst && test.suppressForLegacyImpl) {
				continue
			}
			t.Run(fmt.Sprintf("%d:%s/(ast:%v)", testNo, test.expectedOut, useAst), func(t *testing.T) {
				in, err := contents(test.in)
				require.NoError(t, err)
				expectedOut, err := contents(test.expectedOut)
				require.NoError(t, err)

				bytes, err := parse.Generics(
					test.filename,
					test.pkgName,
					strings.NewReader(in),
					test.types,
					test.imports,
					test.tag,
					useAst)

				// check the error
				if test.expectedErr == nil {
					assert.NoError(t, err, "(%d: %s) No error was expected but got: %s", testNo, test.filename, err)
				} else {
					assert.NotNil(t, err, "(%d: %s) No error was returned by one was expected: %s", testNo, test.filename, test.expectedErr)
					assert.IsType(t, test.expectedErr, err, "(%d: %s) Generate should return object of type %v", testNo, test.filename, test.expectedErr)
				}

				// assert the response
				if !assert.Equal(t, expectedOut, string(bytes), "Parse didn't generate the expected output.") {
					log.Println("EXPECTED: " + expectedOut)
					log.Println("ACTUAL: " + string(bytes))
				}
			})
		}

	}

}

func contents(s string) (string, error) {
	if strings.HasSuffix(s, "go") || strings.HasSuffix(s, "go.nobuild") {
		bytes, err := ioutil.ReadFile(s)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return s, nil
}
