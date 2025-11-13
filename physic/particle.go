package physic

import (
	"fmt"

	"github.com/gass-ita/go-physics-engine/common"
	"gonum.org/v1/gonum/mat"
)

type Particle struct {
	Position mat.VecDense
	Velocity mat.VecDense

	Mass   float64
	Radius float64

	ForceAccum  mat.VecDense
	Static      bool
	ForceFields []*mat.VecDense

	OldPosition mat.VecDense

	PredictedPosition mat.VecDense
	PredictedVelocity mat.VecDense
}

func (p *Particle) ClampPosition(limits [4]float64) {
	if p.Static {
		return
	}

	xmin, xmax, ymin, ymax := limits[0], limits[1], limits[2], limits[3]

	// Posizione X
	if p.Position.AtVec(0) < xmin+p.Radius {
		p.Position.SetVec(0, xmin+p.Radius)
		p.Velocity.SetVec(0, -p.Velocity.AtVec(0)) // rimbalzo
	} else if p.Position.AtVec(0) > xmax-p.Radius {
		p.Position.SetVec(0, xmax-p.Radius)
		p.Velocity.SetVec(0, -p.Velocity.AtVec(0)) // rimbalzo
	}

	// Posizione Y
	if p.Position.AtVec(1) < ymin+p.Radius {
		p.Position.SetVec(1, ymin+p.Radius)
		p.Velocity.SetVec(1, -p.Velocity.AtVec(1))
	} else if p.Position.AtVec(1) > ymax-p.Radius {
		p.Position.SetVec(1, ymax-p.Radius)
		p.Velocity.SetVec(1, -p.Velocity.AtVec(1))
	}
}

// Crea una nuova particella
func NewParticle(pos_m, velocity *mat.VecDense, mass, radius float64, isStatic bool) *Particle {
	dim, _ := pos_m.Dims()

	Position := mat.NewVecDense(dim, nil)
	Position.CopyVec(pos_m)
	Velocity := mat.NewVecDense(dim, nil)
	Velocity.CopyVec(velocity)

	// OldPosition = pos - v*dt
	oldPos := mat.NewVecDense(dim, nil)
	oldPos.ScaleVec(common.DT_PHYSIC, velocity)
	oldPos.SubVec(pos_m, oldPos)

	forceAccum := mat.NewVecDense(dim, nil)
	forceAccum.Zero()

	PredictedPosition := mat.NewVecDense(dim, nil)
	PredictedPosition.Zero()
	PredictedVelocity := mat.NewVecDense(dim, nil)
	PredictedVelocity.Zero()

	if isStatic {
		PredictedPosition.CopyVec(pos_m)
		PredictedVelocity.Zero()
	}

	if isStatic {
		mass = 1
		velocity.Zero()
	}

	p := &Particle{
		Position:          *Position,
		Velocity:          *Velocity,
		Mass:              mass,
		Radius:            radius,
		Static:            isStatic,
		ForceAccum:        *forceAccum,
		OldPosition:       *oldPos,
		ForceFields:       []*mat.VecDense{},
		PredictedPosition: *PredictedPosition,
		PredictedVelocity: *PredictedVelocity,
	}
	return p
}

// Aggiunge una forza alla particella
func (p *Particle) AddForce(force *mat.VecDense) {
	p.ForceAccum.AddVec(&p.ForceAccum, force)
}

// Applica le forze esterne accumulate
func (p *Particle) ApplyExternalForces() {
	p.ForceAccum.Zero()
	if p.Static {
		return
	}

	for _, f := range p.ForceFields {
		p.ForceAccum.AddVec(&p.ForceAccum, f)
	}
}

// Aggiorna lo stato predetto (posizione e velocità) usando Verlet
func (p *Particle) UpdatePredictedState(dt float64) {
	if p.Static {
		return
	}

	acc := mat.NewVecDense(p.ForceAccum.Len(), nil)
	acc.CopyVec(&p.ForceAccum)
	acc.ScaleVec(1/p.Mass, acc)

	// predictedPosition = 2*Position - OldPosition + a*dt^2
	p.PredictedPosition.ScaleVec(2, &p.Position)
	p.PredictedPosition.SubVec(&p.PredictedPosition, &p.OldPosition)
	acc.ScaleVec(dt*dt, acc)
	p.PredictedPosition.AddVec(&p.PredictedPosition, acc)

	// predictedVelocity = (predictedPosition - OldPosition) / (2*dt)
	p.PredictedVelocity.SubVec(&p.PredictedPosition, &p.OldPosition)
	p.PredictedVelocity.ScaleVec(1/(2*dt), &p.PredictedVelocity)
}

// Aggiorna lo stato reale della particella usando Verlet
func (p *Particle) Update(dt float64) {
	if p.Static {
		p.Velocity.Zero()
		p.ForceAccum.Zero()
		return
	}

	acc := mat.NewVecDense(p.ForceAccum.Len(), nil)
	acc.CopyVec(&p.ForceAccum)
	acc.ScaleVec(1/p.Mass, acc)

	// calcola nuova posizione: x_new = 2*x - x_old + a*dt^2
	newPos := mat.NewVecDense(p.Position.Len(), nil)
	newPos.ScaleVec(2, &p.Position)
	newPos.SubVec(newPos, &p.OldPosition)
	acc.ScaleVec(dt*dt, acc)
	newPos.AddVec(newPos, acc)

	// aggiorna velocità: v = (x_new - x_old) / (2*dt)
	p.Velocity.SubVec(newPos, &p.OldPosition)
	p.Velocity.ScaleVec(1/(2*dt), &p.Velocity)

	// aggiorna posizioni
	p.OldPosition.CopyVec(&p.Position)
	p.Position.CopyVec(newPos)

	// reset forza accumulata
	p.ForceAccum.Zero()
}

// Stampa lo stato della particella
func (p *Particle) String() string {
	return fmt.Sprintf("Position: %v, Velocity: %v, Mass: %v, Radius: %v",
		p.Position, p.Velocity, p.Mass, p.Radius)
}
