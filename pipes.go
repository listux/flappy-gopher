package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

type pipes struct {
	mu sync.RWMutex

	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

type pipe struct {
	mu sync.RWMutex

	x        int32
	h        int32
	w        int32
	inverted bool
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	p, err := img.LoadTexture(r, "res/img/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	ps := &pipes{
		texture: p,
		speed:   2,
	}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return ps, nil
}

func newPipe() *pipe {

	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}

}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy pipe: %v", err)
	}

	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}

}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b.touch(p)
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.mu.Lock()
		p.x -= ps.speed
		p.mu.Unlock()

		if p.x+p.w > 0 { // clear some memory, when the pipe is gone.
			rem = append(rem, p)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}
