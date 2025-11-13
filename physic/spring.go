package physic

type Spring struct {
	P1, P2     *Particle
	K          float64
	RestLength float64
}

func NewSpring(p1, p2 *Particle, k float64, restLength float64) *Spring {
	if restLength < 0 {
		diff := p2.Position.Sub(p1.Position)
		restLength = diff.Len()
	}
	return &Spring{P1: p1, P2: p2, K: k, RestLength: restLength}
}

func (s *Spring) Update(dt float64) {
	delta := s.P2.PredictedPosition.Sub(s.P1.PredictedPosition)
	dist := delta.Len()
	if dist < 1e-8 {
		return
	}

	dir := delta.Mul(1.0 / dist)
	displacement := dist - s.RestLength
	force := dir.Mul(s.K * displacement)

	s.P1.AddForce(force)
	s.P2.AddForce(force.Mul(-1))
}
