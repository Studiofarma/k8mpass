package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os/exec"
	"runtime"
	"slices"
	"strings"
)

func NewCronjobOperation(name string, cronjob string, output string) INamespaceOperation {
	return NamespaceOperation{
		Name:      name,
		Condition: CronjobCondition(cronjob),
		Command:   CronjobCommand(cronjob, output),
	}
}

func NewIngressOperation(name string, ingressName string, output string) INamespaceOperation {
	return NamespaceOperation{
		Name:      name,
		Condition: IngressCondition(ingressName),
		Command:   IngressCommand(ingressName, output),
	}
}

func CronjobCondition(cronjob string) K8mpassCondition {
	return func(cs *kubernetes.Clientset, namespace string) bool {
		_, err := cs.BatchV1().CronJobs(namespace).Get(context.TODO(), cronjob, metav1.GetOptions{})
		if err != nil {
			return false
		}
		return true
	}
}

func CronjobCommand(cronjob string, output string) K8mpassCommand {
	return func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			err := triggerCronjob(clientset, namespace, cronjob)
			if err != nil {
				return NoOutputResultMsg{false, err.Error()}
			}
			return NoOutputResultMsg{true, output}
		}
	}
}

func triggerCronjob(clientset *kubernetes.Clientset, namespace string, cronjobName string) error {
	cronjobs := clientset.BatchV1().CronJobs(namespace)
	cronjob, err := cronjobs.Get(context.TODO(), cronjobName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	newUUID := uuid.New()
	h := md5.New()
	h.Write([]byte(newUUID.String()))
	hash := hex.EncodeToString(h.Sum(nil))
	jobSpec := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("k8mpass-%s-%s", cronjobName, hash[:10]),
			Namespace: namespace,
		},
		Spec: cronjob.Spec.JobTemplate.Spec,
	}
	jobs := clientset.BatchV1().Jobs(namespace)

	_, err = jobs.Create(context.TODO(), jobSpec, metav1.CreateOptions{})

	return err
}

func IngressCondition(ingress string) K8mpassCondition {
	return func(cs *kubernetes.Clientset, namespace string) bool {
		ingresses, err := cs.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false
		}

		return slices.IndexFunc(ingresses.Items, func(i v1.Ingress) bool {
			host := i.Spec.Rules[0].Host
			return strings.HasPrefix(host, ingress)
		}) != -1
	}
}

func IngressCommand(ingress string, output string) K8mpassCommand {
	return func(clientset *kubernetes.Clientset, namespace string) tea.Cmd {
		return func() tea.Msg {
			ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return NoOutputResultMsg{false, err.Error()}
			}

			idx := slices.IndexFunc(ingresses.Items, func(i v1.Ingress) bool {
				host := i.Spec.Rules[0].Host
				return strings.HasPrefix(host, ingress)
			})

			if idx == -1 {
				return NoOutputResultMsg{false, "Ingress not found"}
			}
			url := ingresses.Items[idx].Spec.Rules[0].Host
			openbrowser("https://" + url)

			return NoOutputResultMsg{
				true,
				output,
			}
		}
	}
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
