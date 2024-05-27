package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func updateHPAs(kubeconfig string, targetNamespace string) {
	// 유효한 네임스페이스가 제공되었는지 확인합니다.
	if targetNamespace == "" {
		log.Fatal("No target namespace provided")
	}

	// 사용하는 kubeconfig 파일을 기반으로 config를 생성합니다.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 클라이언트셋 생성
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var hpaClient = clientset.AutoscalingV1().HorizontalPodAutoscalers(targetNamespace)
	var hpaList *v1.HorizontalPodAutoscalerList

	// if targetNamespace is "all", get all namespaces, else get the specified namespace
	if targetNamespace == "all" {
		hpaList, err = clientset.AutoscalingV1().HorizontalPodAutoscalers("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
	} else {
		hpaList, err = hpaClient.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
	}

	// HPA MAX Replicas 3개 이상인 HPA 변경
	for i, hpa := range hpaList.Items {
		if targetNamespace == "all" {
			hpaClient = clientset.AutoscalingV1().HorizontalPodAutoscalers(hpa.Namespace)
		}
		if hpa.Spec.MaxReplicas > 99 {
			originalReplicas := hpa.Spec.MaxReplicas
			hpa.Spec.MaxReplicas = 10
			_, err := hpaClient.Update(context.Background(), &hpa, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("%s 네임스페이스의 HPA %s 업데이트 실패: %s", hpa.Namespace, hpa.Name, err)
				continue
			}
			fmt.Printf("%d: HPA 이름: %s, 네임스페이스: %s, Max 복제 수 %d에서 %d로 조정됨\n", i+1, hpa.Name, hpa.Namespace, originalReplicas, hpa.Spec.MaxReplicas)
		}
	}

}

func main() {

	var kubeconfig *string
	var namespace *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	namespace = flag.String("namespace", "", "HPA를 조회할 네임스페이스")
	flag.Parse()

	updateHPAs(*kubeconfig, *namespace)
}
