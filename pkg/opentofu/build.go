package opentofu

import (
	"context"
	"text/template"

	"get.porter.sh/porter/pkg/exec/builder"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

const dockerfileLines = `
ENV PORTER_OPENTOFU_MIXIN_USER_AGENT_OPT_OUT="{{ .UserAgentOptOut}}"
ENV AZURE_HTTP_USER_AGENT="{{ .AzureUserAgent }}"
RUN --mount=type=cache,target=/var/cache/apt --mount=type=cache,target=/var/lib/apt \
 apt-get update && apt-get install -y wget && \
 wget --secure-protocol=TLSv1_2 --https-only https://get.opentofu.org/install-opentofu.sh -O install-opentofu.sh --progress=dot:giga && \
 chmod +x install-opentofu.sh && \
 ./install-opentofu.sh --install-method standalone && \
 rm install-opentofu.sh && \
COPY {{.WorkingDir}}/{{.InitFile}} $BUNDLE_DIR/{{.WorkingDir}}/
RUN cd $BUNDLE_DIR/{{.WorkingDir}} && \
 tofu init -backend=false && \
 rm -fr .tofu/providers && \
 tofu providers mirror /usr/local/share/tofu/plugins
`

// BuildInput represents stdin passed to the mixin for the build command.
type BuildInput struct {
	Config *MixinConfig
}

// MixinConfig represents configuration that can be set on the tofu mixin in porter.yaml
// mixins:
//   - opentofu:
//     version: v0.0.0
type MixinConfig struct {
	// ClientVersion is the version of the tofu CLI to install
	ClientVersion string `yaml:"clientVersion,omitempty"`

	// UserAgentOptOut allows a bundle author to opt out from adding porter and the mixin's version to the tofu user agent string.
	UserAgentOptOut bool `yaml:"userAgentOptOut,omitempty"`

	InitFile   string `yaml:"initFile,omitempty"`
	WorkingDir string `yaml:"workingDir,omitempty"`
}

type buildConfig struct {
	MixinConfig

	// AzureUserAgent is the contents of the azure user agent environment variable
	AzureUserAgent string
}

func (m *Mixin) Build(ctx context.Context) error {
	input := BuildInput{
		Config: &m.config, // Apply config directly to the mixin
	}
	err := builder.LoadAction(ctx, m.RuntimeConfig, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &input)
		return &input, err
	})
	if err != nil {
		return err
	}

	tmpl, err := template.New("Dockerfile").Parse(dockerfileLines)
	if err != nil {
		return errors.Wrapf(err, "error parsing opentofu mixin Dockerfile template")
	}

	cfg := buildConfig{MixinConfig: *input.Config}
	if !input.Config.UserAgentOptOut {
		cfg.AzureUserAgent = m.userAgent
	}

	return tmpl.Execute(m.Out, cfg)
}
