package main

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"os"
	"os/user"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	spb "github.com/ebastos/shell-history/history"
	"github.com/gobuffalo/packr"
)

const (
	address = "localhost:50051"
	timeout = 300 * time.Millisecond
)

// Config
type Config struct {
	Enabled   boolean `json:"enabled"`
}

func initConfig() Config {
        config := Config{true}
	
        // Open shell-history.json
        jsonFile, err := os.Open("~/.config/shell-history.json")

        // if we os.Open returns an error log it.
        if err != nil {
                log.Fatal(err)		
        } else {
                //Read file.
                byteValue, _ := ioutil.ReadAll(jsonFile)

                //Convert json to Config struct.
                json.Unmarshal(byteValue, &config)
        }
        
        jsonFile.Close()

        return config
}

func connect(address string) (*grpc.ClientConn, error) {
	// Let's embed the certificate
	box := packr.NewBox("./certs")
	cert := box.Bytes("localhost.crt")
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(cert)

	creds := credentials.NewClientTLSFromCert(roots, "")

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds), grpc.WithTimeout(timeout))
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
	h.Timestamp = time.Now().UTC().Unix()
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
        config := config()

        if !strconv.ParseBool(config.Enabled) {
                log.Fatal("shell-history is not enabled.")
                return
        }
	
	commandExitCode := flag.Int64("e", 0, "Exit code of last command")
	flag.Parse()

	// logs go to syslog instead of user terminal.
	logwriter, err := syslog.New(syslog.LOG_NOTICE, "shell-history")
	if err == nil {
		log.SetOutput(logwriter)
	}

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
