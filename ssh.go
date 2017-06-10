package cli

import (
	"golang.org/x/crypto/ssh"
	"fmt"
	"log"
	"io"
	"time"
	"strings"
)

type SSHClient struct {
	io.Reader
	io.WriteCloser
}

func NewSSHClient(user string, password string, host string, port string) (*SSHClient) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "3des-cbc")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	connection, err := ssh.Dial("tcp", host + ":" + port, sshConfig)

	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	// todo: add errors checking
	session, _ := connection.NewSession()
	sshOut, _ := session.StdoutPipe()
	sshIn, _ := session.StdinPipe()

	if err := session.RequestPty("xterm", 0, 200, modes); err != nil {
		session.Close()
	}

	if err := session.Shell(); err != nil {
		session.Close()
	}

	return &SSHClient{sshOut, sshIn}
}


// accept timeout, action map, error map, return output from the device
func (s *SSHClient) SendCommand(command string, excpectedOut string) (string) {
	timeout := time.After(5 * time.Second)
	ch := make(chan string)

	s.writeBuff(command)

	go s.readBuff(excpectedOut, ch)

	for {
		select {
			case <-timeout:
				fmt.Println("Timeout Error!")
				return "time out error"
			case response := <-ch:
				return response
		}
	}
}

func (s *SSHClient) readBuff(excpectedOut string, ch chan string) {
	buf := make([]byte, 1000)
	bufferStr := ""

	for {
		fmt.Println("    .")
		//fmt.Println("Start reading bytes....")
		// sshOut.Read <- this will stuck if there is now any data in the buffer, as a result gorutine will not be deleted, this is a memory leak
		// need to find some command like sshOut.Read(buf, timeout)
		n, _ := s.Read(buf) //this reads the ssh terminal
		//fmt.Println("Finish reading bytes....")

		fmt.Println("read ", n, " Bytes")
		bufferStr += string(buf)
		//fmt.Println("REsult ", bufferStr, " '\n")

		if strings.Contains(bufferStr, excpectedOut) {
			//fmt.Println("Got response, send it to the channel", bufferStr)
			ch <- bufferStr
			return
		}
		buf = make([]byte, 1000) // clear buffer
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *SSHClient) writeBuff(command string) (int, error) {
	returnCode, err := s.Write([]byte(command + "\r"))
	return returnCode, err
}

func (s *SSHClient) Close() {
	//todo: close opened SSH session
}