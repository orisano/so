package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
)

func run() error {
	addr := flag.String("addr", "", "(required)")
	in := flag.Bool("i", false, "show stdin")
	out := flag.Bool("o", false, "show stdout")
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		return fmt.Errorf("connect server: %w", err)
	}
	defer conn.Close()

	args := flag.Args()
	cmd := exec.Command(args[0], args[1:]...)
	var r io.Reader = conn
	if *in {
		r = io.TeeReader(r, os.Stdout)
	}
	var w io.Writer = conn
	if *out {
		w = io.MultiWriter(w, os.Stdout)
	}
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run program: %w", err)
	}
	if *in {
		_, err = io.Copy(os.Stdout, conn)
	} else {
		_, err = io.Copy(ioutil.Discard, conn)
	}
	return err
}

func main() {
	log.SetPrefix("so: ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
