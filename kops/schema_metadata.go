package kops

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func schemaMetadata() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":               schemaStringRequired(),
				"creation_timestamp": schemaStringComputed(),
			},
		},
	}
}
