package ssesub

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
	Event string `md:"Event"`
}

func (this *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Event": this.Event,
	}
}

func (this *Output) FromMap(values map[string]interface{}) error {

	var err error
	this.Event, err = coerce.ToString(values["Event"])
	if err != nil {
		return err
	}

	return nil
}
