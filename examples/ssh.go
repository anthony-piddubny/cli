package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"time"
	"github.com/anthony-piddubny/cli"
)

const (
	HOST = "192.168.42.235"
	USER = "root"
	PASSWORD = "Password1"
	ENABLE_PASSWORD = "Password2"
)


func writeBuff(command string, sshIn io.WriteCloser) (int, error) {
	returnCode, err := sshIn.Write([]byte(command + "\r"))
	return returnCode, err
}

func handleError(e error, fatal bool, customMessage ...string) {
	var errorMessage string
	if e != nil {
		if len(customMessage) > 0 {
			errorMessage = strings.Join(customMessage, " ")
		} else {
			errorMessage = "%s"
		}
		if fatal == true {
			log.Fatalf(errorMessage, e)
		} else {
			log.Print(errorMessage, e)
		}
	}
}


func main() {

	ss := cli.NewSession()
	fmt.Println(ss)

	sshConfig := &ssh.ClientConfig{
		User: USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.Config.Ciphers = append(sshConfig.Config.Ciphers, "3des-cbc")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	connection, err := ssh.Dial("tcp", HOST+":22", sshConfig)

	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}

	session, err := connection.NewSession()
	handleError(err, true, "Failed to create session: %s")
	sshOut, err := session.StdoutPipe()
	handleError(err, true, "Unable to setup stdin for session: %v")
	sshIn, err := session.StdinPipe()
	handleError(err, true, "Unable to setup stdout for session: %v")




	if err := session.RequestPty("xterm", 0, 200, modes); err != nil {
		session.Close()
		handleError(err, true, "request for pseudo terminal failed: %s")
	}
	if err := session.Shell(); err != nil {
		session.Close()
		handleError(err, true, "request for shell failed: %s")
	}



	// start commands there

	buf := make([]byte, 1000)
	bufferStr := ""
	excpectedStr := "#sadad"

	//tick := time.Tick(100 * time.Millisecond)
	timeout := time.After(5 * time.Second)
	if _, err := writeBuff("enable", sshIn); err != nil {
		handleError(err, true, "Failed to run: %s")
	}
	for {
		select {
		case <-timeout:
			fmt.Println("Timeout Error!")
			return

		default:
			fmt.Println("    .")

			fmt.Println("Start reading bytes....")
			n, err := sshOut.Read(buf) //this reads the ssh terminal
			if err != nil {
				if err != io.EOF {
					log.Printf("Read error: %s", err)
				}
				break
			}
			fmt.Println("Finish reading bytes....")

			fmt.Println("read ", n, " Bytes")
			bufferStr += string(buf)
			fmt.Println("REsult ", bufferStr, " '\n")

			if strings.Contains(bufferStr, excpectedStr) {
				fmt.Println("GOT RESPONSE !!!!", bufferStr)
				return
			}
			// clear buffer !!s
			buf = make([]byte, 1000)
			time.Sleep(50 * time.Millisecond)
		}
	}








	// WORKING
	//if _, err := writeBuff("enable", sshIn); err != nil {
	//	handleError(err, true, "Failed to run: %s")
	//}
	//
	//waitingString := ""
	//buf := make([]byte, 1000)
	//time.Sleep(time.Second * 5)
	//
	//
	//n, err := sshOut.Read(buf) //this reads the ssh terminal
	//
	//fmt.Println("read ", n, " Bytes")
	//
	//waitingString += string(buf)
	//handleError(err, true, "failed to read from terminal: %s")
	//fmt.Println("read: ", waitingString)
	// END WORKING

	session.Close()
}