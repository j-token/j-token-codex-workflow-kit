param(
    [string]$Version = "latest",
    [string]$InstallDir = ""
)

$ErrorActionPreference = "Stop"
$repository = "j-token/j-token-codex-workflow-kit"

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    $InstallDir = Join-Path $env:LOCALAPPDATA "codex-workflow\bin"
}

$architecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLowerInvariant()
switch ($architecture) {
    "x64" { $assetArchitecture = "x64" }
    "arm64" { $assetArchitecture = "arm64" }
    default { throw "지원하지 않는 아키텍처입니다: $architecture" }
}

$asset = "codex-workflow-windows-$assetArchitecture.exe"
if ($Version -eq "latest") {
    $baseUrl = "https://github.com/$repository/releases/latest/download"
} else {
    if (-not $Version.StartsWith("v")) {
        $Version = "v$Version"
    }
    $baseUrl = "https://github.com/$repository/releases/download/$Version"
}

$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("codex-workflow-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir | Out-Null
$stagedPath = $null
$backupPath = $null

try {
    $binaryPath = Join-Path $tempDir $asset
    $checksumPath = "$binaryPath.sha256"
    Invoke-WebRequest -UseBasicParsing -Uri "$baseUrl/$asset" -OutFile $binaryPath
    Invoke-WebRequest -UseBasicParsing -Uri "$baseUrl/$asset.sha256" -OutFile $checksumPath

    $expected = ((Get-Content -Raw -Encoding utf8 -LiteralPath $checksumPath).Trim() -split "\s+")[0].ToLowerInvariant()
    $actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $binaryPath).Hash.ToLowerInvariant()
    if ($expected -ne $actual) {
        throw "SHA-256 검증에 실패했습니다"
    }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $destination = Join-Path $InstallDir "codex-workflow.exe"
    $operationId = [guid]::NewGuid().ToString("N")
    $stagedPath = "$destination.new.$operationId"
    $backupPath = "$destination.bak.$operationId"
    Copy-Item -LiteralPath $binaryPath -Destination $stagedPath
    if (Test-Path -LiteralPath $destination) {
        [System.IO.File]::Replace($stagedPath, $destination, $backupPath, $true)
        Remove-Item -Force -LiteralPath $backupPath
        $backupPath = $null
    } else {
        [System.IO.File]::Move($stagedPath, $destination)
    }
    $stagedPath = $null

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    $pathParts = @($userPath -split ";" | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
    if ($pathParts -notcontains $InstallDir) {
        $newPath = (($pathParts + $InstallDir) -join ";")
        [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
        $env:Path = "$env:Path;$InstallDir"
        Write-Host "사용자 PATH에 추가했습니다: $InstallDir"
    }

    Write-Host "설치 완료: $destination"
} finally {
    if ($null -ne $stagedPath -and (Test-Path -LiteralPath $stagedPath)) {
        Remove-Item -Force -LiteralPath $stagedPath
    }
    if ($null -ne $backupPath -and (Test-Path -LiteralPath $backupPath)) {
        Remove-Item -Force -LiteralPath $backupPath
    }
    if (Test-Path -LiteralPath $tempDir) {
        Remove-Item -Recurse -Force -LiteralPath $tempDir
    }
}
