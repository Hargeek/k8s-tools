package client

import (
	"errors"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
)

var Helm helm

type helm struct {
	ClientMap   map[string]*action.Configuration // 多集群client
	HelmConfMap map[string]string                // 集群config
}

// GetClient 根据集群名获取client
func (h *helm) GetClient(cluster string) (*action.Configuration, error) {
	client, ok := h.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群:%s 不存在, 无法获取client\n", cluster))
	}
	return client, nil
}

// Init 初始化client
func (h *helm) Init() {

}
