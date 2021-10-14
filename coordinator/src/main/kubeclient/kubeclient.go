package kubeclient

import (
	"bunk8s/coordinator/model"
	"context"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NamespaceIsValid(ctx context.Context, clientset *kubernetes.Clientset, testRunnerPod model.TestRunnerPod) bool {

	_, err := clientset.CoreV1().Namespaces().Get(ctx, testRunnerPod.Namespace, metav1.GetOptions{})
	return err == nil
}

func PodNameIsValid(ctx context.Context, clientset *kubernetes.Clientset, testRunnerPod model.TestRunnerPod) bool {

	pod, err := clientset.CoreV1().Pods(testRunnerPod.Namespace).Get(ctx, testRunnerPod.PodName, metav1.GetOptions{})

	if err != nil && err.Error() != "pods \""+testRunnerPod.PodName+"\" not found" {
		return false
	}
	if pod.GetName() == testRunnerPod.PodName {
		return false
	}

	return true
}

func CreateTestRunnerPods(ctx context.Context, clientset *kubernetes.Clientset, coordinatorConfig model.CoordinatorConfig) ([]*v1.Pod, []error) {

	pods := make([]*v1.Pod, len(coordinatorConfig.TestRunnerPods))
	errors := make([]error, len(coordinatorConfig.TestRunnerPods))

	for i := range coordinatorConfig.TestRunnerPods {

		podObject := getPodObject(coordinatorConfig.TestRunnerPods[i])

		pod, err := clientset.CoreV1().Pods(coordinatorConfig.TestRunnerPods[i].Namespace).Create(ctx, podObject, metav1.CreateOptions{})

		pods[i] = pod
		errors[i] = err

	}

	return pods, errors
}

func WatchTestRunnerPods(ctx context.Context, clientset *kubernetes.Clientset, coordinatorConfig model.CoordinatorConfig) model.WatchResults {

	var wg sync.WaitGroup

	watcherBoolResults := make([]bool, len(coordinatorConfig.TestRunnerPods))
	watcherErrorResults := make([]error, len(coordinatorConfig.TestRunnerPods))

	for i := range coordinatorConfig.TestRunnerPods {
		wg.Add(1)
		go startPodWatcher(ctx, clientset, coordinatorConfig.TestRunnerPods[i], watcherBoolResults, watcherErrorResults, i, &wg)
	}
	wg.Wait()

	results := model.WatchResults{WatcherSuccessful: watcherBoolResults, WatcherErr: watcherErrorResults}

	return results
}

func startPodWatcher(ctx context.Context, clientset *kubernetes.Clientset, testRunnerPod model.TestRunnerPod, watcherBoolResults []bool, watcherErrorResults []error, i int, wg *sync.WaitGroup) {

	defer wg.Done()
	log.Debug().Msgf("Started Go routine %v", i)

	timeout := int64(testRunnerPod.TestTimeout)
	opts := metav1.SingleObject(metav1.ObjectMeta{Name: testRunnerPod.PodName})
	opts.TimeoutSeconds = &timeout

	watch, err := clientset.CoreV1().Pods(testRunnerPod.Namespace).Watch(ctx, opts)
	defer watch.Stop()
	log.Debug().Msgf("Created Watcher %v for go routine %v", watch, i)
	if err != nil {
		watcherBoolResults[i] = false
		watcherErrorResults[i] = err
		log.Err(err).Msgf("Error when creating watcher in go routine %v", i)
		return
	}

	for event := range watch.ResultChan() {
		pod, ok := event.Object.(*v1.Pod)
		if !ok {
			log.Debug().Msg("Pod from event channel not okay!")
			watcherBoolResults[i] = false
			watcherErrorResults[i] = nil
			return
		}

		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				watcherBoolResults[i] = true
				watcherErrorResults[i] = nil
				return
			}
		}
	}
	watcherBoolResults[i] = false
	watcherErrorResults[i] = nil
}

func getPodObject(testRunnerPod model.TestRunnerPod) *v1.Pod {

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testRunnerPod.PodName,
			Namespace: testRunnerPod.Namespace,
			Labels: map[string]string{
				"app": "demo",
			},
			Annotations: map[string]string{
				"linkerd.io/inject": "disabled",
			},
		},
		Spec: v1.PodSpec{
			Volumes: []v1.Volume{
				{
					Name:         "shared-data",
					VolumeSource: v1.VolumeSource{},
				},
			},

			Containers: []v1.Container{

				{
					Name:            testRunnerPod.PodName + "-sidecar-container",
					Image:           "busybox",
					ImagePullPolicy: v1.PullIfNotPresent,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "shared-data",
							MountPath: "/testresults",
						},
					},
					Command: []string{
						"sleep",
						"infinity",
					},
				},
			},

			InitContainers: createTestRunnerInitContainerObjects(testRunnerPod.Containers),

			RestartPolicy: v1.RestartPolicyOnFailure,
		},
	}
}

func createTestRunnerInitContainerObjects(containers []model.Containers) []v1.Container {

	var containerObjects []v1.Container

	for i := range containers {

		commands, args := buildInitiContainerCommand(containers[i].StartupCommands, containers[i].StartupCommandsArgs, containers[i].TestResultPath)

		containerObject := v1.Container{
			Name:            containers[i].ContainerName,
			Image:           containers[i].Image,
			ImagePullPolicy: v1.PullIfNotPresent,
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      "shared-data",
					MountPath: containers[i].TestResultPath,
				},
			},
			Command: commands,
			Args:    args,
		}

		containerObjects = append(containerObjects, containerObject)

	}

	return containerObjects
}

func buildInitiContainerCommand(inputCommands []string, inputArgs []string, testResultPath string) ([]string, []string) {

	var splittedCommands []string

	for _, e := range inputCommands {
		splittedCommands = append(splittedCommands, (strings.Split(e, " "))...)
	}

	var splittedArgs []string

	for _, e := range inputArgs {
		splittedArgs = append(splittedArgs, (strings.Split(e, " "))...)
	}

	commandsWithArgs := append(splittedCommands, splittedArgs...)

	var outputCommands []string
	var outputArgs []string

	if commandsWithArgs[0] == "/bin/sh" {

		if commandsWithArgs[1] == "-c" {

			outputCommands = append(outputCommands, commandsWithArgs[0])

			outputArgs = append(outputArgs, commandsWithArgs[1])
			outputArgs = append(outputArgs, "(("+strings.Join(commandsWithArgs[2:], " ")+" > "+testResultPath+"/testlog && echo $? > "+testResultPath+"/exitcode) || echo $? > "+testResultPath+"/exitcode)")

			return outputCommands, outputArgs

		}

		outputCommands = append(outputCommands, commandsWithArgs[0])
		outputArgs = append(outputArgs, "-c")
		outputArgs = append(outputArgs, "(("+strings.Join(commandsWithArgs[1:], " ")+" > "+testResultPath+"/testlog && echo $? > "+testResultPath+"/exitcode) || echo $? > "+testResultPath+"/exitcode)")

		return outputCommands, outputArgs
	}

	outputCommands = append(outputCommands, "/bin/sh")
	outputArgs = append(outputArgs, "-c")

	outputArgs = append(outputArgs, "(("+strings.Join(commandsWithArgs, " ")+" > "+testResultPath+"/testlog && echo $? > "+testResultPath+"/exitcode) || echo $? > "+testResultPath+"/exitcode)")

	return outputCommands, outputArgs

}
