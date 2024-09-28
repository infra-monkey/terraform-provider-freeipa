// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package freeipa

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	ipa "github.com/RomanButsiy/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure freeipaProvider satisfies various provider interfaces.
var _ provider.Provider = &freeipaProvider{}

//var _ provider.ProviderWithFunctions = &freeipaProvider{}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &freeipaProvider{
			version: version,
		}
	}
}

// freeipaProvider defines the provider implementation.
type freeipaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// freeipaProviderModel describes the provider data model.
type freeipaProviderModel struct {
	Host               types.String `tfsdk:"host"`
	Username           types.String `tfsdk:"username"`
	Password           types.String `tfsdk:"password"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure"`
	CaCertificate      types.String `tfsdk:"ca_certificate"`
}

func (p *freeipaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "freeipa"
	resp.Version = p.version
}

func (p *freeipaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The FreeIPA host",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username to use for connection",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password to use for connection",
				Optional:            true,
				Sensitive:           true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether to verify the server's SSL certificate",
				Optional:            true,
			},
			"ca_certificate": schema.StringAttribute{
				MarkdownDescription: "Path to the server's SSL CA certificate",
				Optional:            true,
			},
		},
	}
}

func (p *freeipaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config freeipaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown FreeIPA API",
			"The provider cannot create the FreeIPA API client as there is an unknown configuration value for the FreeIPA host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FREEIPA_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown FreeIPA Username",
			"The provider cannot create the FreeIPA API client as there is an unknown configuration value for the FreeIPA username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FREEIPA_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown FreeIPA Password",
			"The provider cannot create the FreeIPA API client as there is an unknown configuration value for the FreeIPA password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FREEIPA_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("FREEIPA_HOST")
	username := os.Getenv("FREEIPA_USERNAME")
	password := os.Getenv("FREEIPA_PASSWORD")
	insecure := getEnvAsBool("FREEIPA_INSECURE", false)
	ca_certificate := os.Getenv("FREEIPA_CA_CERT")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.InsecureSkipVerify.IsNull() {
		insecure = config.InsecureSkipVerify.ValueBool()
	}

	if !config.CaCertificate.IsNull() {
		ca_certificate = config.CaCertificate.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing FreeIPA Host",
			"The provider cannot create the FreeIPA API client as there is a missing or empty value for the FreeIPA host. "+
				"Set the host value in the configuration or use the FREEIPA_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing FreeIPA Username",
			"The provider cannot create the FreeIPA API client as there is a missing or empty value for the FreeIPA username. "+
				"Set the username value in the configuration or use the FREEIPA_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing FreeIPA Password",
			"The provider cannot create the FreeIPA API client as there is a missing or empty value for the FreeIPA password. "+
				"Set the password value in the configuration or use the FREEIPA_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if insecure {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("insecure"),
			"FreeIPA InsecureSkipVerify set to TRUE",
			"The provider will skip TLS verification for the FreeIPA API client and therefore cannot guaranty the security of the conneciton. ",
		)
	}
	if !insecure && ca_certificate == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("ca_certificate"),
			"Missing FreeIPA CA Certificate Path",
			"The provider cannot create the FreeIPA API client as there is a missing or empty value for the FreeIPA CA Certificate Path. "+
				"Set the CA Certificate path value in the configuration or use the FREEIPA_CA_CERT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new FreeIPA client using the configuration values
	client, err := p.NewFreeIPAClient(ctx, &config)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create FreeIPA API Client",
			"An unexpected error occurred when creating the FreeIPA API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"FreeIPA Client Error: "+err.Error(),
		)
		return
	}

	// Make the FreeIPA client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Client creates a FreeIPA client scoped to the global API
func (c *freeipaProvider) NewFreeIPAClient(ctx context.Context, conf *freeipaProviderModel) (*ipa.Client, error) {
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] freeipa host : %s", conf.Host.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] freeipa username : %s", conf.Username.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] freeipa password : %s", conf.Password.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] freeipa insecure : %s", conf.InsecureSkipVerify.String()))
	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] freeipa cacert path : %s", conf.CaCertificate.ValueString()))
	caCertPool := x509.NewCertPool()

	if !conf.CaCertificate.IsNull() {
		caCert, err := os.ReadFile(conf.CaCertificate.String())
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify.ValueBool(),
			RootCAs:            caCertPool,
		},
	}

	client, err := ipa.Connect(conf.Host.ValueString(), tspt, conf.Username.ValueString(), conf.Password.ValueString())
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] FreeIPA Client configured for host: %s", conf.Host)

	return client, nil
}

func (p *freeipaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUserGroupResource,
		NewUserResource,
		NewUserGroupMembershipResource,
		NewDNSZoneResource,
	}
}

func (p *freeipaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserGroupDataSource,
		NewUserDataSource,
		NewDnsZoneDataSource,
	}
}

// func (p *freeipaProvider) Functions(ctx context.Context) []func() function.Function {
// 	return []func() function.Function{
// 		NewExampleFunction,
// 	}
// }
