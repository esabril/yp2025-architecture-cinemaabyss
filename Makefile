helm-start:
	helm install cinemaabyss ./src/kubernetes/helm --namespace cinemaabyss --create-namespace

helm-stop:
	helm uninstall cinemaabyss -n cinemaabyss

run-docker-tests:
	cd $(PWD)/tests/postman && ./run-tests.sh -d -e docker
