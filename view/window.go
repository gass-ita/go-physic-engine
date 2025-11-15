package view

import (
	"fmt"
	"image/color"
	"log"

	"github.com/gass-ita/go-physics-engine/common"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Window holds state for Ebiten rendering
type Window struct {
	particlePositions []common.ParticlePos
	linkPositions     []common.LinkPos
	info              common.Info

	posChan  chan []common.ParticlePos
	linkChan chan []common.LinkPos
	infoChan chan common.Info
}

func NewWindow(posChan chan []common.ParticlePos, linkChan chan []common.LinkPos, infoChan chan common.Info) *Window {
	// initialize Ebiten window
	ebiten.SetWindowSize(common.WIDTH, common.HEIGHT)
	ebiten.SetWindowTitle("Particle Simulation")
	return &Window{
		posChan:  posChan,
		linkChan: linkChan,
		infoChan: infoChan,
	}
}

// run window
func (g *Window) Run() error {
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (g *Window) Update() error {
	// Non-blocking read from channel
	select {
	case linkPos := <-g.linkChan:
		g.linkPositions = linkPos
	default:
	}

	select {
	case particlePos := <-g.posChan:
		g.particlePositions = particlePos
	default:
	}

	select {
	case info := <-g.infoChan:
		g.info = info
	default:
	}

	return nil
}

func (g *Window) Draw(screen *ebiten.Image) {
	// set background color
	col := color.RGBA{
		R: uint8((common.BACKGROUND_COLOR >> 24) & 0xFF),
		G: uint8((common.BACKGROUND_COLOR >> 16) & 0xFF),
		B: uint8((common.BACKGROUND_COLOR >> 8) & 0xFF),
		A: uint8(common.BACKGROUND_COLOR & 0xFF),
	}
	screen.Fill(col)

	if g == nil || g.particlePositions == nil {
		return
	}
	link_col := color.RGBA{
		R: uint8((common.LINK_COLOR >> 24) & 0xFF),
		G: uint8((common.LINK_COLOR >> 16) & 0xFF),
		B: uint8((common.LINK_COLOR >> 8) & 0xFF),
		A: uint8(common.LINK_COLOR & 0xFF),
	}
	for _, l := range g.linkPositions {
		// TODO: probably this conversion should be done in the controller
		x1 := l.X1 * common.PX_PER_METER
		y1 := l.Y1 * common.PX_PER_METER
		y1 = float64(screen.Bounds().Dy()) - y1
		x2 := l.X2 * common.PX_PER_METER
		y2 := l.Y2 * common.PX_PER_METER
		y2 = float64(screen.Bounds().Dy()) - y2
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, link_col)
	}

	particle_col := color.RGBA{
		R: uint8((common.PARTICLE_COLOR >> 24) & 0xFF),
		G: uint8((common.PARTICLE_COLOR >> 16) & 0xFF),
		B: uint8((common.PARTICLE_COLOR >> 8) & 0xFF),
		A: uint8(common.PARTICLE_COLOR & 0xFF),
	}

	for _, p := range g.particlePositions {
		// TODO: probably this conversion should be done in the controller
		x := p.X * common.PX_PER_METER
		y := p.Y * common.PX_PER_METER
		// flip Y axis
		y = float64(screen.Bounds().Dy()) - y
		radius := p.Radius * common.PX_PER_METER
		ebitenutil.DrawCircle(screen, x, y, radius, particle_col)
	}

	// Draw info text
	ebitenutil.DebugPrintAt(screen, "Frame Time (ms): "+fmt.Sprintf("%.2f", g.info.ElapsedTime), 10, 10)

}

func (g *Window) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.WIDTH, common.HEIGHT
}
