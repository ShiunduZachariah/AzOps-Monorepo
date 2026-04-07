param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$FunctionArgs
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "Set-DevEnvironment.ps1")

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$FunctionsPath = Join-Path $RepoRoot "src\functions\src\AzOps.Functions"

Push-Location $FunctionsPath
try {
    func start @FunctionArgs
}
finally {
    Pop-Location
}
