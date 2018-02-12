package kubernetes

import (
	"github.com/wdxxs2z/service-broker-helm/config"

	k8sclient "k8s.io/client-go/kubernetes"
	k8sclientcmd "k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"code.cloudfoundry.org/lager"
)

type KubeBuilder struct {
	buildConfig            func(string, string) (*rest.Config, error)
	newKubernetesClientSet func(*rest.Config) (*k8sclient.Clientset, error)
}

func DefaultSourceBuilder() *KubeBuilder {
	return NewKubeBuilder(k8sclientcmd.BuildConfigFromFlags, k8sclient.NewForConfig)
}

func NewKubeBuilder(buildConfig func(string, string) (*rest.Config, error),
	newKubernetesClientSet func(*rest.Config) (*k8sclient.Clientset, error)) *KubeBuilder {
	return &KubeBuilder{
		buildConfig: buildConfig,
		newKubernetesClientSet: newKubernetesClientSet,
	}
}

func (builder *KubeBuilder) CreateKubernetesClientSet(cfg *config.Config, logger lager.Logger) (*k8sclient.Clientset, error) {
	kube_cfg,err := builder.buildConfig("", cfg.KubeConfigPath)
	if err != nil {
		logger.Fatal("building config from flags", err)
	}
	kubeClientSet, err := builder.newKubernetesClientSet(kube_cfg)
	if err != nil {
		logger.Fatal("building kube config set from flags", err)
		return nil,err
	}
	return kubeClientSet, nil
}