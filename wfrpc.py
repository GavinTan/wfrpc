import sys
from PyQt5.QtGui import QIcon, QDesktopServices
from PyQt5.QtWidgets import *
from PyQt5.uic import loadUi
from PyQt5 import QtCore
import configparser
import subprocess
import platform
import random
import requests
import tarfile
import shutil
import tempfile
import zipfile
import os
import time


class Thread(QtCore.QThread):
    trigger = QtCore.pyqtSignal(str)

    sys_name = platform.system().lower()
    sys_machine = platform.machine().lower()
    sys_arch_list = {'x86_64': 'amd64', 'aarch64': 'arm64', 'x86': '386'}
    sys_arch = sys_arch_list.get(sys_machine) or sys_machine

    frpc_base_path = os.path.join(os.path.expanduser('~'), '.wfrpc')
    frpc_bin_path = os.path.join(frpc_base_path, 'frpc.exe')
    frpc_config_path = os.path.join(frpc_base_path, 'frpc.ini')

    frp_api_github_url = 'https://api.github.com/repos/fatedier/frp/releases/latest'
    frp_release_github_url = 'https://github.com/fatedier/frp/releases'

    run_type = None
    run_cmd = None

    def __init__(self, run_type):
        super().__init__()
        self.run_type = run_type

    def run(self):
        if self.run_type == 'download':
            self.download_frp()

        if self.run_type == 'run':
            try:
                self.run_cmd = subprocess.Popen([self.frpc_bin_path, '-c', self.frpc_config_path],
                 stdout=subprocess.PIPE, stderr=subprocess.PIPE, stdin=subprocess.DEVNULL, shell=True)
                while self.run_cmd.poll() is None:
                    output = self.run_cmd.stdout.readline().decode()
                    if output.strip():
                        self.trigger.emit(output.strip())
            except Exception as e:
                log = open(os.path.join(self.frpc_base_path, 'wfrpc.log'), 'a')
                print(
                    f"{time.strftime('%Y-%m-%d %H:%M:%S')} {e}",
                    file=log)
                log.close()


    def download_frp(self):
        if not os.path.exists(self.frpc_bin_path):
            os.makedirs(self.frpc_base_path, exist_ok=True)
            try:
                log = open(os.path.join(self.frpc_base_path, 'wfrpc.log'), 'a')
                r = requests.get(self.frp_api_github_url)
                tmp_path = tempfile.gettempdir()
                for i in r.json().get('assets'):
                    self.trigger.emit('download_start')
                    frp_download_url = i.get('browser_download_url')
                    if f'{self.sys_name}_{self.sys_arch}' in frp_download_url:
                        print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} download {frp_download_url}", file=log)
                        file_name = os.path.basename(frp_download_url)
                        file_path = os.path.join(tmp_path, file_name)
                        frp_download = requests.get(frp_download_url)
                        with open(file_path, 'wb') as f:
                            f.write(frp_download.content)

                        if file_name.endswith('.tar.gz'):
                            with tarfile.open(file_path, "r:gz") as f:
                                f.extractall(tmp_path)

                            frp_exe_tmp_path = os.path.join(file_path.split('.tar.gz')[0], 'frpc.exe')
                            shutil.copy2(frp_exe_tmp_path, self.frpc_base_path)

                        if file_name.endswith('.zip'):
                            with zipfile.ZipFile(file_path) as f:
                                f.extractall(tmp_path)

                            frp_exe_tmp_path = os.path.join(file_path.split('.zip')[0], 'frpc.exe')
                            shutil.copy2(frp_exe_tmp_path, self.frpc_base_path)
                        self.trigger.emit('download_complete')
                        return True
                self.trigger.emit(f"当前系统架构 {self.sys_name}_{self.sys_machine} 没有找到匹配的程序包！\
                    请手动前往<a href=\"{self.frp_release_github_url}\">Github</a>下载程序，并把frpc.exe解压到<a href=\"file:///{self.frpc_base_path}\">软件目录</a>下。")

                print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} 当前系统架构 {self.sys_name}_{self.sys_machine} 没有找到匹配的程序包！请手动下载 {self.frp_release_github_url}。",
                    file=log)
                log.close()
            except Exception as e:
                print(
                    f"{time.strftime('%Y-%m-%d %H:%M:%S')} {e}", file=log)
                self.trigger.emit(f'下载frp出错！请手动前往<a href="{self.frp_release_github_url}">Github</a>下载程序，并把frpc.exe解压到<a href="file:///{self.frpc_base_path}">软件目录</a>下。')



