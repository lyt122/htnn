// Copyright The HTNN Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"google.golang.org/protobuf/reflect/protoreflect"

	"mosn.io/htnn/types/pkg/proto"
)

// We don't put this package in api module, because this will force every third-party
// plugin to depend on "istio/api", even the architecture doesn't use any istio code.
// Anyway, the use case to implement a new registry outside HTNN controller is rare.

// RegistryConfig represents the configuration used by the registry
type RegistryConfig interface {
	// The configuration is defined as a protobuf message
	ProtoReflect() protoreflect.Message
	// This method is generated by protoc-gen-validate. We can override it to provide custom validation
	Validate() error
}

type Registry interface {
	Config() RegistryConfig
}

var (
	registryTypes = make(map[string]Registry)
)

// AddRegistryType adds a new registry type
func AddRegistryType(name string, r Registry) {
	registryTypes[name] = r
}

// GetRegistryType gets the added registry type
func GetRegistryType(name string) Registry {
	return registryTypes[name]
}

// ParseConfig parses the given data and returns the configuration according to the registry
func ParseConfig(reg Registry, data []byte) (RegistryConfig, error) {
	conf := reg.Config()

	err := proto.UnmarshalJSON(data, conf)
	if err != nil {
		return nil, err
	}

	err = conf.Validate()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
