// Copyright © 2018 Oleg Dolya <oleg.dolya@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
)

var Verbose int
var BashCompletion bool
var ReleaseVersion string
var ProjectVersion string
var defaultConfigName string
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mbox",
	Short: "Is a micro gRPC framework",
	Long:  `mBox is a micro gRPC framework, distributed systems development for micro services.`,
	Run: func(cmd *cobra.Command, args []string) {
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if BashCompletion {
			if os.Geteuid() == 0 {
				bkFile := fmt.Sprintf("/etc/bash_completion.d/%s.bash", cmd.Use)
				if runtime.GOOS == "darwin" {
					bkFile = fmt.Sprintf("/usr/local/etc/bash_completion.d/%s.bash", cmd.Use)
				}
				fmt.Println("Generate: ", bkFile)
				cmd.GenBashCompletionFile(bkFile)
			} else {
				if runtime.GOOS == "darwin" {
					fmt.Printf("RUN sudo ./%s --bash-completion\n", cmd.Use)
				} else {
					fmt.Printf("RUN sudo %s --bash-completion\n", cmd.Use)
				}
			}
		}
		// },
		// Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(releaseVersion string) {
	ReleaseVersion = releaseVersion
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultConfigName = ".mbox"
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mbox.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	rootCmd.PersistentFlags().BoolVarP(&BashCompletion, "bash-completion", "", false, "Generating Bash Completions")
	rootCmd.PersistentFlags().StringVar(&ProjectVersion, "ver", "0.0.1", "project version")

	rootCmd.PersistentFlags().CountVarP(&Verbose, "verbose", "v", "verbose output")
	// rootCmd.PersistentFlags().Bool("debug", false, "debug mode")
	// viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// rootCmd.Flags().StringP("mode", "m", "production", "NODE_ENV=\"production\"")
	// viper.BindPFlag("NODE_ENV", rootCmd.PersistentFlags().Lookup("mode"))

	rootCmd.PersistentFlags().String("domain", "local", "External domain name")
	viper.BindPFlag("domain", rootCmd.PersistentFlags().Lookup("domain"))
	rootCmd.PersistentFlags().String("name", "", "Service Name (default project name)")
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	rootCmd.PersistentFlags().String("namespace", "", "Namespace for the service") // если указано то будет (namespace-serviceName.service.local)
	viper.BindPFlag("namespace", rootCmd.PersistentFlags().Lookup("namespace"))
	rootCmd.PersistentFlags().String("address", "", "gRPC Server IP address")
	rootCmd.PersistentFlags().Int64P("port", "p", 9000, "gRPC Server Port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd.PersistentFlags().Bool("proxy", false, "gRPC to JSON proxy")
	viper.BindPFlag("proxy", rootCmd.PersistentFlags().Lookup("proxy"))
	rootCmd.PersistentFlags().String("fqdn", "", "FQDN of service (defaults to serviceName.service.local)")
	viper.BindPFlag("fqdn", rootCmd.PersistentFlags().Lookup("fqdn"))
  rootCmd.PersistentFlags().Bool("sqlx", false, "support SqlX")
  viper.BindPFlag("sqlx", rootCmd.PersistentFlags().Lookup("sqlx"))
  rootCmd.PersistentFlags().Bool("gorm", false, "support GORM ORM")
  viper.BindPFlag("gorm", rootCmd.PersistentFlags().Lookup("gorm"))
	rootCmd.PersistentFlags().String("driver", "postgres", "Database driver")
	viper.BindPFlag("driver", rootCmd.PersistentFlags().Lookup("driver"))
	rootCmd.PersistentFlags().String("cf", "", "Custom config file name")
	viper.BindPFlag("cf", rootCmd.PersistentFlags().Lookup("cf"))

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "MIT")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.SetConfigFile(cfgFile)

	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if Verbose > 0 {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		// uncomment to watch changed config file
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			if Verbose > 0 {
				fmt.Println("Config file changed:", e.Name)
			}
		})
	}
}
