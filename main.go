package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	host     string
	port     string
	endpoint string
	interval time.Duration

	currentData *interface{}
)

func init() {
	flag.StringVar(&host, "host", "localhost", "hostname to listen on")
	flag.StringVar(&port, "port", "3000", "port to listen on")
	flag.StringVar(&endpoint, "endpoint", "/data", "default endpoint to serve data on")
	flag.DurationVar(&interval, "interval", 5*time.Second, "interval every which the served data will be updated")
}

func main() {
	flag.Parse()

	finished := make(chan struct{})

	go updateJSON(finished)

	addr := fmt.Sprintf("%s:%s", host, port)
	http.HandleFunc(endpoint, handle)
	http.ListenAndServe(addr, nil)

	<-finished

	fmt.Println("fileserver: served all files, shutdown")
}

func readFile(which int) (*interface{}, error) {
	path := fmt.Sprint("data/", which, "/data.json")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New(fmt.Sprint("fileserver: error reading a json file:", err.Error()))
	}
	var matchData interface{}
	err = json.Unmarshal(data, &matchData)
	if err != nil {
		return nil, err
	}

	return &matchData, nil
}

func updateJSON(finished chan struct{}) {
	i := 1
	for {
		fmt.Printf("fileserver: serving new data after %.f seconds\n", interval.Seconds())
		var err error
		currentData, err = readFile(i)
		if err != nil {
			fmt.Println("fileserver: served all files, you can exit now")
			break
		}

		i++
		time.Sleep(interval)
	}

	close(finished)
}

func handle(writer http.ResponseWriter, req *http.Request) {
	j, err := json.Marshal(currentData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(j)
}
