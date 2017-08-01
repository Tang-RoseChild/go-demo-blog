package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	postData, _ := json.Marshal(map[string]interface{}{
		"no":       "86",
		"mobile":   "13817171378",
		"password": "demo1",
	})
	// // resp, err := http.Post("http://127.0.0.1:9995/api/login", "application/json", bytes.NewBuffer(postData))
	// resp, err := http.Post("http://127.0.0.1:9995/api/register", "application/json", bytes.NewBuffer(postData))
	// fmt.Println("resp status code , ", resp.StatusCode, " err > ", err)
	// defer resp.Body.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// data, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(data))
	worker := make(chan struct{}, 100)

	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func(worker chan struct{}) {
			for {
				_, ok := <-worker
				if !ok {
					return
				}
				resp, err := http.Post("http://127.0.0.1:9995/api/login", "application/json", bytes.NewBuffer(postData))
				if err != nil {
					fmt.Println(err)
				}
				if resp != nil && resp.StatusCode != http.StatusOK {
					data, _ := ioutil.ReadAll(resp.Body)
					fmt.Println(string(data))
				}

				if resp != nil && resp.Body != nil {
					resp.Body.Close()
				}
				// fmt.Println("post")
			}
		}(worker)

		go func(worker chan struct{}) {
			for {
				select {
				case <-done:
					return
				default:
					// fmt.Println("default>>>")
					worker <- struct{}{}
				}
			}
		}(worker)
	}
	fmt.Println("start >>>>> ")
	time.Sleep(60 * time.Second)
	done <- struct{}{}
	fmt.Println("done >>>>>>>>")
}
