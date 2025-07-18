Here's your complete `README.md` in GitHub-compatible markdown format, containing **all the commands** and steps for implementing and testing the **out-of-cluster Kubernetes Go client** for your Book Project:

---

````markdown
# 📚 Book Project: Out-of-Cluster Kubernetes Go Client

This guide provides **step-by-step commands** to implement an **out-of-cluster Go client** for managing Kubernetes resources (`ConfigMap`, `Deployment`, `Service`) for the Book_Project REST API, running in a **Kind cluster** (`kind-book-project-cluster`) at `~/Desktop/Book_Project`.

The Go client uses `k8s.io/client-go` and authenticates using `~/.kube/config`. It deploys a book management API (`book-project:v3`) and allows testing endpoints like `/api/v1/books` exposed at `localhost:8080` (NodePort 30081). CLI is built using Cobra with `k8sdeploy` command.

---

## 🧰 Prerequisites

### ✅ Install Go (v1.21+)

```bash
sudo apt update
sudo apt install golang-go
go version  # Should show go1.21 or higher
````

---

### 🐳 Install Docker

```bash
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
newgrp docker
docker --version
```

---

### 📦 Install Kind

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
kind version
```

---

### ☸️ Install kubectl

```bash
curl -LO "https://dl.k8s.io/release/v1.31.0/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version --client
```

---

### 🗂️ Set Up Project Directory

```bash
cd ~/Desktop/Book_Project
```

---

## ⚙️ Project Setup

### 📦 Initialize Go Module

```bash
go mod init book-project
```

### 📚 Install Dependencies

```bash
go get k8s.io/client-go@v0.31.1
go get k8s.io/api@v0.31.1
go get k8s.io/apimachinery@v0.31.1
go get github.com/spf13/cobra@latest
go mod tidy
```

---

## 🛠️ Update CLI (Out-of-Cluster)

Ensure `cmd/k8sdeploy.go`:

* Uses `~/.kube/config` for authentication.
* Creates the following:

  * ConfigMap: `book-project-config` with key `JWT_SECRET=bolaJabeNah`
  * Deployment: `book-project-deployment`, image `book-project:v3`, command `./main startProject`
  * Service: `book-project-service` with **NodePort: 30081**

Update service port if necessary:

```go
Ports: []corev1.ServicePort{
    {
        Protocol:   corev1.ProtocolTCP,
        Port:       80,
        TargetPort: intstr.FromInt(8080),
        NodePort:   30081,
    },
}
```

---

## 🔧 Build CLI Binary

```bash
go build -o book-cli main.go
```

---

## 🐳 Build and Load Docker Image

```bash
docker build -t book-project:v3 .
kind load docker-image book-project:v3 --name book-project-cluster
```

---

## ☸️ Kubernetes Cluster Setup

### 📝 Create Kind Config

```bash
cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30081
    hostPort: 8080
    protocol: TCP
EOF
```

### ❌ Delete Existing Cluster (if needed)

```bash
kind delete cluster --name book-project-cluster
```

### ✅ Create New Cluster

```bash
kind create cluster --name book-project-cluster --config kind-config.yaml
kubectl config use-context kind-book-project-cluster
kubectl cluster-info --context kind-book-project-cluster
kubectl get nodes
```

---

## 🚀 Deploy Resources

### 🔄 Clean Up Previous Resources

```bash
kubectl delete svc book-project-service
kubectl delete deploy book-project-deployment
kubectl delete configmap book-project-config
```

### 📦 Run Create Command

```bash
./book-cli k8sdeploy --create
```

✅ **Expected Output**:

```
2025/07/17 XX:XX:XX Starting Kubernetes resource management
2025/07/17 XX:XX:XX Creating ConfigMap 'book-project-config'
2025/07/17 XX:XX:XX Created ConfigMap 'book-project-config'
2025/07/17 XX:XX:XX Creating Deployment 'book-project-deployment'
2025/07/17 XX:XX:XX Created Deployment 'book-project-deployment'
2025/07/17 XX:XX:XX Creating Service 'book-project-service'
2025/07/17 XX:XX:XX Created Service 'book-project-service'
```

