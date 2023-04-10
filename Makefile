deploy:
	kubectl create ns gatekeeper-system
	kubectl get secret redis --namespace=default -o yaml | sed 's/namespace: .*/namespace: gatekeeper-system/' | kubectl apply -f -
	kubectl apply -f pub.yaml

build:
	docker build -f Dockerfile -t docker.io/jaydipgabani/publisher:latest .