package resp

import (
	"bytes"
	"errors"
)

func ParseSimpleString(buffer []byte) (string,int,error) {
	
	if len(buffer) == 0 {
		return "", 0, errors.New("buffer is empty")
	}
	if buffer[0] != '+' {
		return "",0,errors.New("-ERR Invalid RESP Type")
	}
	newLineIdx := bytes.IndexByte(buffer,'\n')
	if newLineIdx == -1 {
		return "",0,errors.New("Incomplete Command")
	}
	if(newLineIdx < 1 || buffer[newLineIdx-1] != '\r') {
		return "",0,errors.New("-ERR Protocol error: Missing CR before LF\r\n")
	}
	
	payload := buffer[1:newLineIdx-1]
	command := string(payload)

	bytesConsumed := newLineIdx+1

	return command,bytesConsumed,nil

}