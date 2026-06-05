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

### Authentication

The provider supports two mutually-exclusive auth modes (configure one):

* **Service-account OAuth (recommended):** set `client_id` + `client_secret`
  (or export `JUMPCLOUD_CLIENT_ID` / `JUMPCLOUD_CLIENT_SECRET`) from a JumpCloud
  Service Account. The provider exchanges them for a short-lived Bearer access
  token via the `client_credentials` grant. Not tied to an individual account.
* **Legacy x-api-key:** set `api_key` (or export `JUMPCLOUD_API_KEY`), an
  admin-user API key retrieved via the API or the UI. OAuth takes precedence
  when both are set.

The Jumpcloud "Organization ID" is optional and only needed for
multi-tenant-setups. Export `JUMPCLOUD_ORG_ID` to set it.
