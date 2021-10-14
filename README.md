# Bunk8s

A tool for service integration testing of microservices in Kubernetes

# Development

## Utilized versions

- Go: 1.16.5
- Minikube: 1.21.0
- Helm: 3.6.1

## Next steps in development

### General

- [x] create main.go files and expand the project structure
- [x] create minimal Dockerfiles
- [x] create makefile(s)

### Coordinator

- [x] create Helm chart for the deployment of the "coordinator"
- [x] add RPC server for communication with the launcher
- [x] add functionality to deploy a test runner pod as remote procedure
- [x] add functionality to deploy a sidecar container and a shared volume with the test runner container inside the test runner pod
- [x] define test runner container as init ~~or ephemeral~~ container in order to avoid CrashLoopBackOff from Kubernetes
- ~~[ ] add server for communication with a PreStop Container Hook to get informed when the test runner finishes~~
- ~~[ ] add PreStop Container Lifecycle hook to deplyoment of test runner container~~
- ~~[ ] add configMap, containing the data/files required by the PreStop hook~~
- [x] add watcher in order to get notified when the test runner pod finishes
- [x] add functionality to deploy multiple test runner pods with multiple containers per pod  

### Launcher

- [x] add RPC client
- [x] add functionality to send configuration data to the coordinator
- [x] externalize the configuration data

### Miscellaneous

- [x] add logging
- [x] improve error handleing
- [x] add timeout
- [ ] improve security

## Annotations

- the first version will be developed for local usage with minikube. It will later be adapted for a public cloud Kubernetes Cluster
- If you open the project in VSCode with the root directory of the project as root directory of the workspace in VSCode **AND** use gopls as language server, you have to add the following entry to the ```settings.json``` of the go extension. Otherwise gopls will report an error with the go modules.

    ```json
    "gopls": {
        "experimentalWorkspaceModule": true,
    }
    ```

- A testing process is launched via [launch.sh](launcher/launch.sh). If the application runs locally in minikube with c, the directory with the config file on the host must be mounted into the minikube vm with the following command, before launch.sh can be executed. Alternatively [launch.sh](launcher/launch.sh) can be run with sudo. This is obsolete, if minikube runs bare-metal with --vm-driver=​none

    ```bash
    minikube mount <path/to/config/on/host>:</path/to/config/in/minikube/vm>
    ```

    ```bash
    launch.sh </path/to/config/in/minikube/vm>
    ```
