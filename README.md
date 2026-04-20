# AzOps Monorepo

Sprint 3 of the Azure operations monorepo is now in place. The repo contains a Go CLI for Azure resource-group and virtual-machine operations, an isolated .NET Azure Functions app with feature-based Core and Infrastructure slices, Key Vault-ready secret retrieval, managed-identity-friendly authentication, repo-local restore/build/test scripts, and the top-level folders needed for later infra and observability sprints.

## Current scope

- `src/cli`: Cobra-based Go CLI with `health`, `groups list`, `vm list`, shared response models, script-friendly plain/JSON output, and Azure authentication through Azure CLI or service principal settings.
- `src/functions`: .NET isolated Azure Functions solution containing `AzOps.Functions`, `AzOps.Core`, `AzOps.Infrastructure`, and the HTTP endpoints `GET /api/health`, `GET /api/ping`, `GET /api/resources/resource-groups`, and `GET /api/secrets/{secretName}`.
- `ops/scripts`: PowerShell helpers for restore, build, test, CLI run, and Functions run
- `infra`: placeholder Bicep and environment assets for later sprints
- `.github/workflows/ci.yml`: Windows CI for restore, build, and test

## Repository layout

```text
.
|-- src
|   |-- cli
|   |   |-- cmd
|   |   |-- internal
|   |   |-- third_party
|   |   |-- go.mod
|   |   `-- main.go
|   `-- functions
|       |-- AzOps.Functions.sln
|       |-- src
|       |   |-- AzOps.Functions
|       |   |-- AzOps.Core
|       |   `-- AzOps.Infrastructure
|       `-- tests
|           `-- AzOps.Functions.Tests
|-- infra
|   |-- bicep
|   `-- env
|-- ops
|   |-- docker
|   `-- scripts
`-- .github
    `-- workflows
```

## Prerequisites

- Go 1.26+
- .NET SDK 10.0.x
- Azure Functions Core Tools 4.x
- Azure access through one of the `DefaultAzureCredential` flows

## Environment

The repo automatically loads `infra/env/.env` for the PowerShell scripts and the Go CLI.

The CLI reads:

- `AZOPS_SUBSCRIPTION_ID` or `AZURE_SUBSCRIPTION_ID`
- `AZOPS_TENANT_ID` or `AZURE_TENANT_ID`
- `AZOPS_CLIENT_ID` or `AZURE_CLIENT_ID`
- `AZOPS_CLIENT_SECRET` or `AZURE_CLIENT_SECRET`
- `AZOPS_AUTH_MODE`
- `AZOPS_OUTPUT`

For Azure authentication, the CLI supports either:

- `az login`
- service principal environment variables from `infra/env/.env`

Set `AZOPS_AUTH_MODE=service-principal` when you want to require client-secret auth. Leave it as `auto` to prefer service principal credentials when present and otherwise fall back to Azure CLI, managed identity, and Azure Developer CLI.

The Functions app reads these Azure-facing settings:

- `AZOPS_SUBSCRIPTION_ID` for the `resources/resource-groups` endpoint
- `AZOPS_KEY_VAULT_NAME` or `AZOPS_KEY_VAULT_URI` for the secrets endpoint

The Functions app uses `DefaultAzureCredential`, which means local runs can use `az login` and Azure deployments can use the Function App's managed identity without code changes. The Functions app includes [`local.settings.json.example`](/c:/Users/ZacH/Documents/Personal-Projects/AzOps-Monorepo/src/functions/src/AzOps.Functions/local.settings.json.example) with the minimum local values for the current Function host.

For the secrets endpoint to work in Azure, enable a managed identity on the Function App and grant that identity permission to read secrets from the target Key Vault.

## Go module path

The CLI module now uses the GitHub-style path:

```text
github.com/ShiunduZachariah/azopscli
```

## Common commands

```powershell
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Restore-All.ps1
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Build-All.ps1
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Test-All.ps1
```

```powershell
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 health
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 groups list --subscription-id <subscription-id>
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 vm list --subscription-id <subscription-id>
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 groups list --subscription-id <subscription-id> --output json
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 vm list --subscription-id <subscription-id> --output json
```

```powershell
powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-Functions.ps1
```

When the Functions host is running locally, the current endpoint is:

- `GET http://localhost:7071/api/health`
- `GET http://localhost:7071/api/ping`
- `GET http://localhost:7071/api/resources/resource-groups`
- `GET http://localhost:7071/api/secrets/{secretName}`

## How To Test

1. Restore everything:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Restore-All.ps1`
2. Build everything:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Build-All.ps1`
3. Run all automated tests:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Test-All.ps1`
4. Smoke-test the CLI wiring:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 health`
5. Test the real Azure CLI flow after `az login` or setting service principal env vars:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 groups list --subscription-id <your-subscription-id> --output json`
6. Test the VM command the same way:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-CLI.ps1 vm list --subscription-id <your-subscription-id> --output json`
7. Run the Functions app locally:
   `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-Functions.ps1`
8. In another terminal, test the health endpoint:
   `Invoke-WebRequest http://localhost:7071/api/health`
9. Test the ping endpoint:
   `Invoke-WebRequest http://localhost:7071/api/ping`
10. Test the resource-groups endpoint:
   `Invoke-WebRequest http://localhost:7071/api/resources/resource-groups`
11. Test the secret endpoint after setting `AZOPS_KEY_VAULT_NAME` or `AZOPS_KEY_VAULT_URI`:
   `Invoke-WebRequest http://localhost:7071/api/secrets/<your-secret-name>`

Expected results:
`health` should print `ok`, the automated test script should pass, `groups list` and `vm list` should return Azure data once auth is available, `/api/health` and `/api/ping` should return healthy JSON payloads, `/api/resources/resource-groups` should return a JSON list for the configured subscription, and `/api/secrets/<secretName>` should return a masked Key Vault-backed result or a structured JSON configuration error.

## Sprint 3 verification

The current repo has been verified with:

- `Restore-All.ps1`
- `Build-All.ps1`
- `Test-All.ps1`
- `Run-CLI.ps1 health`
- `Run-CLI.ps1 groups list --output json`
- `Run-CLI.ps1 vm list --output json`
- `Run-Functions.ps1`
- `Invoke-WebRequest http://localhost:7071/api/health`
- `Invoke-WebRequest http://localhost:7071/api/ping`
- `Invoke-WebRequest http://localhost:7071/api/resources/resource-groups`

## Next sprint

Sprint 4 can build on this foundation with metrics, Prometheus, Grafana, and logging/observability cleanup.
