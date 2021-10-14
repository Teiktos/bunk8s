package model

type LauncherConfig struct {
	CoordinatorIp   string `yaml:"coordinatorIp"`
	CoordinatorPort string `yaml:"coordinatorPort"`
	CertFile        string `yaml:"certFile"`
}
