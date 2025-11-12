package physic

import (
	"time"

	"github.com/gass_ita/go-physics-engine/common"
)

type Scene struct {
	Particles   []*Particle
	Springs     []*Spring
	Dampers     []*Damper
	WorldLimits [4]float64 // xmin, ymin, xmax, ymax
}

func NewScene(WorldLimits [4]float64) *Scene {
	return &Scene{
		Particles:   []*Particle{},
		Springs:     []*Spring{},
		Dampers:     []*Damper{},
		WorldLimits: WorldLimits,
	}
}

func (s *Scene) AddParticle(p *Particle) {
	s.Particles = append(s.Particles, p)
}

func (s *Scene) AddSpring(spring *Spring) {
	s.Springs = append(s.Springs, spring)
}

func (s *Scene) AddDamper(damper *Damper) {
	s.Dampers = append(s.Dampers, damper)
}

func (s *Scene) Update(dt float64) {
	for _, p := range s.Particles {
		p.ApplyExternalForces()
		p.UpdatePredictedState(dt)
	}
	for _, spring := range s.Springs {
		spring.Update(dt)
	}
	for _, damper := range s.Dampers {
		damper.Update(dt)
	}

	// Clamp predicted positions to world limits
	for _, p := range s.Particles {

		p.ClampPosition(s.WorldLimits)
	}

	for _, p := range s.Particles {
		p.Update(dt)
	}
}

func (s *Scene) Start(dt float64, posChan chan<- []common.ParticlePos, linkChan chan<- []common.LinkPos) {
	go func() {
		ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
		defer ticker.Stop()

		for {
			<-ticker.C // wait for next tick
			start := time.Now()
			s.Update(dt)
			elapsed := time.Since(start)
			_ = elapsed // currently not used, but could be logged for performance monitoring

			// Send particle positions (non-blocking)
			select {
			case posChan <- func() []common.ParticlePos {
				positions := make([]common.ParticlePos, len(s.Particles))
				for i, p := range s.Particles {
					positions[i] = common.ParticlePos{
						X:      p.Position.AtVec(0),
						Y:      p.Position.AtVec(1),
						Radius: p.Radius,
					}
				}
				return positions
			}():
			default:
			}
			// Send spring positions (non-blocking)
			select {
			case linkChan <- func() []common.LinkPos {
				positions := make([]common.LinkPos, len(s.Springs)+len(s.Dampers))
				for i, s := range s.Springs {
					positions[i] = common.LinkPos{
						X1: s.P1.Position.AtVec(0),
						Y1: s.P1.Position.AtVec(1),
						X2: s.P2.Position.AtVec(0),
						Y2: s.P2.Position.AtVec(1),
					}
				}
				for _, d := range s.Dampers {
					positions = append(positions, common.LinkPos{
						X1: d.P1.Position.AtVec(0),
						Y1: d.P1.Position.AtVec(1),
						X2: d.P2.Position.AtVec(0),
						Y2: d.P2.Position.AtVec(1),
					})
				}

				return positions
			}():
			default:
			}

		}
	}()
}
