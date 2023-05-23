/*
Copyright 2021 The Kubernetes Authors.

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

package envfuncs

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"k8s.io/klog/v2"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

// BuildDockerImage returns an env.Func that is used for
// building a docker image.
func BuildDockerImage(path string, dockerfile string, tags []string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}

		tar, err := archive.TarWithOptions(path, &archive.TarOptions{})
		if err != nil {
			return nil, err
		}

		opts := types.ImageBuildOptions{
			Dockerfile: dockerfile,
			Tags:       tags,
			Remove:     true,
		}
		res, err := cli.ImageBuild(ctx, tar, opts)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		var lastLine string

		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			lastLine = scanner.Text()
			klog.V(3).Infof("Build: %s", scanner.Text())
		}

		errLine := &ErrorLine{}
		json.Unmarshal([]byte(lastLine), errLine)
		if errLine.Error != "" {
			return nil, errors.New(errLine.Error)
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		return ctx, nil
	}
}
