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

	regCmd    *flag.FlagSet
	balCmd    *flag.FlagSet
	importCmd *flag.FlagSet

	dateLastMonth bool
	dateThisMonth bool

	amex     bool
	citibank bool
	ocbc     bool

	argDateQuery string
	argInputFile string
	argNoDry     bool
)

var usage = `Usage of hledger-cmds:
  ./cmd <cmd> [options]

  cmd: one of [reg, bal, import]
  For more help, use
     ./cmd <cmd> -h
`

func main() {
	regCmd = flag.NewFlagSet("reg", flag.ExitOnError)
	regCmd.BoolVar(&dateLastMonth, "last", false, "filter only last month")
	regCmd.BoolVar(&dateThisMonth, "this", false, "filter only this month")
	regCmd.StringVar(&argDateQuery, "date", "", "custom date query")
	regCmd.StringVar(&argDateQuery, "d", "", "custom date query (shorthand)")
	regCmd.BoolVar(&debug, "debug", false, "enable debug logs")

	balCmd = flag.NewFlagSet("bal", flag.ExitOnError)
	balCmd.BoolVar(&dateLastMonth, "last", false, "filter only last month")
	balCmd.BoolVar(&dateThisMonth, "this", false, "filter only this month")
	balCmd.StringVar(&argDateQuery, "date", "", "custom date query")
	balCmd.StringVar(&argDateQuery, "d", "", "custom date query (shorthand)")
	balCmd.BoolVar(&debug, "debug", false, "enable debug logs")

	importCmd = flag.NewFlagSet("import", flag.ExitOnError)
	importCmd.BoolVar(&amex, "amex", false, "use amex.rules to import")
	importCmd.BoolVar(&citibank, "citi", false, "use citibank.rules to import")
	importCmd.BoolVar(&ocbc, "ocbc", false, "use ocbc.rules to import")
	importCmd.StringVar(&argInputFile, "f", "", "file to import")
	importCmd.BoolVar(&argNoDry, "no-dry", false, "disable dry-run")
	importCmd.BoolVar(&debug, "debug", false, "enable debug logs")

	flag.Usage = func() {
		fmt.Printf(usage)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "reg":
		if err := regCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	case "bal":
		if err := balCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	case "import":
		if err := importCmd.Parse(os.Args[2:]); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(0)
	}

	args, err := buildArgs(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if err := execute(args); err != nil {
		fmt.Printf("error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
}

func buildArgs(command string) ([]string, error) {
	if regCmd.Parsed() {
		args := make([]string, 0)

		if command == "reg" {
			args = []string{register}
		} else if command == "bal" {
			args = []string{balance}
		}

		return addDateQuery(args), nil
	}

	if importCmd.Parsed() {
		if argInputFile == "" {
			return nil, errors.New("filename required")
		}

		rulesFile := getRulesFile()
		if rulesFile == "" {
			return nil, errors.New("could not set rules-file")
		}

		args := []string{importValue, argInputFile, "--rules-file", rulesFile}

		// since this param is false by default,
		// app default behaviour will be to use dry-run
		if !argNoDry {
			args = append(args, "--dry-run")
		}

		return args, nil
	}

	return nil, fmt.Errorf("could not build any args")
}

func execute(args []string) error {
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
