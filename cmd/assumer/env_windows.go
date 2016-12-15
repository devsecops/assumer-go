package main

import (
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/aws/aws-sdk-go/service/sts"
)

func execEnv(t *sts.AssumeRoleOutput) {
	creds := "set AWS_ACCESS_KEY_ID=" + *t.Credentials.AccessKeyId + "\nset AWS_SECRET_ACCESS_KEY=" + *t.Credentials.SecretAccessKey + "\nset AWS_SESSION_TOKEN=" + *t.Credentials.SessionToken
	credsBytes := []byte(creds)
	tmpfile, err := ioutil.TempFile("", "assumer")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(credsBytes); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

}

func openGui(t *sts.AssumeRoleOutput) {
	gURL := GUIURL(t)

	err := exec.Command("open", gURL).Start()
	checkErr(err)
}
