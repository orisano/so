package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func run() error {
	addr := flag.String("addr", "", "(required)")
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		return fmt.Errorf("connect server: %w", err)
	}
	defer conn.Close()

	args := flag.Args()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = io.TeeReader(conn, os.Stdout)
	cmd.Stdout = io.MultiWriter(conn, os.Stdout)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run program: %w", err)
	}
	_, err = io.Copy(os.Stdout, conn)
	return err
}

func main() {
	log.SetPrefix("so: ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}

}
