package pod

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/lipgloss"
	"io"
	"math"
	"slices"
	"time"

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

func (n ItemDelegate) Render(w io.Writer, l list.Model, _ int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}
	i.LoadProperties()
	customProperties := ""
	for _, property := range i.ExtendedProperties {
		customProperties += customPropertiesStyle.Render(fmt.Sprintf(property.Value))
	}
	maxLength := 0
	for _, item := range l.Items() {
		maxLength = max(maxLength, len(item.FilterValue()))
	}
	_, _ = fmt.Fprintf(w, "  %s%s", podStyle(i.K8sPod.Status, maxLength).Render(i.K8sPod.Name), customProperties)
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

func (p *Item) LoadProperties() {
	var properties []Property

	properties = append(properties, Property{
		Key:   "age",
		Value: PodAge(p.K8sPod),
		Order: 0,
	})

	properties = append(properties, Property{
		Key:   "version",
		Value: PodVersionSingle(p.K8sPod),
		Order: 0,
	})
	p.ExtendedProperties = properties
}

func PodVersionSingle(pod v1.Pod) string {
	if !slices.Contains(apps, pod.Labels["app"]) {
		return ""
	}
	appVersion := pod.Labels["AppVersion"]
	if appVersion == "" {
		appVersion = pod.Annotations["AppVersion"]
	}
	version, err := semver.NewVersion(appVersion)
	if err != nil {
		return ""
	}
	if version.Major() > 0 {
		return fmt.Sprintf("(v%s)", version.String())
	}
	return fmt.Sprintf("(%s)", version.Prerelease())
}

func PodAge(pod v1.Pod) string {
	time := time.Now().Sub(pod.CreationTimestamp.Time)
	var res float64
	var unit string
	if time.Minutes() < 60 {
		res = time.Minutes()
		unit = "m"
	} else if time.Hours() < 24 {
		res = time.Hours()
		unit = "h"
	} else {
		res = time.Hours() / 24
		unit = "d"
	}
	s := fmt.Sprintf("%0.f%s", math.Floor(res), unit)
	return lipgloss.NewStyle().Width(3).Render(s)
}
