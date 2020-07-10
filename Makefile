test:
	go test ./...

k3s-build:
	docker build -t registry.localhost:5000/parvaeres:latest .

k3s-push:
	docker push registry.localhost:5000/parvaeres:latest

k3s-deploy:
	kubectl apply -f ./deploy/parvaeres-server.yaml -n argocd

k3s-refresh:
	kubectl delete pod -l app=parvaeres-server -n argocd
