package clusters

import (
	"strconv"

	"../../../../../tools-go/config"
	"../../../../../tools-go/errors"
	"../../../../../tools-go/logging"
	"../../../../../tools-go/runtime"
	"../../../../../tools-go/terminal"
	"../../../../../tools-go/timing"
	"../../../../../tools-go/users"
	"../../../../../tools-go/vars"
	"../clients"
	gcloudconfig "../config"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"k8s.io/client-go/rest"
)

type F struct{}

func (f F) Get(logResult bool) ([]*containerpb.Cluster, error) {
	var err error
	var resp *containerpb.ListClustersResponse
	var success bool
	for ok := true; ok; ok = !success {
		ctx, client := clients.F.GetContainer(clients.F{})
		projectID := config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", "")
		zone := config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", "")
		var request = &containerpb.ListClustersRequest{
			ProjectId: projectID,
			Zone:      zone,
		}
		resp, err = client.ListClusters(ctx, request)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to get clusters!", false, false, true) && resp != nil {
			success = true
		} else {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "Trying to recover..", 0)
			logging.Log([]string{"", vars.EmojiWarning, vars.EmojiInfo}, "ProjectID: "+projectID+" / Zone: "+zone+"\n", 0)
			timing.TimeOut(1, "s")
		}
	}
	logging.ProgressSpinner("stop")
	var clusters []*containerpb.Cluster
	if resp != nil {
		clusters = resp.Clusters
	}
	if logResult {
		logging.Log([]string{"", vars.EmojiKubernetes, ""}, "Clusters in your setup:\n", 0)
		if clusters != nil {
			for iC, cluster := range clusters {
				logging.Log([]string{"", "", ""}, "["+strconv.Itoa(iC)+"] "+cluster.Name, 0)
			}
		} else {
			logging.Log([]string{"", vars.EmojiKubernetes, ""}, "There are no clusters deployed.\n", 0)
		}
	}
	return clusters, err
}

func (f F) Create() bool {
	var success bool = false
	clusterID := config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")
	exists, _ := f.Exists(clusterID)
	if !exists {
		ctx, client := clients.F.GetContainer(clients.F{})
		request := &containerpb.CreateClusterRequest{
			ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			Cluster:   gcloudconfig.F.ClusterConfig(gcloudconfig.F{}, config.Setting("get", "infrastructure", "Spec.Gcloud.MachineType", ""), ""),
		}
		resp, err := client.CreateCluster(ctx, request)
		if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to start cluster creation!", false, true, true) {
			if resp.Status == 1 || resp.Status == 2 {
				success = true
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiProcess}, "Started creation of new cluster "+config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")+"..", 0)
			}
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiSuccess}, "The cluster "+clusterID+" is already created.", 0)
	}
	logging.ProgressSpinner("stop")
	return success
}

func (f F) Delete() {
	clusterID := config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")
	exists, _ := f.Exists(clusterID)
	if exists {
		ctx, client := clients.F.GetContainer(clients.F{})
		request := &containerpb.DeleteClusterRequest{
			ProjectId: config.Setting("get", "infrastructure", "Spec.Gcloud.ProjectID", ""),
			Zone:      config.Setting("get", "infrastructure", "Spec.Gcloud.Zone", ""),
			ClusterId: config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""),
			Name:      config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""),
		}
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiWarning}, "Deleting cluster "+clusterID+"!", 0)
		_, err := client.DeleteCluster(ctx, request)
		errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "Unable to delete cluster "+clusterID+"!", false, false, true)
		logging.ProgressSpinner("start")
		var success bool = false
		for ok := true; ok; ok = !success {
			if f.getStatus(clusterID, true) == 4 {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for the cluster "+clusterID+" to be deleted..", 0)
				timing.TimeOut(1, "s")
			} else {
				logging.ProgressSpinner("stop")
				logging.Log([]string{"", vars.EmojiInfra, vars.EmojiSuccess}, "Cluster deleted.\n", 0)
				success = true
			}
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiSuccess}, "The cluster "+clusterID+" is already deleted.", 0)
	}
}

func (f F) Exists(clusterID string) (bool, error) {
	var exists bool = false
	clusters, err := f.Get(false)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true) {
		for _, cluster := range clusters {
			if cluster.Name == clusterID {
				exists = true
				break
			}
		}
	}
	return exists, err
}

func (f F) getStatus(id string, ignoreDoesNotExist bool) containerpb.Cluster_Status {
	cluster := f.getCluster(id, ignoreDoesNotExist)
	var status containerpb.Cluster_Status = 0
	if cluster != nil {
		status = cluster.Status
	}
	return status
}

func (f F) getCluster(id string, ignoreDoesNotExist bool) *containerpb.Cluster {
	var clusterSelection *containerpb.Cluster
	clusters, err := f.Get(false)
	if !errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true) {
		for _, cluster := range clusters {
			if cluster.Name == id {
				clusterSelection = cluster
				break
			}
		}
	} else if !ignoreDoesNotExist {
		f.notExisting()
	}
	return clusterSelection
}

func (f F) notExisting() {
	clusterID := config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")
	if !users.GetClusterRecentlyDeleted() {
		logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWarning}, "The pre-configured cluster "+clusterID+" does not exist yet.", 0)
		action := terminal.GetUserSelection("Do you want to create the cluster?", []string{}, false, true)
		if action == "Yes" {
			provider := vars.InfraProviderActive
			if provider == vars.InfraProviderGcloud {
				f.Create()
			}
		} else {
			terminal.Exit(0, "")
		}
	} else {
		logging.Log([]string{"\n", vars.EmojiKubernetes, vars.EmojiSuccess}, "The specified cluster "+clusterID+" was successfully deleted.", 0)
		terminal.Exit(0, "")
	}
}

var loggedConnected bool = false

func (f F) GetAuthConfig() *rest.Config {
	var success bool = false
	var cluster *containerpb.Cluster
	logging.ProgressSpinner("start")
	for ok := true; ok; ok = !success {
		cluster = f.getCluster(config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""), false)
		if f.getStatus(config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""), false) == 0 {
			f.notExisting()
		} else if f.getStatus(config.Setting("get", "infrastructure", "Spec.Cluster.Name", ""), false) != 2 {
			logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiWaiting}, "Waiting for cluster "+config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")+" to be ready..", 0)
			timing.TimeOut(1, "s")
		} else {
			success = true
			logging.ProgressSpinner("stop")
			if !loggedConnected {
				logging.Log([]string{"", vars.EmojiKubernetes, vars.EmojiSuccess}, "Connected to cluster "+config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")+".", 0)
				loggedConnected = true
			}
		}
	}
	masterAuth := cluster.MasterAuth
	restconfig := &rest.Config{}
	//restconfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	restconfig.Username = masterAuth.Username
	restconfig.Password = masterAuth.Password
	restconfig.Insecure = true // TODO: Remove workaround
	restconfig.Host = cluster.Endpoint
	restconfig.ServerName = config.Setting("get", "infrastructure", "Spec.Cluster.Name", "")
	//restconfig.TLSClientConfig.CAData = []byte(masterAuth.ClusterCaCertificate)
	//restconfig.TLSClientConfig.CertData = []byte(masterAuth.ClientCertificate)
	//restconfig.KeyData = []byte(masterAuth.ClientKey)
	//errors.Check(err, runtime.F.GetCallerInfo(runtime.F{}, false), "", false, false, true)
	return restconfig
}
