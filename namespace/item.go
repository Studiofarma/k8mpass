package namespace

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/termenv"
	"github.com/studiofarma/k8mpass/api"
	"io"
	v1 "k8s.io/api/core/v1"
	"log"
	"runtime"
	"slices"
	"sort"
)

type Item struct {
	K8sNamespace       v1.Namespace
	ExtendedProperties []Property
}

type Property struct {
	Key   string
	Value string
	Order int
}

func (n Item) FilterValue() string {
	return n.K8sNamespace.Name
}

type ItemDelegate struct {
	Pinned []string
}

func (n ItemDelegate) Height() int {
	return 1
}

func (n ItemDelegate) Spacing() int {
	return 0
}

func (n ItemDelegate) Update(tea.Msg, *list.Model) tea.Cmd {
	return nil
}

func (n ItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	namespace := ellipsis(i.K8sNamespace.Name, nameMaxLength)
	properties := ""
	nsStyle := commonStyle.Copy()
	propertiesStyle := customPropertiesStyle.Copy()
	if slices.Contains(n.Pinned, i.K8sNamespace.Name) {
		nsStyle = nsStyle.Inherit(pinnedStyle)
		namespace = "*" + namespace
		nsStyle.MarginLeft(nsStyle.GetMarginLeft() - 1)
	}
	if index == m.Index() {
		nsStyle = nsStyle.Inherit(selectedItemStyle).UnsetMarginBackground()
		propertiesStyle = propertiesStyle.Inherit(selectedItemStyle).UnsetMarginBackground()
	}
	nsStyle = nsStyle.Inherit(namespaceStatusStyle(i))
	for _, property := range i.ExtendedProperties {
		p := ellipsis(property.Value, propertyMaxWidth)
		properties += propertiesStyle.Render(p)
	}
	_, _ = fmt.Fprint(w, nsStyle.Render(namespace)+properties)
}

func ellipsis(s string, length int) string {
	if len(s) > length {
		return truncate.StringWithTail(s, uint(length), "..")
	}
	return s
}

func (n *Item) LoadCustomProperties(extensions ...api.INamespaceExtension) {
	var properties = make([]Property, len(extensions))
	n.ExtendedProperties = make([]Property, 0)
	for _, p := range extensions {
		fn := p.GetExtendSingle()
		if fn == nil {
			log.Println(fmt.Sprintf("Missing extention function for %s", p.GetName()), "namespace:", n.K8sNamespace.Name)
			continue
		}
		value, err := fn(n.K8sNamespace)
		if err != nil {
			log.Println(fmt.Sprintf("Error while computing extension %s", p.GetName()), "namespace:", n.K8sNamespace.Name)
			continue
		}
		properties = append(properties, Property{Key: p.GetName(), Value: value, Order: p.GetOrder()})
	}
	sort.SliceStable(properties, func(i, j int) bool {
		return properties[i].Order > properties[j].Order
	})
	n.ExtendedProperties = properties
}

func namespaceStatusStyle(ns Item) lipgloss.Style {
	style := lipgloss.NewStyle()
	if ns.K8sNamespace.Status.Phase == v1.NamespaceTerminating {
		style = terminatingNamespaceStyle
	}
	return style
}

func slightlyBrighterTerminalColor() lipgloss.Color {
	switch runtime.GOOS {
	case "windows":
		return lipgloss.Color("#444852")
	}
	multiplier := 1.5
	terminalColor := termenv.BackgroundColor()
	rgb := termenv.ConvertToRGB(terminalColor)
	r, g, b := colorful.Color{
		R: rgb.R * multiplier,
		G: rgb.G * multiplier,
		B: rgb.B * multiplier,
	}.RGB255()
	hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return lipgloss.Color(hex)
}
