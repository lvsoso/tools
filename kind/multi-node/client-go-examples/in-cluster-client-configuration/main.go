package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	fmt.Println("----------------- start ...")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("----------------- loop ...")
	count := 0
	for {
		fmt.Printf("----------------- %d ...\n", count)
		count += 1
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%d pods in the cluster\n", len(pods.Items))

		for i, pod := range pods.Items {
			log.Printf("%d -> %s/%s", i+1, pod.Namespace, pod.Name)
		}

		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("node found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}

		<-time.Tick(time.Second * 10)
	}
}
