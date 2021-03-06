package api

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/glog"
	k8sapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"net"
	"strconv"
	"strings"
)

type Cluster struct {
	unversioned.TypeMeta `json:",inline"`
	k8sapi.ObjectMeta    `json:"metadata,omitempty"`

	Spec ClusterSpec `json:"spec,omitempty"`
}

type ClusterSpec struct {
	// The CloudProvider to use (aws or gce)
	CloudProvider string `json:"cloudProvider,omitempty"`

	// The version of kubernetes to install (optional, and can be a "spec" like stable)
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	//
	//// The Node initializer technique to use: cloudinit or nodeup
	//NodeInit                      string `json:",omitempty"`

	// Configuration of zones we are targeting
	Zones []*ClusterZoneSpec `json:"zones,omitempty"`
	//Region                        string        `json:",omitempty"`

	// Project is the cloud project we should use, required on GCE
	Project string `json:"project,omitempty"`

	// MasterPermissions contains the IAM permissions for the masters
	MasterPermissions *CloudPermissions `json:"masterPermissions,omitempty"`
	// NodePermissions contains the IAM permissions for the nodes
	NodePermissions *CloudPermissions `json:"nodePermissions,omitempty"`

	// MasterPublicName is the external DNS name for the master nodes
	MasterPublicName string `json:"masterPublicName,omitempty"`
	// MasterInternalName is the internal DNS name for the master nodes
	MasterInternalName string `json:"masterInternalName,omitempty"`

	// The CIDR used for the AWS VPC / GCE Network, or otherwise allocated to k8s
	// This is a real CIDR, not the internal k8s network
	NetworkCIDR string `json:"networkCIDR,omitempty"`

	// NetworkID is an identifier of a network, if we want to reuse/share an existing network (e.g. an AWS VPC)
	NetworkID string `json:"networkID,omitempty"`

	// SecretStore is the VFS path to where secrets are stored
	SecretStore string `json:"secretStore,omitempty"`
	// KeyStore is the VFS path to where SSL keys and certificates are stored
	KeyStore string `json:"keyStore,omitempty"`
	// ConfigStore is the VFS path to where the configuration (CloudConfig, NodeSetConfig etc) is stored
	ConfigStore string `json:"configStore,omitempty"`

	// DNSZone is the DNS zone we should use when configuring DNS
	// This is because some clouds let us define a managed zone foo.bar, and then have
	// kubernetes.dev.foo.bar, without needing to define dev.foo.bar as a hosted zone.
	// DNSZone will probably be a suffix of the MasterPublicName and MasterInternalName
	DNSZone string `json:"dnsZone,omitempty"`

	// ClusterDNSDomain is the suffix we use for internal DNS names (normally cluster.local)
	ClusterDNSDomain string `json:"clusterDNSDomain,omitempty"`

	//InstancePrefix                string `json:",omitempty"`

	// ClusterName is a unique identifier for the cluster, and currently must be a DNS name
	//ClusterName       string `json:",omitempty"`

	//AllocateNodeCIDRs *bool `json:"allocateNodeCIDRs,omitempty"`

	Multizone *bool `json:"multizone,omitempty"`

	//ClusterIPRange                string `json:",omitempty"`

	// ServiceClusterIPRange is the CIDR, from the internal network, where we allocate IPs for services
	ServiceClusterIPRange string `json:"serviceClusterIPRange,omitempty"`
	//MasterIPRange                 string `json:",omitempty"`

	// NonMasqueradeCIDR is the CIDR for the internal k8s network (on which pods & services live)
	// It cannot overlap ServiceClusterIPRange
	NonMasqueradeCIDR string `json:"nonMasqueradeCIDR,omitempty"`

	//NetworkProvider               string `json:",omitempty"`
	//
	//HairpinMode                   string `json:",omitempty"`
	//
	//OpencontrailTag               string `json:",omitempty"`
	//OpencontrailKubernetesTag     string `json:",omitempty"`
	//OpencontrailPublicSubnet      string `json:",omitempty"`
	//
	//EnableClusterMonitoring       string `json:",omitempty"`
	//EnableL7LoadBalancing         string `json:",omitempty"`
	//EnableClusterUI               *bool  `json:",omitempty"`
	//
	//EnableClusterDNS              *bool  `json:",omitempty"`
	//DNSReplicas                   int    `json:",omitempty"`
	//DNSServerIP                   string `json:",omitempty"`

	//EnableClusterLogging          *bool  `json:",omitempty"`
	//EnableNodeLogging             *bool  `json:",omitempty"`
	//LoggingDestination            string `json:",omitempty"`
	//ElasticsearchLoggingReplicas  int    `json:",omitempty"`
	//
	//EnableClusterRegistry         *bool  `json:",omitempty"`
	//ClusterRegistryDisk           string `json:",omitempty"`
	//ClusterRegistryDiskSize       int    `json:",omitempty"`
	//
	//EnableCustomMetrics           *bool `json:",omitempty"`
	//
	//RegisterMasterKubelet         *bool  `json:",omitempty"`

	//// Image is the default image spec to use for the cluster
	//Image                     string `json:",omitempty"`

	//KubeUser                      string `json:",omitempty"`
	//
	//// These are moved to CAStore / SecretStore
	////KubePassword			string
	////KubeletToken                  string
	////KubeProxyToken                string
	////BearerToken                   string
	////CACert                        []byte
	////CAKey                         []byte
	////KubeletCert                   []byte
	////KubeletKey                    []byte
	////MasterCert                    []byte
	////MasterKey                     []byte
	////KubecfgCert                   []byte
	////KubecfgKey                    []byte
	//
	//AdmissionControl              string `json:",omitempty"`
	//RuntimeConfig                 string `json:",omitempty"`
	//
	//KubeImageTag                  string `json:",omitempty"`
	//KubeDockerRegistry            string `json:",omitempty"`
	//KubeAddonRegistry             string `json:",omitempty"`
	//
	//KubeletPort                   int `json:",omitempty"`
	//
	//KubeApiserverRequestTimeout   int `json:",omitempty"`
	//
	//TerminatedPodGcThreshold      string `json:",omitempty"`
	//
	//EnableManifestURL             *bool  `json:",omitempty"`
	//ManifestURL                   string `json:",omitempty"`
	//ManifestURLHeader             string `json:",omitempty"`
	//
	//TestCluster                   string `json:",omitempty"`
	//
	//DockerOptions                 string `json:",omitempty"`
	//DockerStorage                 string `json:",omitempty"`
	//ExtraDockerOpts               string `json:",omitempty"`
	//
	//E2EStorageTestEnvironment     string `json:",omitempty"`
	//KubeletTestArgs               string `json:",omitempty"`
	//KubeletTestLogLevel           string `json:",omitempty"`
	//DockerTestArgs                string `json:",omitempty"`
	//DockerTestLogLevel            string `json:",omitempty"`
	//ApiserverTestArgs             string `json:",omitempty"`
	//ApiserverTestLogLevel         string `json:",omitempty"`
	//ControllerManagerTestArgs     string `json:",omitempty"`
	//ControllerManagerTestLogLevel string `json:",omitempty"`
	//SchedulerTestArgs             string `json:",omitempty"`
	//SchedulerTestLogLevel         string `json:",omitempty"`
	//KubeProxyTestArgs             string `json:",omitempty"`
	//KubeProxyTestLogLevel         string `json:",omitempty"`

	//NodeUp                        *NodeUpConfig `json:",omitempty"`

	// nodeSets is a list of all the NodeSets in the cluster.
	// It is not exported: we populate it from other files
	//nodeSets                      []*NodeSetConfig `json:",omitempty"`

	//// Masters is the configuration for each master in the cluster
	//Masters []*MasterConfig `json:",omitempty"`

	// EtcdClusters stores the configuration for each cluster
	EtcdClusters []*EtcdClusterSpec `json:"etcdClusters,omitempty"`

	// Component configurations
	Docker                *DockerConfig                `json:"docker,omitempty"`
	KubeDNS               *KubeDNSConfig               `json:"kubeDNS,omitempty"`
	KubeAPIServer         *KubeAPIServerConfig         `json:"kubeAPIServer,omitempty"`
	KubeControllerManager *KubeControllerManagerConfig `json:"kubeControllerManager,omitempty"`
	KubeScheduler         *KubeSchedulerConfig         `json:"kubeScheduler,omitempty"`
	KubeProxy             *KubeProxyConfig             `json:"kubeProxy,omitempty"`
	Kubelet               *KubeletConfig               `json:"kubelet,omitempty"`
	MasterKubelet         *KubeletConfig               `json:"masterKubelet,omitempty"`
}

