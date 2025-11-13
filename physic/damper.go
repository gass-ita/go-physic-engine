package physic

type Damper struct {
	P1, P2 *Particle
	Beta   float64 // damping coefficient [Ns/m]
}

func NewDamper(p1, p2 *Particle, beta float64) *Damper {
	return &Damper{P1: p1, P2: p2, Beta: beta}
}

func (d *Damper) Update(dt float64) {
	delta := d.P2.PredictedPosition.Sub(d.P1.PredictedPosition)
	dist := delta.Len()
	if dist < 1e-8 {
		return
	}

	dir := delta.Mul(1.0 / dist)
	relVel := d.P2.PredictedVelocity.Sub(d.P1.PredictedVelocity)
	velAlongDir := relVel.Dot(dir)

	force := dir.Mul(d.Beta * velAlongDir)

	d.P1.AddForce(force)
	d.P2.AddForce(force.Mul(-1))
}
