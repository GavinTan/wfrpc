pyinstaller -Fw --noupx -i .\icon\2.ico --add-data ".\\ui\\*;ui" --add-data ".\\icon\*;icon"  wfrpc.py