package assumer

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/viper"
)

// TargetPlane represents the AWS Target Plane Account
type TargetPlane struct {
	Plane
}

// Assume assumes a role in the target account and returns the assumed role creds
func (t *TargetPlane) Assume(c *sts.AssumeRoleOutput) (*sts.AssumeRoleOutput, error) {

	if err := t.GetDefaults(); err != nil {
		return nil, errors.New("Error: Could not get Target Plane defaults")
	}

	targetParams := &sts.AssumeRoleInput{
		RoleArn:         aws.String(t.Plane.RoleArn),
		RoleSessionName: aws.String(fmt.Sprintf("%s_AssumedRole", os.Getenv("USER"))),
	}

	controlCreds := credentials.NewStaticCredentials(*c.Credentials.AccessKeyId, *c.Credentials.SecretAccessKey, *c.Credentials.SessionToken)

	stsClient := sts.New(session.New(), aws.NewConfig().WithRegion(t.Region).WithCredentials(controlCreds))
	resp, err := stsClient.AssumeRole(targetParams)
	if err != nil {
		return nil, errors.New("Error: Could not assume Target Plane Role. " + err.Error())
	}

	return resp, nil
}

// GetDefaults will get TargetPlane default values from assumer config file
func (t *TargetPlane) GetDefaults() error {
	if t.AccountNumber == "" {
		t.AccountNumber = viper.GetString("default.account")
	}

	if t.RoleArn == "" {
		t.RoleArn = "arn:aws:iam::" + viper.GetString("default.account") + ":role/" + viper.GetString("default.role")
	} else {
		t.RoleArn = "arn:aws:iam::" + t.AccountNumber + ":role/" + t.RoleArn
	}

	if t.Region == "" {
		t.Region = viper.GetString("default.region")
	}

	return nil
}
