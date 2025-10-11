package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func deserializeSimpleString(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read simple string data: %w", err)
	}

	return &Value{
		Type:	SimpleString,
		String:	string(data),
	}, nil
}

func deserializeSimpleError(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read simple error data: %w", err)
	}

	return &Value{
		Type:	SimpleError,
		String:	string(data),
	}, nil
}

func deserializeInteger(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read integer data: %w", err)
	}

	integer, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("derserialize integer data: %w", err)
	}

	return &Value{
		Type:		Integer,
		Integer:	integer,
	}, nil
}

func deserializeBulkString(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read bulk string data: %w", err)
	}

	strLen, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse bulk string len: %w", err)
	}

	if strLen < -1 {
		return nil, fmt.Errorf("invalid bulk string len: %w", err)
	}

	if strLen == -1 {
		return &Value{
			Type:	BulkString,
			IsNull: true,
		}, nil
	}

	strBytes := make([]byte, strLen)
	readLen, err := io.ReadFull(reader, strBytes)
	if err != nil {
		return nil, fmt.Errorf("read bulk string: %w", err)
	}

	if readLen != int(strLen) {
		return nil, fmt.Errorf("wrong size of bulk string. expected: %d, got: %d", strLen, readLen)
	}

	crlf := make([]byte, 2)
	n, err := io.ReadFull(reader, crlf)
	if err != nil || n != 2 || crlf[0] != '\r' || crlf[1] != '\n' {
		return nil, fmt.Errorf("bulk string not terminated correctly: %c", crlf)
	}

	return &Value{
		Type:	BulkString,
		String: string(strBytes),
	}, nil
}

func deserializeArray(reader *bufio.Reader) (*Value, error) {
	data, err := readUntilCRLF(reader)
	if err != nil {
		return nil, fmt.Errorf("read array data: %w", err)
	}

	numElements, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse number of array elements: %w", err)
	}

	if numElements < -1 {
		return nil, fmt.Errorf("invalid number of array elements: %w", err)
	}

	if numElements == -1 {
		return &Value{
			Type:	Array,
			IsNull: true,
		}, nil
	}

	array := make([]*Value, numElements)
	for i := range numElements {
		element, err := Deserialize(reader)
		if err != nil {
			return nil, fmt.Errorf("error reading element at index %d: %w", i, err)
		}
		array[i] = element
	}

	return &Value{
		Type:	Array,
		Array: 	array,
	}, nil
}