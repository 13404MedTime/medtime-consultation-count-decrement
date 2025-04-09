package main

import (
	"fmt"
	"handler/function"
	"io"
	"os"
)

func main() {
	jsonFile, err := os.Open("request.json")
	byteValue, _ := io.ReadAll(jsonFile)

	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	fmt.Println(function.Handle(byteValue))
}
