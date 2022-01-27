# Bunk8s

A tool for broad integration testing of microservices in Kubernetes.
The following gives an overview of the different components of Bunk8s 

## The Configuration file

The configuration file must be written in YAML and is structured similarly to Kubernetes ressource definition files. It is possible to define multiple test runner pods, which can be deployed in different namespaces, as well as define multiple test runner containers per pod.

Further information can be found in the [README for configuration files](configuration/README.md)

## The Launcher

The Launcher runs outside of Kubernetes as part of the continuous integration process in a pipeline. After the container images of the services under test are built and pushed to a container registry by the pipeline, the Bunk8s launcher is started, which starts the testing process. It loads a configuration file, which provides the information that is required by the launcher in order to connect to the coordinator.

The following gives an overview of the return codes that are sent from he coordinator to the launcher:

- 0 - Test run finished successfully
- 1 - Namespace name invalid
- 2 - Test runner pod name invalid or already exists in given namespace 
- 3 - Failed to create test runner pod
- 4 - Failed to watch test runner pod
- 5 - Test duration timed out

## The Coordinator

The Bunk8s Coordinator handles the deployment of the test runner containers and watches their state until they are finished. It consists of two parts. An RPC server, which allows the synchronization of the test extraction by the launcher, and a Kubernetes clientset, that connects to the kube-apiserver via HTTP and sends a request to create the test runner pods and watches the test runner pod’s state. The coordinator’s RPC server is exposed to the launcher via an Ingress. After the test runner pods finish, the coordinator sends a response to the Launcher, containing the data that is required to extract the test results from the test runner containers.

# Installation 

1. The launcher image, as well as the coordinator image must be built. The Docker image files are written in such a way, that the build contexts root directory must be the root directory of the Bunk8s repository. 
2. The launcher container image must be uploaded to a container registry from which pipeline runners can pull it so that it can be run in the pipeline. 
3. The coordinator container image must be uploaded to a container registry to which Kubernetes has access. 
4. Deploy the coordinator. It is possible to deploy it to an arbitrary namespace without affecting its functionality, however, it is recommended to deploy it into its own namespace. In order to allow for a quick setup of the coordinator, the Bunk8s repository includes a preconfigured helm chart for deployment to Kubernetes. When deployed with Helm, the coordinator is part of a deployment and therefore of a replica set. However, it is only required to set the number of replicas to one, since the gRPC server can handle multiple simultaneous test runs. 
5. A certificate for TLS must be provided to the ingress and the CAs root certificate must be placed in the `bunk8s/launcher/src/cert` directory, before building the container image. Particular attention is to be paid to the Ingress configuration in the helm chart. Depending on the used Ingress controller in the cluster the Ingresses annotations must be changed in order to enable the communication between the gRPC client and the gRPC server. 
```` yaml
nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
````
6. Create a role, that grants access to the core API, namespaces, and pod resources and assign it to the coordinator pod. For the namespace resource, the required verb is ["get"], while for the pod resource the required verbs are ["get", "list", "watch", "create"]. Additionally, a service account, as well as a role binding must be created and be assigned to the coordinator pod. All three parts of the role-based access control, the role, the service account, and the role binding must be created in the namespace of the coordinator.

# Showcase

[![Bunk8s Showcase](http://img.youtube.com/vi/e8wbS25O4Bo/0.jpg)](https://www.youtube.com/watch?v=e8wbS25O4Bo "Bunk8s Showcase")

# Project Setup

## Utilized versions

- Go: 1.16.5
- Kubernetes: 1.21
- Helm: 3.6.1

## VSCode

- If you open the project in VSCode with the root directory of the project as root directory of the workspace in VSCode **AND** use gopls as language server, you have to add the following entry to the ```settings.json``` of the go extension. Otherwise gopls will report an error with the multiple go modules in the repository.

    ```json
    "gopls": {
        "experimentalWorkspaceModule": true,
    }
    ```

# Reference 

<cite>Reile, C., Chadha, C., Hauner, V., Jindal, A., Hofmann, B., Gerndt, M. (2022).  Bunk8s: Enabling Easy Integration Testing of Microservices in Kubernetes. IEEE International Conference on Software Analysis, USA.