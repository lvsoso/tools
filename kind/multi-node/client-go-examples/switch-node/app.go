package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
)

const (
	master  = "master"
	backup  = "backup"
	unknown = "unknown"

	nodeEnvKey = "NODE"
)

var (
	whereTargetSvc    string
	whereTargetClient string

	targetStsMasterAlive bool
	targetStsBackupAlive bool
	targetStsClientAlive bool

	namespace              = "default"
	targetSvcName          = "nginx"
	targetStsMasterName    = "nginx-master"
	targetStsPodMasterName = targetStsMasterName + "-0"
	targetStsBackupName    = "nginx-backup"
	targetStsPodBackupName = targetStsBackupName + "-0"
	targetStsClientName    = "worker"
	targetStsPodClientName = targetStsClientName + "-0"
	masterNodeLabel        = map[string]string{"app-master": "true"}
	backupNodeLabel        = map[string]string{"app-backup": "true"}

	masterNodePvc = ""
	backupNodePvc = ""

	masterSvcSelector = map[string]string{"app": "nginx-master"}
	backupSvcSelector = map[string]string{"app": "nginx-backup"}

	leaseMode          = "false"
	leaseLockName      = "test"
	leaseLockNamespace = namespace
	// leaselNodeId  = ""

)

func str2Map(data string) map[string]string {
	res := make(map[string]string)
	err := json.Unmarshal([]byte(data), &res)
	if err != nil {
		errStage("str2Map failed, ", err.Error())
	}
	return res
}

// init config
func init() {
	namespace = os.Getenv("NAMESPACE")
	targetSvcName = os.Getenv("TARGET_SVC_NAME")
	targetStsMasterName = os.Getenv("TARGET_STS_MASTER_NAME")
	targetStsBackupName = os.Getenv("TARGET_STS_BACKUP_NAME")
	targetStsClientName = os.Getenv("TARGET_STS_CLIENT_NAME")
	masterNodeLabelEnv := os.Getenv("MASTER_NODE_LABEL")
	backupNodeLabelEnv := os.Getenv("BACKUP_NODE_LABEL")
	masterSvcSelectorEnv := os.Getenv("MASTER_SVC_SELECTOR")
	backupSvcSelectorEnv := os.Getenv("BACKUP_SVC_SELECTOR")

	// pvc
	masterNodePvc = os.Getenv("MASTER_NODE_PVC")
	backupNodePvc = os.Getenv("BACKUP_NODE_PVC")

	leaseLockName = os.Getenv("LEASE_LOCK_NAME")
	leaseMode = os.Getenv("LEASE_MODE")

	// leaselNodeId = ""
	// leaselNodeId = os.Getenv("LEASEL_NODE_ID")

	masterNodeLabel = str2Map(masterNodeLabelEnv)
	backupNodeLabel = str2Map(backupNodeLabelEnv)
	masterSvcSelector = str2Map(masterSvcSelectorEnv)
	backupSvcSelector = str2Map(backupSvcSelectorEnv)

	targetStsPodMasterName = targetStsMasterName + "-0"
	targetStsPodBackupName = targetStsBackupName + "-0"
	targetStsPodClientName = targetStsClientName + "-0"
	leaseLockNamespace = namespace
}

func isPodAlive(pod apicorev1.Pod) bool {
	return pod.Status.Phase == apicorev1.PodRunning
}

func stage(msg ...interface{}) {
	klog.Infoln(msg...)
}

func errStage(msg ...interface{}) {
	klog.Errorln(msg...)
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func main() {
	// get args
	klog.InitFlags(nil)

	var id string
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "path to the kubeconfig")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "path to the kubeconfig")
	}

	flag.StringVar(&id, "id", uuid.New().String(), "the holder identity name")

	if leaseLockName == "" {
		klog.Fatal("unable to get lease lock resource name.")
	}
	if leaseLockNamespace == "" {
		klog.Fatal("unable to get lease lock resource namespace.")
	}

	klog.Infoln("my id: ", id)
	klog.Infoln("my kubeconfig: ", *kubeconfig)

	// build config
	fi, err := os.Stat(*kubeconfig)
	if err != nil || fi.IsDir() {
		*kubeconfig = ""
	}

	config, err := buildConfig(*kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}
	client := kubernetes.NewForConfigOrDie(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
	}()

	if strings.ToLower(leaseMode) == "false" {
		klog.Infoln("normal mode")
		run(ctx, client)
	}

	klog.Infoln("lease mode")
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				run(ctx, client)
			},
			OnStoppedLeading: func() {
				klog.Infof("leader lost: %s", id)
				cancel()
			},
			OnNewLeader: func(identity string) {
				if identity == id {
					return
				}
				klog.Infof("new leader elected: %s", identity)
			},
		},
	})
}

func run(ctx context.Context, clientset *kubernetes.Clientset) {
	stage("start loop")

	for {
		select {
		case <-ctx.Done():
			stage("canceld.")
			return
		case <-time.Tick(10 * time.Second):
			check(ctx, clientset)
		}
	}
}

