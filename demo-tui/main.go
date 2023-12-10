/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"demo-tui/cmd"
	"log"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
