SHELL := /bin/bash

run:
	go run main.go

# ==================================================================
# Building containers

VERSION := 1.0

all: build

build:
	docker build \
		-f zarf/docker/dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==================================================================
# Running from within k8s

k3d-update: build k3d-load k3d-restart

k3d-load:
	k3d image import -c local-cluster service-amd64:$(VERSION)

k3d-apply:
	cat zarf/k8s/base/service-pod/base-service.yaml | kubectl apply -f -

k3d-restart:
	kubectl rollout restart deployment service-deploy --namespace=service-system

k3d-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

k3d-status-service:
	kubectl get pods -o wide --watch -n service-system

k3d-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100

k3d-describe:
	kubectl describe pod -l app=service
