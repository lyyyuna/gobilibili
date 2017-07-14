# gobilibili

B 站直播弹幕 Go 版。

## 安装

    go get github.com/lyyyuna/gobilibili

## 简单说明

参考 example/main.go，打印实时直播弹幕。

程序逻辑来自我的 [B 站直播弹幕姬 Python 版](https://github.com/lyyyuna/bilibili_danmu)。逻辑基本和 Python 版本保持一致，可以对着理解。

正如 [B 站直播弹幕姬 Python 版](https://github.com/lyyyuna/bilibili_danmu) 指出的，B 站直播协议会变，不保证向后兼容性。

在 example 的 main.go 中，roomid 没有自动转换，对于部分 UP 主，不要使用短 ID，请使用他的原始 ID。（Python 的版本做了处理，Go 我偷懒了。）

## 写后感

写并发，Go 比 Python 的 asyncio 舒服。但由于我是 Go 新手，总体上，写起来不舒服，Json 的处理也稍显麻烦。