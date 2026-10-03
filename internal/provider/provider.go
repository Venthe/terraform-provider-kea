package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/venthe/terraform-provider-kea/internal/clients"
	"github.com/venthe/terraform-provider-kea/internal/dhcp4"
	"github.com/venthe/terraform-provider-kea/kea"
	"github.com/venthe/terraform-provider-kea/kea/keadhcp4"
)

type KeaProvider struct {
	version string
}

type KeaProviderModel struct {
	DHCP4 types.Object `tfsdk:"dhcp4"`
}

type KeaProviderClientModel struct {
	Address           types.String `tfsdk:"address"`
	HTTPUsername      types.String `tfsdk:"http_username"`
	HTTPPassword      types.String `tfsdk:"http_password"`
	TLSCAFile         types.String `tfsdk:"tls_ca_file"`
	TLSClientCertFile types.String `tfsdk:"tls_client_cert_file"`
	TLSClientKeyFile  types.String `tfsdk:"tls_client_key_file"`
	TLSServerName     types.String `tfsdk:"tls_server_name"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KeaProvider{version}
	}
}

func (p *KeaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kea"
	resp.Version = p.version
}

func (p *KeaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Kea provider enables management of Kea DHCP server resources via the control channel API. " +
			"It connects to Kea servers using UNIX domain sockets or HTTP/HTTPS endpoints.",
		Attributes: map[string]schema.Attribute{
			"dhcp4": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Connection settings for the Kea DHCPv4 server.",
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						Optional: true,
						Description: "Control channel address: `unix:///path/to/socket` or `http(s)://host:port`. " +
							"Falls back to `KEA_DHCP4_ADDRESS` environment variable.",
						Validators: []validator.String{
							IsValidKeaURL(),
						},
					},
					"http_username": schema.StringAttribute{
						Optional:    true,
						Description: "Username for HTTP basic authentication. Falls back to `KEA_DHCP4_HTTP_USERNAME` environment variable.",
					},
					"http_password": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "Password for HTTP basic authentication. Falls back to `KEA_DHCP4_HTTP_PASSWORD` environment variable.",
					},
					"tls_ca_file": schema.StringAttribute{
						Optional:    true,
						Description: "PEM CA bundle used to verify the Kea HTTPS server. Falls back to `KEA_DHCP4_TLS_CA_FILE` environment variable.",
					},

					"tls_client_cert_file": schema.StringAttribute{
						Optional:    true,
						Description: "PEM client certificate used for mutual TLS authentication. Falls back to `KEA_DHCP4_TLS_CLIENT_CERT_FILE` environment variable.",
					},

					"tls_client_key_file": schema.StringAttribute{
						Optional:    true,
						Sensitive:   true,
						Description: "PEM private key corresponding to tls_client_cert_file. Falls back to `KEA_DHCP4_TLS_CLIENT_KEY_FILE` environment variable.",
					},

					"tls_server_name": schema.StringAttribute{
						Optional:    true,
						Description: "TLS server name used for certificate verification and SNI. Falls back to `KEA_DHCP4_TLS_SERVER_NAME` environment variable.",
					},
				},
			},
		},
	}
}

func (p *KeaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KeaProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var dhcp4Config *KeaProviderClientModel
	if !data.DHCP4.IsNull() && !data.DHCP4.IsUnknown() {
		dhcp4Config = &KeaProviderClientModel{}
		resp.Diagnostics.Append(data.DHCP4.As(ctx, dhcp4Config, basetypes.ObjectAsOptions{})...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	dhcp4Client, err := configureDHCP4Client(dhcp4Config)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Kea DHCP4 client: %s", err))
		return
	}

	keaClients := clients.KeaClients{
		DHCP4: dhcp4Client,
	}

	resp.DataSourceData = keaClients
	resp.ResourceData = keaClients
}

func (p *KeaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return &dhcp4.ReservationResource{} },
	}
}

func (p *KeaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &dhcp4.ReservationDataSource{} },
	}
}

