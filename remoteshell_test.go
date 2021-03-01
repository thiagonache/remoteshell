package remoteshell_test

import (
	"bufio"
	"bytes"
	"net"
	"remoteshell"
	"testing"
	"time"
)

var listenAddr string = "127.0.0.1:8999"

func init() {
	go remoteshell.ListenAndServe(&bytes.Buffer{}, listenAddr)
}

func TestProtocol(t *testing.T) {
	conn, err := net.Dial("tcp4", listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.Write([]byte("hello\n"))
	want := "hello yourself\n"
	reader := bufio.NewReader(conn)
	got, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestProtocolMismatch(t *testing.T) {
	go remoteshell.ListenAndServe(&bytes.Buffer{}, listenAddr)
	// Sensei, do not lose your faith on me
	time.Sleep(300 * time.Millisecond)
	conn, err := net.Dial("tcp4", listenAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.Write([]byte("hi\n"))
	want := "Protocol error!\n"
	reader := bufio.NewReader(conn)
	got, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestProtocolActionAuth(t *testing.T) {
	commands := []string{"hello", "auth abc1234"}
	for _, cmd := range commands {
		want := "OK\n"
		got, err := remoteshell.ProtocolAction(cmd)
		if err != nil {
			t.Fatal(err)
		}
		if want != got {
			t.Errorf("want %q, got %q", want, got)
		}
	}
}
