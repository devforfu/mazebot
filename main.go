package main

import (
    "log"
    "maze/generator"
    "maze/renderer/ascii"
    "maze/solver"
    "maze/utils"
    "strings"
    "time"
)

type solution struct {
    maze *generator.Maze
    actions solver.Path
    timeToSolve time.Duration
}

func MazeFactory(ch chan<- *generator.Maze, sizes []int) {
    for _, size := range sizes {
        bot := generator.MazeBot{}
        maze := bot.CreateMazeWithSize(size)
        if maze == nil {
            log.Printf("Failed to create maze with size %d", size)
            continue
        }
        ch <- maze
    }
    close(ch)
}

func Solve(input <-chan *generator.Maze, output chan<- solution) {
    for {
        select {
        case maze := <-input:
            if maze == nil {
                log.Printf("No more mazes to solve")
                close(output)
                return
            }
            s := solution{}
            elapsed := utils.Timer(func() {
                solution := solver.FindPath(maze, solver.ManhattanDistance)
                maze.MarkVisited(solution.Points)
                s.actions = solution.Actions
            })
            s.maze = maze
            s.timeToSolve = elapsed
            output <- s
        case <-time.After(2 * time.Second):
            log.Printf("Timeout...\n")
            close(output)
            return
        }
    }
}

func Submit(input <-chan solution, results chan<- SolvedMaze) {
   bot := generator.MazeBot{}
   for solution := range input {
       actions := convertToString(solution.actions)
       result := bot.SubmitSolution(solution.maze.ID, actions)
       results <- SolvedMaze{
           Maze:solution.maze,
           Actions:actions,
           TimeToSolveLocally: solution.timeToSolve,
           BotResponse:&result}
   }
   close(results)
}

func convertToString(path solver.Path) string {
    var b strings.Builder
    for _, step := range path {
        switch step {
        case solver.North: b.WriteString("N")
        case solver.East:  b.WriteString("E")
        case solver.South: b.WriteString("S")
        case solver.West:  b.WriteString("W")
        }
    }
    return b.String()
}

type SolvedMaze struct {
    Maze *generator.Maze
    Actions string
    TimeToSolveLocally time.Duration
    BotResponse *generator.SubmissionResult
}

func main() {
    mazeStream := make(chan *generator.Maze)
    solutionStream := make(chan solution)
    results := make(chan SolvedMaze)
    sizes := []int{10, 20, 40, 60, 100, 120, 150, 200}

    go MazeFactory(mazeStream, sizes)
    go Solve(mazeStream, solutionStream)
    go Submit(solutionStream, results)

    for result:= range results {
        resp := result.BotResponse
        if resp.Error != "" || resp.Result != "success" {
            var msg string
            if resp.Error != "" {
                msg = resp.Error
            } else {
                msg = resp.Message
            }
            log.Printf("Failed to submit maze solution: %s", msg)
            log.Printf("Proposed solution was: %s", result.Actions)
        } else {
            log.Printf("Maze %s solution accepted!", result.Maze.ID)
            log.Printf("Submitted path length: %d", resp.SubmittedLen)
            log.Printf("Best possible path length: %d", resp.ShortestLen)
            log.Printf("Is shortest? %v", resp.ShortestLen == resp.SubmittedLen)
            log.Printf("Search time: %v", result.TimeToSolveLocally)
            log.Printf("Turnaround time: %v", resp.Elapsed)
        }
        time.Sleep(1*time.Millisecond)
        ascii.DefaultRenderer.Render(result.Maze)
    }
}
