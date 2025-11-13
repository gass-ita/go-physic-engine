package physic

import (
	"fmt"

	"github.com/gass-ita/go-physics-engine/common"
	"github.com/go-gl/mathgl/mgl64"
)

type Particle struct {
	Position   mgl64.Vec2
	Velocity   mgl64.Vec2
	Mass       float64
	Radius     float64
	Static     bool
	ForceAccum mgl64.Vec2

	ForceFields []mgl64.Vec2

	OldPosition       mgl64.Vec2
	PredictedPosition mgl64.Vec2
	PredictedVelocity mgl64.Vec2
	acc               mgl64.Vec2
	newPos            mgl64.Vec2
}

func (p *Particle) ClampPosition(limits [4]float64) {
	if p.Static {
		return
	}

	xmin, xmax, ymin, ymax := limits[0], limits[1], limits[2], limits[3]

	// X
	if p.Position.X() < xmin+p.Radius {
		p.Position[0] = xmin + p.Radius
		p.Velocity[0] = -p.Velocity.X()
	} else if p.Position.X() > xmax-p.Radius {
		p.Position[0] = xmax - p.Radius
		p.Velocity[0] = -p.Velocity.X()
	}

	// Y
	if p.Position.Y() < ymin+p.Radius {
		p.Position[1] = ymin + p.Radius
		p.Velocity[1] = -p.Velocity.Y()
	} else if p.Position.Y() > ymax-p.Radius {
		p.Position[1] = ymax - p.Radius
		p.Velocity[1] = -p.Velocity.Y()
	}
}

func NewParticle(pos, vel mgl64.Vec2, mass, radius float64, isStatic bool) *Particle {
	oldPos := pos.Sub(vel.Mul(common.DT_PHYSIC))
	p := &Particle{
		Position:          pos,
		Velocity:          vel,
		Mass:              mass,
		Radius:            radius,
		Static:            isStatic,
		ForceAccum:        mgl64.Vec2{0, 0},
		ForceFields:       []mgl64.Vec2{},
		OldPosition:       oldPos,
		PredictedPosition: mgl64.Vec2{0, 0},
		PredictedVelocity: mgl64.Vec2{0, 0},
		acc:               mgl64.Vec2{0, 0},
		newPos:            mgl64.Vec2{0, 0},
	}

	if isStatic {
		p.PredictedPosition = pos
		p.PredictedVelocity = mgl64.Vec2{0, 0}
		p.Mass = 1
		p.Velocity = mgl64.Vec2{0, 0}
	}
	return p
}

func (p *Particle) AddForce(force mgl64.Vec2) {
	p.ForceAccum = p.ForceAccum.Add(force)
}

func (p *Particle) ApplyExternalForces() {
	p.ForceAccum = mgl64.Vec2{0, 0}
	if p.Static {
		return
	}

	// Air friction
	p.ForceAccum = p.ForceAccum.Sub(p.Velocity.Mul(common.AIR_FRICTION))

	// External force fields
	for _, f := range p.ForceFields {
		p.ForceAccum = p.ForceAccum.Add(f)
	}
}

func (p *Particle) UpdatePredictedState(dt float64) {
	if p.Static {
		return
	}

	p.acc = p.ForceAccum.Mul(1 / p.Mass)
	p.PredictedPosition = p.Position.Mul(2).Sub(p.OldPosition).Add(p.acc.Mul(dt * dt))
	p.PredictedVelocity = p.PredictedPosition.Sub(p.OldPosition).Mul(1 / (2 * dt))
}

func (p *Particle) Update(dt float64) {
	if p.Static {
		p.Velocity = mgl64.Vec2{0, 0}
		p.ForceAccum = mgl64.Vec2{0, 0}
		return
	}

	acc := p.ForceAccum.Mul(1 / p.Mass)
	newPos := p.Position.Mul(2).Sub(p.OldPosition).Add(acc.Mul(dt * dt))
	p.Velocity = newPos.Sub(p.OldPosition).Mul(1 / (2 * dt))

	p.OldPosition = p.Position
	p.Position = newPos
	p.ForceAccum = mgl64.Vec2{0, 0}
}

func (p *Particle) String() string {
	return fmt.Sprintf("Pos:(%.3f,%.3f) Vel:(%.3f,%.3f) M:%.2f R:%.2f",
		p.Position.X(), p.Position.Y(), p.Velocity.X(), p.Velocity.Y(), p.Mass, p.Radius)
}
