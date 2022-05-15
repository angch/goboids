package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	emptyImage = ebiten.NewImage(3, 3)

	// emptySubImage is an internal sub image of emptyImage.
	// Use emptySubImage at DrawTriangles instead of emptyImage in order to avoid bleeding edges.
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

const (
	screenWidth  = 640
	screenHeight = 480
)

type vec struct {
	X, Y float32
}

type BoidGame struct {
	counter int
	boids   []*Boid
	center  vec
}

type Boid struct {
	x, y                           float32
	dx, dy                         float32
	r, v                           float32 // recalc based on dx, dy
	colorR, colorG, colorB, colorA float32
}

func (g *BoidGame) Update() error {
	g.counter++
	step := float32(0.01)

	for _, b := range g.boids {
		centerOfMass := vec{X: 0, Y: 0}

		for _, b2 := range g.boids {
			if b2 == b {
				continue
			}
			centerOfMass.X += b2.x
			centerOfMass.Y += b2.y
		}
		nboids := len(g.boids) - 1
		centerOfMass.X /= float32(nboids)
		centerOfMass.Y /= float32(nboids)

		// Rotate right
		if false {
			b.r += math.Pi / 180
			if b.r > 2*math.Pi {
				b.r -= 2 * math.Pi
			}
		}

		// Increment position by velocity
		if false {
			x0, y0 := float32(1), float32(0)
			cos0 := float32(math.Cos(float64(b.r - math.Pi/2)))
			sin0 := float32(math.Sin(float64(b.r - math.Pi/2)))
			b.dx += (x0*cos0 - y0*sin0) * b.v
			b.dy += (x0*sin0 + y0*cos0) * b.v
		}

		// Rule 1: Boids try to fly towards the centre of mass of neighbouring boids.
		if true {
			b.dx += (centerOfMass.X - b.x) * float32(step)
			b.dy += (centerOfMass.Y - b.y) * float32(step)
		}

		// Rule 2: Boids try to keep a small distance away from other objects (including other boids).
		if true {
			for _, other := range g.boids {
				if other == b {
					continue
				}
				dx := float32(other.x - b.x)
				dy := float32(other.y - b.y)
				d := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				if d < 10 {
					if d == 0 {
						d = 0.01
					}
					b.dx -= dx / d
					b.dy -= dy / d
				}
			}
		}

		// Rule 3: Boids try to match velocity with near boids.
		if true {
			v := vec{}
			for _, other := range g.boids {
				if other == b {
					continue
				}
				v.X += other.dx
				v.Y += other.dy
			}
			v.X /= float32(len(g.boids) - 1)
			v.Y /= float32(len(g.boids) - 1)
		}

		// Rule 4: Center of screen pull
		if true {
			b.dx += (g.center.X - b.x) * step / 50
			b.dy += (g.center.Y - b.y) * step / 50
		}

		// Rule 6: Bounds
		if true {
			if b.x+b.dx < 0 {
				b.dx *= 0.5
			}
			if b.x+b.dx > screenWidth {
				b.dx *= 0.5
			}
			if b.y+b.dy < 0 {
				b.dy *= 0.5
			}
			if b.y+b.dy > screenHeight {
				b.dy *= 0.5
			}
		}

		// Rule 5: Speed limit
		speedlimit := float32(5)
		if true {
			d := float32(math.Sqrt(float64(b.dx*b.dx + b.dy*b.dy)))
			if d > speedlimit {
				b.dx = b.dx / d * speedlimit
				b.dy = b.dy / d * speedlimit
			}
		}

		// Apply, and update heading and velocity
		b.x += b.dx
		b.y += b.dy
		b.r = float32(math.Atan2(float64(b.dy), float64(b.dx)))
		b.v = float32(math.Sqrt(float64(b.dx*b.dx + b.dy*b.dy)))
	}

	return nil
}

func (g *BoidGame) DrawBoid(screen *ebiten.Image, b *Boid) {
	if b == nil {
		return
	}
	var size = float32(8)
	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: -size},
		{DstX: size / 2, DstY: size / 2},
		{DstX: 0, DstY: 0},
		{DstX: -size / 2, DstY: +size / 2},
	}

	for k := range vertices {
		x0, y0 := vertices[k].DstX, vertices[k].DstY
		cos0 := float32(math.Cos(float64(b.r + math.Pi/2)))
		sin0 := float32(math.Sin(float64(b.r + math.Pi/2)))
		x1 := b.x + x0*cos0 - y0*sin0
		y1 := b.y + x0*sin0 + y0*cos0

		vertices[k].DstX = x1
		vertices[k].DstY = y1
		vertices[k].ColorR = b.colorR
		vertices[k].ColorG = b.colorG
		vertices[k].ColorB = b.colorB
		vertices[k].ColorA = b.colorA
	}

	screen.DrawTriangles(vertices, []uint16{0, 1, 2, 0, 2, 3}, emptySubImage, op)
}

func (g *BoidGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	for _, b := range g.boids {
		g.DrawBoid(screen, b)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}

func (g *BoidGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &BoidGame{
		counter: 0,
		boids:   []*Boid{
			// {100, 100, 0.5, 1, 0, 0},
			// {110, 100, 0.6, 1, 0, 0},
			// {110, 110, 0.7, 1, 0, 0},
			// {120, 100, 0.7, 1, 0, 0},
		},
		center: vec{X: float32(screenWidth) / 2, Y: float32(screenHeight) / 2},
	}
	for n := 0; n < 50; n++ {
		g.boids = append(g.boids, &Boid{
			x:      float32(screenWidth)*rand.Float32()/8 + float32(screenWidth)/2,
			y:      float32(screenHeight)*rand.Float32()/8 + float32(screenHeight)/2,
			dx:     rand.Float32()*10 - 5,
			dy:     rand.Float32()*10 - 5,
			r:      0,
			v:      0,
			colorR: rand.Float32()/2 + 0.5,
			colorG: rand.Float32()/2 + 0.5,
			colorB: rand.Float32()/2 + 0.5,
			colorA: 1,
		})
	}

	ebiten.SetMaxTPS(60)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
