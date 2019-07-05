package generator

import "fmt"

type Point struct { X, Y int }
func (p Point) String() string { return fmt.Sprintf("(%d, %d)", p.GetX(), p.GetY()) }
func (p Point) GetX() int { return p.X }
func (p Point) GetY() int { return p.Y }
func (p Point) XY() (int, int) { return p.X, p.Y }

type CellType uint8
const (
	Empty CellType = iota
	Wall
	Start
	Exit
	Visited
)

type Maze struct {
	ID string
	Size Point
	Start Point
	Exit Point
	Map [][]CellType
}

func (m *Maze) MarkVisited(points []Point) {
	for _, p := range points {
		x, y := p.XY()
		if m.Map[x][y] == Empty {
			m.Map[x][y] = Visited
		}
	}
}