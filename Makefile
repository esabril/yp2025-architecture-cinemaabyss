helm-start:
	helm install cinemaabyss ./src/kubernetes/helm --namespace cinemaabyss --create-namespace

helm-stop:
	helm uninstall cinemaabyss -n cinemaabyss

run-docker-tests:
	cd $(PWD)/tests/postman && ./run-tests.sh -d -e docker

run-kube-tests:
	cd $(PWD)/tests/postman && npm run test:kubernetes

deploy: deploy-config deploy-infra deploy-services deploy-proxy deploy-ingress

deploy-config:
	kubectl apply -f src/kubernetes/namespace.yaml
	kubectl apply -f src/kubernetes/configmap.yaml
	kubectl apply -f src/kubernetes/secret.yaml
	kubectl apply -f src/kubernetes/dockerconfigsecret.yaml

deploy-infra:
	kubectl apply -f src/kubernetes/postgres-init-configmap.yaml
	kubectl apply -f src/kubernetes/postgres.yaml
	kubectl apply -f src/kubernetes/kafka/kafka.yaml

deploy-services:
	kubectl apply -f src/kubernetes/monolith.yaml
	kubectl apply -f src/kubernetes/movies-service.yaml
	kubectl apply -f src/kubernetes/events-service.yaml

deploy-proxy:
	kubectl apply -f src/kubernetes/proxy-service.yaml

deploy-ingress:
	kubectl apply -f src/kubernetes/ingress.yaml

kube-clean:
	kubectl delete all --all -n cinemaabyss
	kubectl delete namespace cinemaabyss

diagrams-generate:
	  $(RM) $(PWD)/diagrams/*.png && java -jar plantuml.jar $(PWD)/diagrams/*.puml