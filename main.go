package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Sims 4 Mod Manager...")
	app := setupApp()
	app.Run()
}