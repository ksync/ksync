package ksync

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // Not sure why this is needed.
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// TODO: is this singleton a good idea?
	kubeClient  *kubernetes.Clientset
	kubeCfg     *rest.Config
	kubeCfgPath string
	namespace   string
)

func getKubeConfig(context string) (*rest.Config, string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{}

	if context != "" {
		overrides.CurrentContext = context
	}

	clientLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		rules,
		overrides)
	config, err := clientLoader.ClientConfig()
	if err != nil {
		return nil, "", fmt.Errorf(
			"could not get config for context (%q): %s", context, err)
	}

	return config, clientLoader.ConfigAccess().GetDefaultFilename(), nil
}

// InitKubeClient creates a new k8s client for use in talking to the k8s api server.
func InitKubeClient(context string, nspace string) error {
	log.WithFields(log.Fields{
		"context":   context,
		"namespace": namespace,
	}).Debug("initializing kubernetes client")
	config, cfgPath, err := getKubeConfig(context)

	// TODO: better error
	if err != nil {
		return err
	}

	client, err := kubernetes.NewForConfig(config)
	log.WithFields(log.Fields{
		"host": config.Host,
	}).Debug("kubernetes client created")

	// TODO: better error
	if err != nil {
		return err
	}

	kubeClient = client
	kubeCfg = config
	kubeCfgPath = cfgPath
	namespace = nspace

	return nil
}
