package test

import (
	"strings"
	"testing"

	"github.com/shivakuppa/Go_Redis/internals/resp"
	"github.com/stretchr/testify/assert"
)

func TestDeserializeRESP3(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr string
		expectedVal *resp.Value
	}{
		// ===== RESP2 Compatible =====
		{
			name:  "simple string",
			input: "+OK\r\n",
			expectedVal: &resp.Value{
				Type:   resp.SimpleString,
				String: "OK",
			},
		},
		{
			name:  "simple error",
			input: "-ERR invalid command\r\n",
			expectedVal: &resp.Value{
				Type:   resp.SimpleError,
				String: "ERR invalid command",
			},
		},
		{
			name:  "integer",
			input: ":123\r\n",
			expectedVal: &resp.Value{
				Type:    resp.Integer,
				Integer: 123,
			},
		},
		{
			name:  "bulk string",
			input: "$5\r\nhello\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				String: "hello",
			},
		},
		{
			name:  "null bulk string",
			input: "$-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				IsNull: true,
			},
		},
		{
			name:  "array",
			input: "*2\r\n:1\r\n+two\r\n",
			expectedVal: &resp.Value{
				Type: resp.Array,
				Array: []*resp.Value{
					{Type: resp.Integer, Integer: 1},
					{Type: resp.SimpleString, String: "two"},
				},
			},
		},
		{
			name:  "null array",
			input: "*-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Array,
				IsNull: true,
			},
		},

		// ===== RESP3 Specific =====
		{
			name:  "boolean true",
			input: "#t\r\n",
			expectedVal: &resp.Value{
				Type: resp.Boolean,
				Bool: true,
			},
		},
		{
			name:  "boolean false",
			input: "#f\r\n",
			expectedVal: &resp.Value{
				Type: resp.Boolean,
				Bool: false,
			},
		},
		{
			name:  "double",
			input: ",3.14159\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Double,
				Double: 3.14159,
			},
		},
		{
			name:  "bignumber",
			input: "(3492890328409238509324850943850943825024385\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BigNumber,
				String: "3492890328409238509324850943850943825024385",
			},
		},
		{
			name:  "null",
			input: "_\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Null,
				IsNull: true,
			},
		},
		{
			name:  "bulk error",
			input: "!21\r\nSYNTAX invalid input\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkError,
				String: "SYNTAX invalid input",
			},
		},
		{
			name:  "verbatim string",
			input: "=15\r\ntxt:Some text\r\n",
			expectedVal: &resp.Value{
				Type:   resp.VerbatimString,
				String: "txt:Some text",
			},
		},
		{
			name:  "set",
			input: "~2\r\n+apple\r\n+banana\r\n",
			expectedVal: &resp.Value{
				Type: resp.Set,
				Set: map[*resp.Value]struct{}{
					{Type: resp.SimpleString, String: "apple"}: {},
					{Type: resp.SimpleString, String: "banana"}: {},
				},
			},
		},
		{
			name:  "map",
			input: "%2\r\n+key1\r\n+value1\r\n+key2\r\n+value2\r\n",
			expectedVal: &resp.Value{
				Type: resp.Map,
				Map: map[string]*resp.Value{
					"key1": {Type: resp.SimpleString, String: "value1"},
					"key2": {Type: resp.SimpleString, String: "value2"},
				},
			},
		},
		{
			name:  "push message",
			input: ">2\r\n+pubsub\r\n+message\r\n",
			expectedVal: &resp.Value{
				Type: resp.Push,
				Array: []*resp.Value{
					{Type: resp.SimpleString, String: "pubsub"},
					{Type: resp.SimpleString, String: "message"},
				},
			},
		},

		// ===== Error Cases =====
		{
			name:        "missing first byte",
			input:       "",
			expectedErr: "read resp first byte",
		},
		{
			name:        "invalid first byte",
			input:       "?unknown\r\n",
			expectedErr: "invalid resp type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := resp.Deserialize(strings.NewReader(tc.input))

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedVal, val)
			}
		})
	}
}

func TestSerializeRESP3(t *testing.T) {
	testCases := []struct {
		name               string
		value              *resp.Value
		expectedSerialized string
		expectedErr        string
	}{
		// ===== RESP2 Compatible =====
		{
			name: "simple string",
			value: &resp.Value{
				Type:   resp.SimpleString,
				String: "OK",
			},
			expectedSerialized: "+OK\r\n",
		},
		{
			name: "integer",
			value: &resp.Value{
				Type:    resp.Integer,
				Integer: 123,
			},
			expectedSerialized: ":123\r\n",
		},
		{
			name: "array",
			value: &resp.Value{
				Type: resp.Array,
				Array: []*resp.Value{
					{Type: resp.Integer, Integer: 1},
					{Type: resp.SimpleString, String: "two"},
				},
			},
			expectedSerialized: "*2\r\n:1\r\n+two\r\n",
		},

		// ===== RESP3 Types =====
		{
			name: "boolean true",
			value: &resp.Value{
				Type: resp.Boolean,
				Bool: true,
			},
			expectedSerialized: "#t\r\n",
		},
		{
			name: "double",
			value: &resp.Value{
				Type:   resp.Double,
				Double: 3.14,
			},
			expectedSerialized: ",3.14\r\n",
		},
		{
			name: "big number",
			value: &resp.Value{
				Type:   resp.BigNumber,
				String: "12345678901234567890",
			},
			expectedSerialized: "(12345678901234567890\r\n",
		},
		{
			name: "null",
			value: &resp.Value{
				Type:   resp.Null,
				IsNull: true,
			},
			expectedSerialized: "_\r\n",
		},
		{
			name: "bulk error",
			value: &resp.Value{
				Type:   resp.BulkError,
				String: "ERR invalid",
			},
			expectedSerialized: "!11\r\nERR invalid\r\n",
		},
		{
			name: "verbatim string",
			value: &resp.Value{
				Type:   resp.VerbatimString,
				String: "txt:Some text",
			},
			expectedSerialized: "=15\r\ntxt:Some text\r\n",
		},
		{
			name: "map",
			value: &resp.Value{
				Type: resp.Map,
				Map: map[string]*resp.Value{
					"key1": {Type: resp.SimpleString, String: "value1"},
					"key2": {Type: resp.SimpleString, String: "value2"},
				},
			},
			expectedSerialized: "%2\r\n+key1\r\n+value1\r\n+key2\r\n+value2\r\n",
		},
		{
			name: "set",
			value: &resp.Value{
				Type: resp.Set,
				Set: map[*resp.Value]struct{}{
					{Type: resp.SimpleString, String: "apple"}: {},
					{Type: resp.SimpleString, String: "banana"}: {},
				},
			},
			// Order may vary; you might use custom set comparison in Deserialize
			expectedSerialized: "~2\r\n+apple\r\n+banana\r\n",
		},
		{
			name: "push",
			value: &resp.Value{
				Type: resp.Push,
				Array: []*resp.Value{
					{Type: resp.SimpleString, String: "pubsub"},
					{Type: resp.SimpleString, String: "message"},
				},
			},
			expectedSerialized: ">2\r\n+pubsub\r\n+message\r\n",
		},

		// ===== Error Case =====
		{
			name: "invalid type",
			value: &resp.Value{
				Type: 99,
			},
			expectedErr: "invalid resp type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := resp.Serialize(tc.value)

			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedSerialized, encoded)
			}
		})
	}
}
