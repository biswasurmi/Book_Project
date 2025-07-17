Book Project: Out-of-Cluster Kubernetes Go Client
This guide provides step-by-step commands to implement an out-of-cluster Go client for managing Kubernetes resources (ConfigMap, Deployment, Service) for the Book_Project REST API, running in a Kind cluster (kind-book-project-cluster) at ~/Desktop/Book_Project. The client uses k8s.io/client-go with ~/.kube/config for authentication, deploys a book management API (book-project:v3), and supports testing endpoints like /api/v1/books. The setup assumes a Go application with Cobra CLI (k8sdeploy command), JWT_SECRET=bolaJabeNah, and NodePort: 30081 mapped to localhost:8080.
Prerequisites

Install Go (version 1.21+):
sudo apt update
sudo apt install golang-go
go version  # Should show go1.21 or higher


Install Docker:
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER
newgrp docker
docker --version


Install Kind:
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.23.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
kind version


Install kubectl:
curl -LO "https://dl.k8s.io/release/v1.31.0/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version --client


Set Up Project Directory:
cd ~/Desktop/Book_Project



Project Setup

Initialize Go Module (if not already done):
go mod init book-project


Install Dependencies:
go get k8s.io/client-go@v0.31.1
go get k8s.io/api@v0.31.1
go get k8s.io/apimachinery@v0.31.1
go get github.com/spf13/cobra@latest
go mod tidy


Ensure k8sdeploy.go:Verify cmd/k8sdeploy.go matches the provided code (out-of-cluster version). Key details:

Uses ~/.kube/config for authentication.
Creates ConfigMap (book-project-config, jwt-secret: bolaJabeNah).
Creates Deployment (book-project-deployment, image: book-project:v3, command: ./main startProject).
Creates Service (book-project-service, NodePort: 30081).

If needed, update cmd/k8sdeploy.go to use NodePort: 30081:
// In createResources, Service definition
Ports: []corev1.ServicePort{
    {
        Protocol:   corev1.ProtocolTCP,
        Port:       80,
        TargetPort: intstr.FromInt(8080),
        NodePort:   30081,
    },
}


Build the CLI:
go build -o book-cli main.go


Build and Load Docker Image:Ensure book-project:v3 is available:
docker build -t book-project:v3 .
kind load docker-image book-project:v3 --name book-project-cluster



Kubernetes Cluster Setup

Create Kind Config (kind-config.yaml):
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


Delete Existing Cluster (if needed):
kind delete cluster --name book-project-cluster


Create Cluster:
kind create cluster --name book-project-cluster --config kind-config.yaml
kubectl config use-context kind-book-project-cluster
kubectl cluster-info --context kind-book-project-cluster


Verify Cluster:
kubectl get nodes



Deploy Resources

Clean Up Existing Resources:
kubectl delete svc book-project-service
kubectl delete deploy book-project-deployment
kubectl delete configmap book-project-config


Run Create Command:
./book-cli k8sdeploy --create

Expected output:
2025/07/17 17:XX:XX Starting Kubernetes resource management
2025/07/17 17:XX:XX Creating Kubernetes resources
2025/07/17 17:XX:XX Creating ConfigMap 'book-project-config'
2025/07/17 17:XX:XX Created ConfigMap 'book-project-config'
2025/07/17 17:XX:XX Creating Deployment 'book-project-deployment'
2025/07/17 17:XX:XX Created Deployment 'book-project-deployment'
2025/07/17 17:XX:XX Creating Service 'book-project-service'
2025/07/17 17:XX:XX Created Service 'book-project-service'


Verify Resources:
kubectl get configmap book-project-config -o yaml
kubectl get deployment book-project-deployment
kubectl describe deployment book-project-deployment
kubectl get svc book-project-service
kubectl get pods


Check Pod Logs:
kubectl logs book-project-deployment-xxxx

Expected:
2025/07/17 03:XX:XX Starting Book Server on port 8080
2025/07/17 03:XX:XX Server listening on :8080



Test the API

Register a User:
curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password123"}'

Expected: 201 Created

Login to Get JWT Token:
token=$(curl -s -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"password123"}' | jq -r '.token')
echo $token


Create a Book (based on Book struct):
curl -X POST http://localhost:8080/api/v1/books -H "Content-Type: application/json" -H "Authorization: Bearer $token" -d '{"name":"Sample Book1","authorList":["John Doe","joh1"],"publishDate":"2025-07-17","isbn":"123-4567890123"}'

Expected: 201 Created

List Books:
curl -X GET http://localhost:8080/api/v1/books -H "Authorization: Bearer $token"


Test with Postman:

Open Postman, create a collection Book_Project.
Add requests: Register, Login, Create Book, List Books.
For Login, save token:pm.environment.set("token", pm.response.json().token);


Use Authorization: Bearer {{token}} for protected endpoints.



Update or Delete Resources

Update Deployment (e.g., scale to 2 replicas):
./book-cli k8sdeploy --update --replicas=2 --image=book-project:v3
kubectl get pods  # Should show 2 pods


Delete Resources:
./book-cli k8sdeploy --delete
kubectl get pods  # Should show no pods
kubectl get svc
kubectl get configmap


Delete Cluster (optional):
kind delete cluster --name book-project-cluster



Troubleshooting

NodePort Conflict:If error: spec.ports[0].nodePort: Invalid value: 30080: provided port is already allocated:
kubectl get svc --all-namespaces -o wide
kubectl delete svc book-project-service

Update cmd/k8sdeploy.go to use NodePort: 30081, rebuild:
go build -o book-cli main.go
./book-cli k8sdeploy --create


Image Pull Error:If pod shows ImagePullBackOff:
kubectl describe pod book-project-deployment-xxxx
kind load docker-image book-project:v3 --name book-project-cluster


API Errors:

Connection Refused:docker ps | grep kind-book-project-cluster  # Verify 0.0.0.0:8080->30081/tcp
kubectl port-forward svc/book-project-service 8080:80
curl http://localhost:8080/api/v1/register


401 Unauthorized:Ensure Authorization: Bearer $token header and JWT_SECRET=bolaJabeNah:kubectl get configmap book-project-config -o yaml


400 Bad Request:Verify JSON matches Book struct:{"name":"Sample Book1","authorList":["John Doe","joh1"],"publishDate":"2025-07-17","isbn":"123-4567890123"}




Kubeconfig Issues:
ls -l ~/.kube/config
kubectl config view


Timezone (logs show UTC):Add to k8sdeploy.go:
Env: []corev1.EnvVar{
    {Name: "TZ", Value: "Asia/Dhaka"},
    // ... JWT_SECRET ...
},


