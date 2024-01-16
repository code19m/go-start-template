package http

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func (srv *HttpHandler) registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(getNameByPriority)
	}
}

func getNameByPriority(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
	}
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("env"), ",", 2)[0]
	}
	if name == "" {
		name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
	}
	return name
}
