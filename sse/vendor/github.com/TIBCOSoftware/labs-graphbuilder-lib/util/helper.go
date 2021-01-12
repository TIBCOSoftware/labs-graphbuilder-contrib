/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/project-flogo/core/data/property"
)

const (
	GRAPH_ID = "$GRAPH_ID"
)

func SplitFilename(filename string) (string, string) {
	if "" != filename {
		indexSlash := strings.LastIndex(filename, "/")
		indexBackslash := strings.LastIndex(filename, "\\")
		var index int
		if indexSlash > indexBackslash {
			index = indexSlash
		} else {
			index = indexBackslash
		}
		return filename[:index], filename[index+1:]
	}
	return "", ""
}

func ReplaceParameter(origString string, parameter string, replacement string) string {
	if "" != origString && "" != parameter && "" != replacement {
		index := strings.Index(origString, parameter)
		if -1 != index {
			len := len(parameter)
			return fmt.Sprintf("%s%s%s", origString[:index], replacement, origString[index+len:])
		}
	}
	return origString
}

func CastGenMap(mapData interface{}) map[string]interface{} {
	if nil != mapData {
		return mapData.(map[string]interface{})
	}
	return make(map[string]interface{})
}

func CastGenArray(arrayData interface{}) []interface{} {
	if nil != arrayData {
		return arrayData.([]interface{})
	}
	return make([]interface{}, 0)
}

func CastString(stringData interface{}) string {
	if nil != stringData {
		return stringData.(string)
	}
	return ""
}

func TypeString(data interface{}) string {
	return reflect.TypeOf(data).String()
}

func ConvertToInteger(data interface{}) (interface{}, error) {
	//fmt.Println("Convert to Integer ....")
	switch data.(type) {
	case int:
		return int32(data.(int)), nil
	case int32:
		return data, nil
	case int64:
		return int32(data.(int64)), nil
	case float32:
		return int32(data.(float32)), nil
	case float64:
		return int32(data.(float64)), nil
	case string:
		return strconv.ParseInt(data.(string), 10, 32)
	}
	return data, fmt.Errorf("Unable to convert to Interger : ", reflect.TypeOf(data).String())
}

func ConvertToLong(data interface{}) (interface{}, error) {
	//fmt.Println("Convert to Long ....")
	switch data.(type) {
	case int:
		return int64(data.(int)), nil
	case int32:
		return int64(data.(int32)), nil
	case int64:
		return data, nil
	case float32:
		return int64(data.(float32)), nil
	case float64:
		return int64(data.(float64)), nil
	case string:
		return strconv.ParseInt(data.(string), 10, 64)
	}
	return data, fmt.Errorf("Unable to convert to Long : ", reflect.TypeOf(data).String())
}

func ConvertToString(data interface{}, dateTimeSample string) (interface{}, error) {
	//fmt.Println("Convert to String ....")

	switch data.(type) {
	case string:
		return data, nil
	case int:
		return strconv.Itoa(data.(int)), nil
	case int32:
		intData := int(data.(int32))
		return strconv.Itoa(intData), nil
	case int64:
		intData := int(data.(int64))
		return strconv.Itoa(intData), nil
	case float32:
		floatData := data.(float64)
		return strconv.FormatFloat(floatData, 'f', -1, 32), nil
	case float64:
		doubleData := data.(float64)
		return strconv.FormatFloat(doubleData, 'f', -1, 64), nil
	case bool:
		booeanlData := data.(bool)
		return strconv.FormatBool(booeanlData), nil
	case time.Time:
		dateTimeData := data.(time.Time)
		return dateTimeData.Format(dateTimeSample), nil
	}
	return data, fmt.Errorf("Unable to convert to string : ", data)
}

func ConvertToDouble(data interface{}) (interface{}, error) {
	//fmt.Println("Convert to Double ....")
	switch data.(type) {
	case int:
		return float64(data.(int)), nil
	case int32:
		return float64(data.(int32)), nil
	case float32:
		return float64(data.(float32)), nil
	case float64:
		return data, nil
	case string:
		return strconv.ParseFloat(data.(string), 64)
	}
	return data, fmt.Errorf("Unable to convert to Double : ", reflect.TypeOf(data).String())
}

func ConvertToBoolean(data interface{}) (interface{}, error) {
	//fmt.Println("Convert to Boolean ....")
	switch data.(type) {
	case bool:
		return data, nil
	case string:
		return strconv.ParseBool(data.(string))
	}
	return data, fmt.Errorf("Unable to convert to Boolean : ", reflect.TypeOf(data).String())
}

