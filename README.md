# Assumer
Assume roles between AWS Control Plane accounts and Target accounts safely and securely.

## Install
`go get -d github.com/pmbenjamin/assumer`

## Usage
### CLI
```
assumer -h
assumer -a <target-account-number> -r <target-account-role> -A <control-account-number> -R <control-account-role>
```
#### Required Flags
```
  -A, --control-account Control Account Number
  -R, --control-role    Control Account Role
  -a, --target-account  Target Account Number
  -r, --target-role     Target Account Role
```
#### Optional Flags
```
  -g, --gui             AWS Console GUI
  --profile             AWS Profile
  --region              AWS Region
```
### Library
```go
package main

import "github.com/pmbenjamin/assumer"

func main() {

  // construct the control plane object
  controlPlane := assumer.Plane{AccountNumber: "123456789012", RoleArn: "arn:aws:iam::123456789012:role/control-role", Region: "us-west-2"} 

  // construct the target plane object
  targetPlane := assumer.Plane{AccountNumber: "123123123123", RoleArn: "arn:aws:iam::123123123123:role/target-plane"}

  // ... get MFA Token

  controlCreds, err := assumer.AssumeControlPlane(controlPlane, mfaToken)
  if err != nil {
    fmt.Println(err)
  }

  targetCreds, err := assumer.AssumeTargetPlane(targetPlane, controlCreds)
  if err != nil {
    fmt.Println(err)
  }

  // Now you have Target Plane Credentials...
  targetCreds.Credentials.AccessKey
  targetCreds.Credentials.SecretKey
  targetCreds.Credentials.Region
}
```

## Configuration
Assumer expects the config file to be called `config` and supports multiple configuration formats (e.g. [`TOML`](https://github.com/toml-lang/toml), `YAML`, & `JSON`).  
Assumer expects the configuration file to be located in `$HOME/.assumer/config.xyz` or in the current working directory.  
The config file is only used if the user assumes role via `assumer [target-account-name]`, rather than `assumer -a <target-account-number> -r <target-account-role> -A <control-account-number> -R <control-account-role>`

### Example
```
[myControlAccount]
account = 123456789012
role = "my/control/iam/role"
region = "us-west-2"

[myTarget]
  [myTarget.prod.da]
  account = 123456789012
  region = "us-west-2"
  role = "my/target/iam/role"

  [myTarget.prod.ro]
  account = 123456789012
  region = "us-west-2"
  role = "my/target/iam/role"
```

## Upcoming Features
- [ ] Open AWS Console in browser with `-g` or `--gui` flag
- [ ] Assume into target accounts with a simple command: `assumer <target-account-name>`
- [x] Support different configuration formats (e.g. `JSON`, `YAML`)
- [ ] Distribute binary via Homebrew, so users can `brew install assumer`