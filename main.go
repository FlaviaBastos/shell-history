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
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	spb "github.com/ebastos/shell-history/history"
	"github.com/gobuffalo/packr"
)

const (
	timeout = 300 * time.Millisecond
)

// Config
type Config struct {
	RemoteHost string `json:"remote_host"`
	RemotePort int    `json:"remote_port"`
	Disabled   bool   `json:"disabled"`
	Secure     bool   `json:"secure"`
}

func initConfig() Config {
	config := Config{}
	config.Disabled = false
	config.RemoteHost = "localhost"
	config.RemotePort = 50051

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// Open shell-history.json
	jsonFile, err := os.Open(currentUser.HomeDir + "/.config/shell-history.json")

	// if we os.Open returns an error log it.
	if err != nil {
		if strings.ContainsAny("no such file or directory", err.Error()) {
			fmt.Println(err)
			return config
		}
		log.Fatal(err)
	} else {
		// Read file.
		byteValue, _ := ioutil.ReadAll(jsonFile)

		// Convert json to Config struct.
		json.Unmarshal(byteValue, &config)
	}

	jsonFile.Close()
	return config
}

// Transforms a source string into a new (and possibly different) output string.
type Transformer interface {
	transform(input string) (output string)
}

// Represents a filter that maps a regex key to a regex transform.
type Redactor map[string]string

// Transforms source string to and output string when it matches a defined
// redaction.
func (redactor Redactor) transform(source string) (result string) {
	result = source

	for key, value := range redactor {
		regex, err := regexp.Compile(key)
		if err != nil {
			log.Fatalf("Redactor key %q is an invalid regexp", key)
		}
		result = regex.ReplaceAllString(result, value)
	}

	return
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

func connectInsecure(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil

}

func getinformation(argsWithoutProg []string, commandExitCode int64) spb.Command {
	var h spb.Command
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	h.Hostname, _ = os.Hostname()
	h.Timestamp = time.Now().UTC().Unix()
	h.Username = currentUser.Username
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
	config := initConfig()
	address := config.RemoteHost + ":" + strconv.Itoa(config.RemotePort)

	if config.Disabled {
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

	conn, err := connectInsecure(address)
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
