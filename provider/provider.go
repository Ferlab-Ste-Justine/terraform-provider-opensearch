package provider

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Description: "Name of the opensearch user that will be used to access opensearch. Can alternatively be set with the OPENSEARCH_USERNAME environment variable. Can also be omitted if tls certificate authentication will be used instead as the username will be infered from the certificate.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_USERNAME", ""),
			},
			"password": &schema.Schema{
				Description: "Password of the opensearch user that will be used to access opensearch. Can alternatively be set with the OPENSEARCH_PASSWORD environment variable. Can also be omitted if tls certificate authentication will be used instead.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_PASSWORD", ""),
			},
			"ca_cert": &schema.Schema{
				Description: "File that contains the CA certificate that signed the opensearch servers' certificates. Can alternatively be set with the OPENSEARCH_CACERT environment variable. Can also be omitted.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_CACERT", ""),
			},
			"cert": &schema.Schema{
				Description: "File that contains the client certificate used to authentify the user. Can alternatively be set with the OPENSEARCH_CERT environment variable. Can be omitted if password authentication is used.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_CERT", ""),
			},
			"key": &schema.Schema{
				Description: "File that contains the client encryption key used to authentify the user. Can alternatively be set with the OPENSEARCH_KEY environment variable. Can be omitted if password authentication is used.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_KEY", ""),
			},
			"endpoints": &schema.Schema{
				Description: "Endpoints of the opensearch servers. The entry of each server should follow the http|https://ip:port format and be coma separated. Can alternatively be set with the OPENSEARCH_ENDPOINTS environment variable.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OPENSEARCH_ENDPOINTS", ""),
			},
			"connection_timeout": &schema.Schema{
				Description: "Timeout to establish the opensearch servers connection in golang duration format. Defaults to 10 seconds.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10s",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := time.ParseDuration(v)
					if err != nil {
						return []string{}, []error{errors.New("connection_timeout must be a value golang duration string value")}
					}

					return []string{}, []error{}
				},
			},
			"request_timeout": &schema.Schema{
				Description: "Timeout for individual requests the provider makes on the opensearch servers in golang duration format. Defaults to 10 seconds.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "10s",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := time.ParseDuration(v)
					if err != nil {
						return []string{}, []error{errors.New("request_timeout must be a value golang duration string value")}
					}

					return []string{}, []error{}
				},
			},
			"retries": &schema.Schema{
				Description: "Number of times operations that result in retriable errors should be re-attempted. Defaults to 10.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"opensearch_role": resourceOpensearchRole(),
			"opensearch_user": resourceOpensearchUser(),
			"opensearch_role_mapping": resourceOpensearchRoleMapping(),
			"opensearch_ism_policy": resourceOpensearchIsmPolicy(),
		},
		DataSourcesMap: map[string]*schema.Resource{
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	endpoints, _ := d.Get("endpoints").(string)
	username, _ := d.Get("username").(string)
	password, _ := d.Get("password").(string)
	caCert, _ := d.Get("ca_cert").(string)
	cert, _ := d.Get("cert").(string)
	key, _ := d.Get("key").(string)
	connectionTimeout, _ := d.Get("connection_timeout").(string)
	requestTimeout, _ := d.Get("request_timeout").(string)
	retries, _ := d.Get("retries").(int)
	tlsConf := &tls.Config{}

	if cert != "" {
		certData, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		(*tlsConf).Certificates = []tls.Certificate{certData}
		(*tlsConf).InsecureSkipVerify = false
	}

	if caCert != "" {
		caCertContent, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to read root certificate file: %s", err.Error()))
		}
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCertContent)
		if !ok {
			return nil, errors.New("Failed to parse root certificate authority")
		}
		(*tlsConf).RootCAs = roots
	}

	pConnectionTimeout, _ := time.ParseDuration(connectionTimeout)
	pRequestTimeout, _ := time.ParseDuration(requestTimeout)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConf,
			TLSHandshakeTimeout: pConnectionTimeout,
			ResponseHeaderTimeout: pRequestTimeout,
		},
	}

	arrEndpoints := strings.Split(endpoints, ",")
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arrEndpoints), func(i, j int) { arrEndpoints[i], arrEndpoints[j] = arrEndpoints[j], arrEndpoints[i] })

	return OpensearchClient{
		Client: client,
		Endpoints: arrEndpoints,
		Username: username,
		Password: password,
		Retries: retries,
	}, nil
}
