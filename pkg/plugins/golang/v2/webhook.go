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

package v2

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"

	newconfig "sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins/golang/v2/scaffolds"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugins/internal/cmdutil"
)

type createWebhookSubcommand struct {
	config newconfig.Config
	// For help text.
	commandName string

	options *Options
}

var (
	_ plugin.CreateWebhookSubcommand = &createWebhookSubcommand{}
	_ cmdutil.RunOptions             = &createWebhookSubcommand{}
)

func (p *createWebhookSubcommand) UpdateContext(ctx *plugin.Context) {
	ctx.Description = `Scaffold a webhook for an API resource. You can choose to scaffold defaulting,
validating and (or) conversion webhooks.
`
	ctx.Examples = fmt.Sprintf(`  # Create defaulting and validating webhooks for CRD of group ship, version v1beta1
  # and kind Frigate.
  %s create webhook --group ship --version v1beta1 --kind Frigate --defaulting --programmatic-validation

  # Create conversion webhook for CRD of group shio, version v1beta1 and kind Frigate.
  %s create webhook --group ship --version v1beta1 --kind Frigate --conversion
`,
		ctx.CommandName, ctx.CommandName)

	p.commandName = ctx.CommandName
}

func (p *createWebhookSubcommand) BindFlags(fs *pflag.FlagSet) {
	p.options = &Options{}
	fs.StringVar(&p.options.Group, "group", "", "resource Group")
	p.options.Domain = p.config.GetDomain()
	fs.StringVar(&p.options.Version, "version", "", "resource Version")
	fs.StringVar(&p.options.Kind, "kind", "", "resource Kind")
	fs.StringVar(&p.options.Plural, "resource", "", "resource irregular plural form")

	p.options.WebhookVersion = "v1beta1"
	fs.BoolVar(&p.options.DoDefaulting, "defaulting", false,
		"if set, scaffold the defaulting webhook")
	fs.BoolVar(&p.options.DoValidation, "programmatic-validation", false,
		"if set, scaffold the validating webhook")
	fs.BoolVar(&p.options.DoConversion, "conversion", false,
		"if set, scaffold the conversion webhook")
}

func (p *createWebhookSubcommand) InjectConfig(c newconfig.Config) {
	p.config = c
}

func (p *createWebhookSubcommand) Run(fs afero.Fs) error {
	return cmdutil.Run(p, fs)
}

func (p *createWebhookSubcommand) Validate() error {
	if err := p.options.Validate(); err != nil {
		return err
	}

	if !p.options.DoDefaulting && !p.options.DoValidation && !p.options.DoConversion {
		return fmt.Errorf("%s create webhook requires at least one of --defaulting,"+
			" --programmatic-validation and --conversion to be true", p.commandName)
	}

	// check if resource exist to create webhook
	if !p.config.HasResource(p.options.GVK()) {
		return fmt.Errorf("%s create webhook requires an api with the group,"+
			" kind and version provided", p.commandName)
	}

	return nil
}

func (p *createWebhookSubcommand) GetScaffolder() (cmdutil.Scaffolder, error) {
	// Create the resource from the options
	res := p.options.NewResource(p.config)
	return scaffolds.NewWebhookScaffolder(p.config, res), nil
}

func (p *createWebhookSubcommand) PostScaffold() error {
	return nil
}
