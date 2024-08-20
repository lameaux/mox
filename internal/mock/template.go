package mock

import (
	"github.com/Masterminds/sprig/v3"
	"text/template"
)

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	return f
}
