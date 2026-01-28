# Cloud Application Lifecycle Governance Platform

This project is a lightweight, enterprise-style **cloud platform engineering tool** designed to
demonstrate application lifecycle governance, Kubernetes deployment automation, policy enforcement,
diagnostics, observability, and CI/CD readiness in a private cloud–like environment.

The platform simulates **internal cloud platform workflows** commonly found in large-scale
organizations, including:
- Application onboarding via standardized specifications
- Governance and security policy validation
- Automated Kubernetes deployments
- Platform-level diagnostics and troubleshooting
- Observability through metrics and health endpoints
- CI/CD-based build and validation pipelines

---

## Project Structure

```
cloud-app-platform/
├── cmd/
│   └── platformctl/
│       └── main.go
├── internal/
│   ├── cli/
│   │   ├── deploy.go
│   │   ├── validate.go
│   │   ├── status.go
│   │   ├── metrics.go
│   │   └── root.go
│   ├── spec/
│   │   └── spec.go
│   ├── validator/
│   │   └── validator.go
│   ├── render/
│   │   └── render.go
│   ├── kube/
│   │   └── kubectl.go
│   └── observability/
│       └── store.go
├── .github/
│   └── workflows/
│       └── ci.yml
├── app.yaml
├── go.mod
├── go.sum
└── README.md
```

---

## Prerequisites

- Go 1.21+
- Docker
- kubectl
- kind (or minikube)
- macOS / Linux / Windows

Verify installations:
```bash
go version
docker version
kubectl version --client
kind version
```

---

## Kubernetes Cluster Setup

Create a local Kubernetes cluster using kind:
```bash
kind create cluster --name platform-dev
kubectl get nodes
```

---

## Application Specification (`app.yaml`)

```yaml
name: demo-nginx
namespace: demo
image: nginx:1.25
port: 80
replicas: 2
resources:
  cpu: "250m"
  memory: "256Mi"
```

---

## 🚀 Commands to Run the Platform

### Validate application spec
```bash
go run ./cmd/platformctl validate app.yaml
```

### Deploy application
```bash
go run ./cmd/platformctl deploy app.yaml
```

### Status & diagnostics
```bash
go run ./cmd/platformctl status -n demo demo-nginx
```

---

## Observability & Metrics

```bash
go run ./cmd/platformctl metrics --port 9090
curl http://localhost:9090/metrics | egrep "platform_"
```

---

## CI/CD Pipeline

GitHub Actions workflow for vetting, testing, and building the CLI:
```
.github/workflows/ci.yml
```

---

## Key Engineering Concepts

- Application lifecycle governance
- Kubernetes-based private cloud operations
- Cloud-native automation
- Policy enforcement & compliance
- Observability-first platform design
- CI/CD pipelines

---

## Author

Hardik Kotangale  
Indiana University Bloomington
