package api

import v1 "k8s.io/api/core/v1"

type Extension struct {
	Name         string
	ExtendSingle ExtendSingleFunc
	ExtendList   ExtendListFunc
}

type Name string
type ExtensionValue string

type ExtendSingleFunc func(ns v1.Namespace) (ExtensionValue, error)
type ExtendListFunc func(ns []v1.Namespace) map[Name]ExtensionValue
