package assumer

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// AssumeTargetAccount assumes a role in the target account and returns the assumed role creds
func AssumeTargetAccount(t *Plane, c *sts.AssumeRoleOutput) (*sts.AssumeRoleOutput, error) {
	targetParams := &sts.AssumeRoleInput{
		RoleArn:         aws.String(t.RoleArn),
		RoleSessionName: aws.String("AssumedRole"),
	}

	controlCreds := credentials.NewStaticCredentials(*c.Credentials.AccessKeyId, *c.Credentials.SecretAccessKey, *c.Credentials.SessionToken)

	stsClient := sts.New(session.New(), aws.NewConfig().WithRegion(t.Region).WithCredentials(controlCreds))
	resp, err := stsClient.AssumeRole(targetParams)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
