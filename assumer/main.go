package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
	flag "github.com/ogier/pflag"
	"github.com/pmbenjamin/assumer"
	"github.com/spf13/viper"
)

// Flags
var (
	cfgFile      string
	ctrlAcctNum  string
	ctrlAcctRole string
	tgtAcctNum   string
	tgtAcctRole  string
	region       string
	profile      string
	gui          bool
	env          bool
	version      bool
	debug        bool
	help         bool
)

// Vars
var (
	token     string
	ctrlPlane *assumer.Plane
	tgtPlane  *assumer.Plane
)

const semver = "0.0.1"

func main() {

	if help {
		printHelp()
	}

	ctrlPlane := &assumer.Plane{AccountNumber: ctrlAcctNum, RoleArn: ctrlAcctRole, Region: region} // construct the control plane object
	ctrlPlane = getControlPlane(ctrlPlane)                                                         // get defaults if some args are skipped. This requires a config file.

	// get token interactively, or can be passed as a flag
	if token == "" {
		token = mfa(&token)
	}

	tgtPlane := &assumer.Plane{AccountNumber: tgtAcctNum, RoleArn: tgtAcctRole, Region: region} // construct the target plane object
	tgtPlane = getTargetPlane(tgtPlane)                                                         // get defaults if some args are skipped. This requires a config file.

	fmt.Printf("Control Plane: ")
	color.Yellow("%s", *ctrlPlane)
	fmt.Printf("Target Plane: ")
	color.Yellow("%s", *tgtPlane)

	if debug {
		log.Println("Control Plane:", *ctrlPlane)
		log.Println("Target Plane:", *tgtPlane)
	}

	ctrlCreds, err := assumer.AssumeControlPlane(ctrlPlane, &token) // assume into control plane
	if err != nil {
		color.Red("Assume Control Plane Failed!")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	tgtCreds, err := assumer.AssumeTargetAccount(tgtPlane, ctrlCreds) // assume into target plane
	if err != nil {
		color.Red("Assume Target Plane Failed!")
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		color.Green("\nSUCCESS!\n")
		fmt.Println(*tgtCreds.Credentials)
	}

	if debug {
		log.Println("Control Creds:", ctrlCreds)
		log.Println("Target Creds:", tgtCreds)
	}

	if env {
		exportEnv(tgtCreds)
	}

	if gui {
		openGui(tgtCreds)
	}
}

// initialize flags & config files
func init() {
	flag.BoolVarP(&debug, "debug", "d", false, "Debug Mode")
	flag.StringVarP(&ctrlAcctNum, "control-account", "A", "", "Control Account Number")
	flag.StringVarP(&ctrlAcctRole, "control-role", "R", "", "Control Account Role")
	flag.StringVarP(&tgtAcctNum, "target-account", "a", "", "Target Account Number")
	flag.StringVarP(&tgtAcctRole, "target-role", "r", "", "Target Account Role")
	flag.StringVar(&region, "region", "", "AWS Region")
	flag.StringVar(&profile, "profile", "", "AWS Profile")
	flag.StringVar(&cfgFile, "config", "", "Path to config file")
	flag.BoolVarP(&help, "help", "h", false, "Print help message")
	flag.Bool("version", false, "Print Version")
	flag.BoolVarP(&gui, "gui", "g", false, "AWS Console GUI")
	flag.BoolVarP(&env, "env", "e", false, "Export Credentials as Environment Variables")
	flag.Parse()

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.assumer")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		color.Yellow("Warning! Config file not found: %s \n", viper.ConfigFileUsed())
	}
}

func printHelp() {
	fmt.Println("Assumer ", semver)
	fmt.Println("Using config file:", viper.ConfigFileUsed())
	flag.Usage()
	os.Exit(0)
}

func openGui(tgtCreds *sts.AssumeRoleOutput) {
	fmt.Println("Generating AWS Console URL")

	// issuerUrl := "assumer"
	// consoleUrl := "https://console.aws.amazon.com/"
	// signinUrl := "https://signin.aws.amazon.com/federation"

	// sessionJson := ""
}

func exportEnv(tgtCreds *sts.AssumeRoleOutput) {
	os.Setenv("AWS_ACCESS_KEY_ID", *tgtCreds.Credentials.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *tgtCreds.Credentials.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", *tgtCreds.Credentials.SessionToken)

	color.Green("\nExecute the following in your terminal session:\n")
	fmt.Printf("AWS_ACCESS_KEY_ID=%s\n", os.Getenv("AWS_ACCESS_KEY_ID"))
	fmt.Printf("AWS_SECRET_ACCESS_KEY=%s\n", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	fmt.Printf("AWS_SESSION_TOKEN=%s\n\n", os.Getenv("AWS_SESSION_TOKEN"))
}
