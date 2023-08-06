package flexmessage

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TemplateVars struct {
	Name string
	Age  int
}

func TestTemplating(t *testing.T) {
	tmpl, err := template.ParseFiles("templates/test/sample.json.tmpl")
	assert.Nil(t, err)

	v := TemplateVars{"Tokihide Nashiichiro", 14}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, v)
	assert.Nil(t, err)

	assert.JSONEq(t, `{"name":"Tokihide Nashiichiro","age":14}`, string(b.Bytes()))
}
