package kops

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kops/pkg/client/simple/vfsclientset"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,
		Exists: resourceClusterExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: diffJSON,
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	return setResourceData(d, m)
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceClusterExists(d *schema.ResourceData, m interface{}) (bool, error) {
	clientset := m.(*vfsclientset.VFSClientset)
	_, err := clientset.GetCluster(d.Id())
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func setResourceData(d *schema.ResourceData, m interface{}) error {
	// get cluster
	clientset := m.(*vfsclientset.VFSClientset)
	cluster, err := clientset.GetCluster(d.Id())
	if err != nil {
		return err
	}

	clusterJSON, err := json.Marshal(cluster)
	if err != nil {
		return err
	}

	if err := d.Set("content", string(clusterJSON)); err != nil {
		return err
	}
	return nil
}

func diffJSON(k, old, new string, d *schema.ResourceData) bool {
	var o interface{}
	var n interface{}
	var err error

	err = json.Unmarshal([]byte(old), &o)
	if err != nil {
		return false
	}
	err = json.Unmarshal([]byte(new), &n)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(o, n)
}
