# Bunk8s

A tool for broad integration testing of microservices in Kubernetes

# Development

## Utilized versions

- Go: 1.16.5
- Kubernetes: 1.21
- Minikube: 1.21.0
- Helm: 3.6.1

## Annotations

- If you open the project in VSCode with the root directory of the project as root directory of the workspace in VSCode **AND** use gopls as language server, you have to add the following entry to the ```settings.json``` of the go extension. Otherwise gopls will report an error with the go modules.

    ```json
    "gopls": {
        "experimentalWorkspaceModule": true,
    }
    ```