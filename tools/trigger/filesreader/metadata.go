package filesreader

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
}
type HandlerSettings struct {
	Filename        string `md:"Filename"`
	EmitPerLine     bool   `md:"EmitPerLine"`
	Asynchronous    bool   `md:"Asynchronous"`
	MaxNumberOfLine int    `md:"MaxNumberOfLine"`
}

type Output struct {
	MessageId    string `md:"MessageId"`
	FileContent  string `md:"FileContent"`
	ModifiedTime int64  `md:"ModifiedTime"`
	LineNumber   int    `md:"LineNumber"`
}

func (this *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"MessageId":    this.MessageId,
		"FileContent":  this.FileContent,
		"ModifiedTime": this.ModifiedTime,
		"LineNumber":   this.LineNumber,
	}
}

func (this *Output) FromMap(values map[string]interface{}) error {

	var err error
	this.MessageId, err = coerce.ToString(values["MessageId"])
	this.FileContent, err = coerce.ToString(values["FileContent"])
	this.ModifiedTime, err = coerce.ToInt64(values["ModifiedTime"])
	this.LineNumber, err = coerce.ToInt(values["LineNumber"])
	if err != nil {
		return err
	}

	return nil
}
