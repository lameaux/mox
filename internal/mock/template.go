package mock

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	return f
}
