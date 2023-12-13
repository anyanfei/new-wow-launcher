#主要目的

1、练习一下go-fyne的用法以及其中的小技巧

2、让登录器支持mp3播放以及重写io.copy中断播放

请clone项目：

[gitee](https://gitee.com/87066062/new-wow-launcher1)

[github](https://github.com/anyanfei/new-wow-launcher)



#操作步骤

1、安装golang环境，版本能多高就多高，至少1.18以上（笔者是1.20）

2、首先安装gcc https://www.msys2.org/   跟着官网步骤来，把gcc加入到环境变量中哦，记得打开 CGO_ENABLED=1 (命令：go env -w CGO_ENABLED=1)

3、如果是vs code，记得启动vs code以管理员方式启动，不然就算把gcc放到环境变量了，也找不到的

4、go mod tidy  拉依赖

5、go run main.go 跑起来吧


请参阅以下[指南]
- [go-fyne](https://github.com/fyne-io/fyne)
- [go-mp3](github.com/hajimehoshi/go-mp3)
- [play-music](github.com/hajimehoshi/oto)