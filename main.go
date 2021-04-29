package main

import (
	"test/proall"
)

func main() {
	proall.NewProcessor().Register().Boot()
}