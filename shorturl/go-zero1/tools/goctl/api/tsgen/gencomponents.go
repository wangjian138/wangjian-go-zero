package tsgen

import (
	"path"
	"strings"
	"text/template"

	"shorturl/go-zero/tools/goctl/api/spec"
	apiutil "shorturl/go-zero/tools/goctl/api/util"
	"shorturl/go-zero/tools/goctl/util"
)

const (
	componentsTemplate = `// Code generated by goctl. DO NOT EDIT.

{{.componentTypes}}
`
)

func genComponents(dir string, api *spec.ApiSpec) error {
	types := api.Types
	if len(types) == 0 {
		return nil
	}

	val, err := buildTypes(types)
	if err != nil {
		return err
	}

	outputFile := apiutil.ComponentName(api) + ".ts"
	filename := path.Join(dir, outputFile)
	if err := util.RemoveIfExist(filename); err != nil {
		return err
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, ".", outputFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	t := template.Must(template.New("componentsTemplate").Parse(componentsTemplate))
	return t.Execute(fp, map[string]string{
		"componentTypes": val,
	})
}

func buildTypes(types []spec.Type) (string, error) {
	var builder strings.Builder
	first := true
	for _, tp := range types {
		if first {
			first = false
		} else {
			builder.WriteString("\n")
		}
		if err := writeType(&builder, tp); err != nil {
			return "", apiutil.WrapErr(err, "Type "+tp.Name()+" generate error")
		}
	}

	return builder.String(), nil
}
