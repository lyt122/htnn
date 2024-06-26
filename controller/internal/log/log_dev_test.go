// Copyright The HTNN Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestDevLogConvertAnyToStr(t *testing.T) {
	InitLogger("json")
	patches := gomonkey.ApplyMethodFunc(devLogger, "Error", func(err error, msg string, keysAndValues ...any) {
		assert.Equal(t, "1", msg)
	})
	patches.ApplyMethodFunc(devLogger, "Info", func(msg string, keysAndValues ...any) {
		assert.Equal(t, "1", msg)
	})
	defer patches.Reset()

	Error(1)
	Info(1)
	Errorf("%d", 1)
	Infof("%d", 1)
}
