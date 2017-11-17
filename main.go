package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"time"

	spb "github.com/ebastos/shell-history/history"
)

func main() {
	commandExitCode := flag.Int64("e", 0, "Exit code of last command")
	flag.Parse()

	if len(os.Args) < 4 {
		return
	}

	argsWithoutProg := os.Args[3:]
	var h spb.Command

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	h.Hostname, _ = os.Hostname()
	h.Timestamp = time.Now().Unix()
	h.Username = user.Username
	h.Cwd, err = os.Getwd()
	h.Oldpwd = os.Getenv("OLDPWD")
	h.Command = argsWithoutProg
	h.Exitcode = *commandExitCode
	if os.Geteuid() == 0 {

		h.Altusername = os.Getenv("SUDO_USER")
	}

}
