/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package table

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"

	"log"
	"os"
	"sync"
)

type TableManager struct {
	tables map[string]*Table
}

var (
	instance *TableManager
	once     sync.Once
	mux      sync.Mutex
)

func GetTableManager() *TableManager {
	once.Do(func() {
		instance = &TableManager{tables: make(map[string]*Table)}
	})
	return instance
}

func (this *TableManager) GetTable(tablename string) *Table {
	return this.tables[tablename]
}

func (this *TableManager) CreateTable(
	pKey []string,
	tablename string,
	tableSchema *Schema) *Table {

	table := this.tables[tablename]
	if nil == table {
		mux.Lock()
		defer mux.Unlock()
		table = this.tables[tablename]
		if nil == table {
			table = &Table{
				pKey:          pKey,
				indices:       make([][]string, 0),
				theMap:        make(map[CompositKey]*Record),
				theIndexedKey: make(map[CompositKey]map[CompositKey][]CompositKey),
				tableSchema:   tableSchema,
			}
			this.tables[tablename] = table
		}
	}

	return table
}

type Table struct {
	pKey          []string
	indices       [][]string
	tableSchema   *Schema
	theMap        map[CompositKey]*Record
	theIndexedKey map[CompositKey]map[CompositKey][]CompositKey
}

func (this *Table) AddIndex(keyName []string) bool {
	key, _ := ConstructKey(keyName, nil)
	this.indices = append(this.indices, keyName)
	this.theIndexedKey[key] = make(map[CompositKey][]CompositKey)
	log.Println("table after add index: indices = ", this.indices, ", table = ", this.theMap)
	return true
}

func (this *Table) RemoveIndex(keyName []string) bool {
	key, _ := ConstructKey(keyName, nil)
	this.theMap[key] = nil
	return true
}

func (this *Table) GetPkeyNames() []string {
	return this.pKey
}

func (this *Table) Get(searchKey []string, data map[string]interface{}) ([]*Record, bool) {
	log.Println("searchKey : ", searchKey)
	log.Println("data : ", data)
	log.Println("pKey : ", this.pKey)
	log.Println("theMap : ", this.theMap)

	records := make([]*Record, 0)
	pKeyHash, pKeyValueHash := ConstructKey(this.pKey, data)
	searchKeyHash, searchKeyValueHash := ConstructKey(searchKey, data)

	log.Println("pKeyValueHash : ", pKeyValueHash)
	log.Println("searchKeyValueHash : ", searchKeyValueHash)

	searchByPKey := true
	if searchKeyHash == pKeyHash {
		log.Println("Get by primary key !")
		record := this.theMap[pKeyValueHash]
		if nil != record {
			records = append(records, record)
		}
	} else {
		log.Println("Get by indexed search key !")
		searchByPKey = false
		pKeyValueHashs := this.theIndexedKey[searchKeyHash][searchKeyValueHash]
		if nil != pKeyValueHashs {
			newPKeyValueHashs := make([]CompositKey, 0)
			for _, pKeyValueHash := range pKeyValueHashs {
				if nil != this.theMap[pKeyValueHash] {
					records = append(records, this.theMap[pKeyValueHash])
					newPKeyValueHashs = append(newPKeyValueHashs, pKeyValueHash)
				}
			}
			this.theIndexedKey[searchKeyHash][searchKeyValueHash] = newPKeyValueHashs
		}
	}

	return records, searchByPKey
}