// configValue returns the configured attribute, falling back to the environment variable.
func configValue(data *KeaProviderClientModel, get func(*KeaProviderClientModel) types.String, env string) string {
	if data != nil {
		if v := get(data); !v.IsNull() && !v.IsUnknown() {
			return v.ValueString()
		}
	}
	return os.Getenv(env)
}

func configureDHCP4Client(data *KeaProviderClientModel) (*keadhcp4.Client, error) {
	address := configValue(data, func(m *KeaProviderClientModel) types.String { return m.Address }, "KEA_DHCP4_ADDRESS")

	// No address configured - skip DHCP4 client setup
	if address == "" {
		return nil, nil
	}

	http_username := configValue(data, func(m *KeaProviderClientModel) types.String { return m.HTTPUsername }, "KEA_DHCP4_HTTP_USERNAME")
	http_password := configValue(data, func(m *KeaProviderClientModel) types.String { return m.HTTPPassword }, "KEA_DHCP4_HTTP_PASSWORD")
	tls_ca_file := configValue(data, func(m *KeaProviderClientModel) types.String { return m.TLSCAFile }, "KEA_DHCP4_TLS_CA_FILE")
	tls_client_cert_file := configValue(data, func(m *KeaProviderClientModel) types.String { return m.TLSClientCertFile }, "KEA_DHCP4_TLS_CLIENT_CERT_FILE")
	tls_client_key_file := configValue(data, func(m *KeaProviderClientModel) types.String { return m.TLSClientKeyFile }, "KEA_DHCP4_TLS_CLIENT_KEY_FILE")
	tls_server_name := configValue(data, func(m *KeaProviderClientModel) types.String { return m.TLSServerName }, "KEA_DHCP4_TLS_SERVER_NAME")

	transport, err := newKeaTransport(
		address,
		http_username,
		http_password,
		tls_ca_file,
		tls_client_cert_file,
		tls_client_key_file,
		tls_server_name,
	)
	if err != nil {
		return nil, err
	}

	client := keadhcp4.NewClient(transport)
	return client, nil
}

func newTLSHTTPClient(
	caFile string,
	clientCertFile string,
	clientKeyFile string,
	serverName string,
) (*http.Client, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: serverName,
	}

	if caFile != "" {
		caPEM, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("read TLS CA file: %w", err)
		}

		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(caPEM) {
			return nil, fmt.Errorf("TLS CA file contains no valid certificates")
		}

		tlsConfig.RootCAs = roots
	}

	switch {
	case clientCertFile != "" && clientKeyFile == "":
		return nil, fmt.Errorf("tls_client_key_file is required with tls_client_cert_file")

	case clientCertFile == "" && clientKeyFile != "":
		return nil, fmt.Errorf("tls_client_cert_file is required with tls_client_key_file")

	case clientCertFile != "" && clientKeyFile != "":
		cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return nil, fmt.Errorf("load TLS client certificate: %w", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	baseTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("unexpected default HTTP transport type")
	}

	transport := baseTransport.Clone()
	transport.TLSClientConfig = tlsConfig

	return &http.Client{
		Transport: transport,
	}, nil
}

func newKeaTransport(
	uri string,
	httpUsername string,
	httpPassword string,
	caFile string,
	clientCertFile string,
	clientKeyFile string,
	serverName string,
) (kea.Transport, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "unix":
		return &kea.UnixTransport{
			SocketPath: u.Path,
		}, nil

	case "http":
		return &kea.HTTPTransport{
			Endpoint: uri,
			Username: httpUsername,
			Password: httpPassword,
		}, nil

	case "https":
		client, err := newTLSHTTPClient(
			caFile,
			clientCertFile,
			clientKeyFile,
			serverName,
		)
		if err != nil {
			return nil, err
		}

		return &kea.HTTPTransport{
			Endpoint: uri,
			Client:   client,
			Username: httpUsername,
			Password: httpPassword,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
	}
}
