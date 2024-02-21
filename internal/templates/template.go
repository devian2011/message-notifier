package templates

import (
	"bytes"
	"os"
	"text/template"
)

func File(filePath string, params map[string]interface{}) (string, error) {
	body, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return String(string(body), params)
}

func String(body string, params map[string]interface{}) (string, error) {
	buf := bytes.Buffer{}

	tmpl, parseErr := template.New("parse").Parse(body)
	if parseErr != nil {
		return "", parseErr
	}

	execErr := tmpl.Execute(&buf, params)
	if execErr != nil {
		return "", execErr
	}

	return buf.String(), nil
}
