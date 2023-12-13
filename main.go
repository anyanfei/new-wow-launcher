package main

import (
	"context"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"log"
	"net/url"
	"runtime"
	"time"

	"new-wow-launcher/lib"
	"new-wow-launcher/lib/play_mp3"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	_ "github.com/lengzhao/font/autoload"
)

//go:embed "resources/Icon.png"
var iconPng []byte

//go:embed "resources/close.png"
var closePng []byte // 关闭按钮图片

//go:embed "resources/leftPng.png"
var leftPng []byte // 左版面图片

//go:embed "resources/rightBottomBanner.png"
var rightBottomBanner []byte

//go:embed "resources/register.png"
var buttonRegister []byte

//go:embed "resources/cash.png"
var buttonCash []byte

//go:embed "resources/QQqun.png"
var buttonQQ []byte

//go:embed "resources/fcmPng.png"
var fcmPng []byte

//go:embed "resources/sltsPng.png"
var sltsPng []byte

//go:embed "resources/startButton.png"
var startButtonPng []byte

//go:embed "resources/updatingButton.png"
var updatingButtonPng []byte

const (
	ClientWidth  = 990
	ClientHeight = 550
	IconWidth    = 20
	BorderWidth  = 4
)

var TextColor = color.RGBA{R: 187, G: 187, B: 182, A: 255} // 统一文字颜色

