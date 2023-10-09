package api

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/studiofarma/k8mpass/api"
	v1 "k8s.io/api/core/v1"
	"slices"
)

var apps = []string{"backend", "sf-full", "spring-batch-ita", "spring-batch-deu"}

var PodVersion = api.PodExtension{
	Name:         "version",
	ExtendSingle: PodVersionSingle,
	ExtendList:   PodVersionList,
}

func PodVersionSingle(pod v1.Pod) (string, error) {
	if !slices.Contains(apps, pod.Labels["app"]) {
		return "", nil
	}
	appVersion := pod.Labels["AppVersion"]
	if appVersion == "" {
		appVersion = pod.Annotations["AppVersion"]
	}
	version, err := semver.NewVersion(appVersion)
	if err != nil {
		return "", nil
	}
	if version.Major() > 0 {
		return fmt.Sprintf("(v%s)", version.String()), nil
	}
	return fmt.Sprintf("(%s)", version.Prerelease()), nil
}
func PodVersionList(ns []v1.Pod) map[string]string {
	res := make(map[string]string, len(ns))
	for _, n := range ns {
		ext, err := PodVersionSingle(n)
		if err != nil {
			continue
		}
		res[n.Name] = ext
	}
	return res
}
