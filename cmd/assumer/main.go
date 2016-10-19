package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	flag "github.com/ogier/pflag"
	"github.com/pmbenjamin/assumer-go"
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
	token        string
	gui          bool
	version      bool
	debug        bool
	help         bool
)

// Vars
var (
	ctrlPlane *assumer.Plane
	tgtPlane  *assumer.Plane
)

const semver string = "0.0.5"

func main() {

	flag.Parse()

	if help || flag.NFlag() == 0 {
		printHelp()
	}

	if version {
		printVersion(semver)
		os.Exit(0)
	}

	ctrlPlane := &assumer.Plane{AccountNumber: ctrlAcctNum, RoleArn: ctrlAcctRole, Region: region} // construct the control plane object
	ctrlPlane = getControlPlane(ctrlPlane)                                                         // get defaults if some args are skipped. This requires a config file.

	tgtPlane := &assumer.Plane{AccountNumber: tgtAcctNum, RoleArn: tgtAcctRole, Region: region} // construct the target plane object
	tgtPlane = getTargetPlane(tgtPlane)                                                         // get defaults if some args are skipped. This requires a config file.

	fmt.Println(os.Getenv("USER"), "is assuming into:")
	color.Yellow("Target Plane: %s", *tgtPlane)
	fmt.Println("via")
	color.Yellow("Control Plane: %s", *ctrlPlane)

	if debug {
		log.Println("Control Plane:", *ctrlPlane)
		log.Println("Target Plane:", *tgtPlane)
	}

	// get token interactively, or can be passed via "-t" or "--token" flag
	if token == "" {
		token = mfa(&token)
	}

	ctrlCreds, err := assumer.AssumeControlPlane(ctrlPlane, &token) // assume into control plane
	checkErr(err)

	tgtCreds, err := assumer.AssumeTargetAccount(tgtPlane, ctrlCreds) // assume into target plane
	checkErr(err)

	if debug {
		log.Println("Control Creds:", ctrlCreds)
		log.Println("Target Creds:", tgtCreds)
	}

	color.Green("\nSUCCESS!")
	if gui {
		openGui(tgtCreds)
	} else {
		execEnv(tgtCreds)
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
	flag.StringVarP(&token, "token", "t", "", "MFA Token")
	flag.StringVar(&cfgFile, "config", "", "Path to config file")
	flag.BoolVarP(&help, "help", "h", false, "Print help message")
	flag.BoolVarP(&version, "version", "v", false, "Print Version")
	flag.BoolVarP(&gui, "gui", "g", false, "AWS Console GUI")

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.assumer")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		color.Yellow("Warning! Config file not found: %s \n", viper.ConfigFileUsed())
	}
}

func printHelp() {
	printVersion(semver)
	fmt.Printf("CONFIG: %s\n\n", viper.ConfigFileUsed())
	fmt.Printf("USAGE: %s\n\n", "assumer [options]")
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	os.Exit(0)
}
