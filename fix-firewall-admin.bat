@echo off
echo === Fixing Windows Firewall for ZeroTier ===

REM Remove old rules
netsh advfirewall firewall delete rule name="Vite 5196" >nul 2>&1
netsh advfirewall firewall delete rule name="Backend 8082" >nul 2>&1

REM Add rules for ALL profiles (domain, private, public)
netsh advfirewall firewall add rule name="Vite 5196" dir=in action=allow protocol=TCP localport=5196 profile=any
netsh advfirewall firewall add rule name="Backend 8082" dir=in action=allow protocol=TCP localport=8082 profile=any

REM Set ZeroTier adapter to Private network profile via PowerShell
powershell -Command "$a = Get-NetAdapter | Where-Object { $_.InterfaceDescription -like '*ZeroTier*' }; if ($a) { Set-NetConnectionProfile -InterfaceIndex $a.ifIndex -NetworkCategory Private; Write-Host 'ZeroTier set to Private network' } else { Write-Host 'ZeroTier adapter not found' }"

echo.
echo === Done! ===
echo Ports 5196 and 8082 are now open for ALL network profiles.
pause
