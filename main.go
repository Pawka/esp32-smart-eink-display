package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Pawka/esp32-eink-smart-display/service"
)

const serverAddr string = ":3000"

func main() {
	http.HandleFunc("/", serve)
	fmt.Printf("Serving on %s\n", serverAddr)
	http.ListenAndServe(serverAddr, nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := service.NewContent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(c)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
