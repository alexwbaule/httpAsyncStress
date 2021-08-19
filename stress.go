package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	httpasync "github.com/alexwbaule/httpasyncstress/v2/http"
)

/*
# Testes Sem parar
#  Report [Numero do Loop - Codigo HTTP - Tempo Levou - Cidade - IN - Out] - OK
#  Randomico de Cidade
#  Randomico Data: Hoje ate final do ano (de 5 a 7 dias IN Out)
#  Timeout do HTTP de 10 segundos. - OK
*/
var workers int
var etern bool
var runtimes int
var timeout int
var file string
var out_header bool
var out_body bool

//var configfile string

func init() {
	flag.IntVar(&workers, "workers", 5, "Number of concurrent workers itens")
	flag.BoolVar(&etern, "etern", false, "Run the test for a eternity in a brust of 1000 times each ! (don't worry, CTRL + C stop it)")
	flag.IntVar(&runtimes, "times", 1, "If you don't like to run this by the eternity, choose how many times you want it.")
	flag.IntVar(&timeout, "timeout", 20, "HTTP Request Timeout")
	flag.StringVar(&file, "file", "test.json", "File with all configurations")
	flag.BoolVar(&out_header, "header", false, "Print the http Response")
	flag.BoolVar(&out_body, "body", false, "Print the http Response")
}

func main() {
	flag.Parse()

	var rLimit syscall.Rlimit
	rLimit.Max = math.MaxUint64
	rLimit.Cur = math.MaxUint64
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
	fmt.Printf("Limits %v\n", rLimit)

	summary := make(chan bool)
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	go func() {
		for _ = range sig {
			fmt.Println("Wait, getting the summary....\n")
			summary <- true
		}
	}()

	fmt.Println(strings.Repeat("#", 80))

	if etern {
		runtimes = 10 * workers
		fmt.Printf("Workers\t\t-->\t[%d]\nEtern\t\t-->\t[%v]\nBrust\t\t-->\t[%d]\nHTTPTimeout\t-->\t[%d]\nOutput Header\t-->\t[%v]\nOutput Body\t-->\t[%v]\n", workers, etern, runtimes, timeout, out_header, out_body)
	} else {
		fmt.Printf("Workers\t\t-->\t[%d]\nEtern\t\t-->\t[%v]\nTimes\t\t-->\t[%d]\nHTTPTimeout\t-->\t[%d]\nOutput Header\t-->\t[%v]\nOutput Body\t-->\t[%v]\n", workers, etern, runtimes, timeout, out_header, out_body)
	}

	fmt.Println(strings.Repeat("#", 80))

	httptimeout := time.Duration(time.Duration(timeout) * time.Second)
	httpasync.WorkerAsyncHttpGets(workers, runtimes, httptimeout, etern, summary, file, out_header, out_body)
}
