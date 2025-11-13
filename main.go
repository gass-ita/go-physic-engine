package main

import (
	"github.com/gass-ita/go-physics-engine/common"
	"github.com/gass-ita/go-physics-engine/physic"
	"github.com/gass-ita/go-physics-engine/view"
	"gonum.org/v1/gonum/mat"
)

// LinkPos holds screen coordinates of a particle

func main() {
	WorldLimits := [4]float64{0, 13.4, 0, 10} // xmin, xmax, ymin, ymax
	/* // --- Initialize physics ---
	pos1 := mat.NewVecDense(2, []float64{0, 6})
	pos2 := mat.NewVecDense(2, []float64{0, 7})
	vel1 := mat.NewVecDense(2, []float64{0, 0})
	vel2 := mat.NewVecDense(2, []float64{0, 0})

	p1 := physic.NewParticle(pos1, vel1, 1.0, 0.05, false)
	p2 := physic.NewParticle(pos2, vel2, 1.0, 0.05, false)
	//p2 := physic.NewParticle(pos2, vel2, 1.0, 1.0, true)
	spring := physic.NewSpring(p1, p2, 200, 1.0)
	//damper := physic.NewDamper(p1, p2, 4)
	*/

	// make a soft body made of 9 particles connected by springs in a 3x3 grid
	scene := physic.NewScene(WorldLimits)
	gravity := mat.NewVecDense(2, []float64{0, -9.81})

	// make a simple pendulum with a fixed particle and a free particle connected by a spring
	p1 := physic.NewParticle(mat.NewVecDense(2, []float64{6, 10}), mat.NewVecDense(2, []float64{0, 0}), 1.0, 0.05, true)
	p2 := physic.NewParticle(mat.NewVecDense(2, []float64{6, 8}), mat.NewVecDense(2, []float64{5, 0}), 1.0, 0.05, false)
	p2.ForceFields = append(p2.ForceFields, gravity)
	scene.AddParticle(p1)
	scene.AddParticle(p2)
	spring := physic.NewSpring(p1, p2, 1000, 2.0)
	scene.AddSpring(spring)

	xOffset := 1.0
	gridSize := 5
	distanceBetweenParticles := 0.4
	particleMass := 1.0
	springK := 150.0
	damperBeta := 0.0
	particleRadius := 0.05
	particles := make([][]*physic.Particle, gridSize)
	for n := range 100 {
		xOffset += float64(n) * 1
		for i := range gridSize {
			particles[i] = make([]*physic.Particle, gridSize)
			for j := 0; j < gridSize; j++ {
				pos := mat.NewVecDense(2, []float64{float64(i)*distanceBetweenParticles + xOffset, float64(j)*distanceBetweenParticles + 5})
				vel := mat.NewVecDense(2, []float64{0, 0})
				particles[i][j] = physic.NewParticle(pos, vel, particleMass, particleRadius, false)
				particles[i][j].ForceFields = append(particles[i][j].ForceFields, gravity)
				scene.AddParticle(particles[i][j])
			}
		}

		// connect particles with springs 1 to every other forming a complete graph
		for i := range gridSize {
			for j := range gridSize {
				for k := 0; k < gridSize; k++ {
					for l := range gridSize {
						if i == k && j == l {
							continue
						}
						p1 := particles[i][j]
						p2 := particles[k][l]
						restLength := mat.Norm(mat.NewVecDense(2, []float64{
							p2.Position.AtVec(0) - p1.Position.AtVec(0),
							p2.Position.AtVec(1) - p1.Position.AtVec(1),
						}), 2)
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
	// --- Channel to communicate positions ---
	posChan := make(chan []common.ParticlePos, 1)
	linkChan := make(chan []common.LinkPos, 1)
	infoChan := make(chan common.Info, 1)

	// --- Start physics simulation ---
	scene.Start(common.DT_PHYSIC, posChan, linkChan, infoChan)

	// --- Run Ebiten game loop ---
	Window := view.NewWindow(posChan, linkChan, infoChan)
	Window.Run()

}
