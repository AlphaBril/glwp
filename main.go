package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	bmpImages = [4]string{
		"./wallpapers/default/1.bmp",
		"./wallpapers/default/2.bmp",
		"./wallpapers/default/3.bmp",
		"./wallpapers/default/4.bmp",
	}
)

const (
	width  = 1200
	height = 700
)

type layer struct {
	sensitivityX float32
	sensitivityY float32
	tex          *sdl.Texture
}

type wallpaper struct {
	layers    [4]layer
	tex       *sdl.Texture
	originalW int32
	originalH int32
}

func lerp(a int, b int32, t float64) int {
	if t > 1 {
		t = 1
	}
	return int(float64(a) + float64(t)*(float64(b)-float64(a)))
}

func run() (err error) {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var wallpaper wallpaper
	var monitor *sdl.Texture
	var info *sdl.SysWMInfo
	var subsystem string

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	if window, err = sdl.CreateWindow("Parallax", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_OPENGL|sdl.WINDOW_HIDDEN); err != nil {
		return
	}
	defer window.Destroy()
	if info, err = window.GetWMInfo(); err != nil {
		return
	}
	fmt.Printf("info: %v\n", info.Subsystem)
	switch info.Subsystem {
	case sdl.SYSWM_UNKNOWN:
		subsystem = "An unknown system!"
	case sdl.SYSWM_WINDOWS:
		subsystem = "Microsoft Windows(TM)"
	case sdl.SYSWM_X11:
		subsystem = "X Window System"
	case sdl.SYSWM_DIRECTFB:
		subsystem = "DirectFB"
	case sdl.SYSWM_COCOA:
		subsystem = "Apple OS X"
	case sdl.SYSWM_UIKIT:
		subsystem = "UIKit"
	case sdl.SYSWM_WAYLAND:
		subsystem = "Wayland"
	}

	fmt.Printf("This program is running SDL version %d.%d.%d on %s\n",
		info.Version.Major,
		info.Version.Minor,
		info.Version.Patch,
		subsystem)

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return
	}
	defer renderer.Destroy()

	if wallpaper, err = loadWallpaper(renderer, wallpaper, err); err != nil {
		return
	}
	if monitor, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_TARGET, width, height); err != nil {
		return
	}

	var (
		mx    int32  = 0
		my    int32  = 0
		state uint32 = 0
	)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}
		var currentX = 0
		var currentY = 0

		var lastTicks uint32 = 0
		ticks := sdl.GetTicks()
		dT := float64(ticks-lastTicks) / float64(1000.0)
		lastTicks = ticks

		mx, my, state = sdl.GetGlobalMouseState()
		fmt.Println(state)

		currentX = lerp(currentX, mx, dT*float64(8.0))
		currentY = lerp(currentY, my, dT*float64(8.0))

		renderer.SetRenderTarget(monitor)
		renderer.Clear()
		for i := 0; i < 4; i++ {
			src := sdl.Rect{
				X: 0,
				Y: 0,
				W: wallpaper.originalW,
				H: wallpaper.originalH,
			}

			x := -(float64(currentX-width/2) *
				float64(0.05))
			y := -(float64(currentY-height/2) *
				float64(0.05))
			for k := -0; k <= 0; k++ {
				for j := -0; j <= 0; j++ {
					dest := sdl.Rect{
						X: int32(x + float64(j*height)),
						Y: int32(y + float64(k*height)),
						W: width,
						H: height,
					}
					if err = renderer.Copy(wallpaper.layers[i].tex, &src, &dest); err != nil {
						return
					}
				}
			}
		}
		renderer.SetRenderTarget(nil)
		src := sdl.Rect{
			X: 0,
			Y: 0,
			W: wallpaper.originalW,
			H: wallpaper.originalH,
		}
		dest := sdl.Rect{
			X: 0,
			Y: 0,
			W: wallpaper.originalW,
			H: wallpaper.originalH,
		}
		if err = renderer.Copy(monitor, &src, &dest); err != nil {
			return
		}
		renderer.Present()
		sdl.WaitEvent()
	}

	return
}

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}
