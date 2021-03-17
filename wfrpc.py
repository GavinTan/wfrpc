import os
import requests
import platform
import tarfile
import zipfile
import shutil
import configparser
import argparse
import random
import sys
import logging


log = logging.getLogger()
console = logging.StreamHandler()
formatter = logging.Formatter('%(asctime)s %(levelname)s %(message)s', datefmt='%Y/%m/%d %H:%M:%S')
console.setFormatter(formatter)
log.addHandler(console)
log.setLevel(logging.INFO)


parser = argparse.ArgumentParser(description=__doc__, add_help=False)
parser.add_argument('--server_addr', required=False, metavar='', help='frp服务器地址, 默认nb33.3322.org', default='nb33.3322.org')
parser.add_argument('--server_port', required=False, metavar='', help='frp服务器端口, 默认7000', default='7000')
parser.add_argument('--host', required=False, metavar='', help='转发目标IP地址, 默认127.0.0.1', type=str, action='append')
parser.add_argument('--lport', required=False, metavar='', help='转发端口, 支持指定多个, 默认22 80', type=str, action='append')
parser.add_argument('--rport', required=False, metavar='', help='指定转发后的端口, 不指定将随机生成, 指定后当与lport数量不对应时将根据指定的端口自动递增', type=str, action='append')
parser.add_argument('--type', required=False, metavar='', help='指定端口类型 tcp or udp', default='tcp', type=str)
parser.add_argument('--clear', required=False, help='清理配置', action='store_true')
parser.add_argument('-h', '--help', required=False, help='帮助', action='store_true')

args = parser.parse_args()


class FrpBase:
    def __init__(self):
        self.sys_name = platform.system().lower()
        self.sys_machine = platform.machine().lower()
        self.sys_arch_list = {'x86_64': 'amd64', 'aarch64': 'arm64', 'X86': '386'}
        self.sys_arch = self.sys_arch_list.get(self.sys_machine) or self.sys_machine

        self.frpc_bin = '/usr/local/bin/frpc'
        self.frp_config_path = f"{os.path.expanduser('~')}/.frp/frpc.ini"

        self.frp_api_github_url = 'https://api.github.com/repos/fatedier/frp/releases/latest'
        self.frp_release_github_url = 'https://github.com/fatedier/frp/releases'

        self.frp_config_host_list = args.host or ['127.0.0.1']
        self.frp_config_lport_list = args.lport or ['22', '80']
        self.frp_config_rport_list = args.rport or []
        self.frp_config_type = args.type
        self.frp_config_server_addr = args.server_addr
        self.frp_config_server_port = args.server_port

        self.clear = args.clear

        self.config = []

    def download_frp(self):
        if not os.path.exists(self.frpc_bin):
            r = requests.get(self.frp_api_github_url)
            tmp_path = '/tmp'
            for i in r.json().get('assets'):
                frp_download_url = i.get('browser_download_url')
                if f'{self.sys_name}_{self.sys_arch}' in frp_download_url:
                    log.info(f'download {frp_download_url}')
                    file_name = os.path.basename(frp_download_url)
                    file_path = f'{tmp_path}/{file_name}'
                    frp_download = requests.get(frp_download_url)
                    with open(f'{tmp_path}/{file_name}', 'wb') as f:
                        f.write(frp_download.content)
                
                    if file_name.endswith('.tar.gz'):
                        with tarfile.open(file_path, "r:gz") as f:
                            f.extractall(tmp_path)
                        shutil.copy2(f"{tmp_path}/{file_name.split('.tar.gz')[0]}/frpc", '/usr/local/bin/')

                    if file_name.endswith('.zip'):
                        with zipfile.ZipFile(file_path) as f:
                            f.extractall(tmp_path)

                        shutil.copy2(f"{tmp_path}/{file_name.split('.zip')[0]}/frpc", '/usr/local/bin/')
                    return True
            
            log.error(f'当前系统架构 {self.sys_name}_{self.sys_machine} 没有找到匹配的包，请手动下载 {self.frp_release_github_url}')

    def frp_config(self):
        frp_config_parser = configparser.ConfigParser()
        frp_config_parser.read(self.frp_config_path)

        os.path.exists(self.frp_config_path) or os.makedirs(f"{os.path.expanduser('~')}/.frp", exist_ok=True)
        
        if not frp_config_parser.has_section('common'):
            log.info('generate frp configuration file')
            frp_config_parser.add_section('common')
            frp_config_parser.set('common', 'server_addr ', self.frp_config_server_addr)
            frp_config_parser.set('common', 'server_port ', self.frp_config_server_port)

        n = abs(len(self.frp_config_lport_list) - len(self.frp_config_rport_list))

        if self.frp_config_rport_list and n !=0:
            for i in range(1, n + 1):
                self.frp_config_rport_list.append(str(int(self.frp_config_rport_list[-1]) + 1))

        for h in self.frp_config_host_list:
            for lp in self.frp_config_lport_list:
                section = f'{h}_{lp}'
                if not frp_config_parser.has_section(section):
                    frp_config_parser.add_section(section)

                frp_config_parser.set(section, 'type', self.frp_config_type)
                frp_config_parser.set(section, 'local_ip', h)
                frp_config_parser.set(section, 'local_port', lp)

                if self.frp_config_rport_list:
                    frp_config_parser.set(section, 'remote_port', self.frp_config_rport_list[self.frp_config_lport_list.index(lp)])
                else:
                    frp_config_parser.set(section, 'remote_port', str(random.randint(10000,20000)))
                
        frp_config_parser.write(open(self.frp_config_path, 'w')) 

        for i in frp_config_parser.sections():
            if i != 'common':
                self.config.append({k:v for k,v in frp_config_parser.items(i)})

    def main(self):
        try:
            if args.help:
                parser.print_help()
                sys.exit(0)

            if self.clear:
                os.path.isfile(self.frp_config_path) and os.remove(self.frp_config_path)
                # os.path.isfile(self.frpc_bin) and os.remove(self.frpc_bin)
                print('clear config success')
                sys.exit(0)

            self.download_frp()
            self.frp_config()
            print('\t\t+' + '-' * 27 + 'config' + '-' * 27 + '+')
            for i in self.config:
                msg = f"{i.get('local_ip')}:{i.get('local_port')} --> {self.frp_config_server_addr}:{i.get('remote_port')}"
                n = round((60 - len(msg)) / 2)
                print(f"\t\t{' ' * n}{msg}{' ' * n}")
            print('\t\t+' + '-' * 60 + '+')

            os.system(f'{self.frpc_bin} -c {self.frp_config_path}')
        except KeyboardInterrupt:
            print('\nExit!')
            sys.exit(0)
        
        
if __name__ == '__main__':
    f = FrpBase()
    f.main()
