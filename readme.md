🧠 Go Physics Engine
====================

A simple **2D physics engine** written in Go, featuring **particles, springs, dampers**, and a real-time **Ebiten-based visualizer**.\
It uses **Verlet integration** for stable and smooth physics simulation.

* * * * *

🚀 Features
-----------

- 🟢 **Particles** with position, velocity, radius, and mass

- ⚙️ **Verlet integration** for numerically stable motion

- 🧲 **Springs and dampers** connecting particles dynamically

- 🌍 **World boundaries** --- particles stay inside the universe (`ClampPosition`)

- 🧮 **External forces** (gravity, custom fields, etc.)

- 🎨 **Real-time visualization** with [Ebiten](https://ebiten.org/)

- ⚡ Channel-based data sharing between physics and rendering threads

* * * * *

⚙️ Installation
---------------

1. **Clone the repo**

    `git clone https://github.com/gass-ita/go-physics-engine.git
    cd go-physics-engine`

2. **Install dependencies**

    `go mod tidy`

3. **Run the simulation**

    `go run ./...`

* * * * *

🧱 Dependencies
---------------

- [Gonum](https://gonum.org/) -- for vector and matrix operations

- [Ebiten](https://ebiten.org/) -- for real-time 2D rendering

- [Go 1.22+](https://go.dev/)

Install them automatically via:

`go mod tidy`

* * * * *

💡 Future Improvements
----------------------

- 🌪️ Add friction and air resistance

- ⚡ Collision detection between particles

- 🧭 User interaction (drag/move particles)

- 🖼️ GUI controls for simulation parameters
