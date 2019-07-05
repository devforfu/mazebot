# MazeBot

Using the MazeBot API to pull the mazes and find the path from A to B.

## Project Structure

- `main.go` - package entry point that pulls mazes and solves them with A-star search
- `generator` - creates `Maze` objects from the JSON responses returned by API
- `renderer/ascii` - prints `Maze` objects into `stdout`
- `solver` - A-star based maze solver
- `utils` - helper functions

## How to run?

```bash
git clone https://github.com/devforfu/mazebot $GOPATH/src/mazebot
cd $GOPATH/src/mazebot
go run main.go
```
