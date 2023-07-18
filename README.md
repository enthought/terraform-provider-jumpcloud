# JumpCloud Terraform Provider

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13+
- [Go](https://golang.org/doc/install) 1.20

## Building The Provider

Clone repository to: `$GOPATH/src/github.com/enthought/terraform-provider-jumpcloud`

```sh
mkdir -p $GOPATH/src/github.com/enthought
cd $GOPATH/src/github.com/enthought
git clone git@github.com:enthought/terraform-provider-jumpcloud
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/enthought/terraform-provider-jumpcloud
make build
```

## Releasing the provider

Use goreleaser

```
git tag -a $NEW_VERSION -m "REL: release $NEW_VERSION of the jumpcloud provider"
git push --tags
goreleaser release --rm-dist
```

Once done and thoroughly tested, update the `deploy_providers.sh script in the main terraform repository
and deploy the new version to the users.
## Using the provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

The Jumpcloud API key needs to be set before using the provider. It can either be retrieved via the API or through the UI : When selecting a resource, the ID is part of URL.
Export `JUMPCLOUD_API_KEY` to set it.

The Jumpcloud "Organization ID" is optional as only needed for multi-tenant-setups.
Export `JUMPCLOUD_ORG_ID` to set it.
