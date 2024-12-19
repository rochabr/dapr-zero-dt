.PHONY: build deploy clean

build:
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t rochabr/ping-service:v5 --push ./ping-service
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t rochabr/pong-service:v5 --push ./pong-service

bd-pong:
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t rochabr/pong-service:v6 --push ./pong-service
	kubectl apply -f pong-service/k8s/deployment.yaml


deploy:
	kubectl apply -f ping-service/k8s/deployment.yaml
	kubectl apply -f pong-service/k8s/deployment.yaml

clean:
	kubectl delete -f ping-service/k8s/deployment.yaml
	kubectl delete -f pong-service/k8s/deployment.yaml
