package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "init" {
		if err := runInit(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "pswitch init: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := runServe(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "pswitch: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `pswitch

Usage:
  pswitch init [--config PATH] [--force]
  pswitch [--config PATH] [--listen ADDR] [--mode sequential|round_robin|least_failures] [--failure-threshold N] [--cooldown DURATION] [--health-check-interval DURATION] [--health-check-timeout DURATION] [--log-color[=true|false]]`)
}
