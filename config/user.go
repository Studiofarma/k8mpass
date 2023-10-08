package config

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/log"
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
	Persist(context string)
	Load(context string)
}

type UserService struct {
	user string
	data userData
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

func (u *UserService) Persist(context string) {
	err := writeLines(context, u.user, u.data.PinnedNamespaces)
	if err != nil {
		log.Errorf("error persisting", err)
	}
}
func (u *UserService) Load(context string) {
	u.data = loadFromFile(context, u.user)
}

func New(user string) *UserService {
	return &UserService{user: user}
}

func loadFromFile(context string, user string) userData {
	userDir, err := os.UserCacheDir()
	if err != nil {
		return userData{}
	}
	preferencesPath := filepath.Join(userDir, "k8mpass", "users", user, "namespaces")
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

func writeLines(currentContext string, user string, lines []string) error {
	userDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	filePath := filepath.Join(userDir, "k8mpass", "users", user, "namespaces", fmt.Sprintf("%s_ns.txt", currentContext))
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
