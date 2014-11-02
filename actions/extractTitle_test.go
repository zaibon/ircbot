package actions

import (
	"reflect"
	"strings"

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
</html>`, "Hello world"},
}

var te *TitleExtract = NewTitleExtract()

func TestTitleExtract(t *testing.T) {
	for _, tt := range testTable {
		r := strings.NewReader(tt.input)
		actual, err := cssSelectHTML(r, te.selector)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(actual, tt.expect) {
			t.Errorf("title expected %s ,actual %s\n", tt.expect, actual)
		}
	}
}
