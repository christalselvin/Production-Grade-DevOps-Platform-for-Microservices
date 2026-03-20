# 🚀 Production-Grade Cloud Native DevOps Platform for Microservices

<div align="center">

![Platform Status](https://img.shields.io/badge/Status-Active-brightgreen)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.28-326CE5?logo=kubernetes&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-IaC-7B42BC?logo=terraform&logoColor=white)
![Jenkins](https://img.shields.io/badge/Jenkins-CI%2FCD-D24939?logo=jenkins&logoColor=white)
![ArgoCD](https://img.shields.io/badge/ArgoCD-GitOps-EF7B4D?logo=argo&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-Cloud-FF9900?logo=amazonaws&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue)

**An end-to-end production-style DevOps ecosystem simulating real-world cloud infrastructure used in modern tech companies.**

[Architecture](#architecture) · [Microservices](#microservices) · [DevOps Stack](#devops-stack) · [Getting Started](#getting-started) · [CI/CD Pipeline](#cicd-pipeline) · [Monitoring](#monitoring) · [Lessons Learned](#lessons-learned)

</div>

---

## 📌 Overview

This project demonstrates how to design, provision, and operate a **fully automated DevOps platform** for a microservices-based application on AWS. It covers the complete software delivery lifecycle — from infrastructure provisioning with Terraform, containerization with Docker, orchestration with Kubernetes, automated pipelines with Jenkins, GitOps deployments with ArgoCD, to full-stack observability with Prometheus and Grafana.

> **Goal:** Replicate the kind of infrastructure you'd find inside a real product company — not a toy example, but a platform you could actually run production workloads on.

---

## 🏗️ Architecture

```
                          ┌─────────────────────────────────────────────┐
                          │                  AWS Cloud                   │
                          │                                              │
  Developer ──push──▶  GitHub ──webhook──▶  Jenkins (CI)               │
                          │                      │                       │
                          │                  Build & Test                │
                          │                      │                       │
                          │              Docker Image ──push──▶ JFrog   │
                          │                      │           Artifactory │
                          │                  ArgoCD (GitOps)             │
                          │                      │                       │
                          │              ┌───────▼────────┐              │
                          │              │  Kubernetes     │              │
                          │              │  Cluster (EKS)  │              │
                          │              │                 │              │
                          │              │  User Service   │              │
                          │              │  Product Svc    │──▶ RabbitMQ │
                          │              │  Order Service  │              │
                          │              │  Payment Svc    │              │
                          │              └───────┬─────────┘              │
                          │                      │                        │
                          │              Prometheus ──▶ Grafana           │
                          └─────────────────────────────────────────────┘
```

**Key architectural decisions:**
- **Microservices** communicate asynchronously via RabbitMQ to decouple services and handle traffic spikes
- **GitOps** (ArgoCD) is the single source of truth for deployments — no manual `kubectl apply`
- **Infrastructure is fully code-driven** — zero manual AWS console clicks after initial bootstrap
- **Observability-first** — Prometheus scrapes all services; Grafana dashboards are provisioned as code

---

## 🧩 Microservices

| Service | Language | Port | Responsibility |
|---------|----------|------|----------------|
| **User Service** | Go | 3001 | Authentication, user management |
| **Product Service** | Go | 3002 | Product catalog, inventory |
| **Order Service** | Go | 3003 | Order lifecycle management |
| **Payment Service** | Go | 3004 | Payment processing, transaction records |

Each service is independently deployable, has its own Kubernetes `Deployment` and `Service`, and exposes `/health` and `/metrics` endpoints.

---

## 🛠️ DevOps Stack

| Layer | Tool | Purpose |
|-------|------|---------|
| **Source Control** | GitHub | Code hosting, PR workflows, webhooks |
| **CI/CD** | Jenkins | Build, test, image build on every commit |
| **Containerization** | Docker | Consistent build environments |
| **Orchestration** | Kubernetes (EKS) | Container scheduling, scaling, self-healing |
| **IaC** | Terraform | AWS infrastructure provisioning |
| **Config Mgmt** | Ansible | Server configuration, bootstrap scripts |
| **Artifact Storage** | JFrog Artifactory | Docker image registry |
| **GitOps** | ArgoCD | Declarative, Git-driven deployments |
| **Messaging** | RabbitMQ | Async inter-service communication |
| **Metrics** | Prometheus | Metrics collection and alerting |
| **Dashboards** | Grafana | Visualization and observability |

---

## ⚙️ Getting Started

### Prerequisites

- AWS CLI configured with appropriate IAM permissions
- Terraform >= 1.5
- kubectl >= 1.28
- Docker >= 24
- Helm >= 3.12

### 1. Provision AWS Infrastructure

```bash
cd terraform/
terraform init
terraform plan -out=tfplan
terraform apply tfplan
```

This provisions: VPC, subnets, EKS cluster, IAM roles, ECR repositories, and RDS (if applicable).

### 2. Configure Kubernetes Access

```bash
aws eks update-kubeconfig --region ap-south-1 --name devops-cluster
kubectl get nodes  # verify cluster is healthy
```

### 3. Deploy Core Platform Services

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Install Prometheus + Grafana (via Helm)
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Deploy RabbitMQ
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install rabbitmq bitnami/rabbitmq -n messaging --create-namespace
```

### 4. Deploy Microservices via ArgoCD

```bash
kubectl apply -f argocd/application.yaml
# ArgoCD will sync all microservices from the GitOps repo
```

---

## 🔁 CI/CD Pipeline

The Jenkins pipeline runs on every `git push` to `main`:

```
git push
   │
   ▼
[Stage 1] Checkout & Unit Tests
   │
   ▼
[Stage 2] Docker Build
   │
   ▼
[Stage 3] Push to JFrog Artifactory
   │
   ▼
[Stage 4] Update Helm chart values (image tag)
   │
   ▼
[Stage 5] ArgoCD detects Git change → Syncs to Kubernetes
   │
   ▼
[Stage 6] Health checks & deployment verification
```

**Jenkinsfile snippet:**
```groovy
pipeline {
  agent any
  stages {
    stage('Build & Test') {
      steps {
        sh 'go test ./...'
      }
    }
    stage('Docker Build') {
      steps {
        sh "docker build -t ${ARTIFACTORY_URL}/user-service:${BUILD_NUMBER} ."
        sh "docker push ${ARTIFACTORY_URL}/user-service:${BUILD_NUMBER}"
      }
    }
    stage('Update Helm Values') {
      steps {
        sh "sed -i 's/tag:.*/tag: ${BUILD_NUMBER}/' helm/values.yaml"
        sh "git commit -am 'ci: bump image to ${BUILD_NUMBER}' && git push"
      }
    }
  }
}
```

---

## 📊 Monitoring

Grafana dashboards are provisioned automatically via the Helm chart values. Key dashboards:

| Dashboard | What it shows |
|-----------|--------------|
| **Cluster Overview** | Node CPU/memory, pod count, restarts |
| **Service Metrics** | Request rate, error rate, latency (RED method) |
| **RabbitMQ** | Queue depth, message throughput, consumer lag |
| **Jenkins** | Build success rate, pipeline duration |

**Prometheus alerting rules** are defined in `monitoring/alerts/` and cover: pod crash-looping, high error rates, disk pressure, and queue backlog.

---

## 📁 Repository Structure

```
├── services/
│   ├── user-service/
│   ├── product-service/
│   ├── order-service/
│   └── payment-service/
├── terraform/              # AWS infrastructure as code
├── kubernetes/
│   ├── deployments/        # K8s manifests per service
│   └── namespaces/
├── helm/                   # Helm chart for all services
├── argocd/                 # ArgoCD Application definitions
├── jenkins/                # Jenkinsfile and shared libraries
├── ansible/                # Server bootstrap playbooks
├── monitoring/
│   ├── prometheus/         # Scrape configs, alert rules
│   └── grafana/            # Dashboard JSON exports
└── docs/
    └── architecture.md
```

---

## 🧠 Lessons Learned

Building this platform surfaced several real-world DevOps challenges:

1. **GitOps drift** — ArgoCD's self-healing is powerful but requires strict branch protection; a direct `kubectl apply` will be reverted silently.
2. **Secret management** — Kubernetes secrets should not be stored in Git. Integrated AWS Secrets Manager / Sealed Secrets as a follow-up.
3. **Resource limits matter** — Initial deployments had no resource limits set. Services OOM-killed each other under load. Tuning `requests` and `limits` was critical.
4. **Prometheus cardinality** — High-cardinality labels (like user IDs) in metrics caused memory spikes. Label design needs upfront planning.
5. **Terraform state** — Remote state in S3 with DynamoDB locking is non-negotiable for team use.

---

## 🚧 Roadmap

- [ ] Add Istio service mesh for mTLS and traffic management
- [ ] Implement Sealed Secrets for GitOps-safe secret management
- [ ] Add HPA (Horizontal Pod Autoscaler) based on RabbitMQ queue depth
- [ ] Implement canary deployments via ArgoCD Rollouts
- [ ] Add distributed tracing with Jaeger / OpenTelemetry
- [ ] Add a `docker-compose.yml` for local development

---

## 🤝 Connect

**Christal Selvin** — DevOps Engineer | Platform Engineering | SRE

[![LinkedIn](https://img.shields.io/badge/LinkedIn-christalcs4-0A66C2?logo=linkedin)](https://www.linkedin.com/in/christalcs4)
[![GitHub](https://img.shields.io/badge/GitHub-christalselvin-181717?logo=github)](https://github.com/christalselvin)
[![LeetCode](https://img.shields.io/badge/LeetCode-christal4-FFA116?logo=leetcode)](https://leetcode.com/christal4/)
[![Email](https://img.shields.io/badge/Email-christalselvin5@gmail.com-EA4335?logo=gmail)](mailto:christalselvin5@gmail.com)

---

