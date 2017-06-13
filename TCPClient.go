package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"io"

	"github.com/maanur/escmailer/tui"
	"github.com/maanur/testtcpmail/tcpmprobe"
)

func main() {
	var wg sync.WaitGroup
	logfile, err := os.Create("testtcpmail_" + time.Now().Format("060102-150405") + ".log")
	if err != nil {
		log.Fatal(err)
	}
	defer logfile.Close()
	q, wait, step, addr := promptVariables()
	for i := 0; i < q; i++ {
		wg.Add(1)
		go func() {
			tcpmprobe.HelloRun(i, wait, addr, io.MultiWriter(logfile, os.Stdout))
			defer wg.Done()
		}()
		time.Sleep(time.Duration(step) * time.Millisecond)
	}
	wg.Wait()
	println("Завершено")
	_ = tui.Prompt("Нажми Ентер чтоб закрыть", "")
	os.Exit(0)
}

func promptVariables() (q int, w int, s int, a string) {
	var err error
	for {
		q, err = strconv.Atoi(tui.Prompt("Количество соединений : 3", "3"))
		if err != nil {
			log.Println(err)
			println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		w, err = strconv.Atoi(tui.Prompt("Время до отключения (сек) : 20", "20"))
		if err != nil {
			log.Println(err)
			println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		s, err = strconv.Atoi(tui.Prompt("Интервал между инициализацией соединений (мс) : 100", "100"))
		if err != nil {
			log.Println(err)
			println("Некорректный ввод, повтори")
		} else {
			break
		}
	}
	for {
		a = tui.Prompt("Адрес TCPServer (в формате IP:port) : 192.168.51.122:6508", "192.168.51.122:6508")
		_, err := net.ResolveTCPAddr("tcp", a)
		if err != nil {
			log.Println(err)
			println("Некорректный адрес, повтори")
		} else {
			break
		}
	}
	return
}
