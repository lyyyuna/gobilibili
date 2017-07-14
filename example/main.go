package main

import (
	"fmt"
)

func main() {
	str := "速度环"
	bytes := []byte(str)
	fmt.Println(bytes)
	message := string(bytes)
	fmt.Println(message)
}
