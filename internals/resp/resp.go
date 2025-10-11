package resp

import (
	"bufio"
	"fmt"
	"io"
)

// RESP data type	Minimal protocol version	Category	First byte
// Simple strings	RESP2	Simple	+
// Simple Errors	RESP2	Simple	-
// Integers	RESP2	Simple	:
// Bulk strings	RESP2	Aggregate	$
// Arrays	RESP2	Aggregate	*
// Nulls	RESP3	Simple	_
// Booleans	RESP3	Simple	#
// Doubles	RESP3	Simple	,
// Big numbers	RESP3	Simple	(
// Bulk errors	RESP3	Aggregate	!
// Verbatim strings	RESP3	Aggregate	=
// Maps	RESP3	Aggregate	%
// Attributes	RESP3	Aggregate	|
// Sets	RESP3	Aggregate	~
// Pushes	RESP3	Aggregate	>

// DataType represents the first byte of a RESP or RESP3 message type.
type RESPDataType byte

const (
	// RESP2 types
	SimpleString 	RESPDataType = '+' // e.g., +OK\r\n
	SimpleError  	RESPDataType = '-' // e.g., -Error message\r\n
	Integer      	RESPDataType = ':' // e.g., :1000\r\n
	BulkString   	RESPDataType = '$' // e.g., $6\r\nfoobar\r\n
	Array        	RESPDataType = '*' // e.g., *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n

	// RESP3 types
	Null            RESPDataType = '_' // Null value
	Boolean         RESPDataType = '#' // e.g., #t\r\n or #f\r\n
	Double          RESPDataType = ',' // e.g., ,1.23\r\n
	BigNumber       RESPDataType = '(' // e.g., (3492890328409238509324850943850943825024385\r\n
	BulkError       RESPDataType = '!' // e.g., !21\r\nSYNTAX invalid syntax\r\n
	VerbatimString  RESPDataType = '=' // e.g., =15\r\ntxt:Some string\r\n
	Map             RESPDataType = '%' // e.g., %2\r\n+key1\r\n+val1\r\n+key2\r\n+val2\r\n
	Attribute       RESPDataType = '|' // metadata key/value pairs
	Set             RESPDataType = '~' // e.g., ~2\r\n+orange\r\n+apple\r\n
	Push            RESPDataType = '>' // e.g., >4\r\n+pubsub\r\n+message\r\n+chan\r\n+hello\r\n
)

type Value struct {
	Type 	RESPDataType
	String 	string
	Integer int64
	Array 	[]*Value
	IsNull	bool
}

func readUntilCRLF(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read line bytes: %w", err)
	}
	if len(line) < 2 || line[len(line) - 2] != '\r' {
		return nil, fmt.Errorf("line not terminated correctly")
	}

	return line[:len(line) - 2], nil
}

func Deserialize(reader io.Reader) (*Value, error) {
	bufreader := bufio.NewReader(reader)
	
	respType, err := bufreader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read resp first byte: %w", err)
	}

	switch RESPDataType(respType) {
	case SimpleString:
		return deserializeSimpleString(bufreader)

	case SimpleError:
		return deserializeSimpleError(bufreader)

	case Integer:
		return deserializeInteger(bufreader)

	case BulkString:
		return deserializeBulkString(bufreader)

	case Array:
		return deserializeArray(bufreader)

	// RESP3 types
	// case Null:
	// 	return deserializeNull(bufreader)

	// case Boolean:
	// 	return deserializeBoolean(bufreader)

	// case Double:
	// 	return deserializeDouble(bufreader)

	// case BigNumber:
	// 	return deserializeBigNumber(bufreader)

	// case BulkError:
	// 	return deserializeBulkError(bufreader)

	// case VerbatimString:
	// 	return deserializeVerbatimString(bufreader)

	// case Map:
	// 	return deserializeMap(bufreader)

	// case Attribute:
	// 	return deserializeAttribute(bufreader)

	// case Set:
	// 	return deserializeSet(bufreader)

	// case Push:
	// 	return deserializePush(bufreader)

	default:
		return nil, fmt.Errorf("invalid RESP type: %c", respType)
	}
}

func Serialize(value *Value) (string, error) {
	if value == nil {
		return "", fmt.Errorf("value is nil")
	}

	switch value.Type {
	case SimpleString:
		return serializeSimpleString(value)

	case SimpleError:
		return serializeSimpleError(value)

	case Integer:
		return serializeInteger(value)

	case BulkString:
		return serializeBulkString(value)

	case Array:
		return serializeArray(value)

	// RESP3 types
	// case Null:
	// 	return serializeNull(bufreader)

	// case Boolean:
	// 	return serializeBoolean(bufreader)

	// case Double:
	// 	return serializeDouble(bufreader)

	// case BigNumber:
	// 	return serializeBigNumber(bufreader)

	// case BulkError:
	// 	return serializeBulkError(bufreader)

	// case VerbatimString:
	// 	return serializeVerbatimString(bufreader)

	// case Map:
	// 	return serializeMap(bufreader)

	// case Attribute:
	// 	return serializeAttribute(bufreader)

	// case Set:
	// 	return serializeSet(bufreader)

	// case Push:
	// 	return serializePush(bufreader)

	default:
		return "", fmt.Errorf("invalid RESP type: %c", value.Type)
	}
}