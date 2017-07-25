package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"

	"github.com/golang/glog"
)

var templates = map[string]string{
	"utils.PackageDiffResult":             "utils/output_templates/singleVersionOutput.txt",
	"utils.MultiVersionPackageDiffResult": "utils/output_templates/multiVersionOutput.txt",
	"utils.HistDiffResult":                "utils/output_templates/historyOutput.txt",
	"utils.DirDiffResult":                 "utils/output_templates/fsOutput.txt",
}

func JSONify(diff interface{}) error {
	diffBytes, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return err
	}
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(diffBytes)
	return nil
}

func getTemplatePath(diff interface{}) (string, error) {
	diffType := reflect.TypeOf(diff).String()
	if path, ok := templates[diffType]; ok {
		return path, nil
	}
	return "", fmt.Errorf("No available template")
}

func TemplateOutput(diff interface{}) error {
	tempPath, err := getTemplatePath(diff)
	if err != nil {
		glog.Error(err)

	}
	tmpl, err := template.ParseFiles(tempPath)
	if err != nil {
		glog.Error(err)
		return err
	}
	err = tmpl.Execute(os.Stdout, diff)
	if err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
