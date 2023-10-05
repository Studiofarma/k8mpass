package api

import v1 "k8s.io/api/core/v1"

type INamespaceExtension interface {
	GetName() string
	GetExtendSingle() NamespaceExtendSingleFunc
	GetExtendList() NamespaceExtendListFunc
}

type NamespaceExtension struct {
	Name         string
	ExtendSingle NamespaceExtendSingleFunc
	ExtendList   NamespaceExtendListFunc
}

func (e NamespaceExtension) GetName() string {
	return e.Name
}

func (e NamespaceExtension) GetExtendSingle() NamespaceExtendSingleFunc {
	return e.ExtendSingle
}

func (e NamespaceExtension) GetExtendList() NamespaceExtendListFunc {
	return e.ExtendList
}

type NamespaceExtendSingleFunc func(ns v1.Namespace) (string, error)
type NamespaceExtendListFunc func(ns []v1.Namespace) map[string]string
type IPodExtension interface {
	GetName() string
	GetExtendSingle() PodExtendSingleFunc
	GetExtendList() PodExtendListFunc
}

type PodExtension struct {
	Name         string
	ExtendSingle PodExtendSingleFunc
	ExtendList   PodExtendListFunc
}

func (e PodExtension) GetName() string {
	return e.Name
}

func (e PodExtension) GetExtendSingle() PodExtendSingleFunc {
	return e.ExtendSingle
}

func (e PodExtension) GetExtendList() PodExtendListFunc {
	return e.ExtendList
}

type PodExtendSingleFunc func(pod v1.Pod) (string, error)
type PodExtendListFunc func(pods []v1.Pod) map[string]string
