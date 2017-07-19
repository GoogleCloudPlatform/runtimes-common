package utils

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"strings"
)

func JSONify(diff interface{}) (string, error) {
	diffBytes, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return "", err
	}
	return string(diffBytes), nil
}

func Output(diff PackageDiff) error {
	const master = `Packages found only in {{.Image1}}:{{range $name, $value := .Packages1}}{{"\n"}}{{print "-"}}{{$name}}{{"\t"}}{{$value}}{{end}}{{"\n"}}
Packages found only in {{.Image2}}:{{range $name, $value := .Packages2}}{{"\n"}}{{print "-"}}{{$name}}{{"\t"}}{{$value}}{{end}}
Version differences:{{"\n"}}	(Package:	{{.Image1}}{{"\t\t"}}{{.Image2}}){{range .InfoDiff}}
	{{.Package}}:	{{.Info1.Version}}	{{.Info2.Version}}
	{{end}}{{"\n"}}`

	funcs := template.FuncMap{"join": strings.Join}

	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}

	if err := masterTmpl.Execute(os.Stdout, diff); err != nil {
		log.Fatal(err)
	}
	return nil
}

func OutputMulti(diff MultiVersionPackageDiff) error {
	const master = `Packages found only in {{.Image1}}:{{range $name, $value := .Packages1}}{{"\n"}}{{print "-"}}{{$name}}{{end}}{{"\n"}}
Packages found only in {{.Image2}}:{{range $name, $value := .Packages2}}{{"\n"}}{{print "-"}}{{$name}}{{end}}
Version differences:{{"\n"}}	(Package:	{{.Image1}}{{"\t\t"}}{{.Image2}}){{range .InfoDiff}}
	{{.Package}}:	{{range .Info1}}{{.Version}}{{end}}	{{"\t"}} {{range .Info2}}{{.Version}}{{end}}
	{{end}}{{"\n"}}`

	funcs := template.FuncMap{"join": strings.Join}

	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}

	if err := masterTmpl.Execute(os.Stdout, diff); err != nil {
		log.Fatal(err)
	}
	return nil
}
