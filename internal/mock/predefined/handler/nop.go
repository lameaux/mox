package handler

import (
	"net/http"
)

func Nop(_ http.ResponseWriter, _ *http.Request) {
}
