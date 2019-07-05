package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.noopschallenge.com"
type Endpoint uint8
const (
	GenerateRandom = iota
	SubmitRandomSolution
	StartRace
	Unknown
)

func endpointName(endpoint Endpoint) string {
	var name string
	switch endpoint {
	case GenerateRandom: 	   name = "mazebot/random"
	case SubmitRandomSolution: name = "mazebot/mazes"
	case StartRace:			   name = "mazebot/race/start"
	case Unknown: 			   name = ""
	}
	return fmt.Sprintf("%s/%s", baseURL, name)
}

type MazeBot struct {}

type MazeBotResponse struct {
	Result string    `json:"result"`
	Message string   `json:"message"`
	ShortestLen int  `json:"shortestSolutionLength"`
	SubmittedLen int `json:"yourSolutionLength"`
	Elapsed int 	 `json:"elapsed"`
	NextMaze string  `json:"nextMaze"`
	Certificate string `json:"certificate"`
	Error string
}

type Certificate struct {
	Message string
	Elapsed float32
	Completed time.Time
}

func (m *MazeBot) CreateMaze() *Maze {
	return m.requestMaze(endpointName(GenerateRandom))
}

func (m *MazeBot) CreateMazeWithSize(size int) *Maze {
	url := fmt.Sprintf("%s?minSize=%d&maxSize=%d", endpointName(GenerateRandom), size, size)
	return m.requestMaze(url)
}

func (m *MazeBot) SubmitRandomMazeSolution(mazeID string, solution string) MazeBotResponse {
	url := fmt.Sprintf("%s/%s", endpointName(SubmitRandomSolution), mazeID)
	data := []byte(fmt.Sprintf(`{"directions": "%s"}`, solution))
	return m.sendRequest(url, data)
}

func (m *MazeBot) StartRace(username string) MazeBotResponse {
	return m.sendRequest(endpointName(StartRace), []byte(fmt.Sprintf(`{"login": "%s"}`, username)))
}

func (m *MazeBot) GetNextMaze(mazePath string) *Maze {
	return m.requestMaze(fmt.Sprintf("%s%s", baseURL, mazePath))
}

func (m *MazeBot) SubmitRaceSolution(mazePath string, solution string) MazeBotResponse {
	url := fmt.Sprintf("%s%s", baseURL, mazePath)
	data := []byte(fmt.Sprintf(`{"directions": "%s"}`, solution))
	return m.sendRequest(url, data)
}

func (m *MazeBot) FetchCertificate(path string) Certificate {
	url := fmt.Sprintf("%s%s", baseURL, path)
	resp, _ := http.Get(url)
	var cert Certificate
	_ = json.NewDecoder(resp.Body).Decode(&cert)
	return cert
}

func (m *MazeBot) sendRequest(url string, data []byte) (result MazeBotResponse) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		result.Error = fmt.Sprintf("Submission failed: %s", err.Error())
		return result
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		result.Error = fmt.Sprintf("Cannot read response: %s", err.Error())
		return result
	}

	return result
}

type mazeResponse struct {
	Name string    `json:"name"`
	Path string    `json:"mazePath"`
	Start [2]int   `json:"startingPosition"`
	End [2]int     `json:"endingPosition"`
	Map [][]string `json:"map"`
}

var mazeNameRegex = regexp.MustCompile("Maze[\\d\\w\\s]*#\\d+ \\((\\d+)x(\\d+)\\)")

func (m *MazeBot) requestMaze(URL string) *Maze {
	resp, err := http.Get(URL)
	if err != nil { return nil }
	if resp.StatusCode != http.StatusOK { return nil }

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var mr mazeResponse
	err = decoder.Decode(&mr)
	if err != nil { return nil }

	matches := mazeNameRegex.FindStringSubmatch(mr.Name)
	w, _ := strconv.Atoi(matches[1])
	h, _ := strconv.Atoi(matches[2])
	parts := strings.Split(mr.Path, "/")

	maze := Maze{
		ID:parts[len(parts)-1],
		Size:Point{w, h},
		Start:Point{mr.Start[0], mr.Start[1]},
		Exit:Point{mr.End[0], mr.End[1]},
		Map:parseCellStrings(mr.Map, w, h),
	}

	return &maze
}

func parseCellStrings(mazeMap [][]string, nRows, nCols int) [][]CellType {
	var parse = func(x string) CellType {
		switch x {
		case "X": return Wall
		case " ": return Empty
		case "A": return Start
		case "B": return Exit
		}
		log.Printf("Unexpected input char: %s. Replacing with empty cell", x)
		return Empty
	}

	var newMap = make([][]CellType, nRows)
	for j := 0; j < nRows; j++ {
		row := make([]CellType, nCols)
		for i := 0; i < nCols; i++ {
			row[i] = parse(mazeMap[i][j])
		}
		newMap[j] = row
	}

	return newMap
}