package kops

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func schemaClusterSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
				"project":                 schemaStringOptional(),
				"secret_store":            schemaStringOptional(),
				"service_cluster_iprange": schemaStringOptional(),
				"sshkey_name":             schemaStringOptional(),
				"network_id":              schemaStringOptional(),
				"network_cidr":            schemaCIDRStringOptional(),
				"non_masquerade_cidr":     schemaCIDRStringOptional(),
				"ssh_access":              schemaStringSliceOptional(),
				"kubernetes_api_access":   schemaStringSliceOptional(),
				"additional_policies":     schemaStringMap(),
				"subnet":                  schemaClusterSubnet(),
				"topology":                schemaClusterTopology(),
				"etcd_cluster":            schemaClusterEtcdCluster(),
			},
		},
	}
}

func schemaClusterSubnet() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"cidr": schemaCIDRStringRequired(),
				"name": schemaStringRequired(),
				"type": schemaStringInSliceRequired([]string{"Public", "Private", "Utility"}),
				"zone": schemaStringRequired(),
			},
		},
	}
}

func schemaClusterTopology() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"masters": schemaStringOptional(),
				"nodes":   schemaStringOptional(),
				"bastion": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"bastion_public_name":  schemaStringOptional(),
							"idle_timeout_seconds": schemaIntOptional(),
						},
					},
				},
				"dns": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": schemaStringInSliceOptional([]string{"Public", "Private"}),
						},
					},
				},
			},
		},
	}
}

func schemaClusterEtcdCluster() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": schemaStringRequired(),
				//"provider": {
				//	Type:         schema.TypeString,
				//	Optional:     true,
				//	ValidateFunc: validation.StringInSlice([]string{"Manager", "Legacy"}, false),
				//},
				"etcd_member": {
					Type: schema.TypeList, Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"name":             schemaStringRequired(),
							"instance_group":   schemaStringRequired(),
							"volume_type":      schemaStringOptional(),
							"volume_iops":      schemaIntOptional(),
							"volume_size":      schemaIntOptional(),
							"kms_key_id":       schemaStringOptional(),
							"encrypted_volume": schemaBoolOptional(),
						},
					},
				},
				"enable_etcd_tls":         schemaBoolOptional(),
				"enable_tls_auth":         schemaBoolOptional(),
				"version":                 schemaStringOptional(),
				"leader_election_timeout": schemaIntOptional(),
				"heartbeat_interval":      schemaIntOptional(),
				"image":                   schemaStringOptional(),
				//"backups": {
				//	Type:     schema.TypeList,
				//	Optional: true,
				//	MaxItems: 1,
				//	Elem: &schema.Resource{
				//		Schema: map[string]*schema.Schema{
				//			"store": schemaStringOptional(),
				//			"image": schemaStringOptional(),
				//		},
				//	},
				//},
				//"manager": {
				//	Type:     schema.TypeList,
				//	Optional: true,
				//	MaxItems: 1,
				//	Elem: &schema.Resource{
				//		Schema: map[string]*schema.Schema{
				//			"image": schemaStringOptional(),
				//		},
				//	},
				//},
			},
		},
	}
}