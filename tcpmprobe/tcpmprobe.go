package tcpmprobe

import (
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/text/encoding/charmap"
)

func HelloRun(wait time.Duration, addr string, ccode int, log *log.Logger) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Panic(err)
	}
	conn, err := openConn(server, log)
	if err != nil {
		log.Println("Ошибка соединения:")
		log.Panic(err)
	}
	defer closeConn(conn, log)
	msg := "Привет, TCPServer"
	err = sendMsg(msg, conn)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Отправили: " + string(msg))
	JOB := getMsg(conn, log)
	log.Println("%JOB: " + JOB)
	log.Println("Ждем " + wait.String() + " сек...")
	time.Sleep(wait)
	log.Println("Посылаю код и отключаюсь")
	err = sendMsg(strconv.Itoa(ccode), conn)
	if err != nil {
		log.Panic(err)
	}
}

func MonRun(addr string, log *log.Logger) error {
	defer log.Println("Мониторинг завершен")
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Println(addr + " : Адрес НЕ найден")
		return err
	}
	log.Println(addr + " : Адрес найден, IP " + server.IP.String())

	conn, err := openConn(server, log)
	if err != nil {
		log.Println(addr + " : Соединение НЕ установлено")
		return err
	}
	log.Println(addr + " : Соединение установлено")
	defer closeConn(conn, log)

	err = sendMsg("Мониторинг", conn)
	if err != nil {
		log.Println(addr + " : Сообщение НЕ доставлено")
	} else {
		log.Println(addr + " : Сообщение доставлено")
	}
	return err
}

func openConn(addr *net.TCPAddr, log *log.Logger) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err == nil {
		log.Println("Соединяемся...")
	}
	return conn, err
}

func closeConn(conn *net.TCPConn, log *log.Logger) {
	err := conn.Close()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Отключился...")
}

//str2cp866 Конвертирует строку str в байт-код и перекодирует в cp866 и возвращает байт-код
func str2cp866(str string) []byte {
	input := []byte(str)
	cp866 := charmap.CodePage866.NewEncoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Panic(err)
	}
	return output
}

//cp8662str Перекодирует input из cp866 и возвращает строку
func cp8662str(input []byte) string {
	cp866 := charmap.CodePage866.NewDecoder()
	output, err := cp866.Bytes(input)
	if err != nil {
		log.Panic(err)
	}
	return string(output)
}

func sendMsg(msg string, conn *net.TCPConn) error {
	_, err := conn.Write(append(str2cp866(msg), 10))
	return err

}

func getMsg(conn *net.TCPConn, log *log.Logger) string {
	input := make([]byte, 64)
	_, err := conn.Read(input)
	if err != nil {
		log.Panic(err)
	}
	output := cp8662str(input)
	log.Println("Получили: " + output)
	return output
}
