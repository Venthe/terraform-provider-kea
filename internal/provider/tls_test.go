package provider

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/venthe/terraform-provider-kea/kea"
)

// writeTestPKI writes a self-signed CA/cert pair and returns the cert and key paths.
func writeTestPKI(t *testing.T) (certPath, keyPath string) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "test"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")
	if err := os.WriteFile(certPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keyPath, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600); err != nil {
		t.Fatal(err)
	}
	return certPath, keyPath
}

func TestNewTLSHTTPClient(t *testing.T) {
	cert, key := writeTestPKI(t)
	bad := filepath.Join(t.TempDir(), "bad.pem")
	if err := os.WriteFile(bad, []byte("not pem"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name               string
		ca, cert, key, sni string
		wantErr            string
	}{
		{name: "none"},
		{name: "ca only", ca: cert},
		{name: "mtls", ca: cert, cert: cert, key: key, sni: "edge01.home.arpa"},
		{name: "missing ca file", ca: "/nonexistent", wantErr: "read TLS CA file"},
		{name: "invalid ca", ca: bad, wantErr: "no valid certificates"},
		{name: "cert without key", cert: cert, wantErr: "tls_client_key_file is required"},
		{name: "key without cert", key: key, wantErr: "tls_client_cert_file is required"},
		{name: "bad keypair", cert: bad, key: key, wantErr: "load TLS client certificate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := newTLSHTTPClient(tt.ca, tt.cert, tt.key, tt.sni)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("want error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if client == nil {
				t.Fatal("nil client")
			}
		})
	}
}

func TestNewKeaTransport(t *testing.T) {
	cert, key := writeTestPKI(t)

	tr, err := newKeaTransport("https://127.0.0.1:8000", "", "", cert, cert, key, "edge01.home.arpa")
	if err != nil {
		t.Fatal(err)
	}
	ht, ok := tr.(*kea.HTTPTransport)
	if !ok || ht.Client == nil {
		t.Fatalf("expected HTTPTransport with TLS client, got %#v", tr)
	}

	if _, err := newKeaTransport("https://127.0.0.1:8000", "", "", "", cert, "", ""); err == nil {
		t.Fatal("expected error for cert without key")
	}
	if _, err := newKeaTransport("ftp://x", "", "", "", "", "", ""); err == nil {
		t.Fatal("expected unsupported scheme error")
	}
}

func TestConfigureDHCP4ClientTLSEnv(t *testing.T) {
	cert, key := writeTestPKI(t)
	t.Setenv("KEA_DHCP4_ADDRESS", "https://127.0.0.1:8000")
	t.Setenv("KEA_DHCP4_TLS_CA_FILE", cert)
	t.Setenv("KEA_DHCP4_TLS_CLIENT_CERT_FILE", cert)
	t.Setenv("KEA_DHCP4_TLS_CLIENT_KEY_FILE", key)
	t.Setenv("KEA_DHCP4_TLS_SERVER_NAME", "edge01.home.arpa")

	c, err := configureDHCP4Client(nil)
	if err != nil || c == nil {
		t.Fatalf("env config failed: %v", err)
	}

	// Explicit config overrides env: key without cert is an error.
	_, err = configureDHCP4Client(&KeaProviderClientModel{
		Address:           types.StringValue("https://127.0.0.1:8000"),
		TLSCAFile:         types.StringNull(),
		TLSClientCertFile: types.StringValue(""),
		TLSClientKeyFile:  types.StringValue(key),
		TLSServerName:     types.StringNull(),
		HTTPUsername:      types.StringNull(),
		HTTPPassword:      types.StringNull(),
	})
	if err == nil {
		t.Fatal("expected error for key without cert")
	}
}
