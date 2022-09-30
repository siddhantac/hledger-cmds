package main

import "errors"

func buildArgsForImport() ([]string, error) {
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
		args = append(args, "|", hledger, "-f-", "-I", "reg")
	}

	return args, nil
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
