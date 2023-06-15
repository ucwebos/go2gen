# go2gen 命令行工具

## 安装

    go install github.com/ucwebos/go2gen@latest

## 使用帮助

    go2gen help

## 配置文件

    可以在项目根目录下配置go2gen.yaml文件

## 代码生成工具注释标识

    @IGNORE 忽略 用于generate忽略对应实体
    @GI     自动注册单列方法

## 代码生成工具字段解析标签
    db 生成数据库层相关

## 代码范围
    默认generate命令只会解析domain目录