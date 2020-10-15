package main

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
)

const PatchRetain = `{"spec":{"persistentVolumeReclaimPolicy":"` + corev1.PersistentVolumeReclaimRetain + `"}}`

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.Lmsgprefix)

	var err error
	defer exit(&err)

	var cfg *rest.Config
	if cfg, err = rest.InClusterConfig(); err != nil {
		return
	}
	var client *kubernetes.Clientset
	if client, err = kubernetes.NewForConfig(cfg); err != nil {
		return
	}

	var pvs *corev1.PersistentVolumeList
	if pvs, err = client.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{}); err != nil {
		return
	}

	for _, pv := range pvs.Items {
		if pv.Spec.PersistentVolumeReclaimPolicy == corev1.PersistentVolumeReclaimRetain {
			log.Printf("PV %s: OK", pv.Name)
			continue
		}
		if _, err = client.CoreV1().PersistentVolumes().Patch(
			context.Background(),
			pv.Name,
			types.StrategicMergePatchType,
			[]byte(PatchRetain),
			metav1.PatchOptions{},
		); err != nil {
			log.Printf("PV %s: Failed(%s)", pv.Name, err.Error())
			return
		}
		log.Printf("PV %s: Patched", pv.Name)
	}
}
