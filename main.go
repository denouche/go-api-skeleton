package main

import "github.com/denouche/go-api-skeleton/cmd"

//go:generate go run scripts/includeopenapi.go

func main() {
	cmd.Execute()
}
