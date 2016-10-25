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
	ctrlPlane *assumer.ControlPlane
	tgtPlane  *assumer.TargetPlane
)

func main() {

	flag.Parse()

	if help || flag.NFlag() == 0 {
		printHelp()
	}

	if version {
		printVersion()
		os.Exit(0)
	}

	requestedControlPlane := &assumer.Plane{AccountNumber: ctrlAcctNum, RoleArn: ctrlAcctRole, Region: region}
	requestedTargetPlane := &assumer.Plane{AccountNumber: tgtAcctNum, RoleArn: tgtAcctRole, Region: region}

	fmt.Println(os.Getenv("USER"), "is assuming into:")
	color.Yellow("Target Plane: %+v\n", *requestedTargetPlane)
	fmt.Println("via")
	color.Yellow("Control Plane: %+v\n", *requestedControlPlane)

	// get mfa token from user
	if err := assumer.CheckMfa(token); err != nil {
		mfa(&token)
	}

	ctrlPlane = &assumer.ControlPlane{Plane: *requestedControlPlane, MfaToken: token} // construct the control plane object

	if debug {
		log.Printf("Control Plane Values: %+v\n", *ctrlPlane)
	}

	tgtPlane = &assumer.TargetPlane{Plane: *requestedTargetPlane} // construct the target plane object

	if debug {
		log.Printf("Target Plane Values: %+v\n", *tgtPlane)
	}

	ctrlCreds, err := ctrlPlane.Assume()
	checkErr(err)

	if debug {
		log.Printf("Control Role Assumed Successfully: %+v\n", ctrlCreds)
	}

	tgtCreds, err := tgtPlane.Assume(ctrlCreds)
	checkErr(err)

	if debug {
		log.Printf("Target Role Assumed Successfully: %+v\n", tgtCreds)
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

	if err := assumer.Config(); err != nil {
		checkErr(err)
	}
}

func printHelp() {
	printVersion()
	fmt.Printf("CONFIG: %s\n\n", viper.ConfigFileUsed())
	fmt.Printf("USAGE: %s\n\n", "assumer [options]")
	fmt.Println("OPTIONS:")
	flag.PrintDefaults()
	os.Exit(0)
}
