# Content for configuration

The configuration consists out of the following fields:

```yaml
launcherConfig:
  coordinatorIp: 
  coordinatorPort:
coordinatorConfig:
  testRunnerPods:
    - podName:
      namespace: 
      testTimeout: 
      containers:
        - containerName:
          image:
          startupCommands:
          startupCommandsArgs:
          testResultPath:
```

The following is an example configuration, in which two test runners are deployed:

```yaml  
launcherConfig:
  coordinatorIp: 192.168.49.2
  coordinatorPort: 30000
coordinatorConfig:
  testRunnerPods:
    - podName: bunk8s-test-runner-one
      namespace: default 
      testTimeout: 40 
      containers:
        - containerName: test-runner-container-one
          image: busybox
          startupCommands: ['sleep','10']
          startupCommandsArgs:
          testResultPath: /testresults-one
    - podName: bunk8s-test-runner-two
      namespace: default
      testTimeout: 40
      containers:
        - containerName: test-runner-container-one
          image: busybox
          startupCommands: ['sleep','10']
          startupCommandsArgs: 
          testResultPath: /testresults-two
```
