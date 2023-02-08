# Build Status

## Dev
| Step | Status |
| -- | -- |
| Image build | [![Build Status](https://github-private.visualstudio.com/azure/_apis/build/status/Azure.aad-auth-proxy?branchName=main&jobName=Build%20image)](https://github-private.visualstudio.com/azure/_build/latest?definitionId=449&branchName=main) |
| Helm chart | [![Build Status](https://github-private.visualstudio.com/azure/_apis/build/status/Azure.aad-auth-proxy?branchName=main&jobName=Package%20helm%20chart)](https://github-private.visualstudio.com/azure/_build/latest?definitionId=440&branchName=main)

# Project

This project is Azure AD proxy, which can used as a forward proxy to authenticate requests to backend services (example: Azure Monitor workspaces). Clients can use system identity, or user identity (kublet identity) or Azure AD app to fetch tokens whihc will be added as bearer tokens to all forwarded requests.

# Limitations
Only helm v3 is supported.
