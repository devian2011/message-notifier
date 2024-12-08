package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	const tmpl = "{{.name}} <br /> {{.lastname}}"
	const result = "John <br /> Jones"

	parsed, parsedErr := Parse(tmpl, map[string]interface{}{"name": "John", "lastname": "Jones"})
	if parsedErr != nil {
		t.Error(parsedErr)
	}

	assert.Equal(t, result, parsed)
}