type KubeDNSConfig struct {
	Replicas int    `json:"replicas,omitempty"`
	Domain   string `json:"domain,omitempty"`
	ServerIP string `json:"serverIP,omitempty"`
}

//
//type MasterConfig struct {
//	Name string `json:",omitempty"`
//
//	Image       string `json:",omitempty"`
//	Zone        string `json:",omitempty"`
//	MachineType string `json:",omitempty"`
//}
//

type EtcdClusterSpec struct {
	// Name is the name of the etcd cluster (main, events etc)
	Name string `json:"name,omitempty"`

	// EtcdMember stores the configurations for each member of the cluster (including the data volume)
	Members []*EtcdMemberSpec `json:"etcdMembers,omitempty"`
}

type EtcdMemberSpec struct {
	// Name is the name of the member within the etcd cluster
	Name string `json:"name,omitempty"`
	Zone string `json:"zone,omitempty"`

	VolumeType string `json:"volumeType,omitempty"`
	VolumeSize int    `json:"volumeSize,omitempty"`
}

type ClusterZoneSpec struct {
	Name string `json:"name,omitempty"`
	CIDR string `json:"cidr,omitempty"`
}

//type NodeUpConfig struct {
//	Source     string `json:",omitempty"`
//	SourceHash string `json:",omitempty"`
//
//	Tags       []string `json:",omitempty"`
//
//	// Assets that NodeUp should use.  This is a "search-path" for resolving dependencies.
//	Assets     []string `json:",omitempty"`
//}

