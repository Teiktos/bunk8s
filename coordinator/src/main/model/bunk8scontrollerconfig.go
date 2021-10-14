package model

// type CoordinatorConfig struct {
// 	TestTimeout int
// 	Namespace   string
// 	TestRunner  TestRunner
// }

// type TestRunner struct {
// 	PodName        string
// 	Image          string
// 	StartupCommand string
// 	TestResultPath string
// }

type CoordinatorConfig struct {
	TestRunnerPods []TestRunnerPod `json:"testRunnerPods,omitempty"`
}

type TestRunnerPod struct {
	PodName     string       `json:"podName,omitempty"`
	Namespace   string       `json:"namespace,omitempty"`
	TestTimeout int          `json:"testTimeout,omitempty"`
	Containers  []Containers `json:"containers,omitempty"`
}

type Containers struct {
	ContainerName       string   `json:"containerName,omitempty"`
	Image               string   `json:"image,omitempty"`
	StartupCommands     []string `json:"startupCommands,omitempty"`
	StartupCommandsArgs []string `json:"startupCommandsArgs,omitempty"`
	TestResultPath      string   `json:"testResultPath,omitempty"`
}

func (c CoordinatorConfig) GetNamespaces() []string {
	var ns []string
	for i := range c.TestRunnerPods {
		ns = append(ns, c.TestRunnerPods[i].Namespace)
	}
	return ns
}

func (c CoordinatorConfig) GetPodNames() []string {
	var ps []string
	for i := range c.TestRunnerPods {
		ps = append(ps, c.TestRunnerPods[i].PodName)
	}
	return ps
}
