package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func loadWallpaper(renderer *sdl.Renderer, wallpaper wallpaper, err error) (wallpaper, error) {
	for i := 0; i < 4; i++ {
		var surface *sdl.Surface
		var tex *sdl.Texture
		if surface, err = sdl.LoadBMP(bmpImages[i]); err != nil {
			return wallpaper, err
		}
		if tex, err = renderer.CreateTextureFromSurface(surface); err != nil {
			return wallpaper, err
		}
		if i == 0 {
			wallpaper.originalW = surface.W
			wallpaper.originalH = surface.H
		}
		surface.Free()
		wallpaper.layers[i].tex = tex
	}
	if wallpaper.tex, err = renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_TARGET, width, height); err != nil {
		return wallpaper, err
	}
	return wallpaper, err
}
