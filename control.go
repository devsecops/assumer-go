package assumer

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AssumeControlPlane assumes a Control Plane Role and returns the assumed role credentials
func AssumeControlPlane(c *Plane, mfaToken *string) (*sts.AssumeRoleOutput, error) {
	// username := viper.GetString("default.username")
	username := os.Getenv("USER")
	serialNumber := "arn:aws:iam::" + c.AccountNumber + ":mfa/" + username

	// Get sts
	svc := sts.New(session.New(), aws.NewConfig().WithRegion(c.Region))

	controlParams := &sts.AssumeRoleInput{
		RoleArn:         aws.String(c.RoleArn),     // Required
		RoleSessionName: aws.String("AssumedRole"), // Required
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(*mfaToken),
	}

	resp, err := svc.AssumeRole(controlParams)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
