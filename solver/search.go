package solver

import (
	"fmt"
	"math"
	"maze/generator"
	"maze/solver/collections"
)

type Direction int

type Path []Direction

const (
	North Direction = iota
	East
	South
	West
)

type Solution struct {
	Actions Path
	Points []generator.Point
}

func FindPath(m *generator.Maze, h Heuristic) *Solution {
	goal := Solve(m, h)
	if goal == nil { return nil }
	solution := ReconstructPath(goal)
	return solution
}

func Solve(m *generator.Maze, h Heuristic) *Vertex {
	queue := collections.PriorityQueue{NewVertex(m.Start)}
	goal := NewVertex(m.Exit)
	visited := map[generator.Point]bool{}
	for !queue.Empty() {
		curr := queue.Pop().(*Vertex)
		visited[curr.Point] = true
		if curr.Equal(goal) {
			return curr
		}
		for _, neighbour := range getNeighbours(m, curr) {
			_, wasVisited := visited[neighbour.Point]
			if wasVisited { continue }
			neighbour.Score = h(neighbour, goal)
			queue.Push(neighbour)
		}
	}
	return nil
}

func ReconstructPath(last *Vertex) *Solution {
	curr, prev := last, last.Prev
	reversedPath := make(Path, 1)
	points := make([]generator.Point, 1)
	points[0] = curr.Point
	for prev != nil {
		x1, y1 := curr.Point.XY()
		x0, y0 := prev.Point.XY()
		if (x1 - x0) == 1 {
			reversedPath = append(reversedPath, East)
		} else if (x1 - x0) == -1 {
			reversedPath = append(reversedPath, West)
		} else if (y1 - y0) == 1 {
			reversedPath = append(reversedPath, South)
		} else if (y1 - y0) == -1 {
			reversedPath = append(reversedPath, North)
		} else {
			panic("Invalid direction!")
		}
		points = append(points, prev.Point)
		curr, prev = prev, prev.Prev
	}
	n := len(reversedPath)
	path := make(Path, n)
	for i, x := range reversedPath {
		path[n-i-1] = x
	}
	path = path[0:n-1]
	return &Solution{path, points}
}

func getNeighbours(m *generator.Maze, v *Vertex) (vs []*Vertex) {
	maxW, maxH := m.Size.XY()
	i, j := v.GetX(), v.GetY()
	steps := [][]int{{0, 1}, {1, 0}, {-1, 0}, {0, -1}}
	for _, pair := range steps {
		dx, dy := pair[0], pair[1]
		x, y := i+dx, j+dy
		if x >= 0 && x < maxW && y >= 0 && y < maxH {
			p := m.Map[x][y]
			if p == generator.Exit || p == generator.Empty {
				q := NewVertex(generator.Point{x,y})
				q.Prev = v
				q.Mileage = v.Mileage + 1
				vs = append(vs, q)
			}
		}
	}
	return vs
}

func intAbs(x int) int {
	return int(math.Abs(float64(x)))
}

func showQueue(pq collections.PriorityQueue) {
	for _, item := range pq {
		v := item.(*Vertex)
		fmt.Printf("(i=%d, p=%d, pos=%v)\n", item.GetIndex(), item.Priority(), v.Point)
	}
}
