package actions

import (
	"strings"

	"github.com/stretchr/testify/assert"

	"testing"
)

var testTable = []struct {
	input  string
	expect string
}{
	{`<!DOCTYPE html>
<html>
<head>
	<title>Hello world</title>
</head>
<body>

</body>
</html>`, "Hello world"},
	{`<!DOCTYPE html>
<html>
<head>
	<title>
	Hello world
	</title>
</head>
<body>

</body>
</html>`, `
	Hello world
	`},
}

var te *TitleExtract = NewTitleExtract()

func TestTitleExtract(t *testing.T) {
	for _, tt := range testTable {
		r := strings.NewReader(tt.input)
		actual, err := cssSelectHTML(r, te.selector)

		assert.NoError(t, err)
		assert.Equal(t, tt.expect, actual)
	}
}
