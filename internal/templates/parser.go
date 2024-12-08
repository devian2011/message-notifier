package templates

import (
	"bytes"
	"text/template"
)

func Parse(body string, params map[string]interface{}) (string, error) {
	buf := bytes.Buffer{}

	tmpl, parseErr := template.New("config").Parse(body)
	if parseErr != nil {
		return "", parseErr
	}

	execErr := tmpl.Execute(&buf, params)
	if execErr != nil {
		return "", execErr
	}

	return buf.String(), nil
}