---

## 🔎 Verify Resources

```bash
kubectl get configmap book-project-config -o yaml
kubectl get deployment book-project-deployment
kubectl describe deployment book-project-deployment
kubectl get svc book-project-service
kubectl get pods
```

### 📝 Pod Logs

```bash
kubectl logs book-project-deployment-xxxx
```

✅ **Expected Log**:

```
Starting Book Server on port 8080
Server listening on :8080
```

---

## 🧪 Test the API

### 👤 Register a User

```bash
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{"email":"test@example.com","password":"password123"}'
```

### 🔐 Login to Get JWT Token

```bash
token=$(curl -s -X POST http://localhost:8080/api/v1/login \
-H "Content-Type: application/json" \
-d '{"email":"test@example.com","password":"password123"}' | jq -r '.token')
echo $token
```

### 📘 Create a Book

```bash
curl -X POST http://localhost:8080/api/v1/books \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $token" \
-d '{"name":"Sample Book1","authorList":["John Doe","joh1"],"publishDate":"2025-07-17","isbn":"123-4567890123"}'
```

### 📚 List Books

```bash
curl -X GET http://localhost:8080/api/v1/books \
-H "Authorization: Bearer $token"
```

---

## 📬 Postman Testing

1. Create a new **collection**: `Book_Project`.

2. Add requests:

   * `POST /api/v1/register`
   * `POST /api/v1/login`
   * `POST /api/v1/books`
   * `GET /api/v1/books`

3. In **Login request**, save token in **Tests tab**:

   ```js
   pm.environment.set("token", pm.response.json().token);
   ```

4. Use `Authorization: Bearer {{token}}` header for protected endpoints.

---

## 🔄 Update or Delete Resources

### 🔁 Update Deployment (e.g., Scale to 2 Replicas)

```bash
./book-cli k8sdeploy --update --replicas=2 --image=book-project:v3
kubectl get pods
```

### ❌ Delete All Resources

```bash
./book-cli k8sdeploy --delete
kubectl get pods
kubectl get svc
kubectl get configmap
```

### 🧹 Delete Kind Cluster

```bash
kind delete cluster --name book-project-cluster
```

---

## 🧑‍🔧 Troubleshooting

### ❗ NodePort Conflict

```bash
kubectl get svc --all-namespaces -o wide
kubectl delete svc book-project-service
```

Update `NodePort` in `cmd/k8sdeploy.go` to `30081`, rebuild and redeploy:

```bash
go build -o book-cli main.go
./book-cli k8sdeploy --create
```

---

### ❗ Image Pull Error

```bash
kubectl describe pod book-project-deployment-xxxx
kind load docker-image book-project:v3 --name book-project-cluster
```

---

### ❗ API Issues

#### 🔌 Connection Refused

```bash
docker ps | grep kind-book-project-cluster
kubectl port-forward svc/book-project-service 8080:80
curl http://localhost:8080/api/v1/register
```

#### 🔐 401 Unauthorized

```bash
kubectl get configmap book-project-config -o yaml
```

Ensure:

* `Authorization: Bearer $token`
* `JWT_SECRET=bolaJabeNah`

#### 📉 400 Bad Request

Ensure JSON format matches `Book` struct:

```json
{
  "name": "Sample Book1",
  "authorList": ["John Doe", "joh1"],
  "publishDate": "2025-07-17",
  "isbn": "123-4567890123"
}
```

---

## 🌐 Kubeconfig Issues

```bash
ls -l ~/.kube/config
kubectl config view
```

---

## 🌍 Timezone (Optional)

In `k8sdeploy.go`, add:

```go
Env: []corev1.EnvVar{
    {Name: "TZ", Value: "Asia/Dhaka"},
    {Name: "JWT_SECRET", Value: "bolaJabeNah"},
},
```

---


