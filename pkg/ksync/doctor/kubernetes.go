package doctor

import (
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/spf13/viper"
	authorizationapi "k8s.io/client-go/pkg/apis/authorization/v1"

	"github.com/vapor-ware/ksync/pkg/ksync/cluster"
)

var (
	kubeConfigError      = `Unable to get config for context (%s). Does 'kubectl --context=%s cluster-info' work?`
	kubeConnectError     = `Unable to contact the cluster for context (%s). Does 'kubectl --context=%s cluster-info' work?`
	kubeDaemonSetError   = `Unable to create the ksync daemonset in namespace (%s) for context (%s). You can test with 'kubectl --namespace=%s --context=%s auth can-i create daemonset'.`
	kubePortforwardError = `Unable to setup port forwarding for the ksync pods in namespace (%s) for context (%s). You can test with 'kubectl --namespace=%s --context=%s auth can-i get pods --subresource=portforward'.`
	kubeVersionError     = `Your cluster version (%s) does not fall within the acceptible range: %s. Please upgrade to a compatible version.`
)

// IsClusterVersionSupported verifies that the remote cluster's API version
// falls within the acceptable range.
func IsClusterVersionSupported() error {
	clusterInfo, err := cluster.Client.Discovery().ServerVersion()
	if err != nil {
		return err
	}

	clusterVersion, err := semver.Make(
		strings.Replace(clusterInfo.String(), "v", "", -1))
	if err != nil {
		return err
	}

	versionRange, err := semver.ParseRange(KubernetesRange)
	if err != nil {
		return err
	}

	if !versionRange(clusterVersion) {
		return fmt.Errorf(
			kubeVersionError,
			clusterInfo.String(),
			KubernetesRange)
	}

	return nil
}

// IsClusterConfigValid verifies that the local configuration can be loaded
// and potentially used to connect to the cluster.
func IsClusterConfigValid() error {
	ctx := viper.GetString("context")
	_, _, err := cluster.GetKubeConfig(ctx)

	if err != nil {
		return fmt.Errorf(kubeConfigError, ctx, ctx)
	}

	// This is a bit of a hack. The singleton should be setup elsewhere, but
	// doctor will just run through all the checks. If the cluster config is
	// valid, create the client.
	if err := cluster.InitKubeClient(ctx); err != nil {
		return err
	}

	return nil
}

// CanConnectToCluster verifies, naively, that the cluster can be connected to
// and return a basic result.
func CanConnectToCluster() error {
	ctx := viper.GetString("context")
	client := cluster.Client.DiscoveryClient.RESTClient()

	if client.Get().Timeout(1*time.Second).Do().Error() != nil {
		return fmt.Errorf(kubeConnectError, ctx, ctx)
	}

	return nil
}

// HasClusterPermissions verifies that the current context/user is able to
// do everything required on the remote cluster.
func HasClusterPermissions() error {
	ctx := viper.GetString("context")

	service := cluster.NewService()
	reviews := cluster.Client.AuthorizationV1().SelfSubjectAccessReviews()

	createReview := &authorizationapi.SelfSubjectAccessReview{
		Spec: authorizationapi.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationapi.ResourceAttributes{
				Namespace: service.Namespace,
				Verb:      "create",
				Resource:  "daemonset",
			},
		},
	}

	if resp, err := reviews.Create(createReview); err != nil {
		return err
	} else if !resp.Status.Allowed {
		return fmt.Errorf(
			kubeDaemonSetError, service.Namespace, ctx, service.Namespace, ctx)
	}

	portforwardReview := &authorizationapi.SelfSubjectAccessReview{
		Spec: authorizationapi.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationapi.ResourceAttributes{
				Namespace:   service.Namespace,
				Verb:        "get",
				Resource:    "pods",
				Subresource: "portforward",
			},
		},
	}

	if resp, err := reviews.Create(portforwardReview); err != nil {
		return err
	} else if !resp.Status.Allowed {
		return fmt.Errorf(
			kubePortforwardError, service.Namespace, ctx, service.Namespace, ctx)
	}

	return nil
}
