package main

import (
	"log"
	"os"
	"remoteshell"
)

func main() {
	log.Fatal(remoteshell.ListenAndServe(os.Stdout, ":2020"))
}
