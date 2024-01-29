package opentofu

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"

	"github.com/crosbygw/opentofu/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMixin_Build(t *testing.T) {
	testcases := []struct {
		name              string
		inputFile         string
		expectedVersion   string
		expectedUserAgent string
	}{
		{name: "build with custom config", inputFile: "testdata/build-input-with-config.yaml", expectedVersion: "wget --secure-protocol=TLSv1_2 --https-only https://get.opentofu.org/install-opentofu.sh -O install-opentofu.sh", expectedUserAgent: "ENV PORTER_OPENTOFU_MIXIN_USER_AGENT_OPT_OUT=\"true\"\nENV AZURE_HTTP_USER_AGENT=\"\""},
		{name: "build with the default OpenTofu config", expectedVersion: "wget --secure-protocol=TLSv1_2 --https-only https://get.opentofu.org/install-opentofu.sh -O install-opentofu.sh", expectedUserAgent: "ENV PORTER_OPENTOFU_MIXIN_USER_AGENT_OPT_OUT=\"false\"\nENV AZURE_HTTP_USER_AGENT=\"github.com/crosbygw/opentofu/pkg/opentofu"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Set a fake version of the mixin and porter for our user agent
			pkg.Version = "v1.2.3"

			var data []byte
			var err error
			if tc.inputFile != "" {
				data, err = ioutil.ReadFile(tc.inputFile)
				require.NoError(t, err)
			}

			m := NewTestMixin(t)
			m.In = bytes.NewReader(data)

			err = m.Build(context.Background())
			require.NoError(t, err, "build failed")

			gotOutput := m.TestContext.GetOutput()
			assert.Contains(t, gotOutput, tc.expectedVersion)
			assert.Contains(t, gotOutput, tc.expectedUserAgent)
			assert.NotContains(t, "{{.", gotOutput, "Not all of the template values were consumed")
		})
	}
}
