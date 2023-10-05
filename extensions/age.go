package main

import (
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	"math"
	"time"
)

func Age(ns v1.Namespace) (string, error) {
	res := fmt.Sprintf("Age: %0.f minutes", time.Since(ns.CreationTimestamp.Time).Minutes())
	return res, nil
}

func AgeList(ns []v1.Namespace) map[string]string {
	values := make(map[string]string)
	for _, n := range ns {
		age, _ := Age(n)
		values[n.Name] = age
	}
	return values
}

var NamespaceAgeProperty = api.NamespaceExtension{
	Name:         "age",
	ExtendSingle: Age,
	ExtendList:   AgeList,
}
var PodAgeProperty = api.PodExtension{
	Name:         "age",
	ExtendSingle: PodAgeSingle,
	ExtendList:   PodAgeList,
}

func PodAgeSingle(pod v1.Pod) (string, error) {
	t := time.Now().Sub(pod.CreationTimestamp.Time)
	var res float64
	var unit string
	if t.Minutes() < 60 {
		res = t.Minutes()
		unit = "m"
	} else if t.Hours() < 24 {
		res = t.Hours()
		unit = "h"
	} else {
		res = t.Hours() / 24
		unit = "d"
	}
	s := fmt.Sprintf("%0.f%s", math.Floor(res), unit)
	return s, nil
}

func PodAgeList(pods []v1.Pod) map[string]string {
	res := make(map[string]string, len(pods))
	for _, pod := range pods {
		property, err := PodAgeSingle(pod)
		if err != nil {
			continue
		}
		res[pod.Name] = property
	}
	return res
}
