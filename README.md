# cos-tgbot

> 一个可以帮助你快速上传文件到腾讯云cos的telegram机器人

### 功能

* 将你发给Bot的文件(不大于[50MB](https://core.telegram.org/bots/faq#how-do-i-upload-a-large-file))转存到cos
* 快速上传到`/css`、`/img`、`/upload`等文件夹内
* 对上传的文件按照日期进行命名和放置
* 上传后可用bot对文件进行删除
* 无须数据库

### 使用方法

#### 下载

1. 克隆仓库

   ```shell
   git clone https://github.com/rroy233/cos-tgbot.git
   ```

2. 获取可执行文件

    1. 自行编译

       ```shell
       cd cos-tgbot/
       # 自行编译
       # go版本要求：go1.17+
       go build -o cosTgBot
       # 或者使用make交叉编译
       make
       ```

    2. 前往release下载

       下载已编译的[可执行文件](https://github.com/rroy233/cos-tgbot/releases)，重新命名为`cosTgbot`，放于项目文件夹内

#### 找 BotFather 创建Bot

获得`bot_token`,然后设置命令列表

```
help - 获取帮助
keyboard - 获取快捷键盘
upload - 上传到/upload文件夹下
css - 上传到/css文件夹下
js - 上传到/js文件夹下
img - 上传到/img文件夹下
```

#### 创建配置文件

复制`config.example.json`为`config.json`，然后填入配置

```json
{
  "bot_token": "",
  "admin_uid": [114514],
  "cos": {
    "BucketURL": "https://xxxxx.cos.ap-shanghai.myqcloud.com",
    "ServiceURL": "https://cos.ap-shanghai.myqcloud.com",
    "SecretID": "",
    "SecretKey": "",
    "cdnUrlDomain": "https://xxxxx.xxxx.com"
  }
}
```

#### 运行程序

```shell
# 编译并运行(unix)
bash ./buildrun.sh 

# 直接运行(unix)
bash ./run.sh 
```


