package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/user"
	"time"

	"google.golang.org/grpc"

	spb "github.com/ebastos/shell-history/history"
)

const (
	address = "localhost:50051"
)

func main() {
	commandExitCode := flag.Int64("e", 0, "Exit code of last command")
	flag.Parse()

	if len(os.Args) < 4 {
		return
	}

	argsWithoutProg := os.Args[3:]

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := spb.NewHistorianClient(conn)

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
	_, err = c.GetCommand(context.Background(), &h)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// log.Printf("Greeting: %s", r.Response)
	return
}
