package model

type Bunk8sConfig struct {
	LauncherConfig    LauncherConfig    `yaml:"launcherConfig"`
	CoordinatorConfig CoordinatorConfig `yaml:"coordinatorConfig"`
}
