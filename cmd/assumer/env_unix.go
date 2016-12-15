package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
)

func execEnv(t *sts.AssumeRoleOutput) {
	creds := "export AWS_ACCESS_KEY_ID=" + *t.Credentials.AccessKeyId + "\nexport AWS_SECRET_ACCESS_KEY=" + *t.Credentials.SecretAccessKey + "\nexport AWS_SESSION_TOKEN=" + *t.Credentials.SessionToken
	credsBytes := []byte(creds)

	tmpfile, err := ioutil.TempFile("", "assumer")
	checkErr(err)

	if _, err := tmpfile.Write(credsBytes); err != nil {
		checkErr(err)
	}

	if err := tmpfile.Close(); err != nil {
		checkErr(err)
	}

	color.Green("To import environment variables into the shell, execute:\n")
	fmt.Printf("source %s\n\n", tmpfile.Name())
}

func openGui(t *sts.AssumeRoleOutput) {
	gURL := GUIURL(t)

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", gURL).Start()
	case "darwin":
		err = exec.Command("open", gURL).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	checkErr(err)
}
