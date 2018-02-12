package helm

import (
	"github.com/wdxxs2z/service-broker-helm/config"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/tlsutil"
	"code.cloudfoundry.org/lager"
)

type HelmBuilder struct {
	newHelmClient func(opts ...helm.Option) helm.Client
}

func DefaultNewHelmBuilder() *HelmBuilder {
	return &HelmBuilder{
		newHelmClient: helm.NewClient,
	}
}

func NewHelmBuilder(newHelmClient func(opts ...helm.Option) helm.Client) *HelmBuilder{
	return &HelmBuilder{
		newHelmClient: newHelmClient,
	}
}

func (hb *HelmBuilder) CreateHelmClient(cfg *config.Config, logger lager.Logger) helm.Client{
	options := []helm.Option{helm.Host(cfg.HelmOptionConfig.Host)}
	if !cfg.HelmOptionConfig.SkipTlsVerification {
		tlsOpts := tlsutil.Options{KeyFile: cfg.HelmOptionConfig.TlsKeyFile, CertFile: cfg.HelmOptionConfig.TlsCertFile, InsecureSkipVerify: true}
		tlsCfg, err := tlsutil.ClientConfig(tlsOpts)
		if err != nil {
			logger.Fatal("helm tls config from flags", err)
		}
		options = append(options, helm.WithTLS(tlsCfg))
	}
	return hb.newHelmClient(options)
}