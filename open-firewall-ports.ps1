netsh advfirewall firewall delete rule name="Vite 5196" >$null 2>&1
netsh advfirewall firewall delete rule name="Backend 8082" >$null 2>&1
netsh advfirewall firewall add rule name="Vite 5196" dir=in action=allow protocol=TCP localport=5196
netsh advfirewall firewall add rule name="Backend 8082" dir=in action=allow protocol=TCP localport=8082
Write-Host "Done! Ports 5196 and 8082 are now open." -ForegroundColor Green
pause
