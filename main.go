package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/getlantern/systray"
	"gopkg.in/yaml.v2"
)

// Config contains our configuration
type Config struct {
	Item Item
}

type Item struct {
	Extension string `yaml:"extension"`
	From      string `yaml:"from"`
	Dest      string `yaml:"dest"`
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("Hopstarter-Soft-Scraps-Button-Next.ico"))
	systray.SetTitle("Systray move files")
	systray.SetTooltip("Look at me, I'm a tooltip!")

	dat, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	m1 := Item{
		Extension: "torrent",
		From:      "toto",
		Dest:      "tutu",
	}

	m := Config{Item: m1}

	d, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- m dump:\n%s\n\n", string(d))

	config := Config{}
	err = yaml.Unmarshal([]byte(dat), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t:\n%v\n\n", config)

	/*
		go func() {
			for {

				time.Sleep(1 * time.Second)
			}
		}()
	*/

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
