package main

import (
	"compiler/src/app"
	"flag"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "i", "data/source.src", "Specify input file. Defualt is data/source.src")
	flag.Parse()

	app.App(inputFile)
}
