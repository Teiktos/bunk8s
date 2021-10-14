package main

import (
	"bunk8s/coordinator/server"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	logLevel, _ := zerolog.ParseLevel("info")

	zerolog.SetGlobalLevel(logLevel)

	kubeconfig := os.Getenv("KUBECONFIG")

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build kubeConfig")
	}
	log.Debug().Msg("initialized kubernetes config for clientset")

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to build clientset")
	}
	log.Debug().Msg("initialized kubernetes clientset")

	server.RpcServer(clientset)
	log.Debug().Msg("initialized kubernetes clientset")

}
