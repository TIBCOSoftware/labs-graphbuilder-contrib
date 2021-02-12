package sseserver

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
)

type Settings struct {
}

type HandlerSettings struct {
	Connection connection.Manager `md:"sseConnection,required"`
}

type Output struct {
	Request string `md:"Request"`
}

func (this *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Request": this.Request,
	}
}

func (this *Output) FromMap(values map[string]interface{}) error {

	var err error
	this.Request, err = coerce.ToString(values["Request"])
	if err != nil {
		return err
	}

	return nil
}
