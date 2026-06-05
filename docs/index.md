# Jumpcloud Provider

This provider can be used to add users&groups to
Jumpcloud as well as the mappings.

## Authentication

The provider supports two mutually-exclusive auth modes. Configure **one**:

* **Service-account OAuth (recommended)** — `client_id` + `client_secret` from a
  JumpCloud Service Account. The provider performs a `client_credentials` grant
  against `https://admin-oauth.id.jumpcloud.com/oauth2/token` at configure time
  and sends the resulting short-lived (~1h) token as `Authorization: Bearer`.
  Not tied to a person.
* **Legacy x-api-key** — `api_key`, an admin-user API key sent as `x-api-key`.

OAuth takes precedence when both are set. `client_id` and `client_secret` must
be provided together. Each can be supplied via environment variable
(`JUMPCLOUD_CLIENT_ID`, `JUMPCLOUD_CLIENT_SECRET`, `JUMPCLOUD_API_KEY`).

```hcl
provider "jumpcloud" {
  client_id     = var.jc_client_id
  client_secret = var.jc_client_secret
}
```

The minted token lives for the run (~1h) — sufficient for a single plan/apply.
An apply that runs longer than the token lifetime is an accepted edge case.

## ToDo
 * Extended Documentation
 * System Groups
 * Bugfixing
