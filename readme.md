# Target
兼容ubuntu22和ubuntu24自带的lsblk命令输出的json某些字段不统一的问题

# How to use

```shell
go get -u github.com/Runninginsilence1/lsblkjson_parse
```

# 要求
- go 1.21+ 因为唯一的依赖 github.com/spf13/cast 要求 go 1.21


# 发布新的Go包
Go包依赖Git的tag机制。
你应该：
```shell
git tag vx.y.z # 给当前的提交打上一个标记
git push origin vx.y.z # 将标记推送到远程仓库
```

之后调用 go get -u 的时候，goproxy.cn 会自动去获取最新的包。