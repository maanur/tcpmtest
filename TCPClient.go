package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	var wg sync.WaitGroup
	var q = 3
	var wait = 60  // in Seconds
	var step = 100 // in Milliseconds
	var addr = "192.168.51.122:6508"
	q, wait, step, addr = promptVariables()
	for i := 0; i < q; i++ {
		wg.Add(1)
		go func() {
			sampleRun(i, wait, addr)
			defer wg.Done()
		}()
		time.Sleep(time.Duration(step) * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("Done")
}

func promptVariables() (q int, w int, s int, a string) {
	var r = false
	var err error
	for {
		q, err = strconv.Atoi(prompt("Quantity? : 3", "3", r))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Try again!")
		} else {
			fmt.Println("OK")
			break
		}
	}
	for {
		w, err = strconv.Atoi(prompt("Wait time? : 60", "60", r))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Try again!")
			r = true
		} else {
			break
		}
	}
	for {
		s, err = strconv.Atoi(prompt("Escalation step? : 100", "100", r))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Try again!")
			r = true
		} else {
			break
		}
	}
	for {
		a = prompt("Address:Port : 192.168.51.122:6508", "192.168.51.122:6508", r)
		_, err := net.ResolveTCPAddr("tcp", a)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Try again!")
			r = true
		} else {
			break
		}
	}
	return
}

func prompt(ask string, dft string, repeat bool) (output string) {
	consolereader := bufio.NewReader(os.Stdin)
	fmt.Println(ask)
	rn, err := consolereader.ReadBytes('\r') // this will prompt the user for input
	if err != nil {
		log.Fatal(err)
	}
	if !repeat {
		output = string(rn[:len(rn)-1])
	} else {
		output = string(rn[1 : len(rn)-1])
	}
	if output == "" {
		return dft
	}
	return output
}

func sampleRun(i int, wait int, addr string) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	var delim []byte
	delim = append(delim, 10)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := openConn(server)
	if err != nil {
		fmt.Println("Sample " + strconv.Itoa(i+1) + " connection error:")
		log.Fatal(err)
	}
	sendMsg(str2cp866("Привет, TCPServer"), conn, i)
	sendMsg(delim, conn, i)
	JOB := getMsg(conn, i)
	fmt.Println("Sample " + strconv.Itoa(i+1) + ": %JOB: " + JOB)
	fmt.Println("Sample " + strconv.Itoa(i+1) + ": Waiting...")
	time.Sleep(time.Duration(wait) * time.Second)
	fmt.Println("Sample " + strconv.Itoa(i+1) + ": ")
	sendMsg(str2cp866(strconv.Itoa(180020+i*20)), conn, i)
	closeConn(conn)
}

func openConn(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err == nil {
		fmt.Println("Connected...")
	}
	return conn, err
}

func closeConn(conn *net.TCPConn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disonnected...")
}

func str2cp866(str string) []byte {
	input := []byte(str)
	cp866 := charmap.CodePage866.NewEncoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Fatal(err)
	}
	return output
}

func cp8662str(input []byte) string {
	cp866 := charmap.CodePage866.NewDecoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}

func sendMsg(msg []byte, conn *net.TCPConn, i int) {
	_, err := conn.Write(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sample " + strconv.Itoa(i+1) + ": send: " + string(msg))
}

func getMsg(conn *net.TCPConn, i int) string {
	input := make([]byte, 64)
	_, err := conn.Read(input)
	if err != nil {
		log.Fatal(err)
	}
	output := cp8662str(input)
	fmt.Println("Sample " + strconv.Itoa(i+1) + ": read: " + output)
	return output
}
