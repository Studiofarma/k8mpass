package main

import "k8s.io/apimachinery/pkg/types"

type Namespace struct {
	id   types.UID
	name string
}

func (i Namespace) Title() string       { return i.name }
func (i Namespace) Description() string { return "" }
func (i Namespace) FilterValue() string { return i.name }
