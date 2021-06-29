# 一个便捷的frp客户端（UI）


## build

```
git clone https://github.com/GavinTan/wfrpc.git

pip install pyinstaller requests PyQt5

cd wfrpc

pyinstaller -Fw --noupx -i .\icon\2.ico --add-data ".\\ui\\*;ui" --add-data ".\\icon\*;icon"  wfrpc.py
```
