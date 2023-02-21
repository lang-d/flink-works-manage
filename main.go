package main

import "github.com/land-d/flink-works-amage/core"

func main() {
	app := core.NewApp()
	app.Parse()
	app.Run()

}
