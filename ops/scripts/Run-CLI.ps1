param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$CliArgs
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "Set-DevEnvironment.ps1")

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$CliPath = Join-Path $RepoRoot "src\cli"

Push-Location $CliPath
try {
    go run . @CliArgs
}
finally {
    Pop-Location
}
