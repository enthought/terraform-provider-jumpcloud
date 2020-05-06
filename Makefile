install: 
	mkdir -p ~/.terraform.d/plugins
	GO111MODULE=on go build -o ~/.terraform.d/plugins/terraform-provider-jumpcloud
build:
	GO111MODULE=on go build .
testacc: 
	TF_ACC=true GO111MODULE=on go test -v ./...
