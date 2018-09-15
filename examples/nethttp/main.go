package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	globalmux := http.NewServeMux()
	globalmux.HandleFunc("/hello/", hello)

	// wrap the globalmux with the honeycomb middleware to send one event per
	// request
	log.Fatal(http.ListenAndServe("localhost:8090", globalmux))
}

func hello(w http.ResponseWriter, r *http.Request) {
	bigJob(r.Context())
	outboundCall(r.Context())
	// send our response to the caller
	io.WriteString(w, fmt.Sprintf("Hello world!\n"))
}

// bigJob is going to take a long time and do lots of interesting work. It
// should get its own span.
func bigJob(ctx context.Context) {
	// bigJob will take ~300ms
	sleepTime := math.Abs(200.0 + (rand.NormFloat64()*50 + 100))
	time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}

// outboundCall demonstrates wrapping an outbound HTTP client
func outboundCall(ctx context.Context) {
	// let's make an outbound HTTP call
	req, _ := http.NewRequest(http.MethodGet, "http://scooterlabs.com/echo.json", strings.NewReader(""))
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
		bod, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf(string(bod))
	}
}
