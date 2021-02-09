package main

import (
	"Go-000/Week09/proto"
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var NoticeChan chan struct{}

func init()  {
	NoticeChan = make(chan struct{}, 2)
}


func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("收到client发来的数据：", msg)
	}
}

func main() {
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	g.Go(server)
	g.Go(client)

	g.Wait()
}

func server() error{
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err: ", err)
		return err
	}
	defer listen.Close()
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("accept failed, err: ", err)
				continue
			}
			go process(conn)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case sig := <-c:
			fmt.Println("server receive system signal: ", sig, " now return")
			close(NoticeChan)
			return errors.New(sig.String())
		}
	}
}

func client() error{
	time.Sleep(1*time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err: ", err)
		return err
	}
	defer conn.Close()
	go func() {
		for {
			var msg string
			fmt.Scanf("%s", &msg)
			data, err := proto.Encode(msg)
			if err != nil {
				fmt.Println("encode msg failed, err: ", err)
				return
			}
			conn.Write(data)
		}
	}()


	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case sig := <- NoticeChan:
			fmt.Println("client receive server signal: ", sig, " now return")
			return fmt.Errorf("%+v", sig)
		}
	}
}