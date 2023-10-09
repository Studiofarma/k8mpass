package api

import (
	"fmt"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	"math"
	"time"
)

func NamespaceAge(ns v1.Namespace) (string, error) {
	return ResourceAge(ns.CreationTimestamp.Time), nil
}

func NamespaceAgeList(ns []v1.Namespace) map[string]string {
	values := make(map[string]string)
	for _, n := range ns {
		age, _ := NamespaceAge(n)
		values[n.Name] = age
	}
	return values
}

var NamespaceAgeProperty = api.NamespaceExtension{
	Name:         "age",
	ExtendSingle: NamespaceAge,
	ExtendList:   NamespaceAgeList,
}
var PodAgeProperty = api.PodExtension{
	Name:         "age",
	ExtendSingle: PodAge,
	ExtendList:   PodAgeList,
}

func ResourceAge(creation time.Time) string {
	t := time.Now().Sub(creation)
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
	return s

}

func PodAge(pod v1.Pod) (string, error) {
	return ResourceAge(pod.CreationTimestamp.Time), nil
}

func PodAgeList(pods []v1.Pod) map[string]string {
	res := make(map[string]string, len(pods))
	for _, pod := range pods {
		property, err := PodAge(pod)
		if err != nil {
			continue
		}
		res[pod.Name] = property
	}
	return res
}
