package test

import (
	"strings"
	"testing"

	"github.com/shivakuppa/Go_Redis/internals/resp"
	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	decodingTestCases := []struct {
		name        string
		input       string
		expectedErr string
		expectedVal *resp.Value
	}{
		{
			name:  "deserialize array",
			input: "*1\r\n:123\r\n",
			expectedVal: &resp.Value{
				Type:  resp.Array,
				Array: []*resp.Value{{Type: resp.Integer, Integer: 123}},
			},
		},
		{
			name:  "deserialize null array",
			input: "*-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Array,
				IsNull: true,
			},
		},
		{
			name:  "deserialize simple string",
			input: "+OK\r\n",
			expectedVal: &resp.Value{
				Type:   resp.SimpleString,
				String: "OK",
			},
		},
		{
			name:  "deserialize error",
			input: "-ERR invalid command\r\n",
			expectedVal: &resp.Value{
				Type:   resp.SimpleError,
				String: "ERR invalid command",
			},
		},
		{
			name:  "deserialize integer",
			input: ":123\r\n",
			expectedVal: &resp.Value{
				Type:    resp.Integer,
				Integer: 123,
			},
		},
		{
			name:  "deserialize bulk string",
			input: "$5\r\nhello\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				String: "hello",
			},
		},
		{
			name:  "deserialize null bulk string",
			input: "$-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				IsNull: true,
			},
		},
		{
			name:        "no first byte",
			input:       "",
			expectedErr: "read resp first byte",
		},
		{
			name:        "error read line bytes",
			input:       "*invalid",
			expectedErr: "read line bytes",
		},
		{
			name:        "error line not terminated",
			input:       "*invalid\n",
			expectedErr: "line not terminated with",
		},
		{
			name:        "error parse array num",
			input:       "*x\r\n",
			expectedErr: "parse the num elements",
		},

		{
			name:        "error decoding arr el",
			input:       "*1\r\n:x",
			expectedErr: "deserialize array element",
		},
		{
			name:        "error bulk string len",
			input:       "$x\r\nhello\r\n",
			expectedErr: "parse string len",
		},
		{
			name:        "error read bulk string",
			input:       "$5\r\n\r\n",
			expectedErr: "read bulk string",
		},
		{
			name:        "error bulk string not terminated",
			input:       "$5\r\nhell\r\n",
			expectedErr: "bulk string not terminated correctly",
		},
	}

	for _, tc := range decodingTestCases {
		val, err := resp.Deserialize(strings.NewReader(tc.input))

		if tc.expectedErr != "" {
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedErr)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVal, val)
		}
	}
}

func TestSerialize(t *testing.T) {
	serializeTestCases := []struct {
		name               string
		value              *resp.Value
		expectedSerialized string
		expectedErr        string
	}{
		{
			name: "simple string value",
			value: &resp.Value{
				Type:   resp.SimpleString,
				String: "OK",
			},
			expectedSerialized: "+OK\r\n",
		},
		{
			name: "integer value",
			value: &resp.Value{
				Type:    resp.Integer,
				Integer: 123,
			},
			expectedSerialized: ":123\r\n",
		},
		{
			name: "bulk string value",
			value: &resp.Value{
				Type:   resp.BulkString,
				String: "hello",
			},
			expectedSerialized: "$5\r\nhello\r\n",
		},
		{
			name: "array value",
			value: &resp.Value{
				Type: resp.Array,
				Array: []*resp.Value{
					{Type: resp.Integer, Integer: 123},
					{Type: resp.SimpleString, String: "OK"},
				},
			},
			expectedSerialized: "*2\r\n:123\r\n+OK\r\n",
		},
		{
			name: "error value",
			value: &resp.Value{
				Type:   resp.SimpleError,
				String: "ERR invalid command",
			},
			expectedSerialized: "-ERR invalid command\r\n",
		},
		{
			name: "nil bulk string value",
			value: &resp.Value{
				Type:   resp.BulkString,
				IsNull: true,
			},
			expectedSerialized: "$-1\r\n",
		},
		{
			name: "nil value array element",
			value: &resp.Value{
				Type: resp.Array,
				Array: []*resp.Value{
					{Type: resp.BulkString, IsNull: true},
				},
			},
			expectedSerialized: "*1\r\n$-1\r\n",
		},
		{
			name: "invalid type value",
			value: &resp.Value{
				Type: 99,
			},
			expectedErr: "invalid resp type",
		},
		{
			name: "invalid type array element",
			value: &resp.Value{
				Type: resp.Array,
				Array: []*resp.Value{
					{Type: 99},
				},
			},
			expectedErr: "invalid resp type",
		},
	}

	for _, tc := range serializeTestCases {
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