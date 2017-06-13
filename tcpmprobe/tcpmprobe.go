package tcpmprobe

import (
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func HelloRun(i int, wait int, addr string, logger io.Writer) {
	log.SetOutput(logger)
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := openConn(server)
	if err != nil {
		log.Println("Соединение " + strconv.Itoa(i+1) + " Ошибка соединения:")
		log.Fatal(err)
	}
	defer closeConn(conn)
	msg := "Привет, TCPServer"
	err = sendMsg(msg, conn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Соединение " + strconv.Itoa(i+1) + ": Отправили: " + string(msg))
	JOB := getMsg(conn, i)
	log.Println("Соединение " + strconv.Itoa(i+1) + ": %JOB: " + JOB)
	log.Println("Соединение " + strconv.Itoa(i+1) + ": Ждем " + strconv.Itoa(wait) + " сек...")
	time.Sleep(time.Duration(wait) * time.Second)
	log.Println("Соединение " + strconv.Itoa(i+1) + ": Посылаю код и отключаюсь")
	_ = sendMsg(strconv.Itoa(180020+i*20), conn)
}

func MonRun(addr string) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := openConn(server)
	if err != nil {
		log.Println("Соединение " + strconv.Itoa(1) + " Ошибка соединения:")
		log.Fatal(err)
	}
	defer closeConn(conn)
	_ = sendMsg("Мониторинг", conn)
}
func openConn(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err == nil {
		log.Println("Соединяемся...")
	}
	return conn, err
}

func closeConn(conn *net.TCPConn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Отключился...")
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

func sendMsg(msg string, conn *net.TCPConn) error {
	_, err := conn.Write(append(str2cp866(msg), 10))
	return err

}

func getMsg(conn *net.TCPConn, i int) string {
	input := make([]byte, 64)
	_, err := conn.Read(input)
	if err != nil {
		log.Fatal(err)
	}
	output := cp8662str(input)
	log.Println("Соединение " + strconv.Itoa(i+1) + ": Получили: " + output)
	return output
}
