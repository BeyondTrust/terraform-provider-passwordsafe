package provider_framework

import (
	"net/http"
	"regexp"
	"terraform-provider-passwordsafe/providers/constants"
	"terraform-provider-passwordsafe/providers/entities"
	"terraform-provider-passwordsafe/providers/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const totpCodeDataSourceAddress = "data.passwordsafe_secret_totp_code.test"

// totpCodeDataSourceTestConfig returns the test configuration for the TOTP code data source.
func totpCodeDataSourceTestConfig(serverUrl, apiVersion, secretId string) string {
	return utils.TestResourceConfig(entities.PasswordSafeTestConfig{
		ClientID:     constants.FakeClientId,
		ClientSecret: constants.FakeClientSecret,
		URL:          serverUrl,
		APIVersion:   apiVersion,
		Resource: `
		data "passwordsafe_secret_totp_code" "test" {
		secret_id = "` + secretId + `"
		}`,
	})
}

// runTotpCodeDataSourceTest runs a single-step test of the TOTP code data source.
func runTotpCodeDataSourceTest(t *testing.T, step resource.TestStep) {
	t.Helper()
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		PreCheck: func() {},
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"passwordsafe": providerserver.NewProtocol6WithError(NewProvider()),
		},
		Steps: []resource.TestStep{step},
	})
}

func TestSecretTotpCodeDataSource(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusOK, `{"Code": "123456", "ValidFromUtc": "2026-10-19T10:00:00Z", "ValidToUtc": "2026-10-19T10:00:30Z"}`)

	runTotpCodeDataSourceTest(t, resource.TestStep{
		Config: totpCodeDataSourceTestConfig(server.URL, "3.3", totpCodeSecretId),
		ConfigStateChecks: []statecheck.StateCheck{
			statecheck.ExpectKnownValue(totpCodeDataSourceAddress, tfjsonpath.New("secret_id"), knownvalue.StringExact(totpCodeSecretId)),
			statecheck.ExpectKnownValue(totpCodeDataSourceAddress, tfjsonpath.New("code"), knownvalue.StringExact("123456")),
			statecheck.ExpectKnownValue(totpCodeDataSourceAddress, tfjsonpath.New("valid_from_utc"), knownvalue.StringExact("2026-10-19T10:00:00Z")),
			statecheck.ExpectKnownValue(totpCodeDataSourceAddress, tfjsonpath.New("valid_to_utc"), knownvalue.StringExact("2026-10-19T10:00:30Z")),
		},
	})
}

func TestSecretTotpCodeDataSourceNotAvailable(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusNotFound, `"TOTP is not enabled"`)

	runTotpCodeDataSourceTest(t, resource.TestStep{
		Config:      totpCodeDataSourceTestConfig(server.URL, "3.3", totpCodeSecretId),
		ExpectError: regexp.MustCompile("secret not found or TOTP is not enabled for this secret"),
	})
}

func TestSecretTotpCodeDataSourceNotCredentialSecret(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusBadRequest, `"Secret is not a credential"`)

	runTotpCodeDataSourceTest(t, resource.TestStep{
		Config:      totpCodeDataSourceTestConfig(server.URL, "3.3", totpCodeSecretId),
		ExpectError: regexp.MustCompile("secret is not a credential secret"),
	})
}

func TestSecretTotpCodeDataSourceInvalidSecretId(t *testing.T) {
	server := newTotpCodeMockServer(t, http.StatusOK, `{"Code": "123456"}`)

	runTotpCodeDataSourceTest(t, resource.TestStep{
		Config:      totpCodeDataSourceTestConfig(server.URL, "3.3", "not-a-uuid"),
		ExpectError: regexp.MustCompile("must be a valid UUID"),
	})
}
