package model

type ServerReply struct {
	TestRunnerPods []TestRunnerPodReply `json:"testRunnerPods,omitempty"`
}

type TestRunnerPodReply struct {
	ReturnCode                  int                     `json:"returnCode,omitempty"`
	PodName                     string                  `json:"podName,omitempty"`
	Namespace                   string                  `json:"namespace,omitempty"`
	TestRunnerSidecarContainers []SidecarContainerReply `json:"testRunnerSidecarContainers,omitempty"`
}

type SidecarContainerReply struct {
	SidecarContainerName string `json:"sidecarContainerName,omitempty"`
}
