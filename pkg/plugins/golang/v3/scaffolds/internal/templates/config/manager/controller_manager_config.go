/*
Copyright 2020 The Kubernetes Authors.

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

package manager

import (
	"path/filepath"

	"sigs.k8s.io/kubebuilder/v2/pkg/model/file"
)

var _ file.Template = &ControllerManagerConfig{}

// ControllerManagerConfig scaffolds the config file in config/manager folder.
type ControllerManagerConfig struct {
	file.TemplateMixin
	file.QualifiedGroupMixin
	file.RepositoryMixin
}

// SetTemplateDefaults implements input.Template
func (f *ControllerManagerConfig) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", "manager", "controller_manager_config.yaml")
	}

	f.TemplateBody = controllerManagerConfigTemplate

	f.IfExistsAction = file.Error

	return nil
}

const controllerManagerConfigTemplate = `apiVersion: controller-runtime.sigs.k8s.io/v1alpha1
kind: ControllerManagerConfig
health:
  healthProbeBindAddress: :8081
metrics:
  bindAddress: 127.0.0.1:8080
webhook:
  port: 9443
leaderElection:
  leaderElect: true
  resourceName: {{ hashFNV .Repo }}.{{ .QualifiedGroup }}
`
