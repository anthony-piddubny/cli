package main

import (
	"fmt"
	"github.com/anthony-piddubny/cli"
)

const (
	HOST            = "192.168.42.235"
	PORT 		= "22"
	USER            = "root"
	PASSWORD        = "Password1"
)

func main() {
	client := cli.NewSSHClient(USER, PASSWORD, HOST, PORT)
	defer client.Close()

	// start commands there
	// todo: move this to the session
	client.SendCommand("terminal length 0", "#")
	client.SendCommand("no logging console", "#")

	r1 := client.SendCommand("enable", "#")
	r2 := client.SendCommand("show run", "#")

	fmt.Println("Response 1: ", r1)
	fmt.Println("Response 2: ", r2)
}
