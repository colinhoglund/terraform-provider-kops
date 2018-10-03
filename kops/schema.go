package kops

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func clusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name":                    schemaStringRequired(),
		"channel":                 schemaStringOptional(),
		"cloud_provider":          schemaStringOptional(),
		"cluster_dnsdomain":       schemaStringOptional(),
		"config_base":             schemaStringOptional(),
		"config_store":            schemaStringOptional(),
		"dnszone":                 schemaStringOptional(),
		"key_store":               schemaStringOptional(),
		"kubernetes_version":      schemaStringOptional(),
		"master_internal_name":    schemaStringOptional(),
		"master_public_name":      schemaStringOptional(),
		"network_cidr":            schemaCIDRStringOptional(),
		"network_id":              schemaStringOptional(),
		"non_masquerade_cidr":     schemaCIDRStringOptional(),
		"project":                 schemaStringOptional(),
		"secret_store":            schemaStringOptional(),
		"service_cluster_iprange": schemaStringOptional(),
		"sshkey_name":             schemaStringOptional(),
		"creation_timestamp":      schemaStringComputed(),
		"subnet":                  schemaSubnet(),
	}
}

// complex schema objects
func schemaSubnet() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.CIDRNetwork(1, 32),
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"Public", "Private", "Utility"}, false),
				},
				"zone": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

// generic helper schemas
func schemaStringOptional() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
}

func schemaStringRequired() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
}

func schemaStringComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func schemaCIDRStringOptional() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.CIDRNetwork(1, 32),
	}
}
