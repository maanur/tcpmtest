package tcpmprobe

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func SampleRun(i int, wait int, addr string) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := openConn(server)
	if err != nil {
		println("Соединение " + strconv.Itoa(i+1) + " Ошибка соединения:")
		log.Fatal(err)
	}
	sendMsg("Привет, TCPServer", conn, i)
	JOB := getMsg(conn, i)
	println("Соединение " + strconv.Itoa(i+1) + ": %JOB: " + JOB)
	println("Соединение " + strconv.Itoa(i+1) + ": Ждем " + strconv.Itoa(wait) + " сек...")
	time.Sleep(time.Duration(wait) * time.Second)
	println("Соединение " + strconv.Itoa(i+1) + ": Посылаю код и отключаюсь")
	sendMsg(strconv.Itoa(180020+i*20), conn, i)
	closeConn(conn)
}

func openConn(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err == nil {
		println("Соединяемся...")
	}
	return conn, err
}

func closeConn(conn *net.TCPConn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	println("Отключился...")
}

//str2cp866 Конвертирует строку str в байт-код и перекодирует в cp866 и возвращает байт-код
func str2cp866(str string) []byte {
	input := []byte(str)
	cp866 := charmap.CodePage866.NewEncoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Fatal(err)
	}
	return output
}

//cp8662str Перекодирует input из cp866 и возвращает строку
func cp8662str(input []byte) string {
	cp866 := charmap.CodePage866.NewDecoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Fatal(err)
	}
	return string(output)
}

func sendMsg(msg string, conn *net.TCPConn, i int) {
	_, err := conn.Write(append(str2cp866(msg), 10))
	if err != nil {
		log.Fatal(err)
	}
	println("Соединение " + strconv.Itoa(i+1) + ": Отправили: " + string(msg))
}

func getMsg(conn *net.TCPConn, i int) string {
	input := make([]byte, 64)
	_, err := conn.Read(input)
	if err != nil {
		log.Fatal(err)
	}
	output := cp8662str(input)
	println("Соединение " + strconv.Itoa(i+1) + ": Получили: " + output)
	return output
}

func println(this string) {
	fmt.Println(this)
	log.Println(this)
}
