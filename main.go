package main

import "github.com/denouche/go-api-skeleton/cmd"

//go:generate go run scripts/includeopenapi.go
//go:generate scripts/copy-models-to-client.sh

func main() {
	cmd.Execute()
}
