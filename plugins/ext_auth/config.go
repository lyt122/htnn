package ext_auth

import (
	"net/http"
	"time"

	"mosn.io/moe/pkg/filtermanager/api"
	"mosn.io/moe/pkg/plugins"
)

const (
	// We name this plugin as ext_auth to distinguish it from the C++ implementation ext_authz.
	// We may add new feature to this plugin which will make it different from its C++ sibling.
	Name = "ext_auth"
)

func init() {
	plugins.RegisterHttpPlugin(Name, &plugin{})
}

type plugin struct {
	plugins.PluginMethodDefaultImpl
}

func (p *plugin) ConfigFactory() api.FilterConfigFactory {
	return configFactory
}

func (p *plugin) Config() plugins.PluginConfig {
	return &config{}
}

type config struct {
	Config

	client *http.Client
}

func (conf *config) Init(cb api.ConfigCallbackHandler) error {
	du := 200 * time.Millisecond
	timeout := conf.GetHttpService().Timeout
	if timeout != nil {
		du = timeout.AsDuration()
	}

	conf.client = &http.Client{Timeout: du}
	return nil
}