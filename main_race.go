package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "mazebot/generator"
    "mazebot/solver"
    "mazebot/utils"
    "os"
)

var usernamePtr = flag.String("username", "", "GitHub username to authenticate in race")

func init() {
    flag.Parse()
}

func main() {
    bot := generator.MazeBot{}
    username := *usernamePtr
    if username == "" {
        fmt.Printf("GitHub username is not provided!")
        os.Exit(1)
    }

    response := bot.StartRace(username)
    if response.NextMaze == "" {
        fmt.Printf("Cannot start race: make sure the username is valid: %s", username)
        os.Exit(1)
    }

    fmt.Printf("The race has started! 🏎\n")

    var lap int

    for {
        lap += 1
        fmt.Printf("Running lap %d", lap)

        maze := bot.GetNextMaze(response.NextMaze)
        solution := solver.FindPath(maze, solver.ManhattanDistance)
        actions := utils.ConvertToString(solution.Actions)
        response = bot.SubmitRaceSolution(response.NextMaze, actions)

        if response.Result == "success" || response.Result == "finished" {
            fmt.Printf("\rThe lap %d is finished! 🏁\n", lap)

            if response.Result == "finished" {
                fmt.Printf("%s 🎉\n", response.Message)
                cert := bot.FetchCertificate(response.Certificate)
                data, _ := json.Marshal(cert)
                file, _ := os.Create("cert.json")
                _, _ = file.Write(data)
                os.Exit(0)
            }
        } else {
            fmt.Println(response.Message)
            fmt.Println("Failed to finish the race ⛔")
        }
    }
}
