package physic

import (
	"gonum.org/v1/gonum/mat"
)

type Spring struct {
	P1, P2     *Particle
	K          float64 // spring constant [N/m]
	RestLength float64 // rest length of the spring
}

// Create a new spring constraint between two particles
func NewSpring(p1, p2 *Particle, k float64, restLength float64) *Spring {
	// If restLength <= 0, initialize with distance between particles
	if restLength < 0 {
		diff := mat.NewVecDense(p1.Position.Len(), nil)
		diff.SubVec(&p2.Position, &p1.Position)
		restLength = mat.Norm(diff, 2)
	}
	return &Spring{
		P1:         p1,
		P2:         p2,
		K:          k,
		RestLength: restLength,
	}
}

// Update applies the spring force based on predicted positions
func (s *Spring) Update(dt float64) {
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

	// displacement = dist - restLength
	displacement := dist - s.RestLength

	// force = k * displacement / 2 * direction
	force := mat.NewVecDense(delta.Len(), nil)
	force.ScaleVec(s.K*displacement, direction)

	// apply force to particles
	// print the force

	s.P1.AddForce(force)

	// opposite force to P2
	negForce := mat.NewVecDense(force.Len(), nil)
	negForce.ScaleVec(-1.0, force)
	s.P2.AddForce(negForce)
}
