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
	"github.com/maanur/testtcpmail/tcpmprobe"
)

var fq = flag.Int("q", 3, "Количество соединений")
var intervalFlag interval
var fwait = flag.Duration("w", time.Second*20, "Время поддержания соединения перед разрывом")
var fstep = flag.Duration("s", time.Millisecond*100, "Интервал между соединениями")

type interval time.Duration

func (i *interval) String() string {
	return fmt.Sprint(*i)
}

func (i *interval) Set(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*i = interval(duration)
	return nil
}

func main() {
	q, wait, step, addr := promptVariables()
	logger := log.New(os.Stdout, "", log.LstdFlags)
	var wg sync.WaitGroup
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
			tcpmprobe.HelloRun(wait, addr, log)
		}(i, io.MultiWriter(logfile, os.Stdout))
		time.Sleep(step)
	}
	wg.Wait()
	log.Println("Завершено")
}

func promptVariables() (q int, w time.Duration, s time.Duration, a string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println("Использование: tcpmtest addr:port")
			flag.PrintDefaults()
			os.Exit(0)
		}
	}()
	flag.Parse()
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
	/*
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
	*/
	return
}
