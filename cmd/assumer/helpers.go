package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

// mfa prompts user for MFA token
func mfa(token *string) string {
	fmt.Printf("\nEnter MFA: ")
	reader := bufio.NewReader(os.Stdin)
	*token, _ = reader.ReadString('\n')

	if match, _ := regexp.MatchString(`\d{6}`, *token); !match {
		color.Red("Invalid MFA. Please, enter a 6-digit MFA code.")
		mfa(token)
	}

	carriageReturn := regexp.MustCompile("\\n")
	*token = carriageReturn.ReplaceAllString(*token, "")

	return *token
}

func checkErr(err error) {
	if err != nil {
		color.Red("ERROR: %s", err.Error())
		os.Exit(1)
	}
}
