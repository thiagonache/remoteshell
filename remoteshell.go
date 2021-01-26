package remoteshell

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

// I know it should not be hard coded :)
const appToken = "abc1234"

// ProtocolAction takes a string with a command of RFC666
// and take an action depending on the command received
func ProtocolAction(command string) (string, error) {
	args := strings.Split(command, " ")
	switch args[0] {
	case "auth":
		authCode := args[1]
		if authCode != appToken {
			return "", errors.New("invalid auth")
		}
	default:
		return "OK\n", nil
	}

	return "OK\n", nil
}

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	var rfc666 []string = []string{"hello", "auth (\\w+)", "command \\S+", "ciao"}
	reader := bufio.NewReader(conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimRight(input, "\n")

		cmdExpected := rfc666[0]
		rfc666 = rfc666[1:]
		re := regexp.MustCompile(fmt.Sprintf("%s$", cmdExpected))
		find := re.Find([]byte(input))
		// Command expected does not match with the input
		if len(find) == 0 {
			conn.Write([]byte("Protocol mismatch\n"))
			conn.Close()
			break
		}

		output, err := ProtocolAction(input)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
			conn.Close()
			break
		}
		conn.Write([]byte(output))
		if input == "ciao" {
			conn.Close()
			break
		}
	}
	return nil
}

// ListenAndServe takes an io.Writer and a listenAddr in string format.
// It does create a listen, wait for new connection and handles it.
func ListenAndServe(w io.Writer, listenAddr string) error {
	ln, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Listening on %s\n", ln.Addr())
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		err = handleConnection(conn)
		switch err {
		case io.EOF:
			continue
		case nil:
			continue
		default:
			return err
		}
	}
}
