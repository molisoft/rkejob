package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	goworkers "github.com/jrallison/go-workers"
	. "rkejob/config"
)

type Job struct {
	ClassName string
	FuncName  string
	Params    []interface{}
}

func job(msg *goworkers.Msg) {
	url := Config.Job.Url

	string_body, _ := msg.Args().GetIndex(0).String()
	request_body := bytes.NewBuffer([]byte(string_body))
	resp, err := http.Post(url, "application/text", request_body)
	if err != nil {
		fmt.Println("Job err: ", err)
		return
	}
	resp_body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response code: %s\n", resp.Status)
	fmt.Printf("Response result: %s\n", resp_body)
}

func init() {
	goworkers.Configure(map[string]string{
		"server":    fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port),
		"database":  fmt.Sprintf("%d", Config.Queue.Database),
		"pool":      fmt.Sprintf("%d", Config.Queue.Pool),
		"process":   "1",
		"namespace": Config.Queue.Namespace,
	})

	for _, queue := range Config.Queue.Queues {
		goworkers.Process(queue, job, 1)
	}
}

func main() {
	fmt.Println("----START----")
	goworkers.Run()
}
