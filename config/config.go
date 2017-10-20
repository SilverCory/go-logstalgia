package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type LogstalgiaConfig struct {
	StaticFileDirectory   string     `json:"static_file_directory"`
	TemplateFileDirectory string     `json:"template_file_directory"`
	PageConfig            PageConfig `json:"page_config"`
}

type PageConfig struct {
	Speed     int
	Framerate int
	Colours   bool
	Time      bool
	Summarise bool
}

func (l *LogstalgiaConfig) Load() {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		l.Save()
		fmt.Println("The default configuration has been saved. Please edit this and restart!")
		os.Exit(0)
		return
	} else {
		data, err := ioutil.ReadFile("./config.json")
		if err != nil {
			fmt.Println("There was an error loading the config!", err)
			return
		}

		err = json.Unmarshal(data, l)
		if err != nil {
			l.StaticFileDirectory = "./static/"
			l.TemplateFileDirectory = "./templates/"
			l.PageConfig = PageConfig{
				Time:      true,
				Colours:   true,
				Framerate: 30,
				Speed:     15,
				Summarize: true,
			}
			fmt.Println("There was an error loading the config!", err)
			return
		}
	}
}

func (l *LogstalgiaConfig) Save() error {
	data, err := json.MarshalIndent(l, "", "\t")

	if err != nil {
		fmt.Println("There was an error saving the config!", err)
		return err
	}

	err = ioutil.WriteFile("./config.json", data, 0644)
	if err != nil {
		fmt.Println("There was an error saving the config!", err)
		return err
	}

	return nil

}
