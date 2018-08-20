package main

import (
	"net/url"

	maxcdn "github.com/contentflow/go-maxcdn"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return provider()
		},
	})
}

func provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"company_alias": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STACKPATH_ALIAS", nil),
				Description: "Alias for the Stackpath API OAuth app",
			},
			"consumer_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STACKPATH_CONSUMER_KEY", nil),
				Description: "Consumer key for the Stackpath API OAuth app",
			},
			"consumer_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("STACKPATH_CONSUMER_SECRET", nil),
				Description: "Consumer secret for the Stackpath API OAuth app",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.stackpath.com/v1",
				Description: "API endpoint for the Stackpath API",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"stackpath_ssl_certificate": sslCertResource(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	maxcdn.APIHost = d.Get("endpoint").(string)
	return maxcdn.NewMaxCDN(
		d.Get("company_alias").(string),
		d.Get("consumer_key").(string),
		d.Get("consumer_secret").(string),
	), nil
}

func sslCertResource() *schema.Resource {
	return &schema.Resource{
		Create: createSSLCert,
		Read:   readSSLCert,
		Update: updateSSLCert,
		Delete: deleteSSLCert,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use descriptive name (default: auto_{domain}{expiration}{update})",
			},
			"ssl_crt": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate",
			},
			"ssl_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Certificate Private Key",
				Sensitive:   true,
			},
			"ssl_cabundle": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate Authority Intermediate Bundle",
			},
		},
	}
}

type envelope struct {
	SSL struct {
		CompanyID      string  `json:"company_id"`
		DateExpiration string  `json:"date_expiration"`
		Domain         string  `json:"domain"`
		Globalsign     strint  `json:"globalsign"`
		ID             strint  `json:"id"`
		Name           string  `json:"name"`
		SSLCabundle    *string `json:"ssl_cabundle"`
		SSLCrt         string  `json:"ssl_crt"`
		Wildcard       strint  `json:"wildcard"`
	} `json:"ssl"`
}

func createSSLCert(d *schema.ResourceData, m interface{}) error {
	client := m.(*maxcdn.MaxCDN)

	data := url.Values{}
	for _, key := range []string{"name", "ssl_crt", "ssl_key", "ssl_cabundle"} {
		if d.Get(key).(string) != "" {
			data.Set(key, d.Get(key).(string))
		}
	}

	var r envelope
	_, err := client.Post(&r, "/ssl", data)
	if err != nil {
		return err
	}

	d.Set("name", r.SSL.Name)
	d.SetId(r.SSL.ID.String())

	return nil
}

func readSSLCert(d *schema.ResourceData, m interface{}) error {
	client := m.(*maxcdn.MaxCDN)

	var r envelope
	_, err := client.Get(&r, "/ssl/"+d.Id(), nil)
	if err != nil {
		return err
	}

	d.Set("name", r.SSL.Name)
	d.Set("ssl_crt", r.SSL.SSLCrt)
	if r.SSL.SSLCabundle != nil {
		d.Set("ssl_cabundle", *r.SSL.SSLCabundle)
	}

	return nil
}

func updateSSLCert(d *schema.ResourceData, m interface{}) error {
	client := m.(*maxcdn.MaxCDN)

	data := url.Values{}
	for _, key := range []string{"name", "ssl_crt", "ssl_key", "ssl_cabundle"} {
		if d.Get(key).(string) != "" {
			data.Set(key, d.Get(key).(string))
		}
	}

	data.Set("force", "0")

	var r envelope
	_, err := client.Put(&r, "/ssl/"+d.Id(), data)
	if err != nil {
		return err
	}

	d.Set("name", r.SSL.Name)
	d.SetId(r.SSL.ID.String())

	return nil
}

func deleteSSLCert(d *schema.ResourceData, m interface{}) error {
	client := m.(*maxcdn.MaxCDN)

	_, err := client.Delete("/ssl/"+d.Id(), nil)
	return err
}
