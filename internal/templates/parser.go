package templates

import (
	"bytes"
	"os"
	"text/template"
)

func ParseFile(filePath string, params map[string]interface{}) (string, error) {
	body, err := os.ReadFile(filePath)
	if err != nil {
		return "", ErrTemplateFileNotExists
	}

	return Parse(string(body), params)
}

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
