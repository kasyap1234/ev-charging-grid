package main

import (
	"fmt"

	"github.com/kasyap1234/ev-charging-grid/internal/config"
)


func main(){
	config,err :=config.LoadConfig()
	if err !=nil{
		panic(err)
	}
	fmt.Print(
		*config)
	
}
