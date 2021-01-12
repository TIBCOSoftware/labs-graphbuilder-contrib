/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package query

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/TIBCOSoftware/labs-graphbuilder-lib/model"
	"github.com/TIBCOSoftware/labs-graphbuilder-lib/tools"
)

//-====================-//
//  Define Operand
//-====================-//

type Operand struct {
	opType string
	value  interface{}
}

func (this *Operand) SetOpType(opType string) {
	this.opType = opType
}

func (this *Operand) SetValue(value interface{}) {
	this.value = value
}

func (this *Operand) Eval(data map[string]*model.Attribute) interface{} {
	if "*Expression" == this.opType {
		return (*this.value.(*Expression)).Eval(data)
	} else if "Attribute" == this.opType {
		return data[this.value.(string)].GetValue()
	}
	return this.value
}

func NewOperand(operandConf interface{}) Operand {
	operand := Operand{}
	confDataType := reflect.TypeOf(operandConf).String()
	if "map[string]interface {}" == confDataType {
		expr := NewExpression()
		(*expr).Build(operandConf.(map[string]interface{}))
		operand.SetOpType("*Expression")
		operand.SetValue(expr)
	} else if "string" == confDataType {
		strConfData := operandConf.(string)
		if strings.HasPrefix(strConfData, "$") {
			operand.SetOpType("Attribute")
			operand.SetValue(strConfData[1:])
		} else if "now()" == strConfData {
			operand.SetOpType("Const")
			operand.SetValue(tools.GetClock().GetCurrentTime())
		} else {
			operand.SetOpType("Const")
			operand.SetValue(operandConf)
		}
	}

	return operand
}

//-====================-//
//  Define Operator
//-====================-//

type Operator interface {
	GetName() string
	Operate(operands []Operand, data map[string]*model.Attribute) interface{}
}

//-====================-//
//  Define BaseOperator
//-====================-//

type BaseOperator struct {
	name string
}

func (this *BaseOperator) GetName() string {
	return this.name
}

//-==================-//
//    Define and
//-==================-//

type And struct {
	BaseOperator
}

func (this *And) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	return operands[0].Eval(data).(bool) && operands[1].Eval(data).(bool)
}

//-==================-//
//    Define or
//-==================-//

type Or struct {
	BaseOperator
}

func (this *Or) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	return operands[0].Eval(data).(bool) || operands[1].Eval(data).(bool)
}

//-==================-//
//  Define equals
//-==================-//

type Equal struct {
	BaseOperator
}

func (this *Equal) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	//	fmt.Println("Equal -> ", data)
	//	fmt.Println("Equal -> ", operands[0], operands[1])
	//	fmt.Println("Equal -> ", operands[0].Eval(data), operands[1].Eval(data))
	return operands[0].Eval(data) == operands[1].Eval(data)
}

//-==================-//
//  Define greater
//-==================-//

type Greater struct {
	BaseOperator
}

func (this *Greater) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	op0 := operands[0].Eval(data)
	op1 := operands[1].Eval(data)

	//	fmt.Println("Greater -> ", operands[0], operands[1])
	//	fmt.Println("Greater -> ", op0, op1)

	switch op0.(type) {
	case time.Time:
		return op0.(time.Time).Unix() > op1.(time.Time).Unix()
	case int32:
		return op0.(int32) > op1.(int32)
	case int64:
		return op0.(int64) > op1.(int64)
	case int:
		return op0.(int) > op1.(int)
	case float32:
		return op0.(float32) > op1.(float32)
	case float64:
		return op0.(float64) > op1.(float64)
	case string:
		return op0.(string) > op1.(string)
	}
	return false
}

//-==================-//
//  Define greater
//-==================-//

type Less struct {
	BaseOperator
}

func (this *Less) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	op0 := operands[0].Eval(data)
	op1 := operands[1].Eval(data)

	switch op0.(type) {
	case time.Time:
		return op0.(int64) < op1.(int64)
	case int32:
		return op0.(int32) < op1.(int32)
	case int64:
		return op0.(int64) < op1.(int64)
	case int:
		return op0.(int) < op1.(int)
	case float32:
		return op0.(float32) < op1.(float32)
	case float64:
		return op0.(float64) < op1.(float64)
	case string:
		return op0.(string) < op1.(string)
	}
	return false
}

//-==================-//
//  Define nop
//-==================-//

type Nop struct {
	BaseOperator
}

func (this *Nop) Operate(operands []Operand, data map[string]*model.Attribute) interface{} {
	//fmt.Println("Operate -> ", this.GetName())
	return true
}

func CreateOperator(name string) Operator {

	switch name {
	case "eq":
		op := &Equal{}
		op.name = "Equal"
		return op
	case "gt":
		op := &Greater{}
		op.name = "Greater"
		return op
	case "lt":
		op := &Less{}
		op.name = "Less"
		return op
	case "and":
		op := &And{}
		op.name = "And"
		return op
	case "or":
		op := &Or{}
		op.name = "Or"
		return op
	case "nop":
		op := &Nop{}
		op.name = "Nop"
		return op
	}

	return nil
}

//-====================-//
//  Define Expression
//-====================-//

type Expression struct {
	operator Operator
	operands []Operand
}

func (this *Expression) Build(jsonExpr map[string]interface{}) {
	//	fmt.Println("Build expression : ", jsonExpr)
	this.operator = CreateOperator(jsonExpr["operator"].(string))
	this.operands = make([]Operand, 2)
	jsonOperands := jsonExpr["operands"].([]interface{})

	this.operands[0] = NewOperand(jsonOperands[0])
	this.operands[1] = NewOperand(jsonOperands[1])
	//	fmt.Println("Build expression : ", this.ToString())
}

func (this *Expression) Eval(data map[string]*model.Attribute) bool {
	return this.operator.Operate(this.operands, data).(bool)
}

func (this *Expression) ToString() string {
	if "*query.Nop" != reflect.TypeOf(this.operator).String() {
		return fmt.Sprintf("Expr : %s %s %s", this.operator, this.operands[0], this.operands[1])
	} else {
		return fmt.Sprintf("Expr : %s nil nil", this.operator)
	}
}

func NewExpression() *Expression {
	return &Expression{}
}
