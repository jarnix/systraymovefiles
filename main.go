package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/getlantern/systray"
)

// Config contains our configuration
type Config struct {
	Item struct {
		From string `yaml:"from"`
		To   string `yaml:"dest"`
	}
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("Hopstarter-Soft-Scraps-Button-Next.ico"))
	systray.SetTitle("Systray move files")
	systray.SetTooltip("Look at me, I'm a tooltip!")

	go func() {
		for {

			time.Sleep(1 * time.Second)
		}
	}()

}

func onExit() {
	// Cleaning stuff here.
}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
