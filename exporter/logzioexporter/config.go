// Copyright The OpenTelemetry Authors
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

package logzioexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/logzioexporter"

import (
	"errors"

	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config contains Logz.io specific configuration such as Account TracesToken, Region, etc.
type Config struct {
	config.ExporterSettings       `mapstructure:",squash"`          // squash ensures fields are correctly decoded in embedded struct
	confighttp.HTTPClientSettings `mapstructure:",squash"`          // confighttp client settings https://pkg.go.dev/go.opentelemetry.io/collector/config/confighttp#HTTPClientSettings
	exporterhelper.QueueSettings  `mapstructure:"sending_queue"`    // exporter helper queue settings https://pkg.go.dev/go.opentelemetry.io/collector/exporter/exporterhelper#QueueSettings
	exporterhelper.RetrySettings  `mapstructure:"retry_on_failure"` // exporter helper retry settings https://pkg.go.dev/go.opentelemetry.io/collector/exporter/exporterhelper#RetrySettings
	Token                         string                            `mapstructure:"account_token"`    // Your Logz.io Account Token, can be found at https://app.logz.io/#/dashboard/settings/general
	Region                        string                            `mapstructure:"region"`           // Your Logz.io 2-letter region code, can be found at https://docs.logz.io/user-guide/accounts/account-region.html#available-regions
	CustomEndpoint                string                            `mapstructure:"custom_endpoint"`  // **Deprecation** Custom endpoint to ship traces to. Use only for dev and tests.
	DrainInterval                 int                               `mapstructure:"drain_interval"`   // **Deprecation** Queue drain interval in seconds. Defaults to `3`.
	QueueCapacity                 int64                             `mapstructure:"queue_capacity"`   // **Deprecation** Queue capacity in bytes. Defaults to `20 * 1024 * 1024` ~ 20mb.
	QueueMaxLength                int                               `mapstructure:"queue_max_length"` // **Deprecation** Max number of items allowed in the queue. Defaults to `500000`.
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return errors.New("`account_token` not specified")
	}
	return nil
}

// CheckAndWarnDeprecatedOptions Is checking for soon deprecated configuration options (queue_max_length, queue_capacity, drain_interval, custom_endpoint) log a warning message and map to the relevant updated option
func (c *Config) checkAndWarnDeprecatedOptions(logger hclog.Logger) {
	if c.QueueCapacity != 0 {
		logger.Warn("You are using  the deprecated`queue_capacity` option that will be removed in the next release; use exporter helper configuration instead: https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md")
	}
	// Warn and map queue_max_length -> QueueSettings.QueueSize
	if c.QueueMaxLength != 0 {
		logger.Warn("You are using the deprecated `queue_max_length` option that will be removed in the next release; use exporter helper `queue_size` configuration instead: https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md")
		c.QueueSettings.Enabled = true
		logger.Warn("Mapping `queue_max_length` -> `QueueSettings.QueueSize`")
		c.QueueSettings.QueueSize = c.QueueMaxLength
	}
	if c.DrainInterval != 0 {
		logger.Warn("You are using the deprecated `drain_interval` option that will be removed in the next release; use batch processor `timeout` configuration instead: https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor#batch-processor")
	}
	// Warn and map CustomEndpoint -> Endpoint
	if c.CustomEndpoint != "" {
		logger.Warn("You are using the deprecated `custom_endpoint` option that will be removed in the next release; please use `endpoint` configuration instead: https://github.com/open-telemetry/opentelemetry-collector/tree/main/config/confighttp")
		logger.Warn("Mapping `custom_endpoint` -> `Endpoint`")
		c.HTTPClientSettings.Endpoint = c.CustomEndpoint
	}
}