func ConvertToDate(data interface{}, dateTimeSample string) (interface{}, error) {
	//fmt.Println("Convert to Date ....")
	switch data.(type) {
	case time.Time:
		return data, nil
	case string:
		//fmt.Println("parsing time string, template = ", dateTimeSample, ", datetime = ", data)
		dateData, err := time.Parse(dateTimeSample, data.(string))
		intDateData := dateData.Unix()
		//fmt.Println("after parsing time string : ", intDateData, ", type : ", reflect.TypeOf(intDateData).String())
		return intDateData, err
	}
	return data, fmt.Errorf("Unable to convert to Date : ", reflect.TypeOf(data).String())
}

func TypeConversion(data interface{}, dataType string, dateTimeSample string) (interface{}, error) {
	if nil == data {
		return nil, nil
	}
	switch dataType {
	case "String":
		return ConvertToString(data, dateTimeSample)
	case "Integer":
		return ConvertToInteger(data)
	case "Long":
		return ConvertToLong(data)
	case "Boolean":
		return ConvertToBoolean(data)
	case "Double":
		return ConvertToDouble(data)
	case "Date":
		return ConvertToDate(data, dateTimeSample)
	}
	return data, nil
}

func StringToTypes(data string, dataType string, dateTimeSample string) (interface{}, error) {
	if "null" == data {
		return nil, nil
	}
	switch dataType {
	case "String":
		return data, nil
	case "Integer":
		return strconv.ParseInt(data, 10, 32)
	case "Long":
		return strconv.ParseInt(data, 10, 64)
	case "Boolean":
		return strconv.ParseBool(data)
	case "Double":
		return strconv.ParseFloat(data, 64)
	case "Date":
		return time.Parse(dateTimeSample, data)
	}
	return data, nil
}

func ToString(data interface{}, dataType string, dateTimeSample string) (string, error) {
	golangTypeString := reflect.TypeOf(data).String()

	//fmt.Println("data = ", data, ", dataType = ", dataType, ", golangTypeString = ", golangTypeString)

	switch dataType {
	case "String":
		if "string" == golangTypeString {
			stringData := data.(string)
			return fmt.Sprintf("%s", stringData), nil
		} else {
			return "", fmt.Errorf("Not a string type data")
		}
	case "Integer":
		if "int32" == golangTypeString {
			intData := int(data.(int32))
			return strconv.Itoa(intData), nil
		} else {
			return "", fmt.Errorf("Not a Integer(int32) type data")
		}
	case "Long":
		if "int64" == golangTypeString {
			longData := data.(int64)
			return strconv.FormatInt(longData, 10), nil
		} else {
			return "", fmt.Errorf("Not a Long(int64) type data")
		}
	case "Boolean":
		if "bool" == golangTypeString {
			booeanlData := data.(bool)
			return strconv.FormatBool(booeanlData), nil
		} else {
			return "", fmt.Errorf("Not a Boolean(bool) type data")
		}
	case "Double":
		if "float64" == golangTypeString {
			doubleData := data.(float64)
			return strconv.FormatFloat(doubleData, 'f', -1, 64), nil
		} else {
			return "", fmt.Errorf("Not a Double(float64) type data")
		}
	case "Date":
		if "time.Time" == golangTypeString {
			dateTimeData := data.(time.Time)
			return dateTimeData.Format(dateTimeSample), nil
		} else {
			return "", fmt.Errorf("Not a Date(time.Time) type data")
		}
	}
	return data.(string), nil
}

func ReplaceCharacter(str string, targetRegex string, replacement string, doReplace bool) string {
	if doReplace {
		var re = regexp.MustCompile(targetRegex)
		str = re.ReplaceAllString(str, replacement)
	}

	return str
}

func SliceContains(slice []string, targetElement string) bool {
	for _, element := range slice {
		if element == targetElement {
			return true
		}
	}
	return false
}

func IsInteger(data string) bool {
	_, err := strconv.ParseInt(data, 10, 64)
	return err == nil
}

func ReadFile(filename string) (string, error) {

	file, err := os.Open(filename)

	if nil != err {
		return "", err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var buffer bytes.Buffer
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		buffer.WriteString("\r\n")
	}

	file.Close()
	return buffer.String(), nil
}

func Contains(arrays []interface{}, target interface{}) bool {
	for _, element := range arrays {
		if element == target {
			return true
		}
	}
	return false
}

func GetValue(value string) string {

	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		name := value[2 : len(value)-1]
		property, ok := property.DefaultManager().GetProperty(name)
		if ok {
			strProperty, isString := property.(string)
			if isString {
				return strProperty
			}
		}
	}

	return value
}
