@echo off
echo "Gymdl Restart"
taskkill /f /t /fi "imagename eq gymdl-windows-amd64.exe"
echo "Gymdl Stop"
timeout /t 1
exit