Here’s a complete `README.md` for your **Kind-powered local Kubernetes Book Project**, including Docker image setup, cluster creation, deployment, port forwarding, and testing. You can copy this into your project root as `README.md`:

---

````markdown
# 📚 Book Project – Local Kubernetes Deployment with Kind

This project runs a Book API server locally on Kubernetes using [Kind (Kubernetes in Docker)](https://kind.sigs.k8s.io/). It supports full user and book management with JWT and Basic Authentication.

---

## 📦 Prerequisites

- Docker (v28.2.2 or later)
- Kind (v0.16.0 or later)
- kubectl (v1.33.x)
- `book-project:v3` Docker image (build locally or pull from Docker Hub)
- Your app must expose HTTP server at port `8080` inside container

---

## 🛠️ 1. Build & Tag Docker Image

If you're using a local build:
```bash
docker build -t book-project:v3 .
````

If you're using Docker Hub:

```bash
docker pull urmibiswas/book_project:v3
docker tag urmibiswas/book_project:v3 book-project:v3
```

---

## 🧱 2. Create Kind Cluster with Port Mapping

Create a `kind-config.yaml`:

```yaml
# kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 30080
    protocol: TCP
```

Then create the cluster:

```bash
kind create cluster --name book-project-cluster --config kind-config.yaml
```

Check cluster status:

```bash
kubectl cluster-info --context kind-book-project-cluster
kubectl config current-context
```

---

## 🐳 3. Load Docker Image into Kind Cluster

```bash
kind load docker-image book-project:v3 --name book-project-cluster
```

---

## 🚀 4. Deploy Application to Kubernetes

Create a `deployment.yaml` file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: book-project-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: book-project
  template:
    metadata:
      labels:
        app: book-project
    spec:
      containers:
      - name: book-project
        image: book-project:v3
        imagePullPolicy: Never
        ports:
        - containerPort: 8080
        env:
        - name: JWT_SECRET
          valueFrom:
            configMapKeyRef:
              name: book-project-config
              key: jwt-secret
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: book-project-config
data:
  jwt-secret: "bolaJabeNah"
---
apiVersion: v1
kind: Service
metadata:
  name: book-project-service
spec:
  selector:
    app: book-project
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
    nodePort: 30080
  type: NodePort
```

Apply the deployment:

```bash
kubectl apply -f deployment.yaml
kubectl get pods
kubectl logs -l app=book-project
```

---

## 🌐 5. Access API via Port Forwarding

Since Kind doesn't expose NodePorts on `localhost` by default, run:

```bash
kubectl port-forward service/book-project-service 8080:80
```

Now your app is accessible at:

```
http://localhost:8080
```

---

## 🔬 6. Test API in Postman or curl

Here are some quick test examples:

### Register a User

```http
POST http://localhost:8080/api/v1/register
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

### Login (to get token)

```http
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

### Use Token to Access Protected Endpoint

```bash
curl -H "Authorization: Bearer <your_token>" http://localhost:8080/api/v1/books
```

---

## 🧼 Tear Down

To delete the cluster:

```bash
kind delete cluster --name book-project-cluster
```

---

## ✅ Summary

| Component          | Status                                                               |
| ------------------ | -------------------------------------------------------------------- |
| Kubernetes Cluster | ✅ Local Kind cluster                                                 |
| Docker Image       | ✅ Built & loaded                                                     |
| App Running        | ✅ Confirmed via logs                                                 |
| API Accessible     | ✅ On [http://localhost:8080](http://localhost:8080) via port-forward |

---


