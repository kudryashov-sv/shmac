package mocktrans

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func ExampleNewTransport() {
	tr, err := NewTransport()
	if err != nil {
		panic(err)
	}
	defer tr.Close()

	s := http.Server{
		Handler: handler{},
	}
	defer s.Close()
	go func() {
		err := s.Serve(tr)
		if err != nil {
			panic(err)
		}
	}()

	cli := http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return tr.Dial()
			},
		},
	}
	resp, err := cli.Get("http://0.0.0.0/test/my/mock")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:")
	fmt.Println(string(body))
	// output:
	// Status: 200 OK
	// Body:
	// Request "/test/my/mock"
	// All done!
}

type handler struct{}

func (handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Del("Date")
	rw.WriteHeader(http.StatusOK)

	fmt.Fprintf(rw, "Request %q\n", req.URL.Path)
	fmt.Fprintf(rw, "All done!")
}
