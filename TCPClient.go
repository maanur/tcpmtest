package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"io"

	"github.com/maanur/escmailer/tui"
	"github.com/maanur/tcpmtest/tcpmprobe"
)

var fq = flag.Int("q", 3, "Количество соединений")
var fccode = flag.Int("cc", 123456, "Код клиента, которым представляемся")
var fwait = flag.Duration("w", time.Second*20, "Время поддержания соединения перед разрывом")
var fstep = flag.Duration("s", time.Millisecond*100, "Интервал между соединениями")
var fmon = flag.Bool("M", false, "Имитировать мониторинг (только проверка подключения, только одно соединение)")

const usagenote = "При обычном использовании в протоколе TCPMail на сервере отобразится попытка подключения от клиента (123456 по умолчанию).\nПри запуске с флагом -M (мониторинг) на сервере отобразится 'Сервер. Мониторинг'"

func main() {
	q, wait, step, addr, ccode := promptVariables()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logfile, err := os.Create("testtcpmail_" + time.Now().Format("060102-150405") + ".log")
	defer func() {
		logger.Println("Лог записан в " + logfile.Name())
		logfile.Close()
		_ = tui.Prompt("Нажми Ентер чтоб закрыть", "")
	}()
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(io.MultiWriter(logfile, os.Stdout))
	if fmon != nil && *fmon {
		err = tcpmprobe.MonRun(addr, logger)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	var wg sync.WaitGroup
	for i := 0; i < q; i++ {
		wg.Add(1)
		go func(i int, wr io.Writer) {
			log := log.New(wr, "[Соединение "+strconv.Itoa(i+1)+"] ", log.LstdFlags)
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
				wg.Done()
			}()
			tcpmprobe.HelloRun(wait, addr, ccode, log)
		}(i, io.MultiWriter(logfile, os.Stdout))
		time.Sleep(step)
	}
	wg.Wait()
	log.Println("Завершено")
}

func promptVariables() (q int, w time.Duration, s time.Duration, a string, cc int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Ошибка вызова: %v \n\n", err)
			fmt.Println("Использование: tcpmtest addr:port")
			flag.PrintDefaults()
			fmt.Println(usagenote)
			os.Exit(0)
		}
	}()
	flag.Parse()
	switch {
	case len(flag.Args()) == 0:
		panic("Не задан адрес сервера")
	case len(flag.Args()) > 1:
		panic("Слишком много аргументов")
	default:
		_, err := net.ResolveTCPAddr("tcp", flag.Arg(0))
		if err != nil {
			panic("Некорректный адрес сервера")
		}
		a = flag.Arg(0)
	}
	if fq != nil {
		if *fq <= 0 {
			panic("Некорректное кол-во соединений")
		}
		q = *fq
	} else {
		q = 3
	}
	if fwait != nil {
		w = *fwait
	} else {
		w = time.Second * 20
	}
	if fstep != nil {
		s = *fstep
	} else {
		s = time.Millisecond * 100
	}
	if fccode != nil {
		if *fccode <= 0 {
			panic("Некорректный код клиента")
		} else {
			cc = *fccode
		}
	} else {
		cc = 123456
	}
	return
}
