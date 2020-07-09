package jsonnet

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJsonnet(t *testing.T) {
	testCases := []struct {
		jsonInput       string
		expectedJsonnet string
	}{
		{
			jsonInput:       "{}",
			expectedJsonnet: "{}",
		},
		{
			jsonInput: `{
				"foo": "bar",
				"baz": 123
			}`,
			expectedJsonnet: strings.Join([]string{
				`{`,
				`  baz: 123,`,
				`  foo: "bar"`,
				`}`,
			}, "\n"),
		},
		{
			jsonInput: `{
				"foo": "bar",
				"if": "keyword",
				"key-with-dashes": 123,
				"key with spaces": 456
			}`,
			expectedJsonnet: strings.Join([]string{
				`{`,
				`  foo: "bar",`,
				`  "if": "keyword",`,
				`  "key with spaces": 456,`,
				`  "key-with-dashes": 123`,
				`}`,
			}, "\n"),
		},
	}

	for _, testCase := range testCases {
		v := make(map[string]interface{})
		json.Unmarshal([]byte(testCase.jsonInput), &v)
		out, err := MarshalJsonnet(v)
		assert.Equal(t, testCase.expectedJsonnet, string(out))
		assert.NoError(t, err)
	}
}
