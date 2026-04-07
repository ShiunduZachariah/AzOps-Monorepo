Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot "Set-DevEnvironment.ps1")

function Invoke-CheckedCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,

        [Parameter()]
        [string[]]$Arguments = @()
    )

    & $FilePath @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed: $FilePath $($Arguments -join ' ')"
    }
}

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path

Push-Location (Join-Path $RepoRoot "src\cli")
try {
    Invoke-CheckedCommand -FilePath "go" -Arguments @("build", "./...")
}
finally {
    Pop-Location
}

Push-Location (Join-Path $RepoRoot "src\functions")
try {
    Invoke-CheckedCommand -FilePath "dotnet" -Arguments @("build", ".\AzOps.Functions.sln", "--nologo", "--no-restore")
}
finally {
    Pop-Location
}
