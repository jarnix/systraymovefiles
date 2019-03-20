package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
	"gopkg.in/yaml.v2"
)

// Config contains our configuration
type Config struct {
	Items map[string]ConfigItem `yaml:"watched"`
}

// ConfigItem is an element of configuration
type ConfigItem struct {
	Name      string `yaml:"name"`
	Extension string `yaml:"extension"`
	From      string `yaml:"from"`
	Dest      string `yaml:"dest"`
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {

	systray.SetIcon(getIcon("assets" + string(os.PathSeparator) + "Hopstarter-Soft-Scraps-Button-Next.ico"))
	systray.SetTitle("STMF")
	mQuit := systray.AddMenuItem("Quit", "Stop watching and moving")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		err := beeep.Alert("STMF", "config error: "+err.Error(), "assets/Button-Close-icon.png")
		if err != nil {
			panic(err)
		}
	}

	config := Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		err := beeep.Alert("STMF", "config error: "+err.Error(), "assets/Button-Close-icon.png")
		if err != nil {
			panic(err)
		}
	}

	var dests []string
	for _, item := range config.Items {
		dests = append(dests, "- "+item.From+" (*."+item.Extension+")")
	}
	systray.SetTooltip("STMF watching folders: \n" + strings.Join(dests, "\n"))

	go func() {
		for {
			time.Sleep(60 * time.Second)
			for _, item := range config.Items {
				var files []string
				allFiles, err := ioutil.ReadDir(item.From)
				if err != nil {
					log.Fatal(err)
				}

				for _, file := range allFiles {
					if filepath.Ext(file.Name()) == "."+item.Extension {
						fmt.Println(item.From + string(os.PathSeparator) + file.Name())
						files = append(files, file.Name())
					}
				}
				for _, file := range files {
					err := moveFile(item.From+string(os.PathSeparator)+file, item.Dest+string(os.PathSeparator)+file)
					if err != nil {
						err := beeep.Alert("STMF", "moveFile error: "+err.Error(), "assets/Button-Close-icon.png")
						if err != nil {
							panic(err)
						}
					} else {
						err := beeep.Notify("STMF", "File moved: "+file, "assets/Button-Next-icon.png")
						if err != nil {
							panic(err)
						}
					}
				}

			}

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

func moveFile(source, destination string) (err error) {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	fi, err := src.Stat()
	if err != nil {
		return err
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		dst.Close()
		os.Remove(destination)
		return err
	}
	err = dst.Close()
	if err != nil {
		return err
	}
	err = src.Close()
	if err != nil {
		return err
	}
	err = os.Remove(source)
	if err != nil {
		return err
	}
	return nil
}
