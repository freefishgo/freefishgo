package main

import "time"

func main() {
	app := NewFreeFish()
	app.Run()
	time.Sleep(time.Hour)
}
