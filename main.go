package main

import (
	"fmt"

	"github.com/Pawka/esp32-eink-smart-display/service"
)

func main() {
	fmt.Println("vim-go")
	w := service.NewWeather()
	r, e := w.Forecast("vilnius")
	if e != nil {
		panic(e)
	}

	fmt.Printf("%#v", r)
}
