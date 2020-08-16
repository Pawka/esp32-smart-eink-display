package main

import (
	"fmt"

	"github.com/Pawka/esp32-eink-smart-display/gateway/meteolt"
)

func main() {
	fmt.Println("vim-go")
	meteo := meteolt.New()
	r, e := meteo.Forecast("vilnius")
	if e != nil {

		panic(e)
	}

	fmt.Printf("%#v", r)
}
