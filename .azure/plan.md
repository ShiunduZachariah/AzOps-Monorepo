# Azure Preparation Plan

## Status

Approved for Execution

## Scenario

Implement Sprint 3 of the AzOps monorepo by extending the existing .NET isolated Azure Functions app with:

- thin HTTP-triggered endpoints
- Core and Infrastructure feature separation
- Key Vault secret retrieval through Azure SDK clients
- `DefaultAzureCredential`-based authentication that works locally and with managed identity in Azure
- integration-focused tests for the new function flows

## Workspace Mode

MODIFY

## Scope

- Update `src/functions/src/AzOps.Core`
- Update `src/functions/src/AzOps.Infrastructure`
- Update `src/functions/src/AzOps.Functions`
- Add integration-focused tests under `src/functions/tests`
- Update local configuration and README guidance as needed

## Recipe

Bicep-backed Azure Functions application with code-first Function changes in this sprint. Deployment infrastructure remains incremental, while application code is prepared for managed identity and Key Vault access.

## Architecture Decisions

- Keep HTTP triggers thin and feature-based under `AzOps.Functions/Features`
- Keep contracts, use cases, and response models in `AzOps.Core`
- Keep Azure SDK clients, configuration, and Key Vault/resource access in `AzOps.Infrastructure`
- Use `DefaultAzureCredential` so local development can use Azure CLI credentials and Azure deployment can use managed identity without code changes
- Add structured JSON error responses for predictable API behavior

## Deliverables

- `/api/health`
- `/api/ping`
- `/api/resources/resource-groups`
- `/api/secrets/{secretName}`
- Key Vault-backed secret retrieval flow
- resource-group sample operation through the Azure SDK
- integration-focused tests for success and failure paths

## Validation Plan

- `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Restore-All.ps1`
- `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Build-All.ps1`
- `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Test-All.ps1`
- `powershell -ExecutionPolicy Bypass -File .\ops\scripts\Run-Functions.ps1`
- `Invoke-WebRequest http://localhost:7071/api/health`
- `Invoke-WebRequest http://localhost:7071/api/ping`

