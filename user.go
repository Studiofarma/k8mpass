package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"slices"
)

type IPinnedNamespacesService interface {
	GetNamespaces() []string
	Pin(ns string)
	Unpin(ns string)
}

type PinnedNamespaceService struct {
	namespaces []string
}

func (p *PinnedNamespaceService) LoadSavedNamespaces() error {
	namespaces, err := LoadConfigFile()
	if err != nil {
		return err
	}
	p.namespaces = namespaces
	return nil
}

func (p PinnedNamespaceService) GetNamespaces() []string {
	return p.namespaces
}

func (p *PinnedNamespaceService) Pin(ns string) error {
	if slices.Contains(p.namespaces, ns) {
		return nil
	}
	p.namespaces = append(p.namespaces, ns)
	err := WriteLineToConfigFile(ns)
	if err != nil {
		return err
	}
	return nil
}

func (p *PinnedNamespaceService) Unpin(ns string) error {
	p.namespaces = slices.DeleteFunc(p.namespaces, func(s string) bool {
		return s == ns
	})
	err := RemoveLineFromFile(ns)
	if err != nil {
		return err
	}
	return nil
}

func WriteLineToConfigFile(ns string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(userHomeDir, ".k8mpass", "config.txt")

	readFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	_, err = readFile.WriteString(ns + "\n")
	if err != nil {
		return err
	}
	return nil
}

func RemoveLineFromFile(ns string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(userHomeDir, ".k8mpass", "config.txt")

	readFile, err := os.OpenFile(configPath, os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer readFile.Close()
	if err != nil {
		return err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line != ns {
			buf.WriteString(line + "\n")
		}
	}
	err = readFile.Truncate(0)
	if err != nil {
		return err
	}
	_, err = readFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(readFile)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfigFile() ([]string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(userHomeDir, ".k8mpass", "config.txt")

	readFile, err := os.Open(configPath)
	if err != nil {
		readFile, err = os.Create(configPath)
		if err != nil {
			return nil, err
		}
	}
	defer readFile.Close()
	if err != nil {
		return nil, err
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines, nil
}
