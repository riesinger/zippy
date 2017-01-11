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

func getTemplate(name string) (*raymond.Template, error) {
	// First, check if the given path exists directly
	template, exists := templates[name]
	if exists {
		return template, nil
	} else {
		return nil, errors.New(fmt.Sprintf("requested template %s does not exist", name))
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
				templates[name] = tmpl
				logger.Debug("Parsed template", zap.String("name", name))
			}
		}
	}
}
