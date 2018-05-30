package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

var listener = flag.Bool("listener", false, "pass listener if it is listener if it is generator dont pass")
var destination = flag.String("destination", "127.0.0.1", "destination to send data")
var chunkSize = flag.Int("chunkSize", 1048576, "chunk size in bytes")
var totalSizeInMB = flag.Int("totalSizeInMB", 250, "total size in MB")

func main() {
	flag.Parse()
	if *listener {
		StartListener()
	} else {
		StartGenerator(destination)
	}
}

func StartListener() {

	fmt.Println("Starting listener...")
	ln, errListen := net.Listen("tcp", ":8081")
	if errListen != nil {
		fmt.Printf("net.Listen faild with exception")
		fmt.Print(errListen)
		return
	}
	conn, errAccept := ln.Accept()
	defer conn.Close()
	if errAccept != nil {
		fmt.Printf("ln.Accept faild with exception")
		fmt.Print(errAccept)
		return
	}

	numberOfIteration := ((*totalSizeInMB * 1024 * 1024) / *chunkSize)
	fmt.Println(numberOfIteration)

	start := time.Now()
	started := false
	for i := 0; i < numberOfIteration; i++ {
		bufio.NewReader(conn).ReadString('\n')
		if !started {
			start = time.Now()
			started = true
		}
		//fmt.Println("Message Received Length", len(message))
		//conn.Write([]byte("\n"))
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Seconds())
	fmt.Println(float64(*totalSizeInMB) / elapsed.Seconds())
}

func StartGenerator(destination *string) {
	address := *destination + ":8081"
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		fmt.Printf("net.Dial faild with exception")
		fmt.Print(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	numberOfIteration := ((*totalSizeInMB * 1024 * 1024) / *chunkSize)
	fmt.Println(numberOfIteration)

	var data = prepareData(chunkSize)
	start := time.Now()

	for i := 0; i < numberOfIteration; i++ {

		// text, _ := reader.re('\n')
		fmt.Fprintf(conn, *data)
		//bufio.NewReader(conn).ReadString('\n')

		// if err != nil {
		// 	fmt.Println("Message send failed")
		// } else {
		// 	//fmt.Print("Message send successfull time ")

		// }
	}
	elapsed := time.Since(start)
	fmt.Println(elapsed.Seconds())
	fmt.Println(float64(*totalSizeInMB) / elapsed.Seconds())
}

func prepareData(chunkSize *int) *string {

	data := RandString(*chunkSize - 1)
	return &data
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return (string(b) + "\n")
}
