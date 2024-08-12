package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// RESP protocol type identifiers
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value represents a RESP (Redis Serialization Protocol) value
type Value struct {
	Type  string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

// Resp represents a RESP reader
type Resp struct {
	reader *bufio.Reader
}

// NewResp creates a new RESP reader
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Read reads and parses a RESP value
func (r *Resp) Read() (Value, error) {
	typeChar, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeChar {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case STRING:
		return r.readString()
	case INTEGER:
		return r.readInteger()
	case ERROR:
		return r.readError()
	default:
		return Value{}, fmt.Errorf("unknown type: %v", string(typeChar))
	}
}

// readLine reads a line from the RESP stream
func (r *Resp) readLine() ([]byte, error) {
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	return []byte(strings.TrimRight(line, "\r\n")), nil
}

// readArray reads a RESP array
func (r *Resp) readArray() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		value, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		array[i] = value
	}

	return Value{Type: "array", Array: array}, nil
}

// readBulk reads a RESP bulk string
func (r *Resp) readBulk() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	size, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	if size == -1 {
		return Value{Type: "null"}, nil
	}

	bulk := make([]byte, size)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return Value{}, err
	}

	// Read the trailing CRLF
	_, err = r.readLine()
	if err != nil {
		return Value{}, err
	}

	return Value{Type: "bulk", Bulk: string(bulk)}, nil
}

// readString reads a RESP simple string
func (r *Resp) readString() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: "string", Str: string(line)}, nil
}

// readInteger reads a RESP integer
func (r *Resp) readInteger() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	return Value{Type: "integer", Num: num}, nil
}

// readError reads a RESP error
func (r *Resp) readError() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: "error", Str: string(line)}, nil
}

// Marshal converts a Value to its RESP byte representation
func (v Value) Marshal() []byte {
	switch v.Type {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "integer":
		return v.marshalInteger()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

// marshalArray converts an array Value to its RESP byte representation
func (v Value) marshalArray() []byte {
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.Array))...)
	bytes = append(bytes, '\r', '\n')
	for _, item := range v.Array {
		bytes = append(bytes, item.Marshal()...)
	}
	return bytes
}

// marshalBulk converts a bulk string Value to its RESP byte representation
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalString converts a simple string Value to its RESP byte representation
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalInteger converts an integer Value to its RESP byte representation
func (v Value) marshalInteger() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.Num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalError converts an error Value to its RESP byte representation
func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalNull converts a null Value to its RESP byte representation
func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}
