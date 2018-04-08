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

	"github.com/spf13/cobra"
  "os"
  "path/filepath"
  "path"
  "github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize mBox micro service",
	Long: `Initialize mBox micro service.`,
  Run: func(cmd *cobra.Command, args []string) {
    wd, err := os.Getwd()
    if err != nil {
      er(err)
    }

    var project *Project
    if len(args) == 0 {
      project = NewProjectFromPath(wd)
    } else if len(args) == 1 {
      arg := args[0]
      if arg[0] == '.' {
        arg = filepath.Join(wd, arg)
      }
      if filepath.IsAbs(arg) {
        project = NewProjectFromPath(arg)
      } else {
        project = NewProject(arg)
      }
    } else {
      er("please provide only one argument")
    }

    initializeProject(project)

    fmt.Fprintln(cmd.OutOrStdout(), `Your micro service is ready at
`+project.AbsPath()+`.

Give it a try by going there and running
make init 

Add commands to it by running `+"`mbox add [cmdname]`.")
  },
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializeProject(project *Project) {
  if !exists(project.AbsPath()) { // If path doesn't yet exist, create it
    err := os.MkdirAll(project.AbsPath(), os.ModePerm)
    if err != nil {
      er(err)
    }
  } else if !isEmpty(project.AbsPath()) { // If path exists and is not empty don't use it
    er("Cobra will not create a new project in a non empty directory: " + project.AbsPath())
  }

  // We have a directory and it's empty. Time to initialize it.
  createLicenseFile(project.License(), project.AbsPath())
  createMainFile(project)
  createRootCmdFile(project)
  createGitIgnoreFile(project)
  createVodModFile(project)
  createReleaseFile(project)
  createMakeFile(project)
  createEditorConfigFile(project)
}

func createLicenseFile(license License, path string) {
  data := make(map[string]interface{})
  data["copyright"] = copyrightLine()

  // Generate license template from text and data.
  text, err := executeTemplate(license.Text, data)
  if err != nil {
    er(err)
  }

  // Write license text to LICENSE file.
  err = writeStringToFile(filepath.Join(path, "LICENSE"), text)
  if err != nil {
    er(err)
  }
}

func createMainFile(project *Project) {
  mainTemplate := `{{ comment .copyright }}
{{if .license}}{{ comment .license }}{{end}}

package main // import "{{ .vgopath }}"

import "{{ .importpath }}"

var (
  // RELEASE returns the release version
  release = "UNKNOWN"
  // REPO returns the git repository URL
  Repo = "UNKNOWN"
  // COMMIT returns the short sha from git
  Commit = "UNKNOWN"

  BuildTime = "UNKNOWN"
)

func main() {
  cmd.Execute(release)
}
`
  data := make(map[string]interface{})
  data["copyright"] = copyrightLine()
  data["license"] = project.License().Header
  data["importpath"] = path.Join(project.Name(), filepath.Base(project.CmdPath()))
  data["vgopath"] = project.Name()

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), "main.go"), mainScript)
  if err != nil {
    er(err)
  }
}

func createGitIgnoreFile(project *Project) {
  mainTemplate := `# Compiled Object files, Static and Dynamic libs (Shared Objects)
*.o
*.a
*.so
*.orig
*.log
~$*

# Folders
_obj
_test
.idea/*

# Architecture specific extensions/prefixes
*.[568vq]
[568vq].out

*.cgo1.go
*.cgo2.c
_cgo_defun.c
_cgo_gotypes.go
_cgo_export.*

_testmain.go

*.exe
*.test
*.prof
harp.json
.harp
tmp/*
*/_metadata_*
node_modules/
npm-debug.log
/{{ .appName }}
#/vendor/
  `
  data := make(map[string]interface{})
  data["appName"] = path.Base(project.Name())

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), ".gitignore"), mainScript)
  if err != nil {
    er(err)
  }
}

func createVodModFile(project *Project) {
  mainTemplate := `module "{{ .importpath }}"
  `
  data := make(map[string]interface{})
  data["importpath"] = path.Join(project.Name())
  // data["importpath"] = path.Join(project.Name(), filepath.Base(project.CmdPath()))
  // data["appName"] = path.Base(project.Name())

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), "go.mod"), mainScript)
  if err != nil {
    er(err)
  }
}

func createReleaseFile(project *Project) {
  mainTemplate := `{{ .ver }}
  `
  data := make(map[string]interface{})
  data["ver"] = ProjectVersion

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), "RELEASE"), mainScript)
  if err != nil {
    er(err)
  }
}

