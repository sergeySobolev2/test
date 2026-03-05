@echo off
echo === TEMPORARILY DISABLING FIREWALL ===
netsh advfirewall set allprofiles state off
echo.
echo Firewall is OFF. Test from phone now.
echo Press any key to TURN FIREWALL BACK ON...
pause
netsh advfirewall set allprofiles state on
echo Firewall is back ON.
pause
