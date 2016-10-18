package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/pmbenjamin/assumer-go"
	"github.com/spf13/viper"
)

// mfa prompts user for MFA token
func mfa(token *string) string {
	fmt.Println("Enter MFA: ")
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

func getControlPlane(p *assumer.Plane) *assumer.Plane {
	if p.AccountNumber == "" {
		p.AccountNumber = viper.GetString("control.account")
	}
	if p.RoleArn == "" {
		p.RoleArn = "arn:aws:iam::" + viper.GetString("control.account") + ":role/" + viper.GetString("control.role")
	} else {
		p.RoleArn = "arn:aws:iam::" + viper.GetString("control.account") + ":role/" + p.RoleArn
	}
	if p.Region == "" {
		p.Region = viper.GetString("control.region")
	}
	return p
}

func getTargetPlane(p *assumer.Plane) *assumer.Plane {
	if p.AccountNumber == "" {
		p.AccountNumber = viper.GetString("default.account")
	}

	if p.RoleArn == "" {
		p.RoleArn = "arn:aws:iam::" + viper.GetString("default.account") + ":role/" + viper.GetString("default.role")
	} else {
		p.RoleArn = "arn:aws:iam::" + p.AccountNumber + ":role/" + p.RoleArn
	}

	if p.Region == "" {
		p.Region = viper.GetString("default.region")
	}
	return p
}

func checkErr(err error) {
	if err != nil {
		color.Red("ERROR: %s", err.Error())
		os.Exit(1)
	}
}

func printVersion(version string) {
	fmt.Printf("VERSION: %s\n\n", semver)
}
