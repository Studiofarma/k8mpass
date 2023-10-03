package main

import (
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	"time"
)

func Age(ns v1.Namespace) (api.ExtensionValue, error) {
	res := fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
	return api.ExtensionValue(res), nil
}

func AgeList(ns []v1.Namespace) map[api.Name]api.ExtensionValue {
	values := make(map[api.Name]api.ExtensionValue)
	for _, n := range ns {
		age, _ := Age(n)
		values[api.Name(n.Name)] = age
	}
	return values
}

var AgeProperty = api.Extension{
	Name:         "age",
	ExtendSingle: Age,
	ExtendList:   AgeList,
}