func createEditorConfigFile(project *Project) {
  mainTemplate := `# http://editorconfig.org
root = true

[*]
indent_style = space
indent_size = 2
end_of_line = lf
charset = utf-8
trim_trailing_whitespace = true
insert_final_newline = true

[*.md]
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
  `
  data := make(map[string]interface{})
  data["ver"] = ProjectVersion

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), ".editorconfig"), mainScript)
  if err != nil {
    er(err)
  }
}

func createMakeFile(project *Project) {
  mainTemplate := `#
# Makefile for this application
#
-include variable.mak

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m


REPO_URL ?=
IMAGE_NAME ?=
USER_NAME ?= grengojbo
ADMIN_USER ?= grengojbo
TAG_VERSION=$(shell cat RELEASE)

OSNAME=$(shell uname)
GO=$(shell which go)

CUR_TIME=$(shell date '+%Y-%m-%d_%H:%M:%S')
# Program version
VERSION=$(shell cat RELEASE)

# Binary name for bintray
BIN_NAME=$(shell basename $(abspath ./))

# Project name for bintray
PROJECT_NAME=$(shell basename $(abspath ./))
PROJECT_DIR=$(shell pwd)

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com.org

# Grab the current commit
GIT_COMMIT="$(shell git rev-parse HEAD)"

DIST_DIR="${PROJECT_DIR}/dist"
DIST_BIN="${DIST_DIR}/bin"

NAME ?= ${PROJECT_NAME}
TAG=${USER_NAME}/$(PROJECT_NAME)$(IMAGE_NAME):$(TAG_VERSION)

BUILD_TAGS ?= "netgo"
BUILD_TAGS_BINDATA ?= "netgo bindatafs"
BUILD_ENV = GOOS=linux GOARCH=amd64
ENVFLAGS = CGO_ENABLED=1 $(BUILD_ENV)
ifneq ($(GOOS), darwin)
  EXTLDFLAGS = -extldflags "-lm -lstdc++ -static"
else
  EXTLDFLAGS =
endif

GO_LINKER_FLAGS ?= -ldflags '$(EXTLDFLAGS) -s -w \
  -X "main.BuildTime=${CUR_TIME}" \
  -X "main.Version=${VERSION}" \
  -X "main.GitHash=${GIT_COMMIT}" \
  -X "config.Version=${VERSION}"

# Add the godep path to the GOPATH
#GOPATH=$(shell godep path):$(shell echo $$GOPATH)

#ifeq ($(OS),Darwin)
#  URL=$(shell dinghy ip)
#else
#  URL="127.0.0.1"
#endif
URL=$(shell dinghy ip)

default: help

help:
	@echo "..............................................................."
	@echo "Project: $(PROJECT_NAME) | current dir: $(PROJECT_DIR)"
	@echo "version: $(VERSION)\n"
	@echo "make init        - Load project"
	@echo "make protoc      - Generate gRPC"
	@echo "make build       - Build for current OS project"
	@echo "make release     - Build release project"
	@echo "make docs        - Project documentation"
	@echo "make deploy      - Deploy bin files to server"
	@echo "make update      - Update vendor files"
	@echo "make test        - Run all test"
	@echo "make version     - Current project version"
	@echo "make push        - Push Docker image"
	@echo "make serve       - Run local Docker image"
	@echo "...............................................................\n"

init:
	@go get -u github.com/golang/dep/cmd/dep
	@go get -u golang.org/x/vgo
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/proto
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/mwitkow/go-proto-validators/protoc-gen-govalidators

update:
	@dep ensure -update

protoc:
	@echo "Generate gRPC"
	@protoc --proto_path=${GOPATH}/src \
	 --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	 -I types/ \
	 --go_out=plugins=grpc:types \
	 --govalidators_out=types \
	 types/ping/*.proto
	@protoc --proto_path=${GOPATH}/src \
	 --proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	 -I . \
	 --go_out=. \
	 ./grpc/pb/authorize/*.proto


deploy:
	@echo "TODO"

release: clean
	@mkdir -p $(DIST_BIN)
	@echo "building release ${BIN_NAME} ${VERSION}"
	@GOOS=linux GOARCH=amd64 go build -a -tags "netgo bindatafs" -ldflags '-w -X main.BuildTime=${CUR_TIME} -X main.Version=${VERSION} -X main.GitHash=${GIT_COMMIT} -X config.Version=${VERSION}' -o $(DIST_HOME)/$(BIN_NAME) main.go
	@chmod 0755 $(DIST_BIN)/$(BIN_NAME)

clean:
	@test ! -e ./${BIN_NAME} || rm ./${BIN_NAME}
	@test ! -e ${DIST_HOME}/${BIN_NAME} || rm ${DIST_HOME}/${BIN_NAME}
	@#git gc --prune=0 --aggressive
	@find . -name "*.orig" -type f -delete
	@find . -name "*.log" -type f -delete
	@test ! -e ./dist || rm -R ./dist

test:
	@echo "Start test..."

clean-docker:
	docker rmi -f $(REPO_URL)$(TAG)
	docker system prune -f

push:
	docker push $(REPO_URL)$(TAG)

serve:
	@echo "$(OK_COLOR)RUN command line: open http://$(URL):$(PUB_PORT)/$(NO_COLOR)\n\n"
	@docker run --rm \
		--name=$(NAME) \
     	-p=$(PUB_PORT):$(PORT) \
     	-e PORT=$(PORT) \
     	-it $(REPO_URL)$(TAG)

stop-docker:
	docker stop $(NAME)

build-docker:
	docker build --tag=$(REPO_URL)$(TAG) .

# Attach a root terminal to an already running dev shell
shell:
	docker run -it --rm $(REPO_URL)$(TAG) bash

build-linux: clean
	@mkdir -p $(DIST_BIN)
	@echo "building version: ${VERSION} to  ${DIST_HOME}/${BIN_NAME}"
	@DB_NAME=$(DB_TEST) DB_USER=$(DB_USER_TEST) DB_PASS=$(DB_PASS_TEST) DB_HOST=$(DB_HOST)GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -tags netgo -ldflags "-s -w -X main.release=${VERSION} -X main.Commit=${GIT_COMMIT} -X main.BuildTime=${CUR_TIME}" -o $(DIST_BIN)/$(BIN_NAME) main.go
	@echo " "

build: clean
	@mkdir -p $(DIST_BIN)
	@echo "building version: ${VERSION} to  ${DIST_HOME}/${BIN_NAME}"
	@CGO_ENABLED=0 go build -a -installsuffix cgo -tags netgo -ldflags "-s -w -X main.release=${VERSION} -X main.Commit=${GIT_COMMIT} -X main.BuildTime=${CUR_TIME}" -o ./$(BIN_NAME) main.go
	@echo " "

version:
	@echo ${VERSION}

docs:
	godoc -http=:6060 -index

  `
  data := make(map[string]interface{})
  data["ver"] = ProjectVersion

  mainScript, err := executeTemplate(mainTemplate, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.AbsPath(), "Makefile"), mainScript)
  if err != nil {
    er(err)
  }
}

