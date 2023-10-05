package pod

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"io"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
)

type Item struct {
	K8sPod             v1.Pod
	ExtendedProperties []Property
}

type Property struct {
	Key   string
	Value string
	Order int
}

func (n Item) FilterValue() string {
	return n.K8sPod.Name
}

type ItemDelegate struct{}

func (n ItemDelegate) Height() int {
	return 1
}

func (n ItemDelegate) Spacing() int {
	return 0
}

func (n ItemDelegate) Update(_ tea.Msg, m *list.Model) tea.Cmd {
	if len(m.Items()) == 0 {
		m.SetShowStatusBar(false)
	} else {
		m.SetShowStatusBar(true)
	}
	return nil
}

func (n ItemDelegate) Render(w io.Writer, _ list.Model, _ int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	_, _ = fmt.Fprintf(w, "  %s%s", styleString(i.K8sPod.Name, podStyle(i.K8sPod.Status)), customPropertiesStyle.Render(PodVersionSingle(i.K8sPod)))
}

func FindPod(items []list.Item, search Item) int {
	var idx = -1
	for i, item := range items {
		if pod, ok := item.(Item); ok {
			if pod.K8sPod.Name == search.K8sPod.Name {
				idx = i
			}
		}
	}
	return idx
}

var apps = []string{"backend", "sf-full", "spring-batch-ita", "spring-batch-deu"}

func PodVersionSingle(ns v1.Pod) string {
	if !slices.Contains(apps, ns.Labels["app"]) {
		return ""
	}
	appVersion := ns.Labels["AppVersion"]
	if appVersion == "" {
		appVersion = ns.Annotations["AppVersion"]
	}
	version, err := semver.NewVersion(appVersion)
	if err != nil {
		return ""
	}
	if version.Major() > 0 {
		return "v" + version.String()
	}
	return fmt.Sprintf("(%s)", version.Prerelease())
}
