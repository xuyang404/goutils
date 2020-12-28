package main

import (
	"fmt"
	"time"
)

func main() {
	the_time := time.Date(2040, 1, 7, 5, 50, 4, 0, time.Local)
	str_time := time.Unix(2209499404, 0).Format("2006-01-02 15:04:05")
	fmt.Println(the_time.Unix())
	fmt.Println(str_time)
}
