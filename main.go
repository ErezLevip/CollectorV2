package main

import "Collector/Collector"

func main() {
	c := Collector.Collector{}.Make("")
	c.Run()
}