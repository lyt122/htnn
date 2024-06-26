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
	"k8s.io/apimachinery/pkg/types"

	"mosn.io/htnn/controller/internal/log"
	"mosn.io/htnn/controller/pkg/component"
	pkgRegistry "mosn.io/htnn/controller/pkg/registry"
	mosniov1 "mosn.io/htnn/types/apis/v1"
	registrytype "mosn.io/htnn/types/pkg/registry"
)

var (
	registries = map[types.NamespacedName]pkgRegistry.Registry{}
	store      *serviceEntryStore
)

type RegistryManagerOption struct {
	Output component.Output
}

func InitRegistryManager(opt *RegistryManagerOption) {
	store = newServiceEntryStore(opt.Output)
}

func UpdateRegistry(registry *mosniov1.ServiceRegistry, prevServiceRegistry *mosniov1.ServiceRegistry) error {
	if prevServiceRegistry != nil && prevServiceRegistry.Generation == registry.Generation {
		// no change
		return nil
	}

	key := types.NamespacedName{Namespace: registry.Namespace, Name: registry.Name}
	if reg, ok := registries[key]; !ok {
		reg, err := pkgRegistry.CreateRegistry(registry.Spec.Type, store, registry.ObjectMeta)
		if err != nil {
			return err
		}

		conf, err := registrytype.ParseConfig(reg, registry.Spec.Config.Raw)
		if err != nil {
			return err
		}

		log.Infof("start registry %s", key)

		err = reg.Start(conf)
		if err != nil {
			return err
		}

		// only started registry can be put into registries
		registries[key] = reg

	} else {
		conf, err := registrytype.ParseConfig(reg, registry.Spec.Config.Raw)
		if err != nil {
			return err
		}

		log.Infof("reload registry %s", key)

		err = reg.Reload(conf)
		if err != nil {
			return err
		}
	}

	return nil
}

func DeleteRegistry(key types.NamespacedName) error {
	prev, ok := registries[key]
	if !ok {
		// this may happens when deleting an invalid ServiceRegistry
		return nil
	}

	delete(registries, key)
	log.Infof("stop registry %s", key)
	return prev.Stop()
}
