package config

import (
	"fmt"
	"io/ioutil"

	"github.com/pivotal-cf/brokerapi"
	yaml "gopkg.in/yaml.v2"
	"code.cloudfoundry.org/multierror"
)

type ConfigSchema struct {
	ServiceBrokerConfig		ServiceSchema		`yaml:"service_broker_config"`
	HelmConfig			HelmSchema		`yaml:"helm_config"`
	KubeConfigPath			string			`yaml:"kube_config_path"`
}

type ServiceSchema struct {
	ServiceName                 string    			`yaml:"service_name"`
	ServiceID                   string    			`yaml:"service_id"`
	Description                 string    			`yaml:"description"`
	LongDescription             string    			`yaml:"long_description"`
	ProviderDisplayName         string    			`yaml:"provider_display_name"`
	DocumentationURL            string    			`yaml:"documentation_url"`
	SupportURL                  string    			`yaml:"support_url"`
	DisplayName                 string    			`yaml:"display_name"`
	IconImage                   string    			`yaml:"icon_image"`
	Plans                       []brokerapi.ServicePlan 	`yaml:"service_plans"`
}

type HelmSchema struct {
	HelmHost			string		`yaml:"helm_host"`
	HelmNamespace			string		`yaml:"helm_namespace"`
	DefaultNamespace		string 		`yaml:"helm_default_namespace"`
	ReleaseNameHeader		string		`yaml:"helm_release_name_header"`
	SkipTlsVerification		bool		`yaml:"skip_tls_verify"`
	TlsCaCertFile			string		`yaml:"tls_cacert_file"`
	TlsCertFile			string		`yaml:"tls_cert_file"`
	TlsKeyFile			string		`yaml:"tls_key_file"`
}

type HelmConfig struct {
	Host			string
	Namespace		string
	DefaultNamespace	string
	ReleaseNameHeader	string
	SkipTlsVerification	bool
	TlsCaCertFile		string
	TlsCertFile		string
	TlsKeyFile		string
}

type Config struct {
	ServiceBroker	 	brokerapi.Service
	HelmOptionConfig	HelmConfig
	KubeConfigPath		string
}

func NewConfigSchemaFromFile(path string) (*ConfigSchema, error) {
	var schema ConfigSchema
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(raw, &schema)
	if err != nil {
		return nil, err
	}
	return &schema, err
}

func missingOptionError(option, desc string) error {
	return fmt.Errorf("config option: %s, error: %s", option, desc)
}

func (cs *ConfigSchema) ToConfig() (*Config, error) {
	errs := multierror.NewMultiError("config")

	if len(cs.KubeConfigPath) == 0 {
		errs.Add(missingOptionError("kube_config_path", "must not be blank"))
	}

	if errs.Length() > 0 {
		return nil, errs
	}

	return &Config{
		ServiceBroker: 		cs.ServiceBrokerConfig,
		HelmOptionConfig: 	cs.HelmConfig,
		KubeConfigPath:		cs.KubeConfigPath,
	}, nil
}

func (hs *HelmSchema) ToConfig() (*HelmConfig, error) {
	errs := multierror.NewMultiError("config")

	if len(hs.HelmHost) == 0 {
		errs.Add(missingOptionError("helm.host", "must not be blank"))
	}

	if len(hs.DefaultNamespace) == 0 {
		errs.Add(missingOptionError("helm.default_namespace", "must not be blank"))
	}

	if len(hs.HelmNamespace) == 0 {
		errs.Add(missingOptionError("helm.namespace", "must not be blank"))
	}

	if len(hs.SkipTlsVerification) == 0 {
		errs.Add(missingOptionError("helm.skip_tls_verify", "must not be blank"))
	}

	if len(hs.ReleaseNameHeader) == 0 {
		errs.Add(missingOptionError("helm.release_name_header", "must not be blank"))
	}

	if errs.Length() > 0 {
		return &HelmConfig{},errs
	}

	return &HelmConfig{
		Host: 		 	hs.HelmHost,
		Namespace:		hs.HelmNamespace,
		DefaultNamespace: 	hs.DefaultNamespace,
		ReleaseNameHeader: 	hs.ReleaseNameHeader,
		SkipTlsVerification: 	hs.SkipTlsVerification,
		TlsCaCertFile:		hs.TlsCaCertFile,
		TlsCertFile:		hs.TlsCertFile,
		TlsKeyFile:		hs.TlsKeyFile,
	}, nil
}

func (ss *ServiceSchema) ToConfig() (*brokerapi.Service, error) {
	errs := multierror.NewMultiError("config")

	if ss.ServiceID == "" {
		errs.Add(missingOptionError("service.id", "must not be blank"))
	}

	if ss.Description == "" {
		errs.Add(missingOptionError("service.description", "must not be blank"))
	}

	if ss.DisplayName == "" {
		errs.Add(missingOptionError("service.display_name", "must not be blank"))
	}

	if len(ss.Plans) == 0 {
		errs.Add(missingOptionError("service.plans", "must not be blank"))
	}

	if errs.Length() > 0 {
		return &brokerapi.Service{}, errs
	}

	return &brokerapi.Service{
		ID:			ss.ServiceID,
		Name:           	ss.ServiceName,
		Description:		ss.Description,
		Bindable:		true,
		Tags:			[]string{"helm","kubernetes"},
		PlanUpdatable:		true,
		Plans:			ss.Plans,
		Metadata:		brokerapi.ServiceMetadata{DisplayName: ss.DisplayName, ImageUrl: ss.IconImage, LongDescription: ss.LongDescription, ProviderDisplayName: ss.ProviderDisplayName, DocumentationUrl: ss.DocumentationURL, SupportUrl: ss.SupportURL},
		DashboardClient:	nil,
	},nil
}