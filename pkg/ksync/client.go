package ksync

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// TODO: is this singleton a good idea?
	KubeClient *kubernetes.Clientset
	KubeCfg    *rest.Config
	Namespace  string
)

func getKubeConfig(context string) (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{}

	if context != "" {
		overrides.CurrentContext = context
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		overrides).ClientConfig()

	if err != nil {
		return nil, fmt.Errorf(
			"could not get config for context (%q): %s", context, err)
	}

	return config, nil
}

// InitClient creates a new k8s client for use in talking to the k8s api server.
func InitClient(context string, namespace string) error {
	log.WithFields(log.Fields{
		"context":   context,
		"namespace": namespace,
	}).Debug("initializing kubernetes client")
	config, err := getKubeConfig(context)

	// TODO: better error
	if err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	log.WithFields(log.Fields{
		"host": config.Host,
	}).Debug("kubernetes client created")

	// TODO: better error
	if err != nil {
		return err
	}

	KubeClient = kubeClient
	KubeCfg = config
	Namespace = namespace

	return nil
}
