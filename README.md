# 一个便捷带图形界面的frp客户端


## build

```
git clone https://github.com/GavinTan/wfrpc.git

# python3.9及以上版本不支持win7，win7环境运行需要使用python3.9以下版本
# mac下pyqt5 5.15.2以上版本pyinstaller打包无法运行
pip install pyinstaller requests PyQt5==5.15.2

cd wfrpc

# windows
pyinstaller -Fw --noupx --clean -y -i .\icon\2.ico --add-data ".\\ui\\*;ui" --add-data ".\\icon\*;icon" wfrpc.py

# mac
pyinstaller -Fw --noupx --clean -y -i ./icon/2.icns --add-data "./ui/*:ui" --add-data "./icon/*:icon" wfrpc.py
```
