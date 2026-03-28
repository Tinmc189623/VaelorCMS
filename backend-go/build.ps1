# VaelorCMS - 多平台编译脚本 (PowerShell)
# 支持 Windows、Linux、macOS 编译

param(
    [string]$OutputDir = "bin",
    [switch]$Clean,
    [switch]$All
)

$ErrorActionPreference = "Stop"

# 项目配置
$ProjectName = "vaelorcms"
$MainPackage = "cmd/server"

# 目标平台配置
$Targets = @(
    @{ OS = "windows"; Arch = "amd64"; Ext = ".exe" },
    @{ OS = "windows"; Arch = "386"; Ext = ".exe" },
    @{ OS = "linux"; Arch = "amd64"; Ext = "" },
    @{ OS = "linux"; Arch = "386"; Ext = "" },
    @{ OS = "linux"; Arch = "arm64"; Ext = "" },
    @{ OS = "linux"; Arch = "arm"; Ext = "" },
    @{ OS = "darwin"; Arch = "amd64"; Ext = "" },
    @{ OS = "darwin"; Arch = "arm64"; Ext = "" }
)

# 清理输出目录
if ($Clean -and (Test-Path $OutputDir)) {
    Write-Host "正在清理输出目录: $OutputDir" -ForegroundColor Yellow
    Remove-Item -Path $OutputDir -Recurse -Force
}

# 创建输出目录
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# 获取当前平台
$CurrentOS = "windows"
if ($PSVersionTable.Platform -eq "Unix") {
    if ($IsMacOS) {
        $CurrentOS = "darwin"
    } else {
        $CurrentOS = "linux"
    }
}

$CurrentArch = [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture.ToString().ToLower()
if ($CurrentArch -eq "x64") { $CurrentArch = "amd64" }
if ($CurrentArch -eq "x86") { $CurrentArch = "386" }

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  VaelorCMS - 多平台编译脚本" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 如果没有指定 -All，只编译当前平台
if (-not $All) {
    $FilteredTargets = @()
    foreach ($Target in $Targets) {
        if ($Target.OS -eq $CurrentOS -and $Target.Arch -eq $CurrentArch) {
            $FilteredTargets += $Target
        }
    }
    $Targets = $FilteredTargets
    Write-Host "仅编译当前平台: $CurrentOS/$CurrentArch" -ForegroundColor Green
} else {
    Write-Host "编译所有平台..." -ForegroundColor Green
}
Write-Host ""

$SuccessCount = 0
$FailCount = 0

foreach ($Target in $Targets) {
    $OS = $Target.OS
    $Arch = $Target.Arch
    $Ext = $Target.Ext
    
    $OutputFile = "$OutputDir/$ProjectName-$OS-$Arch$Ext"
    Write-Host "正在编译: $OS/$Arch" -ForegroundColor Cyan
    
    try {
        $env:GOOS = $OS
        $env:GOARCH = $Arch
        $env:CGO_ENABLED = "0"
        
        $Start = Get-Date
        go build -ldflags "-s -w" -o $OutputFile $MainPackage
        $End = Get-Date
        $Duration = ($End - $Start).TotalSeconds.ToString("F2")
        
        if (Test-Path $OutputFile) {
            $FileSize = (Get-Item $OutputFile).Length / 1KB
            Write-Host "  ✓ 成功! $([Math]::Round($FileSize, 2)) KB, 耗时 ${Duration}s" -ForegroundColor Green
            $SuccessCount++
        } else {
            Write-Host "  ✗ 失败: 输出文件未生成" -ForegroundColor Red
            $FailCount++
        }
    } catch {
        Write-Host "  ✗ 失败: $_" -ForegroundColor Red
        $FailCount++
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  编译完成!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  成功: $SuccessCount" -ForegroundColor Green
if ($FailCount -gt 0) {
    Write-Host "  失败: $FailCount" -ForegroundColor Red
} else {
    Write-Host "  失败: $FailCount" -ForegroundColor Gray
}
Write-Host "  输出目录: $(Resolve-Path $OutputDir)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

if ($FailCount -gt 0) {
    exit 1
}
