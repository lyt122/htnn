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

package casbin

import (
	"sync"

	"mosn.io/htnn/api/pkg/filtermanager/api"
	"mosn.io/htnn/plugins/pkg/file"
)

func factory(c interface{}, callbacks api.FilterCallbackHandler) api.Filter {
	return &filter{
		callbacks: callbacks,
		config:    c.(*config),
	}
}

type filter struct {
	api.PassThroughFilter

	callbacks api.FilterCallbackHandler
	config    *config
}

var (
	Changed   = false
	ChangedMu sync.RWMutex
)

func (f *filter) DecodeHeaders(headers api.RequestHeaderMap, endStream bool) api.ResultAction {
	conf := f.config
	role, _ := headers.Get(conf.Token.Name) // role can be ""
	url := headers.Url()

	err := file.WatchFiles(func() {
		setChanged(true)
	}, conf.modelFile, conf.policyFile)

	if err != nil {
		api.LogErrorf("failed to watch files: %v", err)
		return &api.LocalResponse{Code: 500}
	}

	if getChanged() {
		conf.reloadEnforcer()
	}

	conf.lock.RLock()
	ok, err := f.config.enforcer.Enforce(role, url.Path, headers.Method())
	conf.lock.RUnlock()

	if !ok {
		if err != nil {
			api.LogErrorf("failed to enforece %s: %v", role, err)
		}
		api.LogInfof("reject forbidden user %s", role)
		return &api.LocalResponse{
			Code: 403,
		}
	}
	return api.Continue
}

func setChanged(change bool) {
	ChangedMu.Lock()
	Changed = change
	ChangedMu.Unlock()
}

func getChanged() bool {
	ChangedMu.RLock()
	changed := Changed
	ChangedMu.RUnlock()
	return changed
}
