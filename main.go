package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	hledger     = "hledger"
	register    = "register"
	balance     = "balance"
	importValue = "import"
	dateQuery   = "-p"
	lastMonth   = `"last month"`
	thisMonth   = `"this month"`

	amexRules = "amex.rules"
	citiRules = "citibank.rules"
	ocbcRules = "ocbc.rules"
)

var (
	debug bool

	cmdReg    bool
	cmdBal    bool
	cmdImport bool

	dateLastMonth bool
	dateThisMonth bool

	amex     bool
	citibank bool
	ocbc     bool

	argDateQuery string
	argInputFile string
	argNoDry     bool
)

func main() {

	flag.BoolVar(&cmdReg, "reg", false, "register")
	flag.BoolVar(&cmdBal, "bal", false, "balance")
	flag.BoolVar(&cmdImport, "import", false, "import")

	flag.BoolVar(&dateLastMonth, "last", false, "filter only last month")
	flag.BoolVar(&dateThisMonth, "this", false, "filter only this month")
	flag.StringVar(&argDateQuery, "date", "", "custom date query")
	flag.StringVar(&argDateQuery, "d", "", "custom date query (shorthand)")

	flag.BoolVar(&amex, "amex", false, "use amex.rules to import")
	flag.BoolVar(&citibank, "citi", false, "use citibank.rules to import")
	flag.BoolVar(&ocbc, "ocbc", false, "use ocbc.rules to import")
	flag.StringVar(&argInputFile, "f", "", "file to import")
	flag.BoolVar(&argNoDry, "no-dry", false, "disable dry-run")

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
		args = addDateQuery(args)
	}

	if cmdBal {
		args = append(args, balance)
		args = addDateQuery(args)
	}

	if cmdImport {
		if argDateQuery != "" {
			return errors.New("date query with import not understood")
		}

		if argInputFile == "" {
			return errors.New("filename required")
		}

		rulesFile := getRulesFile()
		if rulesFile == "" {
			return errors.New("could not set rules-file")
		}

		args = append(args, importValue, argInputFile, "--rules-file", rulesFile)

		// since this param is false by default,
		// app default behaviour will be to use dry-run
		if !argNoDry {
			args = append(args, "--dry-run")
		}
	}

	if len(args) < 1 {
		return errors.New("1 arg required")
	}

	if debug {
		fmt.Println("debug: ", args)
	}

	out, err := exec.Command(hledger, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}

	fmt.Println(string(out))
	return nil
}

func getRulesFile() string {
	if amex {
		return amexRules
	}

	if citibank {
		return citiRules
	}

	if ocbc {
		return ocbcRules
	}

	return ""
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
