package assumer

import (
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/viper"
)

// ControlPlane represents the AWS Control Plane Account
type ControlPlane struct {
	Plane
	SerialNumber string
	MfaToken     string
}

// Assume assumes a Control Plane Role and returns the assumed role credentials
func (c *ControlPlane) Assume() (*sts.AssumeRoleOutput, error) {

	if err := CheckMfa(c.MfaToken); err != nil {
		return nil, err
	}

	c.GetDefaults()

	svc := sts.New(session.New(), aws.NewConfig().WithRegion(c.Region))

	controlParams := &sts.AssumeRoleInput{
		RoleArn:         aws.String(c.RoleArn),
		RoleSessionName: aws.String("AssumedRole"),
		SerialNumber:    aws.String(c.SerialNumber),
		TokenCode:       aws.String(c.MfaToken),
	}

	log.Printf("Control Params: %+v\n", controlParams)

	resp, err := svc.AssumeRole(controlParams)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetDefaults will get ControlPlane default values from assumer config file
func (c *ControlPlane) GetDefaults() error {
	if c.AccountNumber == "" {
		c.AccountNumber = viper.GetString("control.account")
	}

	if c.RoleArn == "" {
		c.RoleArn = "arn:aws:iam::" + viper.GetString("control.account") + ":role/" + viper.GetString("control.role")
	} else {
		c.RoleArn = "arn:aws:iam::" + c.AccountNumber + ":role/" + c.RoleArn
	}

	if c.Region == "" {
		c.Region = viper.GetString("control.region")
	}

	if c.SerialNumber == "" {
		username := os.Getenv("USER")
		c.SerialNumber = "arn:aws:iam::" + c.AccountNumber + ":mfa/" + username
	}

	return nil
}

// CheckMfa checks the presence of an MFA Token.
func CheckMfa(token string) error {
	if token == "" {
		return errors.New("MFA Token cannot be blank")
	}

	if match, _ := regexp.MatchString(`\d{6}`, token); !match {
		return errors.New("Invalid MFA Token. Token must be 6-digits")
	}

	return nil
}
