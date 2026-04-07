Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if ($PSVersionTable.PSVersion.Major -ge 7) {
    $PSNativeCommandUseErrorActionPreference = $true
}

$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..\..")).Path
$originalHome = [Environment]::GetEnvironmentVariable("HOME")
$originalUserProfile = [Environment]::GetEnvironmentVariable("USERPROFILE")

function Import-EnvFile {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    if (-not (Test-Path $Path)) {
        return
    }

    foreach ($line in Get-Content -LiteralPath $Path) {
        $trimmed = $line.Trim()
        if ($trimmed.Length -eq 0 -or $trimmed.StartsWith("#")) {
            continue
        }

        if ($trimmed.StartsWith("export ")) {
            $trimmed = $trimmed.Substring(7).Trim()
        }

        $separatorIndex = $trimmed.IndexOf("=")
        if ($separatorIndex -lt 1) {
            continue
        }

        $key = $trimmed.Substring(0, $separatorIndex).Trim()
        $value = $trimmed.Substring($separatorIndex + 1).Trim()

        if (($value.StartsWith('"') -and $value.EndsWith('"')) -or ($value.StartsWith("'") -and $value.EndsWith("'"))) {
            $value = $value.Substring(1, $value.Length - 2)
        }

        if ([string]::IsNullOrWhiteSpace($key)) {
            continue
        }

        if ([string]::IsNullOrEmpty([Environment]::GetEnvironmentVariable($key))) {
            Set-Item -Path "Env:$key" -Value $value
        }
    }
}

$paths = @{
    DotnetHome = Join-Path $RepoRoot ".dotnet"
    FunctionsHome = Join-Path $RepoRoot ".azurefunctions"
    NuGetPackages = Join-Path $RepoRoot ".nuget\packages"
    GoCache = Join-Path $RepoRoot ".gocache"
    GoModCache = Join-Path $RepoRoot ".gomodcache"
}

foreach ($path in $paths.Values) {
    New-Item -ItemType Directory -Force -Path $path | Out-Null
}

Import-EnvFile -Path (Join-Path $RepoRoot "infra\env\.env")
Import-EnvFile -Path (Join-Path $RepoRoot ".env")

$env:DOTNET_CLI_HOME = $paths.DotnetHome
$env:NUGET_PACKAGES = $paths.NuGetPackages
$env:DOTNET_SKIP_FIRST_TIME_EXPERIENCE = "1"
$env:DOTNET_CLI_TELEMETRY_OPTOUT = "1"
$env:DOTNET_ADD_GLOBAL_TOOLS_TO_PATH = "0"
$env:DOTNET_GENERATE_ASPNET_CERTIFICATE = "false"
$env:DOTNET_NOLOGO = "1"
$env:MSBuildEnableWorkloadResolver = "false"
$env:FUNCTIONS_CORE_TOOLS_TELEMETRY_OPTOUT = "1"
$env:GOCACHE = $paths.GoCache
$env:GOMODCACHE = $paths.GoModCache
$env:GOFLAGS = "-buildvcs=false"

if ([string]::IsNullOrWhiteSpace($env:AZURE_CONFIG_DIR)) {
    $azureConfigRoot = $null

    if (-not [string]::IsNullOrWhiteSpace($originalUserProfile)) {
        $azureConfigRoot = Join-Path $originalUserProfile ".azure"
    }
    elseif (-not [string]::IsNullOrWhiteSpace($originalHome)) {
        $azureConfigRoot = Join-Path $originalHome ".azure"
    }

    if (-not [string]::IsNullOrWhiteSpace($azureConfigRoot)) {
        $env:AZURE_CONFIG_DIR = $azureConfigRoot
    }
}

$isSandboxProfile = ($originalUserProfile -like "*CodexSandboxOffline*") -or ($env:USERNAME -eq "CodexSandboxOffline")
if ($isSandboxProfile) {
    $env:HOME = $paths.DotnetHome
    $env:USERPROFILE = $paths.FunctionsHome
}

$aliasMap = @{
    AZOPS_TENANT_ID = "AZURE_TENANT_ID"
    AZOPS_CLIENT_ID = "AZURE_CLIENT_ID"
    AZOPS_CLIENT_SECRET = "AZURE_CLIENT_SECRET"
    AZOPS_CLIENT_CERTIFICATE_PATH = "AZURE_CLIENT_CERTIFICATE_PATH"
    AZOPS_SUBSCRIPTION_ID = "AZURE_SUBSCRIPTION_ID"
}

foreach ($sourceName in $aliasMap.Keys) {
    $targetName = $aliasMap[$sourceName]
    $sourceValue = [Environment]::GetEnvironmentVariable($sourceName)
    $targetValue = [Environment]::GetEnvironmentVariable($targetName)

    if (-not [string]::IsNullOrWhiteSpace($sourceValue) -and [string]::IsNullOrWhiteSpace($targetValue)) {
        Set-Item -Path "Env:$targetName" -Value $sourceValue
    }
}

Write-Host "Repo development environment configured."
