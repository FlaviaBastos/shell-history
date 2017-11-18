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
	timeout = 300 * time.Millisecond
)

func connect(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getinformation(argsWithoutProg []string, commandExitCode int64) spb.Command {
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
	h.Exitcode = commandExitCode
	if os.Geteuid() == 0 {
		h.Altusername = os.Getenv("SUDO_USER")
	}
	return h
}

func main() {
	commandExitCode := flag.Int64("e", 0, "Exit code of last command")
	flag.Parse()

	if len(os.Args) < 4 {
		return
	}

	argsWithoutProg := os.Args[3:]

	conn, err := connect(address)
	defer conn.Close()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := spb.NewHistorianClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	h := getinformation(argsWithoutProg, *commandExitCode)

	r, err := c.GetCommand(ctx, &h)
	if err != nil {
		log.Fatalf("could not save command: %v", err)
	}

	if r.Status == spb.Status_ERR {
		log.Fatalf("Received error while uploading command")
	}
	return
}
