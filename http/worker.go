package httpasync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alexwbaule/httpasyncstress/v2/config"
	"github.com/alexwbaule/httpasyncstress/v2/structs"
)

func WorkerAsyncHttpGets(workers int, runtimes int, httptimeout time.Duration, etern bool, summary <-chan bool, file string, out_header bool, out_body bool) {
	jobs := make(chan structs.HttpRequest, runtimes)
	results := make(chan structs.HttpResult, runtimes)
	loop := 0

	/* ESTATISTICAS */
	errors := 0
	dones := 0
	found_dones := 0
	found_errors := 0
	var biggest int64 = 0
	var minor int64 = math.MaxInt64
	var middle int64 = 0
	jobtimes := 0
	resulttimes := 0
	var bodyString string
	/* ESTATISTICAS */

	httpData := config.GetFile(file)

	var findregex = regexp.MustCompile(httpData.Grep)

	startTime := time.Now()

	for tt := 0; tt < workers; tt++ {
		go worker(tt, jobs, results)
	}

	for j := 0; j < runtimes; j++ {
		data := structs.HttpRequest{httpData, j, httptimeout}
		jobs <- data
		loop = j
		jobtimes++
	}

	if !etern {
		close(jobs)
	}

	for j := 0; j < runtimes; j++ {
		found := false
		resulttimes++
		r := <-results
		if r.Response != nil {
			bodyBytes, err2 := ioutil.ReadAll(r.Response.Body)
			if err2 == nil {
				bodyString = string(bodyBytes)
				found = findregex.MatchString(bodyString)
			}
			fmt.Printf("Count=%d -- Worker=%d -- Status=%s -- QueryString=[%s] ", r.Count, r.Worker, r.Response.Status, r.Response.Request.URL.RawQuery)
			for v := range r.Soap.Rpl {
				fmt.Printf("-- %s=%s ", r.Soap.Rpl[v].Key, r.Soap.Rpl[v].Value)
			}
			fmt.Printf("-- Took=%v -- Regex=%v -- ContentLength=%d\n", r.Tget, found, r.Response.ContentLength)
			if out_header {
				for name, headers := range r.Response.Header {
					for _, h := range headers {
						fmt.Printf("[%v: %v]\n", name, h)
					}
				}
			}
			if out_body && err2 == nil {
				fmt.Printf("[Response]:\n[%s]\n", bodyString)
			}

			biggest = int64(math.Max(float64(biggest), float64(r.Tget)))
			minor = int64(math.Min(float64(minor), float64(r.Tget)))
			middle += int64(r.Tget)
			if r.Response.StatusCode == 200 {
				dones++
			} else {
				errors++
			}
			if found {
				found_dones++
			} else {
				found_errors++
			}
		} else {
			fmt.Printf("Count=%d -- Worker=%d ", r.Count, r.Worker)
			for v := range r.Soap.Rpl {
				fmt.Printf("-- %s=%s ", r.Soap.Rpl[v].Key, r.Soap.Rpl[v].Value)
			}
			fmt.Printf("-- Took=%v -- Error=%s\n", r.Tget, r.Err)
			errors++
		}
		select {
		case _ = <-summary:
			goto SUMMARY
		default:
		}

		if etern {
			data := structs.HttpRequest{httpData, loop, httptimeout}
			jobs <- data
			loop++
			runtimes++
			jobtimes++
		}
	}
SUMMARY:
	Summary(jobtimes, resulttimes, errors, dones, found_dones, found_errors, time.Duration(biggest), time.Duration(minor), time.Now().Sub(startTime), time.Duration(middle/int64(resulttimes)))
}

func Summary(sended int, received int, errors int, oks int, fdone int, ferror int, biggest time.Duration, minor time.Duration, running time.Duration, avg time.Duration) {
	fmt.Println(strings.Repeat("#", 80))
	fmt.Printf("%-40s %39d\n", "Requests Queued:", sended)
	fmt.Printf("%-40s %39d\n", "Requests Executed:", received)
	fmt.Printf("%-40s %39d\n", "Results with Error:", errors)
	fmt.Printf("%-40s %39d\n", "Results OK:", oks)
	fmt.Printf("%-40s %39d\n", "Grep with Error:", ferror)
	fmt.Printf("%-40s %39d\n", "Grep OK:", fdone)
	fmt.Printf("%-40s %39v\n", "Bigggest Response Time:", biggest)
	if minor == math.MaxInt64 {
		minor = 0
	}
	fmt.Printf("%-40s %39v\n", "Minor Response Time:", minor)
	fmt.Printf("%-40s %39v\n", "Average Response Time:", avg)
	fmt.Printf("%-40s %39v\n", "Running Time:", running)
	fmt.Println(strings.Repeat("#", 80))

	os.Exit(0)
}

func worker(worker int, jobs <-chan structs.HttpRequest, results chan<- structs.HttpResult) {

	for item := range jobs {
		start := time.Now()
		fullUrl := config.DoQueryString(item.HttpData)
		body, city := config.DoBody(item.HttpData)
		client := &http.Client{
			Timeout: item.HttpTimeout,
		}
		tr := &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		}
		client = &http.Client{Transport: tr}

		req, err := http.NewRequest(item.HttpData.Type, fullUrl, bytes.NewBuffer([]byte(body)))
		if err != nil {
			fmt.Println(err)
		}
		for header := range item.HttpData.Headers {
			req.Header.Add(item.HttpData.Headers[header].Name, item.HttpData.Headers[header].Value)
		}
		req.Close = true
		resp, err := client.Do(req)

		end := time.Now().Sub(start)
		results <- structs.HttpResult{resp, err, start, end, worker, item.Count, city}
	}
}
