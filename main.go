package main

import (
	_ "github.com/lib/pq" // postgres driver
	"marketflow/cmd"
)

func main() {
	cmd.RunApp()
}
