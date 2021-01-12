/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/core/data/coerce"
)

var log = logger.GetLogger("activity-aggregate")

func Agg(aggStateOfQuery map[Index]map[string]DataState, data map[string]interface{}) map[string]interface{} {

	var outputTuple map[string]interface{}

	var dataKey []interface{}
	var agg map[string]interface{}
	var groupBy []interface{}
	if nil != data["parameter"] {
		parameter := data["parameter"].(map[string]interface{})
		if nil != parameter["key"] {
			dataKey = parameter["key"].([]interface{})
		}

		if nil != parameter["agg"] {
			agg = parameter["agg"].(map[string]interface{})
		}

		if nil != parameter["index"] {
			groupBy = parameter["index"].([]interface{})
		}

		log.Info("******** Key of output tuple ******** dataKey = ", dataKey)
		log.Info("******** Aggregation function ******** agg = ", agg)
		log.Info("******** Group key for aggrate ******** groupBy = ", groupBy)

		for k, v := range data {
			if "parameter" != k {
				rawTuple := v.(map[string][]interface{})

				keys := make([]string, len(rawTuple))
				counter := 0
				for key, _ := range rawTuple {
					keys[counter] = key
					counter++
				}

				for _, tuple := range SplitTuple(0, keys, rawTuple) {
					/* aggregation group key */
					ok := true
					var groupByValue bytes.Buffer
					for _, keyElement := range groupBy {
						valueElement := tuple[keyElement.(string)]
						if nil == valueElement {
							log.Warn("Invalid groupKey : ", tuple)
							ok = false
							continue
						}
						groupByValue.WriteString(valueElement.(string))
						groupByValue.WriteString("_")
					}

					if !ok {
						continue
					}

					/* build output key for tuple */
					key := make([]interface{}, len(dataKey))
					for i, keyElement := range dataKey {
						valueElement := tuple[keyElement.(string)]
						if nil == valueElement {
							log.Warn("Invalid dataKey : ", tuple)
							ok = false
							continue
						}
						key[i] = valueElement.(string)
					}

					if !ok {
						continue
					}

					keyObj := NewIndex(key)
					aggregatedTuple := aggStateOfQuery[keyObj]
					if nil == aggregatedTuple {
						aggregatedTuple = make(map[string]DataState)
						aggStateOfQuery[keyObj] = aggregatedTuple
					}

					for _, keyElement := range dataKey {
						keyElementStr := keyElement.(string)
						data := aggregatedTuple[keyElementStr]
						if nil == data {
							data = &Data{}
							aggregatedTuple[keyElementStr] = data
						}
						data.Update(tuple[keyElementStr])
					}

					for valueColumn, functionNames := range agg {
						functionNamesArr := functionNames.([]interface{})
						for _, functionName := range functionNamesArr {
							functionNameStr := functionName.(string)
							dataKey, err := DataKey(groupByValue.String(), tuple, functionNameStr, valueColumn)
							if nil != err {
								continue
							}
							function := aggregatedTuple[dataKey]
							if nil == function {
								function = GetFunction(functionNameStr)
								aggregatedTuple[dataKey] = function
							}
							err = function.Update(tuple[valueColumn])
							if nil != err {
								log.Error("Error : ", err)
							}
						}
					}

					log.Debug("--- Tuple - ", tuple, ", aggregatedTuple - ", aggregatedTuple)
					outputTuple = make(map[string]interface{})
					for key, value := range aggregatedTuple {
						outputTuple[key] = value.Value()
					}
					log.Debug("--- Tuple - ", tuple, ", outputTuple - ", outputTuple)
				}
			}
		}
	}

	return outputTuple
}

type DataState interface {
	Update(newData interface{}) error
	Value() interface{}
}

func GetFunction(functionName string) DataState {
	var function DataState
	if "sum" == functionName {
		function = &Sum{}
	} else if "count" == functionName {
		function = &Count{}
	} else if "mean" == functionName {
		function = &Mean{}
	} else if "min" == functionName {
		function = &Min{}
	} else if "max" == functionName {
		function = &Max{}
	}
	return function
}

type Data struct {
	data interface{}
}

func (this *Data) Value() interface{} {
	return this.data
}

func (this *Data) Update(newData interface{}) error {
	this.data = newData
	return nil
}

type Sum struct {
	data float64
}

func (this *Sum) Value() interface{} {
	return this.data
}

func (this *Sum) Update(newData interface{}) error {
	delta, _ := coerce.ToFloat64(newData)
	this.data += delta
	return nil
}

type Count struct {
	counter int
}

