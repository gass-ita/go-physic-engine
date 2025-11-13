package main

import (
	"github.com/gass-ita/go-physics-engine/common"
	"github.com/gass-ita/go-physics-engine/physic"
	"github.com/gass-ita/go-physics-engine/view"
	"github.com/go-gl/mathgl/mgl64"
)

func main() {
	// World boundaries: xmin, xmax, ymin, ymax
	WorldLimits := [4]float64{0, 13.4, 0, 10}

	// --- Initialize physics scene ---
	scene := physic.NewScene(WorldLimits)
	gravity := mgl64.Vec2{0, -9.81}

	// --- Pendulum Example ---
	p1 := physic.NewParticle(mgl64.Vec2{6, 10}, mgl64.Vec2{0, 0}, 1.0, 0.05, true) // Fixed point
	p2 := physic.NewParticle(mgl64.Vec2{6, 8}, mgl64.Vec2{5, 0}, 1.0, 0.05, false) // Free particle
	p2.ForceFields = append(p2.ForceFields, gravity)
	scene.AddParticle(p1)
	scene.AddParticle(p2)
	scene.AddSpring(physic.NewSpring(p1, p2, 1000, 2.0))

	// --- Soft Body Parameters ---
	gridSize := 5
	distance := 0.4
	mass := 1.0
	springK := 150.0
	damperBeta := 0.0
	radius := 0.05
	startX := 2.0
	startY := 5.0
	softBodySpacing := 2.5 // horizontal spacing between soft bodies
	numSoftBodies := 10    // number of soft bodies (n)

	for n := 0; n < numSoftBodies; n++ {
		xOffset := startX + float64(n)*softBodySpacing

		// create the grid of particles
		particles := make([][]*physic.Particle, gridSize)
		for i := 0; i < gridSize; i++ {
			particles[i] = make([]*physic.Particle, gridSize)
			for j := 0; j < gridSize; j++ {
				pos := mgl64.Vec2{
					xOffset + float64(i)*distance,
					startY + float64(j)*distance,
				}
				vel := mgl64.Vec2{0, 0}
				p := physic.NewParticle(pos, vel, mass, radius, false)
				p.ForceFields = append(p.ForceFields, gravity)
				scene.AddParticle(p)
				particles[i][j] = p
			}
		}

		// connect every particle with every other in this soft body
		for i := 0; i < gridSize; i++ {
			for j := 0; j < gridSize; j++ {
				for k := 0; k < gridSize; k++ {
					for l := 0; l < gridSize; l++ {
						if i == k && j == l {
							continue
						}

						p1 := particles[i][j]
						p2 := particles[k][l]
						restLength := p2.Position.Sub(p1.Position).Len()

						if springK != 0 {
							spring := physic.NewSpring(p1, p2, springK, restLength)
							scene.AddSpring(spring)
						}
						if damperBeta != 0 {
							damper := physic.NewDamper(p1, p2, damperBeta)
							scene.AddDamper(damper)
						}
					}
				}
			}
		}
	}

	// --- Channels for rendering ---
	posChan := make(chan []common.ParticlePos, 1)
	linkChan := make(chan []common.LinkPos, 1)
	infoChan := make(chan common.Info, 1)

	// --- Start physics simulation ---
	scene.Start(common.DT_PHYSIC, posChan, linkChan, infoChan)

	// --- Start Ebiten visualization ---
	window := view.NewWindow(posChan, linkChan, infoChan)
	window.Run()
}
