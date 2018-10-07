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
		Schema: map[string]*schema.Schema{
			"metadata": schemaMetadata(),
			"spec":     schemaClusterSpec(),
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

	if err := d.Set("metadata", resourceDataClusterMetadata(cluster)); err != nil {
		return err
	}
	if err := d.Set("spec", resourceDataClusterSpec(cluster)); err != nil {
		return err
	}
	return nil
}

func resourceDataClusterSpec(cluster *kopsapi.Cluster) []map[string]interface{} {
	data := make(map[string]interface{})

	data["channel"] = cluster.Spec.Channel
	data["cloud_provider"] = cluster.Spec.CloudProvider
	data["cluster_dnsdomain"] = cluster.Spec.ClusterDNSDomain
	data["config_base"] = cluster.Spec.ConfigBase
	data["config_store"] = cluster.Spec.ConfigStore
	data["dnszone"] = cluster.Spec.DNSZone
	data["key_store"] = cluster.Spec.KeyStore
	data["kubernetes_version"] = cluster.Spec.KubernetesVersion
	data["master_internal_name"] = cluster.Spec.MasterInternalName
	data["master_public_name"] = cluster.Spec.MasterPublicName
	data["network_cidr"] = cluster.Spec.NetworkCIDR
	data["network_id"] = cluster.Spec.NetworkID
	data["non_masquerade_cidr"] = cluster.Spec.NonMasqueradeCIDR
	data["project"] = cluster.Spec.Project
	data["secret_store"] = cluster.Spec.SecretStore
	data["service_cluster_iprange"] = cluster.Spec.ServiceClusterIPRange
	data["sshkey_name"] = cluster.Spec.SSHKeyName
	data["subnet"] = resourceDataClusterSubnet(cluster.Spec.Subnets)
	data["topology"] = resourceDataClusterTopology(cluster.Spec.Topology)
	data["ssh_access"] = cluster.Spec.SSHAccess
	data["kubernetes_api_access"] = cluster.Spec.KubernetesAPIAccess
	data["additional_policies"] = *cluster.Spec.AdditionalPolicies
	data["etcd_cluster"] = resourceDataClusterEtcdCluster(cluster.Spec.EtcdClusters)

	return []map[string]interface{}{data}
}

func resourceDataClusterMetadata(cluster *kopsapi.Cluster) []map[string]interface{} {
	data := make(map[string]interface{})

	data["name"] = cluster.ObjectMeta.Name
	data["creation_timestamp"] = cluster.ObjectMeta.CreationTimestamp.String()

	return []map[string]interface{}{data}
}

func resourceDataClusterSubnet(subnets []kopsapi.ClusterSubnetSpec) []map[string]interface{} {
	var data []map[string]interface{}
	for _, subnet := range subnets {
		data = append(data, map[string]interface{}{
			"name": subnet.Name,
			"cidr": subnet.CIDR,
			"zone": subnet.Zone,
			"type": string(subnet.Type),
		})
	}
	return data
}

func resourceDataClusterTopology(topology *kopsapi.TopologySpec) []map[string]interface{} {
	data := make(map[string]interface{})

	data["masters"] = topology.Masters
	data["nodes"] = topology.Nodes
	if topology.Bastion != nil {
		data["bastion"] = []map[string]interface{}{
			map[string]interface{}{
				"bastion_public_name":  topology.Bastion.BastionPublicName,
				"idle_timeout_seconds": topology.Bastion.IdleTimeoutSeconds,
			},
		}
	}
	data["dns"] = []map[string]interface{}{
		map[string]interface{}{
			"type": topology.DNS.Type,
		},
	}

	return []map[string]interface{}{data}
}

func resourceDataClusterEtcdCluster(etcdClusters []*kopsapi.EtcdClusterSpec) []map[string]interface{} {
	var data []map[string]interface{}

	for _, cluster := range etcdClusters {
		cl := make(map[string]interface{})

		cl["name"] = cluster.Name

		//if cluster.Provider != nil {
		//	cl["provider"] = cluster.Provider
		//}

		// build etcd_members
		var members []map[string]interface{}
		for _, member := range cluster.Members {
			mem := make(map[string]interface{})
			mem["name"] = member.Name
			mem["instance_group"] = *member.InstanceGroup
			if member.VolumeType != nil {
				mem["volume_type"] = *member.VolumeType
			}
			if member.VolumeIops != nil {
				mem["volume_iops"] = *member.VolumeIops
			}
			if member.VolumeSize != nil {
				mem["volume_size"] = *member.VolumeSize
			}
			if member.KmsKeyId != nil {
				mem["kms_key_id"] = *member.KmsKeyId
			}
			if member.EncryptedVolume != nil {
				mem["encrypted_volume"] = *member.EncryptedVolume
			}
			members = append(members, mem)
		}
		cl["etcd_member"] = members

		cl["enable_etcd_tls"] = cluster.EnableEtcdTLS
		cl["enable_tls_auth"] = cluster.EnableTLSAuth
		cl["version"] = cluster.Version
		if cluster.LeaderElectionTimeout != nil {
			cl["leader_election_timeout"] = cluster.LeaderElectionTimeout
		}
		if cluster.HeartbeatInterval != nil {
			cl["heartbeat_interval"] = cluster.HeartbeatInterval
		}
		cl["image"] = cluster.Image
		//cl["backups"] = []map[string]interface{}{
		//	map[string]interface{}{
		//		"store": cluster.Backups.BackupStore,
		//		"image": cluster.Backups.Image,
		//	},
		//}
		//cl["manager"] = []map[string]interface{}{
		//	map[string]interface{}{
		//		"image": cluster.Manager.Image,
		//	},
		//}

		data = append(data, cl)
	}

	return data
}
