package model

// type CoordinatorConfig struct {
// 	TestTimeout int        `yaml:"testTimeout"`
// 	Namespace   string     `yaml:"namespace"`
// 	TestRunner  TestRunner `yaml:"testRunner"`
// }

// type TestRunner struct {
// 	PodName        string `yaml:"podName"`
// 	Image          string `yaml:"image"`
// 	StartupCommand string `yaml:"startupCommand"`
// 	TestResultPath string `yaml:"testResultPath"`
// }

type CoordinatorConfig struct {
	TestRunnerPods []*TestRunnerPod `yaml:"testRunnerPods" json:"testRunnerPods,omitempty"`
}

type TestRunnerPod struct {
	PodName     string       `yaml:"podName" json:"podName,omitempty"`
	Namespace   string       `yaml:"namespace" json:"namespace,omitempty"`
	TestTimeout int          `yaml:"testTimeout" json:"testTimeout,omitempty"`
	Containers  []Containers `yaml:"containers" json:"containers,omitempty"`
}

type Containers struct {
	ContainerName       string   `yaml:"containerName" json:"containerName,omitempty"`
	Image               string   `yaml:"image" json:"image,omitempty"`
	StartupCommands     []string `yaml:"startupCommands" json:"startupCommands,omitempty"`
	StartupCommandsArgs []string `yaml:"startupCommandsArgs" json:"startupCommandsArgs,omitempty"`
	TestResultPath      string   `yaml:"testResultPath" json:"testResultPath,omitempty"`
}


