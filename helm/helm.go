package helm

import (
	"k8s.io/helm/pkg/helm"
	"code.cloudfoundry.org/lager"
	yaml "gopkg.in/yaml.v2"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"fmt"
	"strings"
	k8s "k8s.io/client-go/kubernetes"
	"github.com/wdxxs2z/service-broker-helm/config"
)

type client struct {
	k8sclientset   	k8s.Interface
	helmclientset 	helm.Interface
}


func (c *client) InstallRelease(instanceID string, chart string, version string, values map[string]interface{}, config	 *config.Config, logger lager.Logger) (*rls.InstallReleaseResponse, error) {
	options := []helm.InstallOption{helm.ReleaseName(instanceID),helm.VersionOption(version)}
	if len(values) > 0 {
		valuesContent, err := yaml.Marshal(values)
		if err != nil {
			logger.Fatal("Error helm values content", err)
			return nil, err
		}
		options = append(options, helm.ValueOverrides([]byte(valuesContent)))
	}
	installResp, err := c.helmclientset.InstallRelease(chart, config.HelmOptionConfig.DefaultNamespace, options)
	if err != nil {
		logger.Fatal("Error install release", err)
		return nil, err
	}
	return installResp, nil
}

func (c *client) UpgradeRelease(instanceID string, chart string, repository string, version string, config *config.Config, logger lager.Logger) error {

	return nil
}

func (c *client) DeleteRelease(instanceID string, config *config.Config, logger lager.Logger) error {

	return nil
}

func (c *client) ReleaseStatus(instanceID string, config *config.Config, logger lager.Logger) (string, string, error) {

	return nil
}

func (c *client) releaseName(instanceID string, config *config.Config) string {
	return fmt.Sprintf("%s-%s", config.HelmOptionConfig.ReleaseNameHeader, strings.Replace(instanceID, "-", "", -1))
}