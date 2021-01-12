package cache

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	cache := NewCache(3)
	datas := []Data{
		NewData("A", "A"),
		NewData("B", "B"),
		NewData("C", "C"),
		NewData("A", "A"),
		NewData("D", "D"),
		NewData("E", "E"),
		NewData("F", "F"),
		NewData("B", "B"),
		NewData("A", "A"),
		NewData("A", "A"),
		NewData("D", "D"),
		NewData("B", "B"),
		NewData("D", "D"),
	}

	for _, data := range datas {
		d := cache.Get(data.id)
		fmt.Println("Get Data = (", data.id, ", ", data.value, ")")
		if nil == d {
			cache.Add(data.id, data)
			fmt.Println("Cache miss !!!!")
		} else {
			fmt.Println("Cache hit, RTN data = (", d.(Data).id, ", ", d.(Data).value, ")")
		}
	}
}

type Data struct {
	id    string
	value string
}

func NewData(id string, value string) Data {
	return Data{id: id, value: value}
}
