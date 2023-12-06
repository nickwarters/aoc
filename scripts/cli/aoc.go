package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	switch os.Args[1] {
	case "run":
		run(os.Args[2:])
	case "setup":
		setup(os.Args[2:])
	default:
		fmt.Println("Expected run, setup or submit")
		os.Exit(1)
	}

}

func run(args []string) {
	cmd := flag.NewFlagSet("run", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Year(), "day")
	test := cmd.Bool("test", false, "test")
	actual := cmd.Bool("actual", false, "actual")
	input := cmd.String("input", "", "input")
	timer := cmd.Bool("time", false, "time")

	cmd.Parse(args)
	fmt.Printf("year=%d\n", *year)
	fmt.Printf("day=%d\n", *day)
	fmt.Printf("test=%v\n", *test)
	fmt.Printf("actual=%v\n", *actual)
	fmt.Printf("input=%s\n", *input)
	fmt.Printf("time=%v\n", *timer)
}

func setup(args []string) {
	cmd := flag.NewFlagSet("setup", flag.ExitOnError)
	year := cmd.Int("year", time.Now().Year(), "year")
	day := cmd.Int("day", time.Now().Year(), "day")

	cmd.Parse(args)
	fmt.Printf("year=%d\n", *year)
	fmt.Printf("day=%d\n", *day)
}
