.PHONY: run-auth run-noauth docker-build docker-run-auth docker-run-noauth docker-push helm-repo-add helm-run helm-port-forward test help

GO_VERSION = 1.24
PORT = 8080
REGISTRY = urmibiswas
BINS = book_project
VERSION = v3
NAMESPACE = default
HELM_REPO_NAME = my-charts
HELM_REPO_URL = https://biswasurmi.github.io/my-helm-charts/book-project

run-auth:
	go run main.go startProject --port=$(PORT) --auth=true

run-noauth:
	go run main.go startProject --port=$(PORT) --auth=false

docker-build:
	docker build -t $(REGISTRY)/$(BINS):$(VERSION) -f Dockerfile .

docker-run-auth:
	docker run -d -p $(PORT):$(PORT) $(REGISTRY)/$(BINS):$(VERSION) ./main startProject --port=$(PORT) --auth=true

docker-run-noauth:
	docker run -d -p $(PORT):$(PORT) $(REGISTRY)/$(BINS):$(VERSION) ./main startProject --port=$(PORT) --auth=false

docker-push:
	docker build -t $(REGISTRY)/$(BINS):$(VERSION) -f Dockerfile . && docker push $(REGISTRY)/$(BINS):$(VERSION)

k8s-run:
	kubectl apply -f deployment.yaml -n $(NAMESPACE)

helm-repo-add:
	helm repo add $(HELM_REPO_NAME) $(HELM_REPO_URL) || true
	helm repo update

helm-run:
	helm upgrade --install book-project my-charts/book-project --namespace $(NAMESPACE) \
		--set image.repository=$(REGISTRY)/$(BINS) \
		--set image.tag=$(VERSION)

helm-port-forward:
	@POD_NAME=$$(kubectl get pods --namespace $(NAMESPACE) -l "app.kubernetes.io/name=book-project,app.kubernetes.io/instance=book-project" -o jsonpath="{.items[0].metadata.name}"); \
	CONTAINER_PORT=$$(kubectl get pod --namespace $(NAMESPACE) $$POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}"); \
	echo "Port forwarding pod $$POD_NAME port $$CONTAINER_PORT to localhost:8080"; \
	kubectl --namespace $(NAMESPACE) port-forward $$POD_NAME 8080:$$CONTAINER_PORT

test:
	docker run --rm -v $(PWD):/app -w /app golang:$(GO_VERSION) go test -v ./test_file

clean:
	rm -rf bin/*
	docker rmi $(REGISTRY)/$(BINS):$(VERSION) || true

help:
	@echo "Available targets:"
	@echo "  run-auth         : Run app locally with auth"
	@echo "  run-noauth       : Run app locally without auth"
	@echo "  docker-build     : Build Docker image ($(REGISTRY)/$(BINS):$(VERSION))"
	@echo "  docker-run-auth  : Run Docker container with auth"
	@echo "  docker-run-noauth: Run Docker container without auth"
	@echo "  docker-push      : Push Docker image to Docker Hub"
	@echo "  k8s-run          : Apply Kubernetes manifests from deployment.yaml"
	@echo "  helm-repo-add    : Add/update Helm chart repository"
	@echo "  helm-run         : Install/upgrade Helm chart"
	@echo "  helm-port-forward: Port forward Helm deployed pod to localhost:8080"
	@echo "  test             : Run Go tests inside Docker"
	@echo "  clean            : Remove build artifacts and Docker image"
