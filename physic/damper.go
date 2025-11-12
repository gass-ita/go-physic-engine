package physic

import (
	"gonum.org/v1/gonum/mat"
)

type Damper struct {
	P1, P2 *Particle
	Beta   float64 // spring constant [Ns/m]
}

// Create a new spring constraint between two particles
func NewDamper(p1, p2 *Particle, beta float64) *Damper {
	return &Damper{
		P1:   p1,
		P2:   p2,
		Beta: beta,
	}
}

// Update applies the spring force based on predicted positions
func (s *Damper) Update(dt float64) {
	// delta = P2_predicted - P1_predicted
	delta := mat.NewVecDense(s.P1.Position.Len(), nil)
	delta.SubVec(&s.P2.PredictedPosition, &s.P1.PredictedPosition)

	dist := mat.Norm(delta, 2)
	if dist < 1e-8 {
		return // avoid division by zero
	}

	// direction = delta / dist
	direction := mat.NewVecDense(delta.Len(), nil)
	direction.ScaleVec(1.0/dist, delta)

	relativeVelocity := mat.NewVecDense(s.P1.Position.Len(), nil)
	relativeVelocity.SubVec(&s.P2.PredictedVelocity, &s.P1.PredictedVelocity)

	velAlongDir := mat.Dot(relativeVelocity, direction)

	// force = k * displacement / 2 * direction
	force := mat.NewVecDense(delta.Len(), nil)
	force.ScaleVec(s.Beta*velAlongDir, direction)

	// apply force to particles
	// print the force

	s.P1.AddForce(force)

	// opposite force to P2
	negForce := mat.NewVecDense(force.Len(), nil)
	negForce.ScaleVec(-1.0, force)
	s.P2.AddForce(negForce)
}
