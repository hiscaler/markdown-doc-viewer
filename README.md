# Markdown doc viewer
Markdown 格式文档浏览程序

## 结构
1. config/config.json:  配置文件
2. docs: 存放文档，一个文档一个目录的形式
3. server.exe|server： 可执行文件（Windows 下为 server.exe）

## config.json 配置文件
###  格式
```json
{
    "port": 8080,
    "documentDir": "./docs",
    "documentDirs": [
        "/www/wwwroot/a.com/docs",
        "/www/wwwroot/b.com/docs",
        "D:\\www\\wwwroot\\c.com\\docs"
    ]
}
```

### 说明
> `port` 为服务开启端口，默认情况下为 8080
>
> `documentDir` 为放置文档的位置，默认情况下为 docs，允许多个目录，程序会自动读取该文件夹以及文件下的所有文件渲染
>  
> `documentDirs` 用来设置放置在其他其他不同位置的零散文档，可以设置多个，需要注意的是，必须是绝对路径的目录。同时不会读取其中子目录的内容。
>
> **在 Windows 操作系统中，路径应该使用两个右斜线进行分隔**

## 使用
1. 进入程序放置目录;
2. 执行 ./server 开启服务；
3. 打开浏览器输入 http://127.0.0.1:8080 访问。
