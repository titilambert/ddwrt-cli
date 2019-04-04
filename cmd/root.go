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
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	//	"github.com/spf13/viper"
)

var cfgFile string
var hostAddress string
var username string
var password string
var dumpFolder string
var loadFolder string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ddwrt-cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ddwrt-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVarP(&hostAddress, "host", "H", "", "DD-WRT Host Address")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "DD-WRT username")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "DD-WRT password")
	//rootCmd.PersistentFlags().StringVarP(&dumpFolder, "dump", "d", "", "Folder where router config will be saved")
	//rootCmd.PersistentFlags().StringVarP(&loadFolder, "load", "l", "", "Folder where router config will be loaded from")
	rootCmd.MarkFlagRequired("host")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if hostAddress == "" {
		fmt.Println("Missing host Address")
		os.Exit(1)
	}
	if username == "" {
		fmt.Println("Missing username")
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("Missing password")
		os.Exit(1)
	}
	if cfgFile == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfgFile = path.Join(home, ".ddwrt-cli.yaml")
	} else {
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	/*

	   	if cfgFile != "" {
	   		// Use config file from the flag.
	   		viper.SetConfigFile(cfgFile)
	   	} else {
	   		// Find home directory.
	   		home, err := homedir.Dir()
	   		if err != nil {
	   			fmt.Println(err)
	   			os.Exit(1)
	   		}

	   		// Search config in home directory with name ".ddwrt-cli" (without extension).
	   		viper.AddConfigPath(home)
	   		viper.SetConfigName(".ddwrt-cli")
	   	}

	   	viper.AutomaticEnv() // read in environment variables that match

	   	// If a config file is found, read it in.
	   	if err := viper.ReadInConfig(); err == nil {
	           viper.SetConfigType("yaml")
	   		fmt.Println("Using config file:", viper.ConfigFileUsed())
	   	}
	*/
}