func (this *Count) Value() interface{} {
	return this.counter
}

func (this *Count) Update(newData interface{}) error {
	this.counter += 1
	return nil
}

type Mean struct {
	sum   float64
	count float64
}

func (this *Mean) Value() interface{} {
	return this.sum / this.count
}

func (this *Mean) Update(newData interface{}) error {
	this.count += 1
	delta, err := coerce.ToFloat64(newData)
	if nil != err {
		return err
	}
	this.sum += delta
	return nil
}

type Min struct {
	min interface{}
}

func (this *Min) Value() interface{} {
	return this.min
}

func (this *Min) Update(newData interface{}) error {
	if nil == this.min {
		this.min = newData
		return nil
	}

	result, err := compare(newData, this.min)

	if nil != err {
		return err
	}

	if 0 > result {
		this.min = newData
	}
	return nil
}

type Max struct {
	max interface{}
}

func (this *Max) Value() interface{} {
	return this.max
}

func (this *Max) Update(newData interface{}) error {
	if nil == this.max {
		this.max = newData
		return nil
	}

	result, err := compare(newData, this.max)

	if nil != err {
		return err
	}

	if 0 < result {
		this.max = newData
	}
	return nil
}

func compare(data1 interface{}, data2 interface{}) (int, error) {

	switch data1.(type) {
	case float64:
		delta1float64, _ := coerce.ToFloat64(data1)
		delta2float64, err := coerce.ToFloat64(data2)
		if nil != err {
			return 0, err
		}
		delta := delta1float64 - delta2float64
		switch {
		case delta > 0:
			return 1, nil
		case delta == 0:
			return 0, nil
		case delta < 0:
			return -1, nil
		}
	case int:
		delta1int, _ := coerce.ToInt(data1)
		delta2int, err := coerce.ToInt(data2)
		if nil != err {
			return 0, err
		}
		delta := delta1int - delta2int
		switch {
		case delta > 0:
			return 1, nil
		case delta == 0:
			return 0, nil
		case delta < 0:
			return -1, nil
		}
	}

	return 0, errors.New("Unable to compare, Uknown type!")
}

type Index struct {
	Id uint64
}

func NewIndex(elements []interface{}) Index {
	keyBytes := []byte{}
	for _, element := range elements {
		elementBytes, _ := json.Marshal(element)
		keyBytes = append(keyBytes, elementBytes...)
	}
	hasher := md5.New()
	hasher.Write(keyBytes)
	return Index{Id: binary.BigEndian.Uint64(hasher.Sum(nil))}
}

func DataKey(
	groupByValue string,
	tuple map[string]interface{},
	functionName string,
	valueColumn string,
) (string, error) {
	var aggKey bytes.Buffer
	aggKey.WriteString(groupByValue)
	aggKey.WriteString(functionName)
	aggKey.WriteString("_")
	aggKey.WriteString(valueColumn)
	return aggKey.String(), nil
}

func SplitTuple(scope int, keys []string, data map[string][]interface{}) []map[string]interface{} {
	var temp []map[string]interface{}
	//fmt.Println("\n\n\n@@@@@@@@@@@@@@@@@@@@@@@@  scope = ", scope)

	if scope < len(keys)-1 {
		//fmt.Println("before down stream split scope = ", scope)
		rtn := SplitTuple(scope+1, keys, data)
		fieldValues := data[keys[scope]]
		//fmt.Println("fieldValues = ", fieldValues)
		size := len(fieldValues)
		if nil == rtn {
			//fmt.Println("1 rtn = ", rtn)
			temp = make([]map[string]interface{}, size)
			for i, fieldValue := range fieldValues {
				tuple := make(map[string]interface{})
				tuple[keys[i]] = fieldValue
				temp[i] = tuple
			}
			//fmt.Println("1 temp = ", temp)
		} else {
			//fmt.Println("2 rtn = ", rtn)
			temp = make([]map[string]interface{}, size*len(rtn))
			for i, fieldValue := range fieldValues {
				for j, rtnTuple := range rtn {
					tuple := make(map[string]interface{})
					for k, v := range rtnTuple {
						tuple[k] = v
						//fmt.Println("2.1 tuple = ", tuple)
					}
					//fmt.Println("2.2 tuple = ", tuple)
					tuple[keys[scope]] = fieldValue
					//fmt.Println("2.3 tuple = ", tuple)
					temp[i*size+j] = tuple
				}
			}
			//fmt.Println("2 temp = ", temp)
		}
	}

	//fmt.Println("before return temp =  ", temp)
	//fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@\n\n\n")

	return temp
}
