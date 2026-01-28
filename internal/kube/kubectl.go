package kube

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runKubectl(args ...string) (string, string, error) {
	cmd := exec.Command("kubectl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func EnsureNamespace(ns string) error {
	// Check if exists
	_, _, err := runKubectl("get", "namespace", ns)
	if err == nil {
		return nil
	}

	// Create
	_, se, err := runKubectl("create", "namespace", ns)
	if err != nil {
		return fmt.Errorf("create namespace %q failed: %v (%s)", ns, err, se)
	}
	return nil
}

func ApplyYAML(yaml string) error {
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewBufferString(yaml)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("kubectl apply failed: %v (%s)", err, stderr.String())
	}
	return nil
}

func PrintPods(namespace, appName string) error {
	so, se, err := runKubectl("get", "pods", "-n", namespace, "-l", "app="+appName, "-o", "wide")
	if err != nil {
		return fmt.Errorf("get pods failed: %v (%s)", err, se)
	}
	fmt.Print(so)
	return nil
}

func PrintService(namespace, appName string) error {
	so, se, err := runKubectl("get", "svc", "-n", namespace, appName, "-o", "wide")
	if err != nil {
		// service might not exist yet; show all services for context
		so2, se2, err2 := runKubectl("get", "svc", "-n", namespace, "-o", "wide")
		if err2 != nil {
			return fmt.Errorf("get service failed: %v (%s)", err, se)
		}
		fmt.Println("Service not found by app name; showing services in namespace:")
		fmt.Print(so2)
		if se2 != "" {
			fmt.Println(se2)
		}
		return nil
	}
	fmt.Print(so)
	return nil
}

func PrintEvents(namespace string) error {
	so, se, err := runKubectl("get", "events", "-n", namespace, "--sort-by=.metadata.creationTimestamp")
	if err != nil {
		return fmt.Errorf("get events failed: %v (%s)", err, se)
	}
	fmt.Print(so)
	return nil
}
