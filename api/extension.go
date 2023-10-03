package api

import v1 "k8s.io/api/core/v1"

type IExtension interface {
	GetName() string
	GetExtendSingle() ExtendSingleFunc
	GetExtendList() ExtendListFunc
}

type Extension struct {
	Name         string
	ExtendSingle ExtendSingleFunc
	ExtendList   ExtendListFunc
}

func (e Extension) GetName() string {
	return e.Name
}

func (e Extension) GetExtendSingle() ExtendSingleFunc {
	return e.ExtendSingle
}

func (e Extension) GetExtendList() ExtendListFunc {
	return e.ExtendList
}

type Name string
type ExtensionValue string

type ExtendSingleFunc func(ns v1.Namespace) (ExtensionValue, error)
type ExtendListFunc func(ns []v1.Namespace) map[Name]ExtensionValue
