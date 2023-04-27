package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	"k8s-tools/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset // 多集群client
	KubeConfMap map[string]string                // 集群config
}

// GetClient 根据集群名获取client
func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群:%s 不存在, 无法获取client\n", cluster))
	}
	return client, nil
}

// Init 初始化client
func (k *k8s) Init() {
	mp := map[string]string{}
	k.ClientMap = map[string]*kubernetes.Clientset{}
	if err := json.Unmarshal([]byte(config.KubeConfigRelativePath), &mp); err != nil {
		panic(fmt.Sprintf("Kubeconfig配置初始化反序列化失败 %v\n", err))
	}
	k.KubeConfMap = mp
	for clusterName, kubeConfigFilePath := range mp {
		kubeConfigFileData, err := config.GetKubeEmbed().ReadFile(kubeConfigFilePath)
		if err != nil {
			panic(fmt.Sprintf("集群 %s: 读取 kubeconfig 文件失败 %v\n", clusterName, err))
		}
		//conf, err := clientcmd.BuildConfigFromFlags("", kubeConfigFilePath)
		conf, err := clientcmd.RESTConfigFromKubeConfig(kubeConfigFileData)
		if err != nil {
			panic(fmt.Sprintf("集群 %s: 获取 K8s 配置失败 %v\n", clusterName, err))
		}
		clientSet, err := kubernetes.NewForConfig(conf)
		if err != nil {
			panic(fmt.Sprintf("集群 %s: 初始化 K8s Client 失败 %v\n", clusterName, err))
		}
		k.ClientMap[clusterName] = clientSet
		logger.Info(fmt.Sprintf("集群 %s: 初始化 K8s Client 成功 ", clusterName))
	}
}
