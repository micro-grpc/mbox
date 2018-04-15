package cmd

import (
	"fmt"
	"github.com/serenize/snaker"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Project contains name, license and paths to projects.
type Project struct {
	AppName      string
	absPath      string
	cmdPath      string
	srcPath      string
	license      License
	name         string
	Namespace    string
	Description  string
	Folder       folder
	ProjectDir   string
	ServiceName  string
	RelativeName string
	PackageName  string
	Copyright    string
	Licenses     string
	ImportPath   string
	Address      string
	Port         int
	Domain       string
	NameLicense  string
	Author       string
	Fqdn         string
}

// NewProject returns Project with specified project name.
func NewProject(projectName string) *Project {
	if projectName == "" {
		er(fmt.Errorf("can't create project with blank name\n"))
	}

	p := new(Project)
	p.name = projectName

	// 1. Find already created protect.
	p.absPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.absPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			er(err)
		}
		for _, srcPath := range srcPaths {
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) {
				p.absPath = filepath.Join(srcPath, projectName)
				break
			}
		}
	}

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.absPath == "" {
		p.absPath = filepath.Join(srcPaths[0], projectName)
	}

	p.Domain = viper.GetString("domain")
	p.ServiceName = path.Base(p.Name())
	if len(viper.GetString("name")) > 0 {
		p.ServiceName = viper.GetString("name")
	}
	n := strings.Replace(p.ServiceName, "-", "_", -1)

	if len(viper.GetString("fqdn")) == 0 {
		p.Fqdn = fmt.Sprintf("%s.service.%s", p.ServiceName, p.Domain)
	} else {
		p.Fqdn = viper.GetString("fqdn")
	}

	p.Address = viper.GetString("address")
	p.Port = viper.GetInt("port")
	p.Namespace = viper.GetString("namespace")
	p.RelativeName = snaker.SnakeToCamel(n)
	p.PackageName = snaker.CamelToSnake(n)
	p.ProjectDir = path.Join(path.Dir(p.Name()), filepath.Base(p.Name()))
	p.ImportPath = path.Join(p.Name(), filepath.Base(p.CmdPath()))
	p.AppName = path.Base(p.Name())

	if viper.GetInt("verbose") > 0 {
		fmt.Println("name:", p.name, "RelativeName:", p.RelativeName, "ServiceName:", p.ServiceName, "PackageName:", p.PackageName)
		fmt.Println("absPath:", p.absPath, "ProjectDir:", p.ProjectDir)
		fmt.Println("ImportPath:", p.ImportPath)
	}

	f := folder{Name: p.RelativeName, AbsPath: p.absPath}
	p.Folder = f

	return p
}

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string {
	if packageName == "" {
		return ""
	}

	for _, srcPath := range srcPaths {
		packagePath := filepath.Join(srcPath, packageName)
		if exists(packagePath) {
			return packagePath
		}
	}

	return ""
}

// NewProjectFromPath returns Project with specified absolute path to
// package.
func NewProjectFromPath(absPath string) *Project {
	if absPath == "" {
		er(fmt.Errorf("can't create project: absPath can't be blank\n"))
	}
	if !filepath.IsAbs(absPath) {
		er(fmt.Errorf("can't create project: absPath is not absolute\n"))
	}

	// If absPath is symlink, use its destination.
	fi, err := os.Lstat(absPath)
	if err != nil {
		er(fmt.Errorf("can't read path info: %v\n", err.Error()))
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		path, err := os.Readlink(absPath)
		if err != nil {
			er(fmt.Errorf("can't read the destination of symlink: %v\n", err.Error()))
		}
		absPath = path
	}

	p := new(Project)
	p.absPath = strings.TrimSuffix(absPath, findCmdDir(absPath))
	p.name = filepath.ToSlash(trimSrcPath(p.absPath, p.SrcPath()))
	if len(viper.GetString("name")) > 0 {
		p.name = viper.GetString("name")
	}
	f := folder{Name: p.name, AbsPath: p.absPath}
	p.Folder = f

	return p
}

// trimSrcPath trims at the beginning of absPath the srcPath.
func trimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		er(err)
	}
	return relPath
}

// License returns the License object of project.
func (p *Project) License() License {
	if p.license.Text == "" && p.license.Name != "None" {
		p.license = getLicense()
	}
	return p.license
}

// Name returns the name of project, e.g. "github.com/spf13/cobra"
func (p Project) Name() string {
	return p.name
}

// CmdPath returns absolute path to directory, where all commands are located.
func (p *Project) CmdPath() string {
	if p.absPath == "" {
		return ""
	}
	if p.cmdPath == "" {
		p.cmdPath = filepath.Join(p.absPath, findCmdDir(p.absPath))
	}
	return p.cmdPath
}

// findCmdDir checks if base of absPath is cmd dir and returns it or
// looks for existing cmd dir in absPath.
func findCmdDir(absPath string) string {
	if !exists(absPath) || isEmpty(absPath) {
		return "cmd"
	}

	if isCmdDir(absPath) {
		return filepath.Base(absPath)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files {
		if isCmdDir(file) {
			return filepath.Base(file)
		}
	}

	return "cmd"
}

// isCmdDir checks if base of name is one of cmdDir.
func isCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

// AbsPath returns absolute path of project.
func (p Project) AbsPath() string {
	return p.absPath
}

// SrcPath returns absolute path to $GOPATH/src where project is located.
func (p *Project) SrcPath() string {
	if p.srcPath != "" {
		return p.srcPath
	}
	if p.absPath == "" {
		p.srcPath = srcPaths[0]
		return p.srcPath
	}

	for _, srcPath := range srcPaths {
		if filepathHasPrefix(p.absPath, srcPath) {
			p.srcPath = srcPath
			break
		}
	}

	return p.srcPath
}

// CamelCaseName returns a CamelCased name of the service
func (p *Project) CamelCaseName() string {
	return snaker.SnakeToCamel(p.ServiceName)
}

// SnakeCaseName returns a snake_cased_type name of the service
func (p *Project) SnakeCaseName() string {
	return snaker.CamelToSnake(p.ServiceName)
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}
	return path[0:len(prefix)] == prefix

}

func (p Project) write(templatePath string) error {
	err := os.MkdirAll(p.ProjectDir, os.ModePerm)
	if err != nil {
		return err
	}

	return p.Folder.render(templatePath, p)
	// return nil
}
