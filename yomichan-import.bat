@echo off
set /p dict_path="Specify dictionary path: "
yomichan-import.exe %dict_path%
pause
