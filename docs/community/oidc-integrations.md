---
layout: default
title: Community-Tested OIDC Integrations
parent: Community
nav_order: 4
---

# OIDC Integratins

**Note** This is community-based content for which the core-maintainers cannot guarantee correctness. The parameters may change over time. If a parameter does not work as documented, please submit a PR to update the list.

## Currently Tested Applications

| Application | Minimal Version                | Notes   |
| :---------: | :----------------------------: | :-----: |
| GitLab      | `13.0.0`                       | |
| Grafana     | `8.0.5`                        | |
| MinIO       | `RELEASE.2021-07-12T02-44-53Z` | must set `MINIO_IDENTITY_OPENID_CLAIM_NAME: groups` in MinIO and set [MinIO policies] as groups in Authelia |

[MinIO policies]: https://docs.min.io/minio/baremetal/security/minio-identity-management/policy-based-access-control.html#minio-policy

## Known Callback URLs

If you do not find the application in the list below, you will need to search for yourself - and maybe come back to open a PR to add your application to this list so others won't have to search for them.

`<DOMAIN>` needs to be substituted with the full URL on which the application runs on. If GitLab, as an example, was reachable under `https://gitlab.example.com`, `<DOMAIN>` would be exactly the same.

| Application | Version                        | Callback URL                                             |
| :---------: | :----------------------------: | :------------------------------------------------------: |
| GitLab      | `14.0.1`                       | `<DOMAIN>/users/auth/openid_connect/callback`            |
| MinIO       | `RELEASE.2021-07-12T02-44-53Z` | `<DOMAIN>/oauth_callback`                                |
