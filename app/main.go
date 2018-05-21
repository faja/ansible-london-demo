package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

var world string
var port int = 12345

func HelloServer(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello %s world!\n", world)
}

func main() {
	log.Printf("starting app server\n")
	viper.SetConfigName("config")
	viper.AddConfigPath("/demo/app/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	world = viper.GetString("world.type")

	log.Printf("server started and listening on :%d\n", port)
	http.HandleFunc("/hello", HelloServer)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