func check(ctx context.Context, clientset *kubernetes.Clientset) {
	stage("new epoch ------------------------------------- ")
	// check service
	targetSvc, err := clientset.CoreV1().Services(namespace).Get(ctx, targetSvcName, metav1.GetOptions{})
	if err != nil {
		errStage("miss", targetStsClientName, err.Error())
		return
	}

	// check where service
	whereTargetSvc = unknown
	for k, v := range targetSvc.Spec.Selector {
		if masterSvcSelector[k] == v {
			whereTargetSvc = master
			break
		}
		if backupSvcSelector[k] == v {
			whereTargetSvc = backup
			break
		}
	}
	stage("running svc", targetSvcName, whereTargetSvc)

	// check master serve
	targetStsPodMaster, err := clientset.CoreV1().Pods(namespace).Get(ctx, targetStsPodMasterName, metav1.GetOptions{})
	if err != nil {
		errStage("miss", targetStsPodMasterName, err.Error())
		targetStsMasterAlive = false
	} else {
		targetStsMasterAlive = isPodAlive(*targetStsPodMaster)
	}
	stage("status", targetStsPodMasterName, targetStsMasterAlive)

	// check backup serve
	targetStsPodBackup, err := clientset.CoreV1().Pods(namespace).Get(ctx, targetStsPodBackupName, metav1.GetOptions{})
	if err != nil {
		errStage("miss", targetStsPodBackupName, err.Error())
		targetStsBackupAlive = false
	} else {
		targetStsBackupAlive = isPodAlive(*targetStsPodBackup)
	}
	stage("status", targetStsPodBackupName, targetStsBackupAlive)

	// check client
	targetStsPodClient, err := clientset.CoreV1().Pods(namespace).Get(ctx, targetStsPodClientName, metav1.GetOptions{})
	if err != nil {
		errStage("miss", targetStsPodClientName, err.Error())
		targetStsClientAlive = false
	} else {
		targetStsClientAlive = isPodAlive(*targetStsPodClient)
	}

	stage("status", targetStsPodClientName, targetStsClientAlive)

	// check where client
	whereTargetClient = unknown
	targetStsClient, err := clientset.AppsV1().StatefulSets(namespace).Get(ctx, targetStsClientName, metav1.GetOptions{})
	if err != nil {
		errStage("miss", targetStsClientName, err.Error())
	} else {
		for k, v := range targetStsClient.Spec.Template.Spec.NodeSelector {
			if masterNodeLabel[k] == v {
				whereTargetClient = master
				break
			}
			if backupNodeLabel[k] == v {
				whereTargetClient = backup
				break
			}
		}
	}
	stage("running", targetStsClientName, whereTargetClient)

	stsClient := clientset.AppsV1().StatefulSets(namespace)
	svcClient := clientset.CoreV1().Services(namespace)

	// choice node
	choiceNode := master
	if (whereTargetSvc == backup && targetStsBackupAlive) || (whereTargetClient == backup && targetStsBackupAlive) || (!targetStsMasterAlive && targetStsBackupAlive) {
		choiceNode = backup
	}

	// update
	switch choiceNode {
	case master:
		if whereTargetClient != master {
			err = updateStatefulSetsClient(ctx, stsClient, targetStsClientName, masterNodeLabel, master, masterNodePvc)
			if err != nil {
				errStage("update client failed", targetStsClientName, err.Error())
			}
			stage("update client ", targetStsClientName, master)
		}

		if whereTargetSvc != master {
			err = updateService(ctx, svcClient, targetSvcName, masterSvcSelector)
			if err != nil {
				errStage("update svc failed", targetSvcName, err.Error())
			}
			stage("update svc ", targetSvcName, master)
		}

	case backup:
		if whereTargetClient != backup {
			err = updateStatefulSetsClient(ctx, stsClient, targetStsClientName, backupNodeLabel, backup, backupNodePvc)
			if err != nil {
				errStage("update client failed", targetStsClientName, err.Error())
			}
			stage("update client ", targetStsClientName, backup)
		}

		if whereTargetSvc != backup {
			err = updateService(ctx, svcClient, targetSvcName, backupSvcSelector)
			if err != nil {
				errStage("update svc failed", targetSvcName, err.Error())
			}
			stage("update svc ", targetSvcName, backup)
		}
	}

}

func updateStatefulSetsClient(ctx context.Context, stsClient v1.StatefulSetInterface, targetName string, nodeSelector map[string]string, node string, pvc string) error {
	sts, err := stsClient.Get(ctx, targetName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	sts.Spec.Template.Spec.NodeSelector = nodeSelector
	env := sts.Spec.Template.Spec.Containers[0].Env

	foundIdx := -1
	for idx, e := range env {
		if e.Name == nodeEnvKey {
			foundIdx = idx
			break
		}
	}
	if foundIdx != -1 {
		env[foundIdx].Value = node
	} else {
		env = append(env, apicorev1.EnvVar{
			Name:  nodeEnvKey,
			Value: node,
		})
	}
	sts.Spec.Template.Spec.Containers[0].Env = env

	volumes := sts.Spec.Template.Spec.Volumes
	volumes[0].PersistentVolumeClaim.ClaimName = pvc
	sts.Spec.Template.Spec.Volumes = volumes

	return retry.RetryOnConflict(
		retry.DefaultRetry, func() error {
			_, err = stsClient.Update(ctx, sts, metav1.UpdateOptions{})
			return err
		},
	)
}

func updateService(ctx context.Context, svcClient corev1.ServiceInterface, targetName string, selector map[string]string) error {
	svc, err := svcClient.Get(ctx, targetName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	svc.Spec.Selector = selector
	return retry.RetryOnConflict(
		retry.DefaultRetry, func() error {
			_, err = svcClient.Update(ctx, svc, metav1.UpdateOptions{})
			return err
		},
	)
}
