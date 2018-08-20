export DRAFT=false
export ARCHS=linux/amd64

jenkins: start_build

start_build:
	curl -sSLo golang.sh https://raw.githubusercontent.com/Luzifer/github-publish/master/golang.sh
	docker run --rm -i \
		-e "GITHUB_TOKEN=$(GITHUB_TOKEN)" \
		-v "$(CURDIR):/go/src/github.com/contentflow/terraform-provider-stackpath" \
		-w "/go/src/github.com/contentflow/terraform-provider-stackpath" \
		golang:latest \
		make publish

publish:
		bash golang.sh
