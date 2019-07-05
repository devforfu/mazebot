package solver

import (
    "fmt"
    "math"
    "maze/generator"
)

type Vertex struct {
    Point generator.Point
    Score Score
    Mileage Score
    Index int
    Prev *Vertex
}
func NewVertex(point generator.Point) *Vertex {
    return &Vertex{Point:point, Score:math.MaxInt64, Index:0, Prev:nil}
}
func (v *Vertex) GetX() int { return v.Point.X }
func (v *Vertex) GetY() int { return v.Point.Y }
func (v *Vertex) Priority() int { return int(v.Score + v.Mileage) }
func (v *Vertex) UpdatePriority(value int) { v.Score = Score(value) }
func (v *Vertex) GetIndex() int { return v.Index }
func (v *Vertex) SetIndex(value int) { v.Index = value }
func (v *Vertex) Equal(other *Vertex) bool { return v.GetX() == other.GetX() && v.GetY() == other.GetY() }
func (v *Vertex) String() string {
    return fmt.Sprintf("Vertex(x=%d, y=%d, score=%d)", v.GetX(), v.GetY(), v.Score)
}
