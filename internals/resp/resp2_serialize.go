package resp

import (
	"strconv"
	"fmt"
)

func serializeSimpleString(v *Value) (string, error) {
	return "+" + v.String + "\r\n", nil
}

func serializeSimpleError(v *Value) (string, error) {
	return "-" + v.String + "\r\n", nil
}

func serializeInteger(v *Value) (string, error) {
	return ":" + strconv.Itoa(int(v.Integer)) + "\r\n", nil
}

func serializeBulkString(v *Value) (string, error) {
	if v.IsNull {
		return "$-1\r\n", nil
	}
	return "$" + strconv.Itoa(len(v.String)) + "\r\n" + v.String + "\r\n", nil
}

func serializeArray(v *Value) (string, error) {
	if v.IsNull {
		return "*-1\r\n", nil
	}

	var serialized string
	for _, elem := range v.Array {
		s, err := Serialize(elem)
		if err != nil {
			return "", fmt.Errorf("serializing error element: %w", err)
		}
		serialized += s
	}

	return "*" + strconv.Itoa(len(v.Array)) + "\r\n" + serialized, nil
}