// PerformAssignments populates values that are required and immutable
// For example, it assigns stable Keys to NodeSets & Masters, and
// it assigns CIDRs to subnets
func (c *Cluster) PerformAssignments() error {
	if c.Spec.NetworkCIDR == "" && !c.SharedVPC() {
		// TODO: Choose non-overlapping networking CIDRs for VPCs?
		c.Spec.NetworkCIDR = "172.20.0.0/16"
	}

	if c.Spec.NonMasqueradeCIDR == "" {
		c.Spec.NonMasqueradeCIDR = "100.64.0.0/10"
	}

	for _, zone := range c.Spec.Zones {
		err := zone.performAssignments(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (z *ClusterZoneSpec) performAssignments(c *Cluster) error {
	if z.CIDR == "" {
		cidr, err := z.assignCIDR(c)
		if err != nil {
			return err
		}
		glog.Infof("Assigned CIDR %s to zone %s", cidr, z.Name)
		z.CIDR = cidr
	}

	return nil
}

func (z *ClusterZoneSpec) assignCIDR(c *Cluster) (string, error) {
	// TODO: We probably could query for the existing subnets & allocate appropriately
	// for now we'll require users to set CIDRs themselves

	lastCharMap := make(map[byte]bool)
	for _, nodeZone := range c.Spec.Zones {
		lastChar := nodeZone.Name[len(nodeZone.Name)-1]
		lastCharMap[lastChar] = true
	}

	index := -1

	if len(lastCharMap) == len(c.Spec.Zones) {
		// Last char of zones are unique (GCE, AWS)
		// At least on AWS, we also want 'a' to be 1, so that we don't collide with the lowest range,
		// because kube-up uses that range
		index = int(z.Name[len(z.Name)-1])
	} else {
		glog.Warningf("Last char of zone names not unique")

		for i, nodeZone := range c.Spec.Zones {
			if nodeZone.Name == z.Name {
				index = i
				break
			}
		}
		if index == -1 {
			return "", fmt.Errorf("zone not configured: %q", z.Name)
		}
	}

	_, cidr, err := net.ParseCIDR(c.Spec.NetworkCIDR)
	if err != nil {
		return "", fmt.Errorf("Invalid NetworkCIDR: %q", c.Spec.NetworkCIDR)
	}
	networkLength, _ := cidr.Mask.Size()

	// We assume a maximum of 8 subnets per network
	// TODO: Does this make sense on GCE?
	// TODO: Should we limit this to say 1000 IPs per subnet? (any reason to?)
	index = index % 8
	networkLength += 3

	ip4 := cidr.IP.To4()
	if ip4 != nil {
		n := binary.BigEndian.Uint32(ip4)
		n += uint32(index) << uint(32-networkLength)
		subnetIP := make(net.IP, len(ip4))
		binary.BigEndian.PutUint32(subnetIP, n)
		subnetCIDR := subnetIP.String() + "/" + strconv.Itoa(networkLength)
		glog.V(2).Infof("Computed CIDR for subnet in zone %q as %q", z.Name, subnetCIDR)
		return subnetCIDR, nil
	}

	return "", fmt.Errorf("Unexpected IP address type for NetworkCIDR: %s", c.Spec.NetworkCIDR)
}

// SharedVPC is a simple helper function which makes the templates for a shared VPC clearer
func (c *Cluster) SharedVPC() bool {
	return c.Spec.NetworkID != ""
}

// CloudPermissions holds IAM-style permissions
type CloudPermissions struct {
	Permissions []*CloudPermission `json:"permissions,omitempty"`
}

// CloudPermission holds a single IAM-style permission
type CloudPermission struct {
	Resource string `json:"resource,omitempty"`
}

// AddS3Bucket adds a bucket if it does not already exist
func (p *CloudPermissions) AddS3Bucket(bucket string) {
	for _, p := range p.Permissions {
		if p.Resource == "s3://"+bucket {
			return
		}
	}

	p.Permissions = append(p.Permissions, &CloudPermission{
		Resource: "s3://" + bucket,
	})
}

// S3Buckets returns each of the S3 buckets in the permission
// TODO: Replace with something generic (probably we should just generate the permission)
func (p *CloudPermissions) S3Buckets() []string {
	var buckets []string
	for _, p := range p.Permissions {
		if strings.HasPrefix(p.Resource, "s3://") {
			buckets = append(buckets, strings.TrimPrefix(p.Resource, "s3://"))
		}
	}

	return buckets
}

//
//// findImage finds the default image
//func (c*NodeSetConfig) resolveImage() error {
//	cloud.(*awsup.AWSCloud).ResolveImage()
//
//	if n.Image == "" {
//		if defaultImage == "" {
//			image, err := c.determineImage()
//			if err != nil {
//				return err
//			}
//			defaultImage = image
//		}
//		n.Image = defaultImage
//	}
//
//
//	return nil
//}
