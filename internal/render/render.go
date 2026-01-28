package render

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/HardikKotangale/Cloud-App-Platform/internal/spec"
)

const deploymentTmpl = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
    managed-by: platformctl
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
        managed-by: platformctl
    spec:
      containers:
        - name: {{.Name}}
          image: {{.Image}}
          ports:
            - containerPort: {{.Port}}
          resources:
            requests:
              cpu: {{.Resources.CPU}}
              memory: {{.Resources.Memory}}
            limits:
              cpu: {{.Resources.CPU}}
              memory: {{.Resources.Memory}}
`

const serviceTmpl = `
---
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
    managed-by: platformctl
spec:
  selector:
    app: {{.Name}}
  ports:
    - name: http
      port: {{.Port}}
      targetPort: {{.Port}}
  type: ClusterIP
`

func RenderManifests(a *spec.AppSpec) (string, error) {
	var out bytes.Buffer

	t1, err := template.New("deploy").Parse(deploymentTmpl)
	if err != nil {
		return "", fmt.Errorf("parse deployment template: %w", err)
	}
	if err := t1.Execute(&out, a); err != nil {
		return "", fmt.Errorf("render deployment: %w", err)
	}

	t2, err := template.New("svc").Parse(serviceTmpl)
	if err != nil {
		return "", fmt.Errorf("parse service template: %w", err)
	}
	if err := t2.Execute(&out, a); err != nil {
		return "", fmt.Errorf("render service: %w", err)
	}

	return out.String(), nil
}
