# 1. tools-api
一个基于golang的开源web服务，专注于提供各类API以便在日常开发中提高生产效率


# 2. 部署
下载二进制程序后直接运行

```bash
./go-tools-api
```

运行参数如下：

* **-host**：监听地址，默认 `127.0.0.1`
* **-port**：监听端口，默认 `8085`
* **-image**：图片API的文件夹
* **-jsonMaxMemorySize**： `anyjson` api接口存储的内存占用量（Mib），默认32Mb
* **-jsonValidityDays**：`anyjson` api接口存储的json有效期（天），默认31天


# 3. 二次开发
## 3.1. 环境依赖
Go版本
* go version go1.21.4 linux/amd64

环境搭建
```shell
git clone <repo>
cd <repo_name>
go get -u
```

## 3.2. Visual Studio Code开发配置
参考启动的配置文件（LAUNCH. JSON）如下
```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program":"main.go",
            "cwd": "${workspaceFolder}",
            "args":["-host=0.0.0.0","-port=8017"]
        }
    ]
}
```
# 4. 感谢
项目技术依赖
* https://github.com/gin-gonic/gin