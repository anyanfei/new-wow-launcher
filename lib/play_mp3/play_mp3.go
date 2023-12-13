package play_mp3

import (
	"context"
	"embed"
	"io"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

// 这里替换mp3文件即可
//
//go:embed "newLogin.mp3"
var mp3File embed.FS

func PlayMusic(ctx context.Context) {
	f, _ := mp3File.Open("newLogin.mp3")
	defer f.Close()
	d, _ := mp3.NewDecoder(f)
	c, _ := oto.NewContext(44100, 2, 2, 8192)
	defer c.Close()
	p := c.NewPlayer()
	defer p.Close()
	copy(ctx, p, d)
}

// 不要动下面的内容
type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// 重写io.Copy方法 ，使其支持取消操作
func copy(ctx context.Context, dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
	return err
}
