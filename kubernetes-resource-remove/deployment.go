package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// 사용하는 kubeconfig 파일을 기반으로 config를 생성합니다.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 클라이언트셋 생성
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// deployment 클라이언트 생성
	deploymentClient := clientset.AppsV1().Deployments("")

	// 모든 네임스페이스에서 deployment 목록 가져오기
	deploymentList, err := deploymentClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// deployment 목록 출력
	fmt.Printf("There are %d Deployments in the cluster\n", len(deploymentList.Items))
	for i, deployment := range deploymentList.Items {
		fmt.Printf("%d: Deployment Name: %s, Namespace: %s, Replicas: %d\n", i+1, deployment.Name, deployment.Namespace, *deployment.Spec.Replicas)
	}
}
