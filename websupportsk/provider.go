package websupportsk

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema:               providerSchema(),
		ResourcesMap:         providerResources(),
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return NewClient(d.Get("api_key").(string), d.Get("api_secret").(string), d.Get("api_url").(string)), nil
}

func providerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("WEBSUPPORTSK_API_KEY", ""),
			Description: "API KEY used to authenticate with",
		},

		"api_secret": {
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc("WEBSUPPORTSK_API_SECRET", ""),
			Description: "API SECRET used to authenticate with",
		},
		"api_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The URL to the API",
			DefaultFunc: schema.EnvDefaultFunc("WEBSUPPORTSK_API_URL", "https://rest.websupport.sk"),
		},
	}
}

func providerResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"websupportsk_dns_record": resourceDnsRecord(),
	}
}
