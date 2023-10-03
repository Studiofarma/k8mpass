package namespace

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
)

type Extension struct {
	Name         string
	ExtendSingle ExtendSingleFunc
	ExtendList   ExtendListFunc
}

type Name string
type ExtensionValue string

type ExtendSingleFunc func(ns v1.Namespace) (ExtensionValue, error)
type ExtendListFunc func(ns []v1.Namespace) map[Name]ExtensionValue

func Age(ns v1.Namespace) (ExtensionValue, error) {
	res := fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
	return ExtensionValue(res), nil
}

func AgeList(ns []v1.Namespace) map[Name]ExtensionValue {
	values := make(map[Name]ExtensionValue)
	for _, n := range ns {
		age, _ := Age(n)
		values[Name(n.Name)] = age
	}
	return values
}

var AgeProperty = Extension{
	Name:         "age",
	ExtendSingle: Age,
	ExtendList:   AgeList,
}
