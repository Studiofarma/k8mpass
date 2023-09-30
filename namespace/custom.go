package namespace

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
)

type NamespaceCustomProperty struct {
	Name string
	Func CustomPropertyFunc
}

type CustomPropertyFunc func(ns *v1.Namespace) string

func NamespaceAge(ns *v1.Namespace) string {
	return fmt.Sprintf("Age: %0.f minutes", time.Now().Sub(ns.CreationTimestamp.Time).Minutes())
}

var NamespaceAgeProperty = NamespaceCustomProperty{
	Name: "age",
	Func: NamespaceAge,
}
