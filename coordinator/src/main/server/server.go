package server

import (
	"bunk8s/coordinator/kubeclient"
	"bunk8s/coordinator/model"
	"bunk8s/coordinator/pb"
	"context"
	"encoding/json"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

type GrpcServerStruct struct {
	pb.UnimplementedTestRunServer
}

var clientset *kubernetes.Clientset

func (s *GrpcServerStruct) DeployTestRunner(ctx context.Context, coordinatorConfigProto *pb.CoordinatorConfig) (*pb.ServerReply, error) {

	coordinatorConfig := coordinatorConfigProtoToStruct(coordinatorConfigProto)

	serverReply := buildServerReply(coordinatorConfig)

	log.Debug().Msg("Created serverReply struct")

	for i := range coordinatorConfig.TestRunnerPods {
		if !kubeclient.NamespaceIsValid(ctx, clientset, coordinatorConfig.TestRunnerPods[i]) {
			serverReply.TestRunnerPods[i].ReturnCode = 1
			log.Debug().Msg("Couldn't get namespaces!")
			serverReplyProto := serverReplyStructToProto(serverReply)
			return &serverReplyProto, nil
		}
	}

	log.Debug().Msg("Namespaces are valid!")

	for i := range coordinatorConfig.TestRunnerPods {
		if !kubeclient.PodNameIsValid(ctx, clientset, coordinatorConfig.TestRunnerPods[i]) {
			serverReply.TestRunnerPods[i].ReturnCode = 2
			serverReplyProto := serverReplyStructToProto(serverReply)
			return &serverReplyProto, nil
		}
	}

	log.Debug().Msg("Pod names are valid!")

	createdPods, errors := kubeclient.CreateTestRunnerPods(ctx, clientset, coordinatorConfig)
	for i := range coordinatorConfig.TestRunnerPods {
		if errors[i] != nil {
			log.Error().Err(errors[i]).Msg("Failed to create test runner pod")
			serverReply.TestRunnerPods[i].ReturnCode = 3
			serverReplyProto := serverReplyStructToProto(serverReply)
			return &serverReplyProto, errors[i]
		}

		log.Info().Msg("created Pod " + createdPods[i].Name + " in namespace " + createdPods[i].Namespace)
	}

	watchResults := kubeclient.WatchTestRunnerPods(ctx, clientset, coordinatorConfig)
	for i := range watchResults.WatcherErr {
		if watchResults.WatcherErr[i] != nil {
			log.Error().Err(watchResults.WatcherErr[i]).Msg("Failed to watch test runner pod")
			serverReply.TestRunnerPods[i].ReturnCode = 4
			serverReplyProto := serverReplyStructToProto(serverReply)
			return &serverReplyProto, watchResults.WatcherErr[i]
		}
	}
	for i := range watchResults.WatcherSuccessful {
		if !watchResults.WatcherSuccessful[i] {
			log.Error().Msg("Test run timed out")
			log.Error().Msgf("%v", watchResults.WatcherSuccessful)
			serverReply.TestRunnerPods[i].ReturnCode = 5
			serverReplyProto := serverReplyStructToProto(serverReply)
			return &serverReplyProto, nil
		}
	}

	serverReplyProto := serverReplyStructToProto(serverReply)

	return &serverReplyProto, nil
}

func coordinatorConfigProtoToStruct(coordinatorConfigProto *pb.CoordinatorConfig) model.CoordinatorConfig {

	coordinatorConfigJson, err := json.Marshal(coordinatorConfigProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from config proto to json")
	}

	coordinatorConfig := model.CoordinatorConfig{}

	err = json.Unmarshal(coordinatorConfigJson, &coordinatorConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from json to config struct")
	}

	log.Debug().Msgf("Marshalled coordinator config to: %+v", coordinatorConfig)

	return coordinatorConfig
}

func serverReplyStructToProto(serverReply model.ServerReply) pb.ServerReply {

	serverReplyJson, err := json.Marshal(&serverReply)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from server reply struct to json")
	}
	log.Debug().Msg("Marshalled from server reply struct to json")

	serverReplyProto := pb.ServerReply{}

	err = json.Unmarshal(serverReplyJson, &serverReplyProto)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshall from json to server reply proto")
	}
	log.Debug().Msg("Marshalled from json to server reply proto")

	return serverReplyProto

}

func buildServerReply(coordinatorConfig model.CoordinatorConfig) model.ServerReply {

	serverReply := model.ServerReply{}

	var pods []model.TestRunnerPodReply

	for i := range coordinatorConfig.TestRunnerPods {

		pod := model.TestRunnerPodReply{
			PodName:                     coordinatorConfig.TestRunnerPods[i].PodName,
			Namespace:                   coordinatorConfig.TestRunnerPods[i].Namespace,
			TestRunnerSidecarContainers: buildServerReplySidecarContainers(coordinatorConfig.TestRunnerPods[i]),
		}

		pods = append(pods, pod)
	}

	serverReply.TestRunnerPods = pods

	return serverReply

}

func buildServerReplySidecarContainers(pod model.TestRunnerPod) []model.SidecarContainerReply {

	var containers []model.SidecarContainerReply

	container := model.SidecarContainerReply{
		SidecarContainerName: pod.PodName + "-sidecar-container",
	}

	containers = append(containers, container)

	return containers

}

func RpcServer(cs *kubernetes.Clientset) {

	clientset = cs

	lis, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on port 80")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTestRunServer(grpcServer, &GrpcServerStruct{})

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to serve gRPC server on port 80")
	}
}
