package kubernetes

import (
	"context"
	"flag"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type ICluster interface {
	Connect() error
}

type Cluster struct {
	cs             *kubernetes.Clientset
	namespaceWatch watch.Interface
	podWatch       watch.Interface
}

func (c *Cluster) Connect() error {
	cs, err := getConnection()
	if err != nil {
		return err
	}
	c.cs = cs
	return nil
}

func getConnection() (*kubernetes.Clientset, error) {
	configPath := getConfigPath()
	// To add a minimum spinner time
	sleep := time.NewTimer(time.Millisecond * 500).C

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}
	cs, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	err = healthCheck(cs)
	if err != nil {
		return nil, err
	}
	<-sleep
	return cs, nil
}

func getConfigPath() string {
	var configPath string
	if f := configPathFromFlag(); f != "" {
		configPath = f
	} else if e := configPathFromEnvVar(); e != "" {
		configPath = e
	} else {
		configPath = defaultKubeConfigFilePath()
	}
	return configPath
}

func healthCheck(cs *kubernetes.Clientset) error {
	res := cs.RESTClient().Get().AbsPath("/healthz").Do(context.TODO())
	if err := res.Error(); err != nil {
		return err
	}
	return nil
}

func defaultKubeConfigFilePath() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(userHomeDir, ".kube", "config")
}

func configPathFromEnvVar() string {
	path, found := os.LookupEnv("KUBECONFIG")
	if !found {
		return ""
	}
	return path
}

func configPathFromFlag() string {
	path := flag.String("kubeconfig", "", "specify kubernetes config file to use")
	flag.Parse()
	if path == nil {
		return *path
	}
	return ""
}
