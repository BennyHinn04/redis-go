package resp

import (
	"bytes"
	"errors"
	"strconv"
)

func ParseRESP(inputBytes []byte) (any, int, error) {
	if len(inputBytes) == 0 {
		return "", 0, errors.New("input buffer is empty")
	}

	dataType := inputBytes[0]

	switch dataType {
	case '+':
		return ParseSimpleString(inputBytes)
	case '$':
		return ParseBulkString(inputBytes)
	case '*':
		return ParseArray(inputBytes)
	default:
		return nil, 0, errors.New("-ERR unknown RESP data type")
	}
}

func ParseSimpleString(inputBytes []byte) (string, int, error) {

	if len(inputBytes) == 0 {
		return "", 0, errors.New("input buffer is empty")
	}
	if inputBytes[0] != '+' {
		return "", 0, errors.New("-ERR RESP Protocol error: expected '+'")
	}

	newLineIdx := bytes.IndexByte(inputBytes, '\n')
	if newLineIdx == -1 {
		return "", 0, errors.New("Incomplete Command")
	}
	if inputBytes[newLineIdx-1] != '\r' {
		return "", 0, errors.New("-ERR Protocol error: expected CRLF")
	}

	payload := inputBytes[1 : newLineIdx-1]
	command := string(payload)

	bytesConsumed := newLineIdx + 1

	return command, bytesConsumed, nil

}

func ParseBulkString(inputBytes []byte) (string, int, error) {

	if len(inputBytes) == 0 {
		return "", 0, errors.New("input buffer is empty")
	}

	if inputBytes[0] != '$' {
		return "", 0, errors.New("-ERR Protocol error: expected '$'")
	}

	newLineIdx := bytes.IndexByte(inputBytes, '\n')
	if newLineIdx == -1 {
		return "", 0, errors.New("Incomplete command")
	}

	bulkStringLenStr := string(inputBytes[1 : newLineIdx-1])
	bulkStringLen, err := strconv.Atoi(bulkStringLenStr)
	if err != nil {
		return "", 0, errors.New("-ERR RESP Protocol error: invalid bulk string length")
	}

	if bulkStringLen == -1 {
		return "", newLineIdx + 1, nil
	}
	payloadStartPos := newLineIdx + 1
	payloadEndPos := payloadStartPos + bulkStringLen

	if len(inputBytes) < payloadEndPos+2 {
		return "", 0, errors.New("Incomplete Command")
	}

	if inputBytes[payloadEndPos] != '\r' || inputBytes[payloadEndPos+1] != '\n' {
		return "", 0, errors.New("ERR Protocol error: expected CRLF after bulk string")
	}

	payload := string(inputBytes[payloadStartPos:payloadEndPos])
	consumedBytes := payloadEndPos + 2
	return payload, consumedBytes, nil
}

func ParseArray(inputBytes []byte) ([]string, int, error) {
	if len(inputBytes) == 0 {
		return nil, 0, errors.New("buffer is empty")
	}

	if inputBytes[0] != '*' {
		return nil, 0, errors.New("-ERR RESP Protocol error: expected '*'")
	}

	newLineIdx := bytes.IndexByte(inputBytes, '\n')
	if newLineIdx == -1 {
		return nil, 0, errors.New("Incomplete Command")
	}

	arrElementCountStr := string(inputBytes[1 : newLineIdx-1])
	arrElementCount, err := strconv.Atoi(arrElementCountStr)
	if err != nil {
		return nil, 0, errors.New("-ERR RESP Protocol error: invalid array length")
	}
	offset := newLineIdx + 1

	var result []string
	for idx := 0; idx < arrElementCount; idx++ {

		str, consumedBytes, err := ParseBulkString(inputBytes[offset:])

		if err != nil {
			return nil, 0, err
		}

		result = append(result, str)

		offset += consumedBytes
	}
	return result, offset, nil
}
