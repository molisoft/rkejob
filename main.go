package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	goworkers "github.com/jrallison/go-workers"
	libcron "github.com/molisoft/cron"
	. "rkejob/config"
)

func job(msg *goworkers.Msg) {
	fmt.Println("------------job--------", msg.OriginalJson())
	url := Config.Job.Url

	string_body, _ := msg.Args().GetIndex(0).String()
	request_body := bytes.NewReader([]byte(string_body))
	resp, err := http.Post(url, "application/text", request_body)
	if err != nil {
		fmt.Println("Job err: ", err)
		return
	}
	resp_body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response code: %s\n", resp.Status)
	fmt.Printf("Response result: %s\n", resp_body)
}

func cron(item CronItemConfig) {
	fmt.Println(time.Now(), "Runing cron ", item)
	request_body := bytes.NewReader([]byte(item.Name))
	_, err := http.Post(item.Url, "application/text", request_body)
	if err != nil {
		fmt.Println(item.Name, " Job err:", err)
		return
	}
	fmt.Println(time.Now(), "Done cron ", item)
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

	init_cron()
}

func init_cron() {
	c := libcron.New()
	for i, item := range Config.Crons {
		fmt.Println("-----", item)
		err := c.AddFunc(item.Spec, func(e *libcron.Entry) {
			cron(Config.Crons[e.Param.(int)])
		}, i)
		if err != nil {
			fmt.Println("cron add job err ", err)
		}
	}
	c.Start()
}

func main() {
	fmt.Println("----START----")
	goworkers.Run()
}
