package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
)

// GUIURL generates the GUI console URL
func GUIURL(t *sts.AssumeRoleOutput) string {
	signinURL := "https://signin.aws.amazon.com/federation"
	issuer := "assumer"
	consoleURL := "https://console.aws.amazon.com/"

	// anonymous struct for marshalling data into JSON
	session := struct {
		SessionID    string `json:"sessionId"`
		SessionKey   string `json:"sessionKey"`
		SessionToken string `json:"sessionToken"`
	}{
		*t.Credentials.AccessKeyId,
		*t.Credentials.SecretAccessKey,
		*t.Credentials.SessionToken,
	}

	// compose URL for getting signin token
	sessionJSON, err := json.Marshal(session)
	checkErr(err)
	signinTokenURL := fmt.Sprintf("%s?Action=getSigninToken&SessionType=json&Session=%s", signinURL, url.QueryEscape(string(sessionJSON)))

	if debug {
		fmt.Printf("Auth Request URL: %s\n\n", signinTokenURL)
	}

	// get signin token
	res, err := http.Get(signinTokenURL)
	checkErr(err)
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	// anonymous struct for unmarshalling the signin token
	token := struct {
		SigninToken string `json:"SigninToken"`
	}{
		"",
	}
	json.Unmarshal(resData, &token)

	if debug {
		fmt.Printf("Token: %s\n\n", token.SigninToken)
	}

	loginURL := fmt.Sprintf("%s?Action=login&SigninToken=%s&Issuer=%s&Destination=%s", signinURL, token.SigninToken, issuer, consoleURL)
	fmt.Printf("Logging into the browser console via: \n%s\n", loginURL)
	return loginURL
}

// mfa prompts user for MFA token
func mfa(token *string) string {
	fmt.Printf("\nEnter MFA: ")
	reader := bufio.NewReader(os.Stdin)
	*token, _ = reader.ReadString('\n')

	if match, _ := regexp.MatchString(`\d{6}`, *token); !match {
		color.Red("Invalid MFA. Please enter a 6-digit MFA code.")
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
