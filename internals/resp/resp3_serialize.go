package resp

import (
	"fmt"
	"math"
	"strconv"
	"strings"

)

func serializeNull() (string, error) {
	return "_\r\n", nil
}

func serializeBoolean(v *Value) (string, error) {
	if v.Bool {
		return "#t\r\n", nil
	}
	return "#f\r\n", nil
}

func serializeDouble(v *Value) (string, error) {
	switch {
	case math.IsInf(v.Double, 1):
		return ",inf\r\n", nil
	case math.IsInf(v.Double, -1):
		return ",-inf\r\n", nil
	case math.IsNaN(v.Double):
		return ",nan\r\n", nil
	default:
		return fmt.Sprintf(",%v\r\n", v.Double), nil
	}
}

func serializeBigNumber(v *Value) (string, error) {
	data := v.String
	if data == "" {
		return "", fmt.Errorf("empty big number string")
	}

	start := 0
	if data[0] == '+' || data[0] == '-' {
		if len(data) == 1 {
			return "", fmt.Errorf("big number cannot be only a sign")
		}
		start = 1
	}
	
	for i := start ; i < len(data) ; i++ {
		if (data[i] < '0' || data[i] > '9') {
			return "", fmt.Errorf("invalid big number")
		}
	}

	return "(" + data + "\r\n", nil
}

func serializeBulkError(v *Value) (string, error) {
	if v.IsNull {
		return "!\r\n", nil
	}
	return "!" + strconv.Itoa(len(v.String)) + "\r\n" + v.String + "\r\n", nil
}

func serializeVerbatimString(v *Value) (string, error) {
	// RESP3 verbatim strings use =<length>\r\n<format>:<data>\r\n
	// format examples: "txt:", "mkd:", "html:"
	if v.IsNull {
		return "$-1\r\n", nil
	}

	formatPrefix := "txt:" // default format
	if !strings.Contains(v.String, ":") {
		v.String = formatPrefix + v.String
	}

	return "=" + strconv.Itoa(len(v.String)) + "\r\n" + v.String + "\r\n", nil
}

func serializeMap(v *Value) (string, error) {
	if v.IsNull {
		return "%-1\r\n", nil
	}

	var sb strings.Builder
	sb.WriteString("%")
	sb.WriteString(strconv.Itoa(len(v.Map)))
	sb.WriteString("\r\n")

	for key, val := range v.Map {
		// Keys are always bulk strings in RESP3
		keyVal := &Value{Type: BulkString, String: key}
		keyStr, err := serializeBulkString(keyVal)
		if err != nil {
			return "", fmt.Errorf("serialize map key: %w", err)
		}
		sb.WriteString(keyStr)

		valStr, err := Serialize(val)
		if err != nil {
			return "", fmt.Errorf("serialize map value: %w", err)
		}
		sb.WriteString(valStr)
	}

	return sb.String(), nil
}

func serializeSet(v *Value) (string, error) {
	if v.IsNull {
		return "~-1\r\n", nil
	}

	var sb strings.Builder
	sb.WriteString("~")
	sb.WriteString(strconv.Itoa(len(v.Set)))
	sb.WriteString("\r\n")

	for el := range v.Set {
		s, err := Serialize(el)
		if err != nil {
			return "", fmt.Errorf("serialize set element: %w", err)
		}
		sb.WriteString(s)
	}

	return sb.String(), nil
}

func serializePush(v *Value) (string, error) {
	if v.IsNull {
		return ">-1\r\n", nil
	}

	var sb strings.Builder
	sb.WriteString(">")
	sb.WriteString(strconv.Itoa(len(v.Array)))
	sb.WriteString("\r\n")

	for _, el := range v.Array {
		s, err := Serialize(el)
		if err != nil {
			return "", fmt.Errorf("serialize push element: %w", err)
		}
		sb.WriteString(s)
	}

	return sb.String(), nil
}

func serializeAttribute(v *Value) (string, error) {
	if v.IsNull {
		return "|-1\r\n", nil
	}

	var sb strings.Builder
	sb.WriteString("|")
	sb.WriteString(strconv.Itoa(len(v.Map)))
	sb.WriteString("\r\n")

	for key, val := range v.Map {
		keyVal := &Value{Type: BulkString, String: key}
		keyStr, err := serializeBulkString(keyVal)
		if err != nil {
			return "", fmt.Errorf("serialize attribute key: %w", err)
		}
		sb.WriteString(keyStr)

		valStr, err := Serialize(val)
		if err != nil {
			return "", fmt.Errorf("serialize attribute value: %w", err)
		}
		sb.WriteString(valStr)
	}

	return sb.String(), nil
}
