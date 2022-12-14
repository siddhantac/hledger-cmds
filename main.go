package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	hledger     = "hledger"
	register    = "register"
	balance     = "balance"
	importValue = "import"
	dateQuery   = "-p"
	lastMonth   = `"..last month"`
	thisMonth   = `"..this month"`

	amexRules = "amex.rules"
	citiRules = "citibank.rules"
	ocbcRules = "ocbc.rules"
	dbsRules  = "dbs.rules"
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
	dbs      bool
	amt      int

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

// TODO
// - add account filters to bal and reg cmds
// - print unknowns
// - print txns above $50
// - print a end-of-month report

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
	importCmd.BoolVar(&dbs, "dbs", false, "use dbs.rules to import")
	importCmd.StringVar(&argInputFile, "f", "", "file to import")
	importCmd.BoolVar(&argNoDry, "no-dry", false, "disable dry-run")
	importCmd.BoolVar(&debug, "debug", false, "enable debug logs")
	importCmd.IntVar(&amt, "amt", 0, "print txns above this amount")

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
		os.Exit(1)
	}
}

func buildArgs(command string) ([]string, error) {
	if regCmd.Parsed() {
		args := []string{register, "-H"}
		return addDateQuery(args), nil
	}

	if balCmd.Parsed() {
		args := []string{balance, "-H"}
		return addDateQuery(args), nil
	}

	if importCmd.Parsed() {
		return buildArgsForImport()
	}

	return nil, fmt.Errorf("could not build any args")
}

func execute(args []string) error {
	command := strings.Join(args, " ")
	command = hledger + " " + command

	if debug {
		fmt.Println("debug: ", command)
	}

	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
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

	return args
}
