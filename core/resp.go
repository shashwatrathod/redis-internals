package core

import (
	"errors"
	"fmt"
)

const (
	RespSimpleStringIdentifier byte = '+'
	RespSimpleErrorIdentifier  byte = '-'
	RespIntegerIdentifier      byte = ':'
	RespBulkStringIdentifier   byte = '$'
	RespArrayIdentifier        byte = '*'
)

// Decodes a RESP encoded simple string
// and returns the decoded string, delta and error(if any)
func decodeRespSimpleString(data []byte, startPos int) (string, int, error) {
	if data[startPos] != RespSimpleStringIdentifier &&
		data[startPos] != RespSimpleErrorIdentifier {
		return "", startPos, errors.New("data is not a simple string")
	}

	startPos += 1
	idx := findCarriageReturnIdx(data, startPos)

	return string(data[startPos:idx]), idx + 2, nil
}

// Decodes a RESP encoded Bulk String.
// Returns the decoded string, delta and errors (if any)
func decodeRespBulkString(data []byte, startPos int) (string, int, error) {
	if data[startPos] != RespBulkStringIdentifier {
		return "", startPos, errors.New("data is not a bulk string")
	}

	startPos += 1
	idx := findCarriageReturnIdx(data, startPos)

	var length int = int(stoi(string(data[startPos:idx])))
	idx += 2
	endIdx := idx + length

	return string(data[idx:endIdx]), endIdx + 2, nil
}

// Decodes a RESP encoded Simple Error.
// Returns the decoded Error as a simple string, delta, and error (if any)
func decodeRespSimpleError(data []byte, startPos int) (string, int, error) {
	return decodeRespSimpleString(data, startPos)
}

// decodes a RESP encoded signed/unsigned integer.
// returns the decoded int64, delta, and error (if any)
func decodeRespInteger(data []byte, startPos int) (int64, int, error) {
	if data[startPos] != RespIntegerIdentifier {
		return -1, startPos, errors.New("data is not of int64 type")
	}

	startPos++

	var sign int64 = 1

	if data[startPos] == '+' || data[startPos] == '-' {
		if data[startPos] == '-' {
			sign = -1
		}

		startPos++
	}

	endIdx := findCarriageReturnIdx(data, startPos)

	return sign * stoi(string(data[startPos:endIdx])), endIdx + 2, nil
}

// Decodes a RESP encoded array.
// Returns the decoded array, delta and error (if any)
func decodeRespArray(data []byte, startPos int) ([]interface{}, int, error) {

	var res []interface{} = []interface{}{}

	if data[startPos] != RespArrayIdentifier {
		return res, startPos, errors.New("data is not of array type")
	}

	idx := findCarriageReturnIdx(data, startPos+1)

	arrayLength := stoi(string(data[startPos+1 : idx]))
	idx += 2

	for iter := 0; iter < int(arrayLength); iter++ {
		d, nextPos, err := DecodeOne(data, idx)
		if err != nil {
			return []interface{}{}, startPos, err
		}

		idx = nextPos
		res = append(res, d)
	}

	return res, idx, nil
}

func DecodeOne(data []byte, startPos int) (interface{}, int, error) {
	if len(data[startPos:]) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[startPos] {
	case RespSimpleStringIdentifier:
		return decodeRespSimpleString(data, startPos)
	case RespSimpleErrorIdentifier:
		return decodeRespSimpleError(data, startPos)
	case RespBulkStringIdentifier:
		return decodeRespBulkString(data, startPos)
	case RespIntegerIdentifier:
		return decodeRespInteger(data, startPos)
	case RespArrayIdentifier:
		return decodeRespArray(data, startPos)
	default:
		return nil, 0, errors.New("unknown datatype: " + string(data[startPos:]))
	}
}

// Decode takes a byte slice as input and attempts to decode it into the given type interface{}.
// It returns the decoded value and any error encountered during the decoding process.
//
// Parameters:
// - data: A byte slice containing the data to be decoded.
//
// Returns:
// - An interface{} representing the decoded value.
// - An error if the decoding process fails or if the input data is empty.
func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeOne(data, 0)
	return value, err
}

// converts the given string into an int64.
// assumes that the numString is a valid integer.
func stoi(numString string) int64 {
	var num int64 = 0
	for i := 0; i < len(numString); i++ {
		num = num*10 + int64(numString[i]-'0')
	}

	return num
}

// Finds and returns the index of the '\r' (carriage return)
// character in the byte array. Returns -1 if \r not found.
func findCarriageReturnIdx(data []byte, startPos int) int {
	crIdx := -1

	idx := startPos
	for {
		if data[idx] == '\r' {
			crIdx = idx
			break
		} else {
			idx++
		}
	}

	return crIdx
}

func Encode(val interface{}, isSimpleStr bool) []byte {
	switch value := val.(type) {
	case string:
		if isSimpleStr {
			return []byte(fmt.Sprintf("%c%s\r\n", RespSimpleStringIdentifier, value))
		} else {
			return []byte(fmt.Sprintf("%c%d\r\n%s\r\n", RespBulkStringIdentifier, len(value), value))
		}
	case error:
		return []byte(fmt.Sprintf("%c%s\r\n", RespSimpleErrorIdentifier, value))
	default:
		return []byte(fmt.Sprintf("%c%s\r\n", RespSimpleErrorIdentifier, value))
	}
}
