/*
 * Copyright © 2019. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package textreplacer

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/util"
)

// activityLogger is the default logger for the Filter Activity
var log = logger.GetLogger("activity-textreplacer")

const (
	sLeftToken      = "leftToken"
	sRightToken     = "rightToken"
	iInputDocument  = "inputDocument"
	iReplacements   = "replacements"
	oOutputDocument = "outputDocument"
)

// Mapping is an Activity that is used to Filter a message to the console
type TextReplacer struct {
	metadata    *activity.Metadata
	initialized bool
	mux         sync.Mutex
	tokenMap    map[string][]string
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	aCSVParserActivity := &TextReplacer{
		metadata: metadata,
		tokenMap: make(map[string][]string),
	}
	return aCSVParserActivity
}

// Metadata returns the activity's metadata
func (a *TextReplacer) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - Filters the Message
func (a *TextReplacer) Eval(ctx activity.Context) (done bool, err error) {

	tokens, err := a.getTokens(ctx)
	if nil != err {
		return false, err
	}

	inputDocument, ok := ctx.GetInput(iInputDocument).(string)
	if !ok {
		return false, errors.New("Invalid document ... ")
	}

	replacements := ctx.GetInput(iReplacements).(*data.ComplexObject)
	replacementMap := replacements.Value.(map[string]interface{})

	mapper := NewKeywordMapper(inputDocument, tokens[0], tokens[1])
	document := mapper.replace("", replacementMap)

	log.Info("document = ", document)

	ctx.SetOutput(oOutputDocument, document)

	return true, nil
}

func (a *TextReplacer) getTokens(context activity.Context) ([]string, error) {
	myId := util.ActivityId(context)
	tokens := a.tokenMap[myId]
	log.Info("tokenMap : ", a.tokenMap, ", myId : ", myId)
	if nil == tokens {
		a.mux.Lock()
		defer a.mux.Unlock()
		tokens = a.tokenMap[myId]
		if nil == tokens {
			tokens = make([]string, 2)
			leftToken, _ := context.GetSetting(sLeftToken)
			if nil != leftToken {
				tokens[0] = leftToken.(string)
			}
			rightToken, _ := context.GetSetting(sRightToken)
			if nil != rightToken {
				tokens[1] = rightToken.(string)
			}
			a.tokenMap[myId] = tokens
		}
	}

	return tokens, nil
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
