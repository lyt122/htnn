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

package integration

import (
	"os"
	"testing"
	"time"

	"mosn.io/htnn/api/plugins/tests/integration/controlplane"
)

var (
	controlPlane *controlplane.ControlPlane
)

func TestMain(m *testing.M) {
	controlPlane = controlplane.NewControlPlane()
	go func() {
		controlPlane.Start()
	}()
	time.Sleep(1 * time.Second)

	os.Exit(m.Run())
}
