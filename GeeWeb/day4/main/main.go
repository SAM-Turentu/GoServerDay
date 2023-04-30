package main

import (
	"context"
	"geerpc"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func startServer_v1(addr chan string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error: ", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

//func main() {
//	addr := make(chan string)
//	go startServer(addr)
//
//	conn, _ := net.Dial("tcp", <-addr)
//	defer func() { _ = conn.Close() }()
//
//	time.Sleep(time.Second)
//
//	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
//	cc := codec.NewGobCodec(conn)
//
//	for i := 0; i < 5; i++ {
//		h := &codec.Header{
//			ServiceMethod: "Foo.Sum",
//			Seq:           uint64(i),
//		}
//		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
//		_ = cc.ReadHeader(h)
//		var reply string
//		_ = cc.ReadBody(&reply)
//		log.Println("reply: ", reply)
//	}
//
//}

//func main() {
//	log.SetFlags(0)
//	addr := make(chan string)
//	go startServer(addr)
//	client, _ := geerpc.Dial("tcp", <-addr)
//	defer func() { _ = client.Close() }()
//
//	time.Sleep(time.Second)
//	// 发送请求， 接受响应
//	var wg sync.WaitGroup
//	for i := 0; i < 5; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			args := fmt.Sprintf("geerpc req %d", i)
//			var reply string
//			if err := client.Call("Foo.Sum", args, &reply); err != nil {
//				log.Fatal("call Foo.Sum error: ", err)
//			}
//			log.Println("reply: ", reply)
//		}(i)
//	}
//	wg.Wait()
//}

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

//func startServer_v2(addr chan string) {
//	var foo Foo
//	if err := geerpc.Register(&foo); err != nil {
//		log.Fatal("register error: ", err)
//	}
//	l, err := net.Listen("tcp", ":0")
//	if err != nil {
//		log.Fatal("network error: ", err)
//	}
//	log.Println("start rpc server on", l.Addr())
//	addr <- l.Addr().String()
//	geerpc.Accept(l)
//}
//
//func main() {
//	log.SetFlags(0)
//	addr := make(chan string)
//	//go startServer_v2(addr)
//	go startServer(addr)
//	client, _ := geerpc.Dial("tcp", <-addr)
//	defer func() { _ = client.Close() }()
//
//	time.Sleep(time.Second)
//	// 发送请求， 接受响应
//	var wg sync.WaitGroup
//	for i := 0; i < 5; i++ {
//		wg.Add(1)
//		go func(i int) {
//			defer wg.Done()
//			args := &Args{Num1: i, Num2: i * i}
//			var reply int
//			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
//				log.Fatal("call Foo.Sum error: ", err)
//			}
//			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
//		}(i)
//	}
//	wg.Wait()
//}

func startServer(addrCh chan string) {
	var foo Foo
	l, _ := net.Listen("tcp", ":9999")
	_ = geerpc.Register(&foo)
	geerpc.HandleHTTP()
	addrCh <- l.Addr().String()
	_ = http.Serve(l, nil)
}

func call(addrCh chan string) {
	client, _ := geerpc.DialHTTP("tcp", <-addrCh)
	defer func() { _ = client.Close() }()
	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{
				Num1: 1,
				Num2: i * i,
			}
			var reply int
			if err := client.Call(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error: ", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func main() {
	log.SetFlags(0)
	ch := make(chan string)
	go call(ch)
	startServer(ch)
}
