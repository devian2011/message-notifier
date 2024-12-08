package swap

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwapAndLoad(t *testing.T) {
	//TODO: Try to test read write data without filesystem
	const fileName = "/tmp/test.tmp"
	defer os.Remove(fileName)

	type in struct {
		Input  string
		Output string
	}

	i1 := &in{
		Input:  "test_input",
		Output: "test_output",
	}

	swapErr := Swap(fileName, i1)
	if swapErr != nil {
		t.Error(swapErr)
	}
	i2 := &in{}

	loadErr := Load(fileName, i2)
	if loadErr != nil {
		t.Error(loadErr)
	}

	assert.Equal(t, i1.Input, i2.Input)
	assert.Equal(t, i1.Output, i2.Output)
}
