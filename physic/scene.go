package physic

import (
	"time"

	"github.com/gass-ita/go-physics-engine/common"
)

type Scene struct {
	Particles   []*Particle
	Springs     []*Spring
	Dampers     []*Damper
	WorldLimits [4]float64 // xmin, xmax, ymin, ymax
}

func NewScene(world [4]float64) *Scene {
	return &Scene{
		Particles:   []*Particle{},
		Springs:     []*Spring{},
		Dampers:     []*Damper{},
		WorldLimits: world,
	}
}

func (s *Scene) AddParticle(p *Particle) { s.Particles = append(s.Particles, p) }
func (s *Scene) AddSpring(sp *Spring)    { s.Springs = append(s.Springs, sp) }
func (s *Scene) AddDamper(dp *Damper)    { s.Dampers = append(s.Dampers, dp) }

func (s *Scene) Update(dt float64) {
	for _, p := range s.Particles {
		p.ApplyExternalForces()
		p.UpdatePredictedState(dt)
	}
	for _, sp := range s.Springs {
		sp.Update(dt)
	}
	for _, dp := range s.Dampers {
		dp.Update(dt)
	}

	for _, p := range s.Particles {
		p.ClampPosition(s.WorldLimits)
		p.Update(dt)
	}
}

func (s *Scene) Start(dt float64, posChan chan<- []common.ParticlePos,
	linkChan chan<- []common.LinkPos, infoChan chan<- common.Info) {

	go func() {
		ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
		defer ticker.Stop()

		for {
			<-ticker.C
			start := time.Now()
			s.Update(dt)
			elapsed := time.Since(start).Milliseconds()

			select {
			case posChan <- func() []common.ParticlePos {
				out := make([]common.ParticlePos, len(s.Particles))
				for i, p := range s.Particles {
					out[i] = common.ParticlePos{
						X:      p.Position.X(),
						Y:      p.Position.Y(),
						Radius: p.Radius,
					}
				}
				return out
			}():
			default:
			}

			select {
			case linkChan <- func() []common.LinkPos {
				links := make([]common.LinkPos, 0, len(s.Springs)+len(s.Dampers))
				for _, sp := range s.Springs {
					links = append(links, common.LinkPos{
						X1: sp.P1.Position.X(), Y1: sp.P1.Position.Y(),
						X2: sp.P2.Position.X(), Y2: sp.P2.Position.Y(),
					})
				}
				for _, dp := range s.Dampers {
					links = append(links, common.LinkPos{
						X1: dp.P1.Position.X(), Y1: dp.P1.Position.Y(),
						X2: dp.P2.Position.X(), Y2: dp.P2.Position.Y(),
					})
				}
				return links
			}():
			default:
			}

			select {
			case infoChan <- common.Info{ElapsedTime: float64(elapsed)}:
			default:
			}
		}
	}()
}
