package cli

import (
	"fmt"
	"io"
	"strings"
)

type Session struct {

}

func NewSession() (*Session) {
	return &Session{}
}

// accept timeout, action map, error map, return output from the device
func (s *Session) SendCommand() (string) {
	fmt.Println("TBD")
	commandOutput := "some output"


	var bufferStr string
	var chunk = make([]byte, bufferStrSize)


	readLen, err := conn.Read(chunk)

	if err == io.EOF {
		fmt.Println("?? the connection from: ", conn.RemoteAddr())
		//conn.Close()
		//break
	}
	fmt.Println("Read lines: ", readLen)
	bufferStr += string(chunk)

	if strings.Contains(bufferStr, excpectedStr) {
		// clean response (remove prompt, etc)
	}

	// clear bufferStr before next usage
	chunk = make([]byte, bufferStrSize)

	return commandOutput
}
