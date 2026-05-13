package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	next := flag.Int("next", 3, "The next n times the cron expression is programmed to run")
	json := flag.Bool("json", false, "Whether to output the results in json format")

	flag.Parse()

	cronExpression := flag.Arg(0)
	if len(cronExpression) == 0 {
		fmt.Fprintln(os.Stderr, "gocron: missing cron expression")
		os.Exit(1)
	}

	fmt.Printf("expression: %s\n", cronExpression)
	fmt.Printf("next: %d\n", *next)
	fmt.Printf("json: %t\n", *json)
}
