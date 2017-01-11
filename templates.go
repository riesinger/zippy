package main

import (
	"errors"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/uber-go/zap"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var templates map[string]*raymond.Template

func loadTemplates(themeDir string) {
	templates = make(map[string]*raymond.Template)
	// Find all template files in the current theme
	compileDir(filepath.Join(themeDir, "templates"), filepath.Join(themeDir, "templates"))
	logger.Info("Compiled all templates")
}

func getTemplate(path string, list bool) (*raymond.Template, error) {
	// First, check if the given path exists directly
	template, exists := templates[path]
	if exists {
		return template, nil
	} else {
		// If not, check for the _single and _list versions
		var name string
		if list == true {
			name = "_list"
		} else {
			name = "_single"
		}
		tmpl, exists := templates[path+"/"+name]
		if exists == false {
			return nil, errors.New(fmt.Sprintf("requested template %s does not exist", path))
		}
		return tmpl, nil

	}

}

func compileDir(dirpath, themepath string) {
	files, _ := ioutil.ReadDir(dirpath)
	for _, f := range files {
		fp := filepath.Join(dirpath, f.Name())
		if f.IsDir() {
			compileDir(fp, themepath)
		} else {
			tmpl, err := raymond.ParseFile(fp)
			if err != nil {
				logger.Error("Could not parse template", zap.String("file", fp), zap.Error(err))
			} else {
				name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
				relpath, _ := filepath.Rel(themepath, dirpath)
				templates[relpath+"/"+name] = tmpl
			}
		}
	}
}
