@echo off
echo "Gymdl Start"
taskkill /f /t /fi "imagename eq gymdl-windows-amd64.exe"
echo "Gymdl Stop"
cscript gymdl_no_window.vbs
echo "Gymdl Start"
timeout /t 1
exit
