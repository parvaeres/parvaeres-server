.PHONY: test
test:
	go test -v ./...

.PHONY: k3s-up
k3s-up:
	./scripts/k3s-test-cluster.sh up

.PHONY: k3s-down
k3s-down:
	./scripts/k3s-test-cluster.sh down

.PHONY: k3s-build
k3s-build:
	docker build -t registry.localhost:5000/parvaeres:latest .

.PHONY: k3s-push
k3s-push:
	docker push registry.localhost:5000/parvaeres:latest

.PHONY: k3s-deploy
k3s-deploy:
	kubectl apply -f ./deploy/parvaeres-server.yaml -n argocd

.PHONY: k3s-refresh
k3s-refresh:
	kubectl delete pod -l app=parvaeres-server -n argocd

.PHONY: k3s-logs
k3s-logs:
	kubectl logs -n argocd -l app=parvaeres-server -f
