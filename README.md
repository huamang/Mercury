# Mercury
一个成长中的扫描器，暂时只有scan模式

## 安装
```
git clone https://github.com/huamang/Mercury.git
cd Mercury/cmd
go build -o Mercury
```

## 用法

```
NAME:
   Mercury - A interesting scanner

USAGE:
   Mercury [global options] command [command options] [arguments...]

VERSION:
   v1.0

COMMANDS:
   scan     scan mode
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

scan模式

```
NAME:
   Mercury scan - scan mode

USAGE:
   Mercury scan [command options] [arguments...]

OPTIONS:
   --target value, -t value  target host, eg: CIDR or "-" or "," to split
   --file value, -f value    target file, eg: /tmp/target.txt
   --thread value, -T value  number of concurrent threads (default: 100)
   --port value, -p value    appoint ports, eg: 80,8080 or 80-88
   --icmp                    icmp mode scan,default ping mode (default: false)
   --check                   just check alive (default: false)
   --help, -h                show help
```

普通扫描模式

![image-20230507053905372](https://tuchuang.huamang.xyz/img/image-20230507053905372.png)

ICMP监听扫描模式，需要高权限

![image-20230507053927090](https://tuchuang.huamang.xyz/img/image-20230507053927090.png)

## 免责声明

本工具仅面向**合法授权**的企业安全建设行为，如您需要测试本工具的可用性，请自行搭建靶机环境。

为避免被恶意使用，本项目所有收录的poc均为漏洞的理论判断，不存在漏洞利用过程，不会对目标发起真实攻击和漏洞利用。

在使用本工具进行检测时，您应确保该行为符合当地的法律法规，并且已经取得了足够的授权。**请勿对非授权目标进行扫描。**

如您在使用本工具的过程中存在任何非法行为，您需自行承担相应后果，我们将不承担任何法律及连带责任。

在安装并使用本工具前，请您**务必审慎阅读、充分理解各条款内容**，限制、免责条款或者其他涉及您重大权益的条款可能会以加粗、加下划线等形式提示您重点注意。 除非您已充分阅读、完全理解并接受本协议所有条款，否则，请您不要安装并使用本工具。您的使用行为或者您以其他任何明示或者默示方式表示接受本协议的，即视为您已阅读并同意本协议的约束。

