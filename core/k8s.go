package core

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
)

/*
*
操作集群资源的工具
*/
type Kubectl struct {
	client *dynamic.DynamicClient
}

/*
*
新建一个kubectl工具
*/
func NewKubectl(configStr string) *Kubectl {

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(configStr))

	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return &Kubectl{client: client}

}

func (this *Kubectl) Apply(resourceYaml string) error {
	/**
	1. 查看是否有该资源
	2. 如果没有则执行创建
	3. 如果有则执行删除
	*/

	resourceStruct, deploymentRes, err := this.parseResourceStruct(resourceYaml)
	if err != nil {
		return err
	}

	getResource, err := this.client.Resource(deploymentRes).Namespace(resourceStruct.GetNamespace()).Get(context.TODO(), resourceStruct.GetName(), metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 资源存在
	if getResource != nil {
		return this.Update(resourceYaml)
	}

	// 资源不存在
	return this.Create(resourceYaml)

}

func (this *Kubectl) Create(resourceYaml string) error {
	resourceStruct, deploymentRes, err := this.parseResourceStruct(resourceYaml)
	if err != nil {
		return err
	}

	_, err = this.client.Resource(deploymentRes).Namespace(resourceStruct.GetNamespace()).Create(context.TODO(), resourceStruct, metav1.CreateOptions{})

	return err
}

func (this *Kubectl) delete(resourceYaml string) error {
	resourceStruct, deploymentRes, err := this.parseResourceStruct(resourceYaml)
	if err != nil {
		return err
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err = this.client.Resource(deploymentRes).Namespace(resourceStruct.GetNamespace()).Delete(context.TODO(), resourceStruct.GetName(), deleteOptions)

	return err

}

func (this *Kubectl) Update(resourceYaml string) error {

	resourceStruct, deploymentRes, err := this.parseResourceStruct(resourceYaml)
	if err != nil {
		return err
	}

	_, err = this.client.Resource(deploymentRes).Namespace(resourceStruct.GetNamespace()).Update(context.TODO(), resourceStruct, metav1.UpdateOptions{})

	return err
}

func (this *Kubectl) parseResourceStruct(resourceYaml string) (*unstructured.Unstructured, schema.GroupVersionResource, error) {
	resourceStruct := &unstructured.Unstructured{}
	err := yaml.Unmarshal([]byte(resourceYaml), resourceStruct)
	if err != nil {
		return nil, schema.GroupVersionResource{}, err
	}

	apiVersion := strings.Split(resourceStruct.GetAPIVersion(), "/")

	deploymentRes := schema.GroupVersionResource{Group: apiVersion[0], Version: apiVersion[1], Resource: resourceStruct.GetKind()}

	return resourceStruct, deploymentRes, err
}
