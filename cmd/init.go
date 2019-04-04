// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"

	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"os"
	"path"
)

type CliConfig struct {
	Host     string
	Sections map[string]Section
}

// initCmd represents the services command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sections, err := getSections()
		if err != nil {
			log.Fatal(err)
		}
		// Write to config file
		cliConfig := CliConfig{Host: hostAddress,
			Sections: sections,
		}
		dataYaml, err := yaml.Marshal(cliConfig)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		err = ioutil.WriteFile(cfgFile, dataYaml, 0644)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
	},
}

func ReadCommands() CliConfig {
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		os.Exit(1)
	}
	cliConfig := CliConfig{}
	err = yaml.Unmarshal([]byte(data), &cliConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return cliConfig
}

func SetSections() {
	cliConfig := ReadCommands()
	for key, section := range cliConfig.Sections {
		var cmd = &cobra.Command{
			Use:   key,
			Short: section.Name,
		}
		for subkey, subSection := range section.SubSections {
			var subCmd = &cobra.Command{
				Use:    subkey,
				Short:  subSection.Name,
				Long:   subSection.Name,
				Run:    RunSection,
				PreRun: PreRunSection,
			}
			subCmd.PersistentFlags().StringVarP(&dumpFolder, "dump", "d", "", "Folder where router config will be saved")
			subCmd.PersistentFlags().StringVarP(&loadFolder, "load", "l", "", "Folder where router config will be loaded from")

			cmd.AddCommand(subCmd)
		}
		rootCmd.AddCommand(cmd)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfgFile = path.Join(home, ".ddwrt-cli.yaml")
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Print("\nWARNING: Sections file doesn't exist, please run `init` command.\n\n")
	} else {
		SetSections()
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// servicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// servicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
