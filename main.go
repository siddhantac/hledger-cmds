package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	hledger   = "hledger"
	register  = "register"
	balance   = "balance"
	dateQuery = "-p"
	lastMonth = `"last month"`
	thisMonth = `"this month"`
)

var (
	debug bool

	cmdReg        bool
	cmdBal        bool
	dateLastMonth bool
	dateThisMonth bool

	argDateQuery string
)

func main() {

	flag.BoolVar(&cmdReg, "reg", false, "register")
	flag.BoolVar(&cmdBal, "bal", false, "balance")

	flag.BoolVar(&dateLastMonth, "last", false, "filter only last month")
	flag.BoolVar(&dateThisMonth, "this", false, "filter only this month")
	flag.StringVar(&argDateQuery, "date", "", "custom date query")
	flag.StringVar(&argDateQuery, "d", "", "custom date query (shorthand)")

	flag.BoolVar(&debug, "debug", false, "enable debug logs")
	flag.Parse()

	err := execute()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func execute() error {
	args := make([]string, 0)

	if cmdReg {
		args = append(args, register)
	}

	if cmdBal {
		args = append(args, balance)
	}

	if len(args) < 1 {
		return errors.New("arg required")
	}

	args = addDateQuery(args)

	if debug {
		fmt.Println("debug: ", args)
	}
	out, err := exec.Command(hledger, args...).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

func addDateQuery(args []string) []string {
	if argDateQuery != "" {
		return append(args, dateQuery, argDateQuery)
	}

	if dateLastMonth {
		return append(args, dateQuery, lastMonth)
	}

	if dateThisMonth {
		return append(args, dateQuery, thisMonth)
	}

	return nil
}
