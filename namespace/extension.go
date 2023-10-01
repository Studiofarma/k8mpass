package namespace

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
)

type NamespaceExtension struct {
	Name         string
	ExtendSingle ExtendSingleFunc
	ExtendList   ExtendListFunc
}

type NamespaceName string
type ExtensionValue string

type ExtendSingleFunc func(ns v1.Namespace) (ExtensionValue, error)
type ExtendListFunc func(ns []v1.Namespace) map[NamespaceName]ExtensionValue

func NamespaceAge(ns v1.Namespace) (ExtensionValue, error) {
	res := fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
	return ExtensionValue(res), nil
}

func NamespaceAgeList(ns []v1.Namespace) map[NamespaceName]ExtensionValue {
	values := make(map[NamespaceName]ExtensionValue)
	for _, n := range ns {
		age, _ := NamespaceAge(n)
		values[NamespaceName(n.Name)] = age
	}
	return values
}

var NamespaceAgeProperty = NamespaceExtension{
	Name:         "age",
	ExtendSingle: NamespaceAge,
	ExtendList:   NamespaceAgeList,
}
