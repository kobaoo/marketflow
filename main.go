package main

import (
	"marketflow/cmd"

	_ "github.com/lib/pq" // postgres driver
)

func main() {
	cmd.RunApp()
}
