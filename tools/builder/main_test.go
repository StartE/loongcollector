// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test_loadPluginConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://test.com/external_plugins_config", func(req *http.Request) (*http.Response, error) {
		data, err := os.ReadFile("testdata/external_plugins.yml")
		assert.Nil(t, err)
		return httpmock.NewBytesResponse(200, data), nil
	})

	cases := []struct {
		configFile string
		wantErr    bool
		wantConf   pluginConfig
	}{
		{
			configFile: "testdata/external_plugins.yml",
			wantErr:    false,
			wantConf: pluginConfig{
				Plugins: pluginCategory{
					Common: []*pluginModule{
						{GoMod: "github.com/mock/ilogtail_plugin1 v1.0.0", Import: "github.com/mock/ilogtail_plugin1"},
						{GoMod: "github.com/mock/ilogtail_plugin2 v1.0.0", Import: "github.com/mock/ilogtail_plugin"},
					},
				},
			},
		},
		{
			configFile: "https://test.com/external_plugins_config",
			wantErr:    false,
			wantConf: pluginConfig{
				Plugins: pluginCategory{
					Common: []*pluginModule{
						{GoMod: "github.com/mock/ilogtail_plugin1 v1.0.0", Import: "github.com/mock/ilogtail_plugin1"},
						{GoMod: "github.com/mock/ilogtail_plugin2 v1.0.0", Import: "github.com/mock/ilogtail_plugin"},
					},
				},
			},
		},
	}

	for _, c := range cases {
		conf, err := loadPluginConfig(c.configFile)
		if c.wantErr {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
		assert.Equal(t, c.wantConf, *conf)
	}
}

func Test_generatePluginSourceCode(t *testing.T) {
	ctx := &buildContext{
		ProjectRoot: t.TempDir(),
		ModFile:     "go.mod",
		GoModContent: `module github.com/alibaba/ilogtail

go 1.23

require (
	github.com/alibaba/ilogtail/base v0.0.0-20230128101543-c0c844084b0e
)`,
		Config: pluginConfig{
			Plugins: pluginCategory{
				Common: []*pluginModule{
					{
						GoMod:  "github.com/mock/common_plugins1 v1.0.0",
						Import: "github.com/mock/common_plugins1",
						Path:   "../",
					},
					{
						GoMod:  "github.com/mock/common_plugins2 v1.0.0",
						Import: "github.com/mock/common_plugins2",
					},
				},
				Windows: []*pluginModule{
					{
						GoMod:  "github.com/mock/windows_plugins1 v1.0.0",
						Import: "github.com/mock/windows_plugins1",
					},
					{
						GoMod:  "github.com/mock/windows_plugins2 v1.0.0",
						Import: "github.com/mock/windows_plugins2",
					},
				},
				Linux: []*pluginModule{
					{
						GoMod:  "github.com/mock/linux_plugins1 v1.0.0",
						Import: "github.com/mock/linux_plugins1",
					},
					{
						GoMod:  "github.com/mock/linux_plugins2 v1.0.0",
						Import: "github.com/mock/linux_plugins2",
					},
				},
				Debug: []*pluginModule{
					{
						GoMod:  "github.com/mock/debug_plugins1 v1.0.0",
						Import: "github.com/mock/debug_plugins1",
					},
					{
						GoMod:  "github.com/mock/debug_plugins2 v1.0.0",
						Import: "github.com/mock/debug_plugins2",
					},
				},
			},
		},
	}

	wants := map[string]string{
		"plugins/all/all.go": `// Code generated by "github.com/alibaba/ilogtail/tools/builder". DO NOT EDIT.
// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package all

import _ "github.com/mock/common_plugins1"
import _ "github.com/mock/common_plugins2"`,
		"plugins/all/all_windows.go": `// Code generated by "github.com/alibaba/ilogtail/tools/builder". DO NOT EDIT.
// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows
// +build windows

package all

import _ "github.com/mock/windows_plugins1"
import _ "github.com/mock/windows_plugins2"`,
		"plugins/all/all_linux.go": `// Code generated by "github.com/alibaba/ilogtail/tools/builder". DO NOT EDIT.
// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build linux
// +build linux

package all

import _ "github.com/mock/linux_plugins1"
import _ "github.com/mock/linux_plugins2"`,
		"plugins/all/all_debug.go": `// Code generated by "github.com/alibaba/ilogtail/tools/builder". DO NOT EDIT.
// Copyright 2021 iLogtail Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build debug
// +build debug

package all

import _ "github.com/mock/debug_plugins1"
import _ "github.com/mock/debug_plugins2"`,
		"go.mod": `module github.com/alibaba/ilogtail

go 1.23

require (
	github.com/alibaba/ilogtail/base v0.0.0-20230128101543-c0c844084b0e
)


require github.com/mock/common_plugins1 v1.0.0
require github.com/mock/common_plugins2 v1.0.0
require github.com/mock/windows_plugins1 v1.0.0
require github.com/mock/windows_plugins2 v1.0.0
require github.com/mock/linux_plugins1 v1.0.0
require github.com/mock/linux_plugins2 v1.0.0
require github.com/mock/debug_plugins1 v1.0.0
require github.com/mock/debug_plugins2 v1.0.0
replace github.com/mock/common_plugins1 v1.0.0 => ../






`,
	}

	err := generatePluginSourceCode(ctx)
	assert.Nil(t, err)

	for path, expect := range wants {
		actual, err := os.ReadFile(filepath.Join(ctx.ProjectRoot, path))
		assert.Nil(t, err)
		assert.Equal(t, expect, string(actual))
	}
}
