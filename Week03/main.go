
package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const HOST1 string = "127.0.0.1:8080"

type IndexHandler struct {
	name string
}

// CloseHandler 可触发http.Server Close
type CloseHandler struct {
	CloseChan chan error
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte(h.name))
}

func (h *CloseHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	_, _ = w.Write([]byte("closing"))

	select {
	default:
		h.CloseChan <- errors.New("api shutdown")
	case <-h.CloseChan:
	}

}

func main() {
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	mux := http.NewServeMux()
	indexHandler := &IndexHandler{name: "index1"}
	mux.Handle("/", indexHandler)
	closeHandler := &CloseHandler{}
	mux.Handle("/close", closeHandler)

	s1Ch := make(chan error, 1)
	server := &http.Server{Addr: HOST1, Handler: mux}
	closeHandler.CloseChan = s1Ch

	closeServerChan := make(chan struct{})
	closeSignalChan := make(chan struct{})

	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			return err
		}
		select {
		case err := <-s1Ch:
			fmt.Println("SSSS")
			closeSignalChan <- struct{}{}
			return err
		case <- closeServerChan:
			fmt.Println("Server: 被动退出")
			return nil
		}
	})

	// 接收信号
	sigs := make(chan os.Signal, 1)
	g.Go(func() error{
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			closeServerChan <- struct{}{}
			err := errors.New(fmt.Sprintf("receive signal: %s", sig))
			return err
		case <- closeSignalChan:
			fmt.Println("Signal: 被动退出")
			return nil
		}
	})


	fmt.Println(fmt.Sprintf("group err: %s", g.Wait()))
}