func createRootCmdFile(project *Project) {
  template := `{{comment .copyright}}
{{if .license}}{{comment .license}}{{end}}

package cmd

import (
	"fmt"
	"os"
  "runtime"
{{if .viper}}
	homedir "github.com/mitchellh/go-homedir"{{end}}
	"github.com/spf13/cobra"{{if .viper}}
  "github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"{{end}}
){{if .viper}}

var Verbose bool
var BashCompletion bool
var ReleaseVersion string
var defaultConfigName string
var cfgFile string{{end}}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{.appName}}",
	Short: "A brief description of your application",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
  },
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

func init() { {{- if .viper}}
  defaultConfigName = ".{{ .appName }}"
	cobra.OnInitialize(initConfig)
{{end}}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.{{ if .viper }}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .appName }}.json)"){{ else }}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .appName }}.json)"){{ end }}
  rootCmd.PersistentFlags().BoolVarP(&BashCompletion, "bash-completion", "", false, "Generating Bash Completions")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}{{ if .viper }}

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
    if Verbose {
      fmt.Println("Using config file:", viper.ConfigFileUsed())
    }

    // uncomment to watch changed config file
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
      if Verbose {
        fmt.Println("Config file changed:", e.Name)
      }
    })
  }
}{{ end }}
`

  data := make(map[string]interface{})
  data["copyright"] = copyrightLine()
  data["viper"] = viper.GetBool("useViper")
  data["license"] = project.License().Header
  data["appName"] = path.Base(project.Name())

  rootCmdScript, err := executeTemplate(template, data)
  if err != nil {
    er(err)
  }

  err = writeStringToFile(filepath.Join(project.CmdPath(), "root.go"), rootCmdScript)
  if err != nil {
    er(err)
  }

}