func main() {
	myApp := app.New()
	myApp.SetIcon(fyne.NewStaticResource("", iconPng))          // 给窗口设置图标
	ctx, cancelFunc := context.WithCancel(context.Background()) // chinese english : When click StartGame button , the mp3 player will be closed
	go play_mp3.PlayMusic(ctx)                                  // 播放一个mp3

	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		// 创建无边框窗口
		w := drv.CreateSplashWindow().(fyne.Window)
		w.SetTitle("战歌峡谷登录器") // 设置标题
		// 设置主页面大小
		w.Resize(fyne.NewSize(ClientWidth, ClientHeight))

		// wow图标
		ico1 := widget.NewIcon(fyne.NewStaticResource("", iconPng))
		ico1.Resize(fyne.NewSize(IconWidth, IconWidth))
		ico1.Move(fyne.NewPos(BorderWidth, BorderWidth))

		text1 := canvas.NewText("战歌登录器 80WLK", TextColor)
		text1.TextSize = 16
		text1.Move(fyne.NewPos(IconWidth+BorderWidth+2, (IconWidth+BorderWidth-10)/2))

		// 关闭按钮
		btn1 := widget.NewButtonWithIcon("", fyne.NewStaticResource("", closePng), func() {
			myApp.Quit()
		})
		btn1.Resize(fyne.NewSize(IconWidth, IconWidth))
		btn1.Move(fyne.NewPos(ClientWidth-IconWidth-BorderWidth, BorderWidth))

		// 画一个矩形框
		rect := canvas.NewRectangle(color.RGBA{R: 0, G: 0, B: 0, A: 0})                               // 这里矩形是没有颜色，中空
		rect.StrokeColor = color.RGBA{R: 155, G: 69, B: 91, A: 20}                                    // 矩形给了边框颜色
		rect.Resize(fyne.NewSize(ClientWidth-2*BorderWidth, ClientHeight-60-IconWidth-2*BorderWidth)) // 设置矩形大小
		rect.Move(fyne.NewPos(BorderWidth, IconWidth+2*BorderWidth))                                  // 移动矩形绝对定位
		rect.StrokeWidth = 1                                                                          // 边框宽度

		// 左版面单图
		leftImage := canvas.NewImageFromResource(fyne.NewStaticResource("", leftPng))
		leftImage.FillMode = canvas.ImageFillOriginal
		leftImage.Resize(fyne.NewSize(480, 460))
		leftImage.Move(fyne.NewPos(BorderWidth+1, IconWidth+2*BorderWidth+1))

		// 右版面 三个按钮

		//1、 注册
		registerImage := canvas.NewImageFromResource(fyne.NewStaticResource("", buttonRegister))
		registerImage.FillMode = canvas.ImageFillContain
		registerImage.Resize(fyne.NewSize(110, 50))
		registerImage.Move(fyne.NewPos(leftImage.Size().Width+IconWidth, IconWidth+6*BorderWidth))
		rightBtn1 := widget.NewButton("", func() {
			username := widget.NewEntry()
			username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]{6,15}$`, "账号请填写在6~15个字符之间")
			password := widget.NewPasswordEntry()
			password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]{6,15}$`, "密码请填写在6~15个字符之间")
			againPassword := widget.NewPasswordEntry()
			againPassword.Validator = func(s string) error { // 自定义校验
				if len(s) == 0 {
					return errors.New("必须再次输入密码")
				}
				if s != password.Text {
					return errors.New("两次密码不一致")
				}
				return nil
			}
			emailEntry := widget.NewEntry()
			emailEntry.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "邮箱地址不正确")
			var captchaIdStr string
			getCodeFromButton := widget.NewButton("获取验证码", func() {
				var imageData string
				imageData, captchaIdStr, _ = lib.GetCode()
				bs, err := base64.StdEncoding.DecodeString(imageData[22:])
				if err != nil {
					log.Fatalln("error loading image data", err)
				}
				img := canvas.NewImageFromResource(fyne.NewStaticResource("验证码", bs))
				img.SetMinSize(fyne.NewSize(160, 60))
				img.Refresh()
				dialog.ShowCustom("验证码", "记住了", img, w)
			})

			validateCodeEntry := widget.NewEntry()
			validateCodeEntry.Validator = func(s string) error {
				if s == "" && len(s) < 4 {
					return errors.New("验证码不正确")
				}
				return nil
			}

			items := []*widget.FormItem{
				widget.NewFormItem("账	号	", username),
				widget.NewFormItem("密	码	", password),
				widget.NewFormItem("确认密码	", againPassword),
				widget.NewFormItem("邮	箱	", emailEntry),
				widget.NewFormItem("点击获取	", getCodeFromButton),
				widget.NewFormItem("验证码	", validateCodeEntry),
			}

			dialog.ShowForm("战歌峡谷注册 Register", "注册提交", "取消", items, func(b bool) {
				if !b {
					return
				}
				err := lib.Register(username.Text, password.Text, againPassword.Text, emailEntry.Text, captchaIdStr, validateCodeEntry.Text)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("注册成功 请截图保存", fmt.Sprintf("请务必截图保存好当前窗口的信息，以后改密码可用\n账号：%s\n密码：%s，邮箱：%s", username.Text, password.Text, emailEntry.Text), w)
			}, w)
		})

		rightBtn1.Resize(fyne.NewSize(110, 50))
		rightBtn1.Move(fyne.NewPos(leftImage.Size().Width+IconWidth, IconWidth+6*BorderWidth+1))

		// 充值
		cashImage := canvas.NewImageFromResource(fyne.NewStaticResource("", buttonCash))
		cashImage.FillMode = canvas.ImageFillContain
		cashImage.Resize(fyne.NewSize(110, 50))
		cashImage.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+110+BorderWidth*2, IconWidth+6*BorderWidth))
		rightBtn2 := widget.NewButton("", func() {
			u, _ := url.Parse("https://www.laghaim.cn")
			fyne.CurrentApp().OpenURL(u)
		})
		rightBtn2.Resize(fyne.NewSize(110, 50))
		rightBtn2.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+110+BorderWidth*2, IconWidth+6*BorderWidth+1))

		// 群聊
		qunImage := canvas.NewImageFromResource(fyne.NewStaticResource("", buttonQQ))
		qunImage.FillMode = canvas.ImageFillContain
		qunImage.Resize(fyne.NewSize(110, 50))
		qunImage.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+2*110+BorderWidth*4, IconWidth+6*BorderWidth))
		rightBtn3 := widget.NewButton("", func() {
			u, _ := url.Parse("http://qm.qq.com/cgi-bin/qm/qr?_wv=1027&k=R77hnihDEdUh7qtuJJVZOpI4YDsEJlC5&authKey=FPnclPsC8jKAPDLzggGpKWGJ0ivsiRFatcR97sHtSjL5%2F%2FQvN10rTusdOzRLlFUM&noverify=0&group_code=474963796")
			fyne.CurrentApp().OpenURL(u)
		})
		rightBtn3.Resize(fyne.NewSize(110, 50))
		rightBtn3.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+2*110+BorderWidth*4, IconWidth+6*BorderWidth+1))

		// 防沉迷
		fcmImage := canvas.NewImageFromResource(fyne.NewStaticResource("", fcmPng))
		fcmImage.FillMode = canvas.ImageFillContain
		fcmImage.Resize(fyne.NewSize(54, 70))
		fcmImage.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+3*110+BorderWidth*6, 8*BorderWidth))

		// 适龄提示
		sltsImage := canvas.NewImageFromResource(fyne.NewStaticResource("", sltsPng))
		sltsImage.FillMode = canvas.ImageFillContain
		sltsImage.Resize(fyne.NewSize(54, 70))
		sltsImage.Move(fyne.NewPos(leftImage.Size().Width+IconWidth+3*110+fcmImage.Size().Width+BorderWidth*7, 8*BorderWidth))

		// 画一根分割线
		line := canvas.NewLine(color.RGBA{R: 155, G: 69, B: 91, A: 20})
		line.Position1 = fyne.NewPos(leftImage.Size().Width+BorderWidth+2, 108) // 线的起点坐标
		line.Position2 = fyne.NewPos(leftImage.Size().Width+504, 108)           // 线的终点坐标
		line.StrokeWidth = 1

		// 战歌公告标题
		rightTitle := canvas.NewText("战歌公告(PatchNotes)", color.White)
		rightTitle.TextSize = 16
		rightTitle.TextStyle.Bold = true
		rightTitle.Move(fyne.NewPos(leftImage.Size().Width+IconWidth, 130))

		// 80*20的方框 160, 44, 26
		rect1 := canvas.NewRectangle(color.RGBA{R: 160, G: 44, B: 26, A: 255})
		rect1.Resize(fyne.NewSize(76, 20))
		rect1.Move(fyne.NewPos(ClientWidth-25*BorderWidth, 127))
		rect1.CornerRadius = 2

		// 最新消息分割线
		rightTitleLine := canvas.NewLine(color.RGBA{R: 106, G: 55, B: 62, A: 255})
		rightTitleLine.Position1 = fyne.NewPos(leftImage.Size().Width+IconWidth, 150) // 线的起点坐标
		rightTitleLine.Position2 = fyne.NewPos(ClientWidth-6*BorderWidth, 150)        // 线的终点坐标
		rightTitleLine.StrokeWidth = 1

		// 获取战歌快讯内容
		newsTitle, newsContent := lib.GetNotice()
		// 时间
		rightTitleTime := canvas.NewText(newsTitle, color.White)
		rightTitleTime.TextSize = 14
		rightTitleTime.Move(fyne.NewPos(ClientWidth-24*BorderWidth-2, 130))
		rightTitleTime.Refresh()

		// 消息内容
		rightContentText := widget.NewLabel(newsContent)
		truncSize := rightContentText.MinSize().SubtractWidthHeight(260, 158)
		rightContentText.Resize(truncSize)
		rightContentText.Wrapping = fyne.TextWrapBreak
		rightContentText.Move(fyne.NewPos(leftImage.Size().Width+IconWidth, 158))
		rightContentText.Refresh()

		// 进度条下的文字222, 202, 183
		progressUnderText := canvas.NewText("更新中", color.RGBA{R: 222, G: 202, B: 183, A: 255})
		progressUnderText.TextSize = 12
		progressUnderTextInfo := make(chan string)
		progressUnderText.Move(fyne.NewPos(ClientWidth/2-20, ClientHeight-20))

		// 进度条
		progress := widget.NewProgressBar()
		progress.Resize(fyne.NewSize(600, 14))
		progress.Move(fyne.NewPos((ClientWidth-600)/2, ClientHeight-40))

		var i float64 = 0
		progress.SetValue(i)
		go func() {
			for {
				if i >= 1.0 {
					progressUnderTextInfo <- "更新完成"
					runtime.Goexit()
				}
				i += 0.01
				time.Sleep(time.Millisecond * 20)
				progress.SetValue(i)
			}
		}()

		// 更新按钮
		updateImage := canvas.NewImageFromResource(fyne.NewStaticResource("", updatingButtonPng))
		updateImage.FillMode = canvas.ImageFillOriginal
		updateImage.Resize(fyne.NewSize(182, 52))
		updateImage.Move(fyne.NewPos(ClientWidth-178-3*BorderWidth, ClientHeight-52-BorderWidth))

		// 开始游戏按钮
		startImage := canvas.NewImageFromResource(fyne.NewStaticResource("", startButtonPng))
		startImage.FillMode = canvas.ImageFillOriginal
		startImage.Resize(fyne.NewSize(182, 52))
		startImage.Move(fyne.NewPos(ClientWidth-178-3*BorderWidth, ClientHeight-52-BorderWidth))
		startImage.Hide()
		startBtn := widget.NewButton("", func() {
			cancelFunc()
			log.Println("点击了开始按钮")
		})
		startBtn.Resize(fyne.NewSize(158, 54))
		startBtn.Move(fyne.NewPos(ClientWidth-170-2*BorderWidth, ClientHeight-52-BorderWidth))
		startBtn.Hide()

		// 进度条改变事件，更新进度条下的文字和旁边按钮
		go processBarChangeEvent(progressUnderText, progressUnderTextInfo, startImage, startBtn, updateImage)

		// 右底部banner
		rightBottomBannerCanvas := canvas.NewImageFromResource(fyne.NewStaticResource("", rightBottomBanner))
		rightBottomBannerCanvas.Resize(fyne.NewSize(470, 120))
		rightBottomBannerCanvas.Move(fyne.NewPos(leftImage.Size().Width+IconWidth, ClientHeight-(180+2*BorderWidth)))

		// 将所有组件都写入到手动布局中
		c := container.NewWithoutLayout(ico1,
			text1,
			btn1,
			rect,
			leftImage,
			rightBottomBannerCanvas,
			rightTitle,
			rightBtn1,
			registerImage,
			rightBtn2,
			cashImage,
			rightBtn3,
			qunImage,
			fcmImage,
			sltsImage,
			line,
			rightTitleLine,
			rect1,
			rightTitleTime,
			rightContentText,
			progress,
			progressUnderText,
			updateImage,
			startBtn,
			startImage,
		)
		w.SetContent(c)
		w.Show()
	}

	myApp.Run()
}

// 进度条改变事件，更新进度条下的文字和旁边按钮
func processBarChangeEvent(thisProcessText *canvas.Text, processTextChanStr chan string, thisImage *canvas.Image, thisButton *widget.Button, updateButtonImg *canvas.Image) {
	thisProcessText.Text = <-processTextChanStr
	thisProcessText.Move(fyne.NewPos(ClientWidth/2-26, ClientHeight-20))
	close(processTextChanStr)
	thisProcessText.Refresh()

	updateButtonImg.Hide()
	thisImage.Show()
	thisButton.Show()

}
