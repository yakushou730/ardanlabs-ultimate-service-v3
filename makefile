SHELL := /bin/bash

# ==================================================================
# Testing running system

# expvarmon -ports=":30191" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# hey -m GET -c 100 -n 10000 http://localhost:30190/v1/test

# To generate a private/public key PEM file
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./admin genkey

# Testing Auth
# curl -il http://localhost:30190/v1/testauth
# curl -il -H "Authorization: Bearer wrong-test-token" http://localhost:30190/v1/testauth

# ==================================================================

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

admin:
	go run app/tooling/admin/main.go

# ==================================================================
# Building containers

VERSION := 1.0

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==================================================================
# Running from within k8s

k3d-update: sales-api k3d-load k3d-restart

k3d-update-apply: sales-api k3d-load k3d-apply

k3d-load:
	cd zarf/k8s/k3d/sales-deploy; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	k3d image import -c local-cluster sales-api-amd64:$(VERSION)

k3d-apply:
	kustomize build zarf/k8s/k3d/sales-deploy | kubectl apply -f -

k3d-restart:
	kubectl rollout restart deployment sales-deploy --namespace=sales-system

k3d-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

k3d-status-sales:
	kubectl get pods -o wide --watch -n sales-system

k3d-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

k3d-describe:
	kubectl describe pod -l app=sales

tidy:
	go mod tidy
