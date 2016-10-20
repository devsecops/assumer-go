package assumer

// Plane represents an AWS Account.
type Plane struct {
	AccountNumber string `min:"12" type:"string" required:"true"`
	RoleArn       string `min:"20" type:"string" required:"true"`
	Region        string
}
