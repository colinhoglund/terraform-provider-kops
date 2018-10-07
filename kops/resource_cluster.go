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
	spec["subnet"] = resourceDataClusterSubnet(cluster.Spec.Subnets)
	spec["topology"] = resourceDataClusterTopology(cluster.Spec.Topology)
	spec["ssh_access"] = cluster.Spec.SSHAccess
	spec["kubernetes_api_access"] = cluster.Spec.KubernetesAPIAccess
	spec["additional_policies"] = *cluster.Spec.AdditionalPolicies
	spec["etcd_cluster"] = resourceDataClusterEtcdCluster(cluster.Spec.EtcdClusters)

	return []map[string]interface{}{spec}
}

func resourceDataClusterMetadata(cluster *kopsapi.Cluster) []map[string]interface{} {
	meta := make(map[string]interface{})

	meta["name"] = cluster.ObjectMeta.Name
	meta["creation_timestamp"] = cluster.ObjectMeta.CreationTimestamp.String()

	return []map[string]interface{}{meta}
}

func resourceDataClusterSubnet(subnets []kopsapi.ClusterSubnetSpec) []map[string]interface{} {
	var ret []map[string]interface{}
	for _, subnet := range subnets {
		ret = append(ret, map[string]interface{}{
			"name": subnet.Name,
			"cidr": subnet.CIDR,
			"zone": subnet.Zone,
			"type": string(subnet.Type),
		})
	}
	return ret
}

func resourceDataClusterTopology(topology *kopsapi.TopologySpec) []map[string]interface{} {
	ret := make(map[string]interface{})

	ret["masters"] = topology.Masters
	ret["nodes"] = topology.Nodes
	if topology.Bastion != nil {
		ret["bastion"] = []map[string]interface{}{
			map[string]interface{}{
				"bastion_public_name":  topology.Bastion.BastionPublicName,
				"idle_timeout_seconds": topology.Bastion.IdleTimeoutSeconds,
			},
		}
	}
	ret["dns"] = []map[string]interface{}{
		map[string]interface{}{
			"type": topology.DNS.Type,
		},
	}

	return []map[string]interface{}{ret}
}

func resourceDataClusterEtcdCluster(etcdClusters []*kopsapi.EtcdClusterSpec) []map[string]interface{} {
	var ret []map[string]interface{}

	for _, cluster := range etcdClusters {
		cur := make(map[string]interface{})

		cur["name"] = cluster.Name

		//if cluster.Provider != nil {
		//	cur["provider"] = cluster.Provider
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
		cur["etcd_member"] = members

		cur["enable_etcd_tls"] = cluster.EnableEtcdTLS
		cur["enable_tls_auth"] = cluster.EnableTLSAuth
		cur["version"] = cluster.Version
		if cluster.LeaderElectionTimeout != nil {
			cur["leader_election_timeout"] = cluster.LeaderElectionTimeout
		}
		if cluster.HeartbeatInterval != nil {
			cur["heartbeat_interval"] = cluster.HeartbeatInterval
		}
		cur["image"] = cluster.Image
		//cur["backups"] = []map[string]interface{}{
		//	map[string]interface{}{
		//		"store": cluster.Backups.BackupStore,
		//		"image": cluster.Backups.Image,
		//	},
		//}
		//cur["manager"] = []map[string]interface{}{
		//	map[string]interface{}{
		//		"image": cluster.Manager.Image,
		//	},
		//}

		ret = append(ret, cur)
	}

	return ret
}
