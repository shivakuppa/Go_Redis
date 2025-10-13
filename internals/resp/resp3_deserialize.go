package resp

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func deserializeNull(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read null: %w", err)
	}

	if string(data) != "_" {
		return nil, fmt.Errorf("invalid null format: %w", err)
	}

	return &Value{
		Type:   Null,
		IsNull: true,
	}, nil
}

func deserializeBoolean(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read boolean: %w", err)
	}

	if len(data) != 1 || (data[0] != 't' && data[0] != 'f') {
		return nil, fmt.Errorf("invalid boolean value: %q", data)
	}

	return &Value{
		Type: Boolean,
		Bool: data[0] == 't',
	}, nil
}

func deserializeDouble(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read double value: %w", err)
	}

	str := string(data)
	switch strings.ToLower(str) {
	case "inf":
		return &Value{Type: Double, Double: math.Inf(1)}, nil
	case "-inf":
		return &Value{Type: Double, Double: math.Inf(-1)}, nil
	case "nan":
		return &Value{Type: Double, Double: math.NaN()}, nil
	}

	double, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil, fmt.Errorf("parse double %q: %w", str, err)
	}

	return &Value{
		Type:   Double,
		Double: double,
	}, nil
}

func deserializeBigNumber(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read big number")
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty big number")
	}

	start := 0
	if data[0] == '+' || data[0] == '-' {
		if len(data) == 1 {
			return nil, fmt.Errorf("big number cannot be only a sign")
		}
		start = 1
	}

	for i := start; i < len(data); i++ {
		if data[i] < '0' || data[i] > '9' {
			return nil, fmt.Errorf("invalid big number")
		}
	}

	return &Value{
		Type:   BigNumber,
		String: string(data),
	}, nil
}

func deserializeBulkError(reader *bufio.Reader) (*Value, error) {
	lengthBytes, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read bulk error length: %w", err)
	}

	length, err := strconv.ParseInt(string(lengthBytes), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse bulk error length: %w", err)
	}

	data := make([]byte, length+2)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, fmt.Errorf("read bulk error data: %w", err)
	}

	return &Value{
		Type:   BulkError,
		String: string(data[:length]),
	}, nil
}

func deserializeVerbatimString(reader *bufio.Reader) (*Value, error) {
	lengthBytes, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read verbatim length: %w", err)
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return nil, fmt.Errorf("parse verbatim length: %w", err)
	}

	data := make([]byte, length+2)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, fmt.Errorf("read verbatim string: %w", err)
	}

	return &Value{
		Type:   VerbatimString,
		String: string(data[:length]),
	}, nil
}

func deserializeMap(reader *bufio.Reader) (*Value, error) {
	countBytes, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read map length: %w", err)
	}
	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		return nil, fmt.Errorf("parse map length: %w", err)
	}

	m := make(map[string]*Value, count)
	for i := 0; i < count; i++ {
		keyVal, err := Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("read map key %d: %w", i, err)
		}
		valVal, err := Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("read map value %d: %w", i, err)
		}
		m[keyVal.String] = valVal
	}

	return &Value{
		Type: Map,
		Map:  m,
	}, nil
}

func deserializeAttribute(reader *bufio.Reader) (*Value, error) {
	// Format is same as Map, just different type
	v, err := deserializeMap(reader)
	if err != nil {
		return nil, fmt.Errorf("read attribute: %w", err)
	}
	v.Type = Attribute
	return v, nil
}

func deserializeSet(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read set length: %w", err)
	}

	numElements, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse set length: %w", err)
	}

	if numElements == -1 {
		return &Value{
			Type:   Set,
			IsNull: true,
		}, nil
	}

	if numElements < 0 {
		return nil, fmt.Errorf("invalid set length: %d", numElements)
	}

	set := make(map[*Value]struct{}, numElements)
	for i := int64(0); i < numElements; i++ {
		elem, err := Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("error reading set element %d: %w", i, err)
		}
		set[elem] = struct{}{}
	}

	return &Value{
		Type: Set,
		Set:  set,
	}, nil
}

func deserializePush(reader *bufio.Reader) (*Value, error) {
	countBytes, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read push length: %w", err)
	}

	count, err := strconv.Atoi(string(countBytes))
	if err != nil {
		return nil, fmt.Errorf("parse push length: %w", err)
	}

	elements := make([]*Value, count)
	for i := 0; i < count; i++ {
		elem, err := Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("read push element %d: %w", i, err)
		}
		elements[i] = elem
	}

	return &Value{
		Type:  Push,
		Array: elements,
	}, nil
}
