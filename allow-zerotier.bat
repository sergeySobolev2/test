@echo off
echo Adding ZeroTier subnet allow rule...
netsh advfirewall firewall add rule name="ZeroTier Allow All Inbound" dir=in action=allow remoteip=10.57.0.0/16 profile=any
echo.
echo Result: %ERRORLEVEL%
echo.
netsh advfirewall firewall show rule name="ZeroTier Allow All Inbound"
echo.
pause