class MainWindow(QMainWindow):

    frpc_base_path = os.path.join(os.path.expanduser('~'), '.wfrpc')
    frpc_bin_path = os.path.join(frpc_base_path, 'frpc.exe')
    frpc_config_path = os.path.join(frpc_base_path, 'frpc.ini')

    frpc_config_server_addr = 'nb33.3322.org'
    frpc_config_server_port = '7000'
    frpc_config_server_passwd = ''

    frpc_config_list = []

    def __init__(self, parent=None):
        super(MainWindow, self).__init__(parent)
        if getattr(sys, 'frozen', False):
            # we are running in a bundle
            bundle_dir = sys._MEIPASS
        else:
            # we are running in a normal Python environment
            bundle_dir = os.path.dirname(os.path.abspath(__file__))

        loadUi(bundle_dir + '/ui/main.ui', self)
        self.add_config_dialog = loadUi( bundle_dir + '/ui/add_config_dialog.ui')

        self.add_config_dialog.setWindowIcon(QIcon(bundle_dir + '/icon/3.ico'))
        self.setWindowIcon(QIcon(bundle_dir + '/icon/3.ico'))
        self.tray_icon = QSystemTrayIcon(self)
        # self.tray_icon.setIcon(self.style().standardIcon(QStyle.SP_ComputerIcon))
        self.tray_icon.setIcon(QIcon(bundle_dir + '/icon/2.ico'))
        self.tray_icon.setToolTip('wfrpc')

        show_action = QAction("打开", self)
        quit_action = QAction("退出", self)
        show_action.triggered.connect(self.showNormal)
        quit_action.triggered.connect(lambda: (self.stop_frpc(), qApp.quit()))
        tray_menu = QMenu()
        tray_menu.addAction(show_action)
        tray_menu.addAction(quit_action)
        self.tray_icon.setContextMenu(tray_menu)
        self.tray_icon.show()
        self.tray_icon.activated[QSystemTrayIcon.ActivationReason].connect(self.iconActivated)

        self.statusbar.showMessage('树叶的一生难道只是为了归根吗？')

        def handle_links(url):
            if not url.scheme():
                url = QtCore.QUrl.fromLocalFile(url.toString())
            QDesktopServices.openUrl(url)
        self.show_log.setOpenLinks(False)
        self.show_log.anchorClicked.connect(handle_links)
        # self.show_log.setOpenExternalLinks(True)

        self.btn_start_frpc.clicked.connect(self.start_frpc)
        self.btn_stop_frpc.clicked.connect(lambda: self.stop_frpc(show_msg=True))
        self.btn_clear_config.clicked.connect(self.clear_frpc_config)
        self.btn_add_config.clicked.connect(self.show_add_config_dialog)
        self.btn_run_backgrounder.clicked.connect(self.run_backgrounder)

        self.add_config_dialog.btn_ok.clicked.connect(self.add_frpc_config)

        
        self.run_frpc = Thread(run_type='run')
        self.run_frpc.trigger.connect(lambda msg: self.show_log.append(msg))

        self.download_frpc = Thread(run_type='download')
        self.download_frpc.trigger.connect(self.alert_message)
        self.download_frpc.start()
        
        self.download_failed = False, '正在后台下载frp中。。。'
        self.start_frpc_status = True

    def alert_message(self, msg):
        if msg == 'download_complete':
            if not self.start_frpc_status:
                self.download_failed = False, msg
                QMessageBox.information(self, '提示', 'frp下载完成!')
        elif msg == 'download_start':
            if not self.start_frpc_status:
                self.download_failed = False, '正在后台下载frp中。。。'
        else:
            if self.start_frpc_status:
                self.download_failed = True, msg
                msgbox = QMessageBox()
                msgbox.setIcon(QMessageBox.Critical)
                msgbox.setText(msg)
                # msgbox.setTextInteractionFlags(QtCore.Qt.TextSelectableByMouse)
                msgbox.setWindowTitle('错误')
                msgbox.exec_()
            
    def show_add_config_dialog(self):
        mw = self.width()
        mx = self.x()
        my = self.y()

        self.add_config_dialog.move(mx + mw + 20, my)
        self.add_config_dialog.setWindowFlags(QtCore.Qt.WindowCloseButtonHint)

        self.add_config_dialog.input_local_ip.clear()
        self.add_config_dialog.input_local_port.clear()
        self.add_config_dialog.input_remote_port.clear()
        self.add_config_dialog.select_port_type.setCurrentIndex(0)

        self.add_config_dialog.exec_()

    def add_frpc_config(self):
        local_ip =  self.add_config_dialog.input_local_ip.text()
        local_port = self.add_config_dialog.input_local_port.text()
        remote_port = self.add_config_dialog.input_remote_port.text() or str(random.randint(20000,30000))
        port_type = self.add_config_dialog.select_port_type.currentText()

        if not local_ip:
            QMessageBox.warning(self, '警告', '未填写目标IP')
        elif not local_port:
            QMessageBox.warning(self, '警告', '未填写目标端口')
        elif not remote_port:
            QMessageBox.warning(self, '警告', '未填写转发端口')
        else:
            self.__dict__['add_frpc_config'] = {
                'local_ip': local_ip,
                'local_port': local_port,
                'remote_port': remote_port,
                'port_type': port_type
            }
            if self.run_frpc.run_cmd and self.run_frpc.run_cmd.poll() is None:
                self.stop_frpc()
                self.start_frpc()
            else:
                self.frpc_config()
            QMessageBox.information(self, '提示', '添加配置成功')
            self.add_config_dialog.close()

    def frpc_config(self):
        self.frpc_config_server_addr = self.input_frp_server_addr.text()
        self.frpc_config_server_port = self.input_frp_server_port.text()
        self.frpc_config_server_passwd = self.input_frp_server_passwd.text()
        add_frpc_config = self.__dict__.get('add_frpc_config')
        frpc_config_parser = configparser.ConfigParser()
        frpc_config_parser.read(self.frpc_config_path)

        os.path.exists(self.frpc_config_path) or os.makedirs(self.frpc_base_path, exist_ok=True)
        log = open(os.path.join(self.frpc_base_path, 'wfrpc.log'), 'a')

        if not frpc_config_parser.has_section('common'):
            print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} add frp common configuration", file=log)
            frpc_config_parser.add_section('common')

        frpc_config_parser.set('common', 'server_addr', self.frpc_config_server_addr)
        frpc_config_parser.set('common', 'server_port', self.frpc_config_server_port)
        frpc_config_parser.set('common', 'token', self.frpc_config_server_passwd)

        if add_frpc_config:
            section = f"{add_frpc_config.get('local_ip')}_{add_frpc_config.get('local_port')}"

            if not frpc_config_parser.has_section(section):
                print(f"{time.strftime('%Y-%m-%d %H:%M:%S')} add frp {section} forwarding configuration", file=log)
                frpc_config_parser.add_section(section)

            frpc_config_parser.set(section, 'local_ip', add_frpc_config.get('local_ip'))
            frpc_config_parser.set(section, 'local_port', add_frpc_config.get('local_port'))
            frpc_config_parser.set(section, 'remote_port', add_frpc_config.get('remote_port'))
            frpc_config_parser.set(section, 'type', add_frpc_config.get('port_type'))
                
        frpc_config_parser.write(open(self.frpc_config_path, 'w')) 

        self.frpc_config_list.clear()
        for i in frpc_config_parser.sections():
            if i != 'common':
                self.frpc_config_list.append({k: v for k, v in frpc_config_parser.items(i)})
        
        log.close()

    def start_frpc(self):
        if self.download_failed[0] or not os.path.exists(self.frpc_bin_path):
            self.download_frpc.start()
            self.show_log.append(self.download_failed[1])
            self.start_frpc_status = False
            return False

        if not self.run_frpc.run_cmd or not self.run_frpc.run_cmd.poll() is None:
            self.frpc_config()

            # QMessageBox.information(self, '提示', 'frp已启动')
            
            self.show_log.clear()
            self.show_log.append(f'{" " * 20}+' + '-' * 27 + 'config' + '-' * 27 + '+')
            for i in self.frpc_config_list:
                output = f"{i.get('local_ip')}:{i.get('local_port')} --> {self.frpc_config_server_addr}:{i.get('remote_port')}"
                n = round((60 - len(output)) / 2)
                self.show_log.append(f"{' ' * 20}{' ' * n}{output}{' ' * n}")
            self.show_log.append(f'{" " * 20}+' + '-' * 60 + '+')
            
            self.run_frpc.start()
            self.start_frpc_status = True

    def stop_frpc(self, show_msg=False):
        # os.system(f'TASKKILL /F /PID {self.run_frpc.run_cmd.pid} /T')
        if self.run_frpc.run_cmd and self.run_frpc.run_cmd.poll() is None:
            si = subprocess.STARTUPINFO()
            si.dwFlags |= subprocess.STARTF_USESHOWWINDOW
            subprocess.call(['taskkill', '/F', '/T', '/PID',  str(self.run_frpc.run_cmd.pid)], startupinfo=si)
            self.run_frpc.run_cmd.terminate()
            self.show_log.append('frp已停止')
        if self.run_frpc:
            self.run_frpc.quit()
        
        if show_msg:
            QMessageBox.information(self, '提示', 'frp已关闭')

    def clear_frpc_config(self):
        os.remove(self.frpc_config_path)
        QMessageBox.information(self, '提示', '配置清理完成')
    
    def closeEvent(self, event):
        self.stop_frpc()
    
    def run_backgrounder(self):
        self.hide()
        self.tray_icon.showMessage(
            "提示",
            "wfrpc已在后台运行，点击系统托盘图标打开！",
            QSystemTrayIcon.Information,
            2000
        )

    def iconActivated(self, reason):
        if reason == self.tray_icon.Trigger:
            self.showNormal()


if __name__ == "__main__":
    app = QApplication(sys.argv)
    w = MainWindow()
    w.show()
    sys.exit(app.exec())
