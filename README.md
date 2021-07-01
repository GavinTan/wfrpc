# 一个便捷的frp客户端


## build

```
git clone https://github.com/GavinTan/wfrpc.git

pip install pyinstaller requests PyQt5

cd wfrpc

# windows
pyinstaller -Fw --noupx --clean -y -i .\icon\2.ico --add-data ".\\ui\\*;ui" --add-data ".\\icon\*;icon" wfrpc.py

# mac
pyinstaller -Fw --noupx --clean -y -i ./icon/2.icns --add-data "./ui/*:ui" --add-data "./icon/*:icon" wfrpc.py
```
