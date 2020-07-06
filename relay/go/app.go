package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/signalwire/signalwire-golang/signalwire"
)

// App environment settings
var (
	// required
	ProjectID      = "" //os.Getenv("ProjectID")
	TokenID        = "" //os.Getenv("TokenID")
	DefaultContext = "FIRESTARTER"  //os.Getenv("DefaultContext")
)

// Contexts needed for inbound calls
var Contexts = []string{DefaultContext}

// PProjectID passed from command-line
var PProjectID string

// PTokenID passed from command-line
var PTokenID string

// PContext passed from command line (just one being passed, although we support many)
var PContext string

/*gopl.io spinner*/
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// MyOnIncomingCall - gets executed when we receive an incoming call
func MyOnIncomingCall(consumer *signalwire.Consumer, call *signalwire.CallObj) {
	fmt.Printf("got incoming call.\n")

	resultAnswer, _ := call.Answer()
	if !resultAnswer.Successful {
		if err := consumer.Stop(); err != nil {
			signalwire.Log.Error("Error occurred while trying to stop Consumer\n")
		}

		return
	}

	signalwire.Log.Info("Playing audio on call..\n")

	go spinner(100 * time.Millisecond)

	// run the blocking PlayAudio() in a go-routine , we could have run PlayAudioAsync(), without the need to put it in a go-routine
	go func() {
		// blocking until media file finishes playing.
		if _, err := call.PlayAudio("https://cdn.signalwire.com/default-music/welcome.mp3"); err != nil {
			signalwire.Log.Error("Error occurred while trying to play audio\n")
		}
	}()

	timer := time.NewTimer(10 * time.Second)

	<-timer.C
	signalwire.Log.Info("Hangup call..\n")

	hangupResult, err := call.Hangup()
	if err != nil {
		// RELAY error
		signalwire.Log.Error("Error occurred while trying to hangup call\n")
	}

	if hangupResult.GetSuccessful() {
		signalwire.Log.Info("Call disconnect result: %s\n", hangupResult.GetReason().String())
	}

	if err := consumer.Stop(); err != nil {
		signalwire.Log.Error("Error occurred while trying to stop Consumer\n")
	}
}

func main() {
	var printVersion bool

	var verbose bool

	flag.BoolVar(&printVersion, "v", false, " Show version ")
	flag.StringVar(&PProjectID, "p", ProjectID, " ProjectID ")
	flag.StringVar(&PTokenID, "t", TokenID, " TokenID ")
	flag.StringVar(&PContext, "c", DefaultContext, " Context ")
	flag.BoolVar(&verbose, "d", false, " Enable debug mode ")
	flag.Parse()

	if printVersion {
		fmt.Printf("%s\n", signalwire.SDKVersion)
		fmt.Printf("Blade version: %d.%d.%d\n", signalwire.BladeVersionMajor, signalwire.BladeVersionMinor, signalwire.BladeRevision)
		fmt.Printf("App built with GO Lang version: " + fmt.Sprintf("%s\n", runtime.Version()))

		os.Exit(0)
	}

	if verbose {
		signalwire.Log.SetLevel(signalwire.DebugLevelLog)
	}

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
		for {
			s := <-interrupt
			switch s {
			case syscall.SIGHUP:
				fallthrough
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGINT:
				signalwire.Log.Info("Exit\n")
				os.Exit(0)
			}
		}
	}()

	Contexts = append(Contexts, PContext)
	consumer := new(signalwire.Consumer)
	// setup the Client
	consumer.Setup(PProjectID, PTokenID, Contexts)
	// register callback
	consumer.OnIncomingCall = MyOnIncomingCall

	signalwire.Log.Info("Wait incoming call..\n")

	// start
	if err := consumer.Run(); err != nil {
		signalwire.Log.Error("Error occurred while starting Signalwire Consumer: %v\n", err)
	}
}
