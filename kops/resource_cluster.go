package kops

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kops/pkg/client/simple/vfsclientset"

	kopsapi "k8s.io/kops/pkg/apis/kops"
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
		Schema: clusterSchema(),
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
		} else {
			return false, err
		}
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

	d.Set("channel", cluster.Spec.Channel)
	d.Set("cloud_provider", cluster.Spec.CloudProvider)
	d.Set("cluster_dnsdomain", cluster.Spec.ClusterDNSDomain)
	d.Set("config_base", cluster.Spec.ConfigBase)
	d.Set("config_store", cluster.Spec.ConfigStore)
	d.Set("creation_timestamp", cluster.ObjectMeta.CreationTimestamp.String())
	d.Set("dnszone", cluster.Spec.DNSZone)
	d.Set("key_store", cluster.Spec.KeyStore)
	d.Set("kubernetes_version", cluster.Spec.KubernetesVersion)
	d.Set("master_internal_name", cluster.Spec.MasterInternalName)
	d.Set("master_public_name", cluster.Spec.MasterPublicName)
	d.Set("name", cluster.ObjectMeta.Name)
	d.Set("network_cidr", cluster.Spec.NetworkCIDR)
	d.Set("network_id", cluster.Spec.NetworkID)
	d.Set("non_masquerade_cidr", cluster.Spec.NonMasqueradeCIDR)
	d.Set("project", cluster.Spec.Project)
	d.Set("secret_store", cluster.Spec.SecretStore)
	d.Set("service_cluster_iprange", cluster.Spec.ServiceClusterIPRange)
	d.Set("sshkey_name", cluster.Spec.SSHKeyName)

	// set subnets
	d.Set("subnet", flattenSubnet(cluster.Spec.Subnets))
	return nil
}

func flattenSubnet(subnets []kopsapi.ClusterSubnetSpec) []map[string]string {
	var out []map[string]string
	for _, subnet := range subnets {
		out = append(out, map[string]string{
			"name": subnet.Name,
			"cidr": subnet.CIDR,
			"zone": subnet.Zone,
			"type": string(subnet.Type),
		})
	}
	return out
}
