package main

import (
	"github.com/Masterminds/semver/v3"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	"slices"
)

var apps = []string{"backend"}

var PodVersion = api.Extension{
	Name:         "version",
	ExtendSingle: PodVersionSingle,
	ExtendList:   PodVersionList,
}

func PodVersionSingle(ns v1.Namespace) (api.ExtensionValue, error) {
	if !slices.Contains(apps, ns.Labels["app"]) {
		return "", nil
	}
	version, err := semver.NewVersion(ns.Labels["AppVersion"])
	if err != nil {
		return "", err
	}
	if version.Major() > 0 {
		return api.ExtensionValue("v" + version.String()), nil
	}
	return api.ExtensionValue("commit " + version.Prerelease()), nil
}
func PodVersionList(ns []v1.Namespace) map[api.Name]api.ExtensionValue {
	res := make(map[api.Name]api.ExtensionValue, len(ns))
	for _, n := range ns {
		ext, err := PodVersionSingle(n)
		if err != nil {
			continue
		}
		res[api.Name(n.Name)] = ext
	}
	return res
}
