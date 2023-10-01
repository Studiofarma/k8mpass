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

type ExtendSingleFunc func(ns v1.Namespace) string
type ExtendListFunc func(ns []v1.Namespace) map[string]string

func NamespaceAge(ns v1.Namespace) string {
	return fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
}

var NamespaceAgeProperty = NamespaceExtension{
	Name:         "age",
	ExtendSingle: NamespaceAge,
}
