# 一个便捷的frp客户端


## build

```
git clone https://github.com/GavinTan/wfrpc.git

# mac下pyqt5 5.15.2以上版本pyinstaller打包无法运行
pip install pyinstaller requests PyQt5==5.15.2

cd wfrpc

# windows
pyinstaller -Fw --noupx --clean -y -i .\icon\2.ico --add-data ".\\ui\\*;ui" --add-data ".\\icon\*;icon" wfrpc.py

# mac
pyinstaller -Fw --noupx --clean -y -i ./icon/2.icns --add-data "./ui/*:ui" --add-data "./icon/*:icon" wfrpc.py
```
