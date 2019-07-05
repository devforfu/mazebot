package solver

type Score int

type Heuristic func(state, goal *Vertex) Score

func ManhattanDistance(state, goal *Vertex) Score {
    dx := intAbs(state.GetX() - goal.GetX())
    dy := intAbs(state.GetY() - goal.GetY())
    return Score(dx + dy)
}
