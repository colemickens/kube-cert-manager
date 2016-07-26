/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudprovider

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/glog"
)

// TODO: credit the code I'm yanking from kubernetes main repo (plugins.go)
//       this could is an exact replica of kubernetes/pkg/cloudprovider/plugins.go, plus
//          s/cloudprovider/dnsprovider/g
//          s/cloud provider/dns provider/g
//          s/CloudProvider/DNSProvider/g

// Factory is a function that returns a cloudprovider.Interface.
// The config parameter provides an io.Reader handler to the factory in
// order to load specific configurations. If no configuration is provided
// the parameter is nil.
type Factory func(config io.Reader) (Interface, error)

// All registered dns providers.
var providersMutex sync.Mutex
var providers = make(map[string]Factory)

// RegisterDNSProvider registers a cloudprovider.Factory by name.  This
// is expected to happen during app startup.
func RegisterDNSProvider(name string, cloud Factory) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	if _, found := providers[name]; found {
		glog.Fatalf("Cloud provider %q was registered twice", name)
	}
	glog.V(1).Infof("Registered dns provider %q", name)
	providers[name] = cloud
}

// GetDNSProvider creates an instance of the named dns provider, or nil if
// the name is not known.  The error return is only used if the named provider
// was known but failed to initialize. The config parameter specifies the
// io.Reader handler of the configuration file for the dns provider, or nil
// for no configuation.
func GetDNSProvider(name string, config io.Reader) (Interface, error) {
	providersMutex.Lock()
	defer providersMutex.Unlock()
	f, found := providers[name]
	if !found {
		return nil, nil
	}
	return f(config)
}

// InitDNSProvider creates an instance of the named dns provider.
func InitDNSProvider(name string, configFilePath string) (Interface, error) {
	var cloud Interface
	var err error

	if name == "" {
		glog.Info("No dns provider specified.")
		return nil, nil
	}

	if configFilePath != "" {
		var config *os.File
		config, err = os.Open(configFilePath)
		if err != nil {
			glog.Fatalf("Couldn't open dns provider configuration %s: %#v",
				configFilePath, err)
		}

		defer config.Close()
		cloud, err = GetDNSProvider(name, config)
	} else {
		// Pass explicit nil so plugins can actually check for nil. See
		// "Why is my nil error value not equal to nil?" in golang.org/doc/faq.
		cloud, err = GetDNSProvider(name, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("could not init dns provider %q: %v", name, err)
	}
	if cloud == nil {
		return nil, fmt.Errorf("unknown dns provider %q", name)
	}

	return cloud, nil
}