func (this *Table) Upsert(data map[string]interface{}) *Record {
	log.Println("data : ", data)
	log.Println("pKey : ", this.pKey)
	log.Println("theMap before : ", this.theMap)

	_, pKeyValueHash := ConstructKey(this.pKey, data)
	record := this.theMap[pKeyValueHash]

	if nil != record {
		for _, fieldInfo := range *this.tableSchema.DataSchemas() {
			fieldName := fieldInfo["Name"].(string)
			fieldValue := data[fieldName]
			if nil != fieldValue {
				(*record)[fieldName] = fieldValue
			}
		}
	} else {
		// Create new record
		record = &Record{}
		for _, fieldInfo := range *this.tableSchema.DataSchemas() {
			(*record)[fieldInfo["Name"].(string)] = data[fieldInfo["Name"].(string)]
		}
		this.theMap[pKeyValueHash] = record

		// Indexing record
		for _, index := range this.indices {
			indexHash, indexValueHash := ConstructKey(index, data)
			pKeyValueHashs := this.theIndexedKey[indexHash][indexValueHash]
			if nil != pKeyValueHashs {
				this.theIndexedKey[indexHash][indexValueHash] = append(this.theIndexedKey[indexHash][indexValueHash], pKeyValueHash)
			} else {
				this.theIndexedKey[indexHash][indexValueHash] = []CompositKey{pKeyValueHash}
			}
		}
	}

	log.Println("theMap after : ", this.theMap)

	return record
}

func (this *Table) Delete(data map[string]interface{}) *Record {
	log.Println("data : ", data)
	log.Println("pKey : ", this.pKey)
	log.Println("theMap before : ", this.theMap)

	_, pKeyValueHash := ConstructKey(this.pKey, data)
	record := this.theMap[pKeyValueHash]

	if nil != record {
		delete(this.theMap, pKeyValueHash)
	}

	log.Println("theMap after : ", this.theMap)

	return record
}

func (this *Table) RowCount() int {
	return len(this.theMap)
}

func (this *Table) Load(file *os.File) {
}

func (this *Table) SaveSchema(file *os.File) {
}

func (this *Table) SaveData(file *os.File) {
}

func (this *Table) GenerateKeys(arr []string, data []string, start int, end int, index int, r int) {
	log.Println("GenerateKeys, index = ", index, ", r = ", r, ", arr", arr)
	if index == r {
		log.Println("GenerateKeys, data = ", data)
		key := make([]string, 0)
		for j := 0; j < r; j++ {
			key = append(key, data[j])
		}
		this.AddIndex(key)
		return
	}

	i := start
	for i <= end && end-i+1 >= r-index {
		data[index] = arr[i]
		this.GenerateKeys(arr, data, i+1, end, index+1, r)
		i += 1
	}
}

type Record map[string]interface{}

func CreateSchema(schema *[]map[string]interface{}) *Schema {
	return &Schema{
		schema: schema,
	}
}

type Schema struct {
	schema *[](map[string]interface{})
}

func (this *Schema) DataSchemas() *[](map[string]interface{}) {
	return this.schema
}

func (this *Schema) Length() int {
	return len(*this.schema)
}

type CompositKey struct {
	Id uint64
}

func KeyFromDataArray(elements []interface{}) CompositKey {
	keyBytes := []byte{}
	for _, element := range elements {
		elementBytes, _ := json.Marshal(element)
		keyBytes = append(keyBytes, elementBytes...)
	}
	hasher := md5.New()
	hasher.Write(keyBytes)
	return CompositKey{Id: binary.BigEndian.Uint64(hasher.Sum(nil))}
}

func ConstructKey(keyNameStrs []string, tuple map[string]interface{}) (CompositKey, CompositKey) {
	log.Println("keyNameStrs : ", keyNameStrs)
	log.Println("tuple : ", tuple)

	/* build key */
	key := make([]interface{}, len(keyNameStrs))
	keyFields := make([]interface{}, len(keyNameStrs))
	for j, keyNameStr := range keyNameStrs {
		key[j] = tuple[keyNameStr]
		keyFields[j] = keyNameStr
	}
	log.Println("keyFields : ", keyFields)
	log.Println("key : ", key)

	return KeyFromDataArray(keyFields), KeyFromDataArray(key)
}

/*
func combination(arr []string, data []string, start int, end int, index int, r int, combs []string) {
    if (index == r) {
        comb := ""
        for j:=0; j<r; j++ {
            if(j!=0)
            	comb = fmt.sprint("%s ", comb)
            comb += data[j]
        }
        combs = append(combs, comb)
        return
    }

    i := start
    while i <= end && (end- i + 1) >= (r - index) {
        data[index] = arr[i]
        combination(arr, data, i + 1, end, index + 1, r, combs)
        i += 1
    }
}*/
