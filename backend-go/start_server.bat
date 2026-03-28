@echo off
chcp 65001 >nul
echo ========================================
echo   VaelorCMS - 完整内容管理系统
echo   版本: 1.0.0
echo   作者: Tinmc189623
echo   团队: Nexlyh
echo ========================================
echo.
echo 正在启动服务器...
cd /d "%~dp0"
go run full_server.go
pause
