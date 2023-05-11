
APP_NAME = uptime

setup:
	go mod download
	go mod vendor

build:
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

apply-example:
	@read -p "This will switch the kubectl context to docker-desktop and apply the manifests in the example folder, are you sure? (y/N): " ans && ans=$${ans:-N} ; \
	if [ $${ans} = y ] || [ $${ans} = Y ]; then \
		kubectl config use-context docker-desktop ; \
		kubectl apply -f example/ ; \
		kubectl create clusterrolebinding serviceaccounts-cluster-admin --clusterrole=cluster-admin --group=system:serviceaccounts ; \
		echo "Done" ; \
	fi

mocks:
	go generate ./...

test:
	go test -v -cover ./...