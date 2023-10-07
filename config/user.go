package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type userData struct {
	PinnedNamespaces []string
}

type IUserService interface {
	GetPinnedNamespaces() []string
	Pin(string)
	Unpin(string)
	Persist()
}

type UserService struct {
	data           userData
	currentContext string
}

func (u *UserService) Pin(ns string) {
	if !slices.Contains(u.data.PinnedNamespaces, ns) {
		u.data.PinnedNamespaces = append(u.data.PinnedNamespaces, ns)
	}
}

func (u *UserService) Unpin(ns string) {
	namespaces := u.data.PinnedNamespaces
	u.data.PinnedNamespaces = slices.DeleteFunc(namespaces, func(s string) bool {
		return s == ns
	})
}

func (u *UserService) GetPinnedNamespaces() []string {
	return u.data.PinnedNamespaces
}

func (u *UserService) Persist() {
	err := writeLines(u.currentContext, u.data.PinnedNamespaces)
	if err != nil {
		return
	}
}

func New(currentContext string) *UserService {
	data := loadFromFile(currentContext)
	return &UserService{data: data, currentContext: currentContext}
}

func loadFromFile(context string) userData {
	userDir, err := os.UserCacheDir()
	if err != nil {
		return userData{}
	}
	preferencesPath := filepath.Join(userDir, ".k8mpass", "namespaces")
	err = os.MkdirAll(preferencesPath, os.ModePerm)
	if err != nil {
		return userData{}
	}
	contextFile := fmt.Sprintf("%s_ns.txt", context)
	file, err := os.OpenFile(filepath.Join(preferencesPath, contextFile), os.O_RDONLY|os.O_CREATE, 0660)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	namespaces := readLines(file)
	if err != nil {
		return userData{}
	}
	return userData{PinnedNamespaces: namespaces}
}

func readLines(file *os.File) []string {
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	return lines
}

func writeLines(currentContext string, lines []string) error {
	userDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(userDir, ".k8mpass", "namespaces", fmt.Sprintf("%s_ns.txt", currentContext))
	readFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0660)
	if err != nil {
		return err
	}
	defer func(readFile *os.File) {
		_ = readFile.Close()
	}(readFile)
	var text string
	for _, line := range lines {
		text += fmt.Sprintln(line)
	}
	_, err = readFile.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
