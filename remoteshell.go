package remoteshell

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	history "github.com/thiagonache/go-history"
)

// I know it should not be hard coded :)
const appToken = "abc1234"

func handleConnection(conn net.Conn, r *history.Recorder) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	if scanner.Text() != "hello" {
		fmt.Fprintln(conn, "Protocol error!")
		return
	}
	fmt.Fprintln(conn, "hello yourself")
	if !scanner.Scan() {
		return
	}
	line := scanner.Text()
	items := strings.Split(line, " ")
	if len(items) < 2 {
		fmt.Fprintln(conn, "Protocol error!")
		return
	}
	if items[0] != "auth" {
		fmt.Fprintln(conn, "Protocol error!")
		return
	}
	if items[1] != appToken {
		fmt.Fprintln(conn, "Protocol error!")
		return
	}
	fmt.Fprintln(conn, "welcome to the VIP lounge")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "ciao" {
			fmt.Fprintln(conn, "Aloha")
			return
		}
		items := strings.Split(line, " ")
		if len(items) < 2 {
			fmt.Fprintln(conn, "Protocol error!")
			return
		}
		if items[0] != "command" {
			fmt.Fprintln(conn, "Protocol error!")
			return
		}
		fmt.Fprintf(conn, "Running command %q with args %q\n", items[1], items[2:])
		r.Execute(strings.Join(items, " "))
	}
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
	r, err := history.NewRecorder()
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		handleConnection(conn, r)
	}
}
