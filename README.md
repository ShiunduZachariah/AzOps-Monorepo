# AzOps Monorepo

Sprint 2 of the Azure operations monorepo is now in place. The repo contains a Go CLI for Azure resource-group and virtual-machine operations, an isolated .NET Azure Functions app with ASP.NET Core integration, repo-local restore/build/test scripts, and the top-level folders needed for later infra and observability sprints.

## Current scope

- `src/cli`: Cobra-based Go CLI with `health`, `groups list`, `vm list`, shared response models, script-friendly plain/JSON output, and Azure authentication through Azure CLI or service principal settings.
- `src/functions`: .NET isolated Azure Functions solution containing `AzOps.Functions`, `AzOps.Core`, `AzOps.Infrastructure`, and `GET /api/health`.
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

The Functions app includes [`local.settings.json.example`](/c:/Users/ZacH/Documents/Personal-Projects/AzOps-Monorepo/src/functions/src/AzOps.Functions/local.settings.json.example) with the minimum local values for the current Function host.

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

Expected results:
`health` should print `ok`, the automated test script should pass, `groups list` and `vm list` should return Azure data once auth is available, and `/api/health` should return a healthy JSON payload.

## Sprint 2 verification

The current repo has been verified with:

- `Restore-All.ps1`
- `Build-All.ps1`
- `Test-All.ps1`
- `Run-CLI.ps1 health`
- `Run-CLI.ps1 groups list --output json`
- `Run-CLI.ps1 vm list --output json`

## Next sprint

Sprint 3 can build on this foundation with Function endpoints beyond health, Key Vault integration, and the managed-identity flow.
