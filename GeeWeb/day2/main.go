package main

import (
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "234",
	"Jack": "883",
	"SAM":  "999",
}

//type server int
//
//func (h *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	log.Println(r.URL.Path)
//	w.Write([]byte("hello world"))
//}

func main() {
	//var s server
	//http.ListenAndServe("0.0.0.0:8082", &s)

	geecache.NewGroup("score", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "0.0.0.0:8082"
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

}
