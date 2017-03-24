package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/maanur/escmailer/tui"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	var wg sync.WaitGroup
	q, wait, step, addr := promptVariables()
	for i := 0; i < q; i++ {
		wg.Add(1)
		go func() {
			sampleRun(i, wait, addr)
			defer wg.Done()
		}()
		time.Sleep(time.Duration(step) * time.Millisecond)
	}
	wg.Wait()
	fmt.Println("Завершено")
	_ = tui.Prompt("Нажми Ентер чтоб закрыть", "")
	os.Exit(0)
}

func promptVariables() (q int, w int, s int, a string) {
	var err error
	for {
		q, err = strconv.Atoi(tui.Prompt("Количество соединений : 3", "3"))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		w, err = strconv.Atoi(tui.Prompt("Время до отключения (сек) : 20", "20"))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		s, err = strconv.Atoi(tui.Prompt("Интервал между инициализацией соединений (мс) : 100", "100"))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		a = tui.Prompt("Адрес TCPServer (в формате IP:port) : 192.168.51.122:6508", "192.168.51.122:6508")
		_, err := net.ResolveTCPAddr("tcp", a)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Некорректный адрес, повтори")
		} else {
			break
		}
	}
	return
}

func sampleRun(i int, wait int, addr string) {
	server, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := openConn(server)
	if err != nil {
		fmt.Println("Соединение " + strconv.Itoa(i+1) + " Ошибка соединения:")
		log.Fatal(err)
	}
	sendMsg("Привет, TCPServer", conn, i)
	JOB := getMsg(conn, i)
	fmt.Println("Соединение " + strconv.Itoa(i+1) + ": %JOB: " + JOB)
	fmt.Println("Соединение " + strconv.Itoa(i+1) + ": Ждем " + strconv.Itoa(wait) + " сек...")
	time.Sleep(time.Duration(wait) * time.Second)
	fmt.Println("Соединение " + strconv.Itoa(i+1) + ": Посылаю код и отключаюсь")
	sendMsg(strconv.Itoa(180020+i*20), conn, i)
	closeConn(conn)
}

func openConn(addr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp4", nil, addr)
	if err == nil {
		fmt.Println("Соединяемся...")
	}
	return conn, err
}

func closeConn(conn *net.TCPConn) {
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Отключился...")
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
	fmt.Println("Соединение " + strconv.Itoa(i+1) + ": Отправили: " + string(msg))
}

func getMsg(conn *net.TCPConn, i int) string {
	input := make([]byte, 64)
	_, err := conn.Read(input)
	if err != nil {
		log.Fatal(err)
	}
	output := cp8662str(input)
	fmt.Println("Соединение " + strconv.Itoa(i+1) + ": Получили: " + output)
	return output
}
