package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	iface := "en0" // default interface; change to "en1", "utun0", etc. if needed
	if len(os.Args) > 1 {
		iface = os.Args[1]
	}
	fmt.Println("Used interface: ", iface)
	fd, err := Open() //https://docs.oracle.com/cd/E36784_01/html/E36884/bpf-7d.html
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	// _ = readFD(fd)

	if err := setImmediate(fd); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	if err := BindInterface(fd, iface); err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	// _ = readFD(fd) //in this point we can see raw network packet bytes
	buffLen, err := GetBuffLen(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}

	buffer := make([]byte, buffLen)

	for {
		n, err := unix.Read(fd, buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
			break
		}
		fmt.Println("Read: ", n)
	}
}

// Just seeing what is inside this fd.
func readFD(fd int) error {
	buf := make([]byte, 4096)
	n, err := unix.Read(fd, buf)
	if err != nil {
		fmt.Println("Error: ", err) //(Error:  device not configured) This is the erro we will get
		// because  the BPF device is opened but never
		// bound it to a network interface.
		// A BPF fd must be attached to an interface before it can capture packets.
		// (//https://docs.oracle.com/cd/E36784_01/html/E36884/bpf-7d.html)
		return err
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%02x ", buf[i])
	}
	fmt.Println()
	return nil
}
