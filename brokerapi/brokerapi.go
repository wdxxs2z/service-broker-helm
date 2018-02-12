package brokerapi

import (
	"github.com/wdxxs2z/service-broker-helm/config"
	"github.com/wdxxs2z/service-broker-helm/helm"
	"github.com/pivotal-cf/brokerapi"
	"github.com/mitchellh/mapstructure"
	"code.cloudfoundry.org/lager"

	"fmt"
)

type ProvisionParameters map[string]interface{}

type HelmServiceBroker struct {
	config			*config.Config
	helmClient		*helm.Client
	logger     		lager.Logger
}

func New(config *config.Config, helmClient *helm.Client, logger lager.Logger) *HelmServiceBroker {
	return &HelmServiceBroker{
		config:     config,
		helmClient: helmClient,
		logger:     logger.Session("helm-broker"),
	}
}

func (hsb *HelmServiceBroker) Service() []brokerapi.Service {
	return hsb.config.ServiceBroker
}

func (hsb *HelmServiceBroker) Provision(instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error){
	provisionedServiceSpec := brokerapi.ProvisionedServiceSpec{IsAsync: true}

	if !asyncAllowed {
		return provisionedServiceSpec, brokerapi.ErrAsyncRequired
	}

	provisionParameters := ProvisionParameters{}
	if err := mapstructure.Decode(serviceDetails.RawParameters, &provisionParameters); err != nil {
		return provisionedServiceSpec, fmt.Errorf("Error parsing provision parameters: %s", err)
	}

	chart := ""
	if v, ok := provisionParameters["chart"]; ok {
		chart = v
	} else {
		return provisionedServiceSpec, fmt.Errorf("chart parameter must be set.")
	}

	version := ""
	if v, ok := provisionParameters["version"]; ok {
		chart = v
	} else {
		return provisionedServiceSpec, fmt.Errorf("version parameter must be set.")
	}

	values := make(map[string]interface{})
	if v, ok := provisionParameters["values"]; ok {
		values = v
	} else {
		return provisionedServiceSpec, fmt.Errorf("values parameter must be set.")
	}

	_, err = hsb.helmClient.InstallRelease(instanceID, chart, version, values, hsb.logger)
	if err != nil {
		return provisionedServiceSpec, err
	}

	return provisionedServiceSpec, nil
}

func (hsb *HelmServiceBroker) Deprovision(instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {

	return nil
}

func (hsb *HelmServiceBroker) Bind(instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {

	return nil
}

func (hsb *HelmServiceBroker) Unbind(instanceID, bindingID string, details brokerapi.UnbindDetails) error {

	return nil
}

func (hsb *HelmServiceBroker) LastOperation(instanceID, operationData string) (brokerapi.LastOperation, error) {

	return nil
}