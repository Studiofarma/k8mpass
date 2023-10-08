package kubernetes

import (
	"context"
	"github.com/studiofarma/k8mpass/config"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type ICluster interface {
	Connect() error
	GetContext() string
	GetUser() string
}

type Cluster struct {
	context        string
	user           string
	cs             *kubernetes.Clientset
	namespaceWatch watch.Interface
	podWatch       watch.Interface
}

func New(user string) Cluster {
	return Cluster{user: user}
}

func (c *Cluster) GetContext() string {
	return c.context
}
func (c *Cluster) GetUser() string {
	return c.user
}

func (c *Cluster) Connect() error {
	configPath := getConfigPath()
	cs, err := getConnection(configPath)
	if err != nil {
		return err
	}
	kubeCtx, err := getContext(configPath)
	if err != nil {
		return err
	}
	c.cs = cs
	c.context = kubeCtx
	return nil
}

func getConnection(configPath string) (*kubernetes.Clientset, error) {
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

func getContext(configPath string) (string, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: configPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()
	if err != nil {
		return "", nil
	}
	return config.CurrentContext, nil
}

func getConfigPath() string {
	var configPath string
	if f := config.Config; f != "" {
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
