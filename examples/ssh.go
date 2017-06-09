package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"time"
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
	now := time.Now()
	fmt.Println(now)
	sshConfig := &ssh.ClientConfig{
		User: USER,
		Auth: []ssh.AuthMethod{
			ssh.Password(PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//		Timeout: time.Duration(time.Second * 0,00000001),
// timeout for the wjole tcp connection?? ----- not what we neeed
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
//	session.SetReadDeadline(now)
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
	sendCommand(sshIn, sshOut, "terminal length 0", "#")
	sendCommand(sshIn, sshOut, "no logging console", "#")

    r1 := sendCommand(sshIn, sshOut, "enable", "#")
    r2 := sendCommand(sshIn, sshOut, "show run", "#")

    fmt.Println("Response 1: ", r1)
    fmt.Println("Response 2: ", r2)

	session.Close()
}


func readBuf(sshOut io.Reader, reSrt string, c chan string){
    buf := make([]byte, 1000)
	bufferStr := ""

    for {
        fmt.Println("    .")
        //fmt.Println("Start reading bytes....")
	// sshOut.Read <- this will stuck if there is now any data in the buffer, as a result gorutine will not be deleted, this is a memory leak
	// need to find some command like sshOut.Read(buf, timeout)
        n, _ := sshOut.Read(buf) //this reads the ssh terminal
        //fmt.Println("Finish reading bytes....")

        fmt.Println("read ", n, " Bytes")
        bufferStr += string(buf)
        //fmt.Println("REsult ", bufferStr, " '\n")

        if strings.Contains(bufferStr, reSrt) {
            //fmt.Println("Got response, send it to the channel", bufferStr)
            c <-bufferStr
            return
        }
        buf = make([]byte, 1000)  // clear buffer
        time.Sleep(500 * time.Millisecond)
    }
}



func sendCommand(sshIn io.WriteCloser, sshOut io.Reader, command string, reSrt string) (string){
    timeout := time.After(5 * time.Second)
    ch := make(chan string)

	if _, err := writeBuff(command, sshIn); err != nil {
		handleError(err, true, "Failed to run: %s")
	}
    go readBuf(sshOut, reSrt, ch)

	for {
		select {
		case <-timeout:
			fmt.Println("Timeout Error!")
			return "time out error"
        case response:= <-ch:
            //fmt.Println("Timeout Error!")
            //fmt.Println("GOT resosnse!")
            //fmt.Println("START>>>\n", response, "\n<<<END")
			return response
		default:
            fmt.Println("do nothing ....")
            time.Sleep(50 * time.Millisecond)
		}
	}
}
