# Contracts

Versioned API and event contracts for the Go microservice decomposition.

## Scope

This repo owns cross-service compatibility for:

- `auth-service`
- `user-service`
- `notification-otp-service`
- `api-gateway`

## Layout

```text
openapi/
  auth-service.yaml
  user-service.yaml
  notification-otp-service.yaml
events/
  user-registered.v1.schema.json
  otp-requested.v1.schema.json
  password-reset-requested.v1.schema.json
  email-verification-requested.v1.schema.json
samples/
  *.json
scripts/
  validate-contracts.sh
```

## Versioning Rules

- Additive fields are backward compatible.
- Removing or renaming fields requires a new version.
- Event consumers must ignore unknown fields.
- Event schemas are strict producer contracts; consumers deserialize permissively for forward compatibility.
- Services test against pinned contract versions and a latest-contract compatibility job.

## Internal Auth

Internal APIs use signed service JWTs with:

- `sub`: calling service name,
- `aud`: target service name,
- `exp`: short expiry,
- `kid`: signing key ID.

Plain user JWTs are rejected on internal endpoints unless a contract explicitly allows them.
