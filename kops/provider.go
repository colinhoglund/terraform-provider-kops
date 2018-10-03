package kops

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/kops/pkg/client/simple/vfsclientset"
	"k8s.io/kops/util/pkg/vfs"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"state_store": {
				Description: "Location of state storage.",
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KOPS_STATE_STORE", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kops_cluster": dataSourceCluster(),
			//"kops_instance_group": dataSourceInstanceGroup(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"kops_cluster": resourceCluster(),
			//"kops_instance_group": resourceInstanceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

const (
	INVALID_STATE_ERROR = `Unable to read state store s3 bucket.
Please use a valid s3 bucket uri on state_store attribute or KOPS_STATE_STORE env var.
A valid value follows the format s3://<bucket>.
Trailing slash will be trimmed.`
)

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	registryPath := d.Get("state_store").(string)
	basePath, err := vfs.Context.BuildVfsPath(registryPath)
	if err != nil {
		return nil, fmt.Errorf("error building path for %q: %v", registryPath, err)
	}

	if !vfs.IsClusterReadable(basePath) {
		return nil, field.Invalid(field.NewPath("State Store"), registryPath, INVALID_STATE_ERROR)
	}

	return vfsclientset.NewVFSClientset(basePath, true), nil
}
