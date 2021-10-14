package model

type ServerReply struct {
	TestRunnerPods []TestRunnerPodReply
}

type TestRunnerPodReply struct {
	ReturnCode                  int
	PodName                     string
	Namespace                   string
	TestRunnerSidecarContainers []SidecarContainerReply
}

type SidecarContainerReply struct {
	SidecarContainerName string
}
