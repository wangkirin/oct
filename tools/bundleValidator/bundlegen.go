// Copyright 2015 The oct Authors
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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/opencontainers/specs"
)

var flags = []cli.Flag{
	cli.StringFlag{Name: "outpath", Usage: "path to the output bundles"},
	cli.StringFlag{Name: "rootfs", Usage: "path to the rootfs"},
	cli.BoolFlag{Name: "read-only", Usage: "make the container's rootfs read-only"},
	cli.IntFlag{Name: "gid", Usage: "gid for the process"},
	cli.StringSliceFlag{Name: "groups", Usage: "supplementary groups for the process"},
}

func generateBundle(c *cli.Context) {
	//init
	configinit := SetInitConfig()
	runtimeinit := SetInitRumtime()
	path := "./"
	if c.String("outpath") != "" {
		path = c.String("outpath")
	}
	//settings
	if c.String("rootfs") != "" {
		configinit, runtimeinit = setRootfs(c.String("rootfs"), &configinit, &runtimeinit)
	}

	//generate json
	err := LinuxSpecToConfig(path+"config.json", &configinit)
	if err != nil {
		fmt.Errorf("generate config.json failed")
	}
	err = LinuxRuntimeToConfig(path+"runtime.json", &runtimeinit)
	if err != nil {
		fmt.Errorf("generate runtime.json failed")
	}
}

func setRootfs(path string, config *specs.LinuxSpec, runtime *specs.LinuxRuntimeSpec) (specs.LinuxSpec, specs.LinuxRuntimeSpec) {
	config.Spec.Root.Path = path
	return *config, *runtime
}

func SetInitConfig() specs.LinuxSpec {
	var linuxSpec specs.LinuxSpec = specs.LinuxSpec{
		Spec: specs.Spec{
			Version: "0.2.0",
			Platform: specs.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Root: specs.Root{
				Path:     "rootfs",
				Readonly: true,
			},
			Process: specs.Process{
				Terminal: false,
				User: specs.User{
					UID:            0,
					GID:            0,
					AdditionalGids: nil,
				},
				Args: []string{""},
			},
			Mounts: []specs.MountPoint{
				{
					Name: "proc",
					Path: "/proc",
				},
			},
		},
	}

	return linuxSpec
}

func SetInitRumtime() specs.LinuxRuntimeSpec {
	var linuxRuntimeSpec specs.LinuxRuntimeSpec = specs.LinuxRuntimeSpec{
		RuntimeSpec: specs.RuntimeSpec{
			Mounts: map[string]specs.Mount{
				"proc": specs.Mount{
					Type:    "proc",
					Source:  "proc",
					Options: []string{""},
				},
			},
		},
		Linux: specs.LinuxRuntime{
			Resources: &specs.Resources{
				Memory: specs.Memory{
					Swappiness: -1,
				},
			},
			Namespaces: []specs.Namespace{
				{
					Type: "mount",
					Path: "",
				},
			},
		},
	}
	return linuxRuntimeSpec
}

//

func LinuxSpecToConfig(filePath string, linuxspec *specs.LinuxSpec) error {
	stream, err := json.Marshal(linuxspec)
	if err != nil {
		return err
	}
	objToJson(stream, filePath)
	return err
}

func objToJson(stream []byte, filePath string) {
	fd, err := os.OpenFile(filePath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		fmt.Errorf(" open file err, %v", err)
	}
	defer fd.Close()
	_, err = fd.Write(stream)
	if err != nil {
		fmt.Errorf(" write file err, %v", err)
	}
}

//write specs.LinuxRuntimeSpec to runtime.json
func LinuxRuntimeToConfig(filePath string, linuxRuntime *specs.LinuxRuntimeSpec) error {
	stream, err := json.Marshal(linuxRuntime)
	if err != nil {
		return err
	}
	objToJson(stream, filePath)
	return err
}
