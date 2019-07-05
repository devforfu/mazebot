package utils

import (
    "mazebot/solver"
    "strings"
)

func ConvertToString(path solver.Path) string {
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

