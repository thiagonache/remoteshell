package remoteshell

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func ListenAndServe(w io.Writer, listenAddr string) error {
	ln, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Listening on %s\n", ln.Addr())
	defer ln.Close()
	for {
		var rfc666 []string = []string{"hello\n", "auth\n", "command\n", "ciao\n"}
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		reader := bufio.NewReader(conn)
		for {
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			cmdExpected := rfc666[0]
			rfc666 = rfc666[1:]
			if cmdExpected != input {
				conn.Write([]byte("Protocol mismatch\n"))
				conn.Close()
				break
			}
			conn.Write([]byte("OK\n"))
			if input == "ciao\n" {
				conn.Close()
				break
			}
		}
	}
}
