/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package filereader

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

const (
	sBaseFolder = "baseFolder"
	iFilename   = "filename"
	oContent    = "content"
)

var log = logger.GetLogger("tibco-activity-filereader")

type FileReaderActivity struct {
	metadata    *activity.Metadata
	baseFolders map[string]string
	mux         sync.Mutex
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &FileReaderActivity{
		metadata:    metadata,
		baseFolders: make(map[string]string),
	}
}

func (a *FileReaderActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *FileReaderActivity) Eval(context activity.Context) (done bool, err error) {
	log.Info("(Eval) Graph to file entering ......... ")

	baseFolder, err := a.getBaseFolder(context)
	if nil != err {
		return false, err
	}

	a.mux.Lock()
	defer a.mux.Unlock()

	filename := context.GetInput(iFilename).(string)
	if "" != baseFolder {
		filename = fmt.Sprintf("%s/%s", baseFolder, filename)
	}

	log.Info("(Eval) File name : ", filename)
	descriptor, err := readFile(filename)
	if nil != err {
		return false, err
	}

	mapper := NewKeywordMapper(descriptor, "{{", "}}")

	context.SetOutput(oContent, mapper.replace("", map[string]interface{}{"temp.project.home": baseFolder}))

	log.Info("(Eval) write object to file exit ......... ")
	return true, nil
}

func (a *FileReaderActivity) getBaseFolder(context activity.Context) (string, error) {

	myId := util.ActivityId(context)
	baseFolder := a.baseFolders[myId]
	log.Info("%%%%%%%%%%%%%%%%%%%%%% getOutputFile : myId = ", myId, ", outputFile = ", baseFolder)

	if "" == baseFolder {
		a.mux.Lock()
		defer a.mux.Unlock()
		baseFolder = a.baseFolders[myId]
		if "" == baseFolder {
			log.Info("Initializing FileReader Service start ...")

			baseFolderSetting, _ := context.GetSetting(sBaseFolder)
			baseFolder = baseFolderSetting.(string)

			log.Info("Initializing FileReader Service end ...")
			a.baseFolders[myId] = baseFolder
		}
	}

	return baseFolder, nil
}

func readFile(filename string) (string, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	//fmt.Println("Contents of file:", string(fileContent))
	return string(fileContent), nil
}

type KeywordReplaceHandler struct {
	result     string
	keywordMap map[string]interface{}
}

func (this *KeywordReplaceHandler) setMap(keywordMap map[string]interface{}) {
	this.keywordMap = keywordMap
}

func (this *KeywordReplaceHandler) startToMap() {
	this.result = ""
}

func (this *KeywordReplaceHandler) replace(keyword string) string {
	if nil != this.keywordMap[keyword] {
		return this.keywordMap[keyword].(string)
	}
	return ""
}

func (this *KeywordReplaceHandler) endOfMapping(document string) {
	this.result = document
}

func (this *KeywordReplaceHandler) getResult() string {
	return this.result
}

func NewKeywordMapper(
	template string,
	lefttag string,
	righttag string) *KeywordMapper {
	mapper := KeywordMapper{
		template:     template,
		keywordOnly:  false,
		slefttag:     lefttag,
		srighttag:    righttag,
		slefttaglen:  len(lefttag),
		srighttaglen: len(righttag),
	}
	return &mapper
}

type KeywordMapper struct {
	template     string
	keywordOnly  bool
	slefttag     string
	srighttag    string
	slefttaglen  int
	srighttaglen int
	document     bytes.Buffer
	keyword      bytes.Buffer
	mh           KeywordReplaceHandler
}

func (this *KeywordMapper) replace(template string, keywordMap map[string]interface{}) string {
	if "" == template {
		template = this.template
		if "" == template {
			return ""
		}
	}

	this.mh.setMap(keywordMap)
	this.document.Reset()
	this.keyword.Reset()

	scope := false
	boundary := false
	skeyword := ""
	svalue := ""

	this.mh.startToMap()
	for i := 0; i < len(template); i++ {
		//log.Infof("template[%d] = ", i, template[i])
		// maybe find a keyword beginning Tag - now isn't in a keyword
		if !scope && template[i] == this.slefttag[0] {
			if this.isATag(i, this.slefttag, template) {
				this.keyword.Reset()
				scope = true
			}
		} else if scope && template[i] == this.srighttag[0] {
			// maybe find a keyword ending Tag - now in a keyword
			if this.isATag(i, this.srighttag, template) {
				i = i + this.srighttaglen - 1
				skeyword = this.keyword.String()[this.slefttaglen:this.keyword.Len()]
				svalue = this.mh.replace(skeyword)
				if "" == svalue {
					svalue = fmt.Sprintf("%s%s%s", this.slefttag, skeyword, this.srighttag)
				}
				//log.Info("value ->", svalue);
				this.document.WriteString(svalue)
				boundary = true
				scope = false
			}
		}

		if !boundary {
			if !scope && !this.keywordOnly {
				this.document.WriteByte(template[i])
			} else {
				this.keyword.WriteByte(template[i])
			}
		} else {
			boundary = false
		}

		//log.Info("document = ", this.document)

	}
	this.mh.endOfMapping(this.document.String())
	return this.mh.getResult()
}

func (this *KeywordMapper) isATag(i int, tag string, template string) bool {
	for j := 0; j < len(tag); j++ {
		if tag[j] != template[i+j] {
			return false
		}
	}
	return true
}
