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

	if err := d.Set("metadata", flattenMetadata(cluster)); err != nil {
		return err
	}
	if err := d.Set("spec", flattenSpec(cluster)); err != nil {
		return err
	}
	return nil
}

func flattenSpec(cluster *kopsapi.Cluster) []map[string]interface{} {
	spec := make(map[string]interface{})

	spec["channel"] = cluster.Spec.Channel
	spec["cloud_provider"] = cluster.Spec.CloudProvider
	spec["cluster_dnsdomain"] = cluster.Spec.ClusterDNSDomain
	spec["config_base"] = cluster.Spec.ConfigBase
	spec["config_store"] = cluster.Spec.ConfigStore
	spec["dnszone"] = cluster.Spec.DNSZone
	spec["key_store"] = cluster.Spec.KeyStore
	spec["kubernetes_version"] = cluster.Spec.KubernetesVersion
	spec["master_internal_name"] = cluster.Spec.MasterInternalName
	spec["master_public_name"] = cluster.Spec.MasterPublicName
	spec["network_cidr"] = cluster.Spec.NetworkCIDR
	spec["network_id"] = cluster.Spec.NetworkID
	spec["non_masquerade_cidr"] = cluster.Spec.NonMasqueradeCIDR
	spec["project"] = cluster.Spec.Project
	spec["secret_store"] = cluster.Spec.SecretStore
	spec["service_cluster_iprange"] = cluster.Spec.ServiceClusterIPRange
	spec["sshkey_name"] = cluster.Spec.SSHKeyName
	spec["subnet"] = flattenSubnet(cluster.Spec.Subnets)
	spec["topology"] = flattenTopology(cluster.Spec.Topology)

	return []map[string]interface{}{spec}
}

func flattenMetadata(cluster *kopsapi.Cluster) []map[string]interface{} {
	meta := make(map[string]interface{})

	meta["name"] = cluster.ObjectMeta.Name
	meta["creation_timestamp"] = cluster.ObjectMeta.CreationTimestamp.String()

	return []map[string]interface{}{meta}
}

func flattenSubnet(subnets []kopsapi.ClusterSubnetSpec) []map[string]interface{} {
	var out []map[string]interface{}
	for _, subnet := range subnets {
		out = append(out, map[string]interface{}{
			"name": subnet.Name,
			"cidr": subnet.CIDR,
			"zone": subnet.Zone,
			"type": string(subnet.Type),
		})
	}
	return out
}

func flattenTopology(topology *kopsapi.TopologySpec) []map[string]interface{} {
	out := make(map[string]interface{})

	out["masters"] = topology.Masters
	out["nodes"] = topology.Nodes
	if topology.Bastion != nil {
		out["bastion"] = []map[string]interface{}{
			map[string]interface{}{
				"bastion_public_name":  topology.Bastion.BastionPublicName,
				"idle_timeout_seconds": topology.Bastion.IdleTimeoutSeconds,
			},
		}
	}
	out["dns"] = []map[string]interface{}{
		map[string]interface{}{
			"type": topology.DNS.Type,
		},
	}

	return []map[string]interface{}{out}
}
