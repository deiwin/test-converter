package convert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/deiwin/test-converter/convert"
)

func TestConvertTest(t *testing.T) {
	reader := strings.NewReader(input)
	writer := new(bytes.Buffer)

	convert.Test(reader, writer)

	if writer.String() != output {
		t.Fatal("Expected something different in the output")
	}
}
