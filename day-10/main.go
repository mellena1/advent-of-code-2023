package main

import (
	"fmt"
	"io"
	"slices"

	"github.com/mellena1/advent-of-code-2023/utils"
)

func main() {
	f := utils.ReadFile("input.txt")
	defer f.Close()

	grid, sLocation := parseGrid(f)

	loop := findLoop(grid, sLocation)

	fmt.Printf("Part one solution: %d\n", loop.furthestFromStart())
	fmt.Printf("Part two solution: %d\n", grid.findAreaInsideLoop(loop))
}

type Node struct {
	attachedNodes []*Node
	coor          Coordinate
	r             rune
}

func NewNode(coor Coordinate, r rune) *Node {
	return &Node{
		attachedNodes: []*Node{},
		coor:          coor,
		r:             r,
	}
}

func (n *Node) nextNode(last *Node) *Node {
	for _, attached := range n.attachedNodes {
		if attached != last {
			return attached
		}
	}
	// should never get here
	return nil
}

func (n *Node) furthestFromStart() int {
	lastLeft, lastRight := n, n
	left, right := n.attachedNodes[0], n.attachedNodes[1]

	steps := 1
	for left != right {
		newLeft := left.nextNode(lastLeft)
		lastLeft, left = left, newLeft

		newRight := right.nextNode(lastRight)
		lastRight, right = right, newRight

		steps++
	}

	return steps
}

type Direction Coordinate

var (
	LEFT = Direction{
		x: -1,
		y: 0,
	}
	RIGHT = Direction{
		x: 1,
		y: 0,
	}
	UP = Direction{
		x: 0,
		y: -1,
	}
	DOWN = Direction{
		x: 0,
		y: 1,
	}
)

var PossibleNextPipes = map[Direction][]rune{
	UP:    {'F', '7', '|'},
	DOWN:  {'L', 'J', '|'},
	LEFT:  {'F', 'L', '-'},
	RIGHT: {'7', 'J', '-'},
}

var PipeAttachments = map[rune][]Direction{
	'|': {UP, DOWN},
	'-': {LEFT, RIGHT},
	'L': {UP, RIGHT},
	'J': {LEFT, UP},
	'7': {LEFT, DOWN},
	'F': {RIGHT, DOWN},
	'S': {UP, DOWN, LEFT, RIGHT},
}

func findLoop(grid Grid, startingLocation Coordinate) *Node {
	sNode := NewNode(startingLocation, grid.get(startingLocation).r)
	foundLoop := false

	var dfs func(curNode *Node, coor Coordinate, lastCoor Coordinate) bool
	dfs = func(curNode *Node, coor Coordinate, lastCoor Coordinate) bool {
		if foundLoop {
			return false
		}

		gridNode := grid.get(coor)

		for _, direction := range PipeAttachments[gridNode.r] {
			newCoor := Coordinate{
				x: coor.x + direction.x,
				y: coor.y + direction.y,
			}

			if newCoor.x < 0 || newCoor.x >= len(grid[0]) || newCoor.y < 0 || newCoor.y >= len(grid) {
				// out of bounds
				continue
			}
			if newCoor.x == lastCoor.x && newCoor.y == lastCoor.y {
				// don't go backwards
				continue
			}

			neighbor := grid.get(newCoor)

			if neighbor.r == 'S' {
				foundLoop = true
				curNode.attachedNodes = append(curNode.attachedNodes, sNode)
				sNode.attachedNodes = append(sNode.attachedNodes, curNode)
				return true
			}

			if slices.Contains(PossibleNextPipes[direction], neighbor.r) {
				newNode := NewNode(newCoor, neighbor.r)
				foundS := dfs(newNode, newCoor, coor)

				if foundS {
					curNode.attachedNodes = append(curNode.attachedNodes, newNode)
					newNode.attachedNodes = append(newNode.attachedNodes, curNode)
					return true
				}
			}
		}

		return false
	}

	dfs(sNode, startingLocation, Coordinate{x: -1, y: -1})

	sNode.r = findPipeUnderS(sNode)

	return sNode
}

func findPipeUnderS(sNode *Node) rune {
	sDirections := []Direction{}

	for _, attachment := range sNode.attachedNodes {
		dir := sNode.coor.findDirection(attachment.coor)
		sDirections = append(sDirections, dir)
	}
OUTER:
	for pipe, directions := range PipeAttachments {
		for _, d := range directions {
			if !slices.Contains(sDirections, d) {
				continue OUTER
			}
		}
		return pipe
	}

	panic("invalid pipe arrangement")
}

func isCorner(r rune) bool {
	switch r {
	case 'L', 'J', '7', 'F':
		return true
	}
	return false
}

type Grid [][]GridNode

func (g Grid) get(coor Coordinate) *GridNode {
	return &g[coor.y][coor.x]
}

func (g Grid) markLoopNodes(startNode *Node) {
	g.get(startNode.coor).inLoop = true

	lastNode := startNode
	curNode := startNode.attachedNodes[0]
	for curNode != startNode {
		g.get(curNode.coor).inLoop = true

		newNode := curNode.nextNode(lastNode)
		lastNode, curNode = curNode, newNode
	}
}

var Rotations = map[rune]map[Direction]Direction{
	'7': {
		DOWN:  LEFT,
		RIGHT: UP,
	},
	'J': {
		LEFT: UP,
		DOWN: RIGHT,
	},
	'L': {
		LEFT: DOWN,
		UP:   RIGHT,
	},
	'F': {
		RIGHT: DOWN,
		UP:    LEFT,
	},
}

func (g Grid) findAreaInsideLoop(startNode *Node) int {
	g.markLoopNodes(startNode)

	// replace S with actual val
	g.get(startNode.coor).r = startNode.r

	// should always be a F
	leftCorner := g.findTopLeft(startNode)

	lastNode := leftCorner
	curNode := leftCorner
	for _, attach := range curNode.attachedNodes {
		if attach.coor.x > curNode.coor.x {
			curNode = attach
			break
		}
	}

	markNodesAsTouched := func(curNode *Node, dir Direction) {
		switch dir {
		case DOWN:
			for i := curNode.coor.y + 1; i < len(g); i++ {
				gNode := g.get(Coordinate{
					x: curNode.coor.x,
					y: i,
				})

				if gNode.inLoop {
					break
				}

				gNode.touched.Down = true
			}
		case UP:
			for i := curNode.coor.y - 1; i >= 0; i-- {
				gNode := g.get(Coordinate{
					x: curNode.coor.x,
					y: i,
				})

				if gNode.inLoop {
					break
				}

				gNode.touched.Up = true
			}
		case LEFT:
			for i := curNode.coor.x - 1; i >= 0; i-- {
				gNode := g.get(Coordinate{
					x: i,
					y: curNode.coor.y,
				})

				if gNode.inLoop {
					break
				}

				gNode.touched.Left = true
			}
		case RIGHT:
			for i := curNode.coor.x + 1; i < len(g[curNode.coor.y]); i++ {
				gNode := g.get(Coordinate{
					x: i,
					y: curNode.coor.y,
				})

				if gNode.inLoop {
					break
				}

				gNode.touched.Right = true
			}
		}
	}

	curInsideDirection := DOWN

	for curNode != leftCorner {
		markNodesAsTouched(curNode, curInsideDirection)

		if isCorner(curNode.r) {
			curInsideDirection = Rotations[curNode.r][curInsideDirection]
			markNodesAsTouched(curNode, curInsideDirection)
		}

		newNode := curNode.nextNode(lastNode)
		lastNode, curNode = curNode, newNode
	}

	area := 0
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[i]); j++ {
			gNode := g.get(Coordinate{
				x: j,
				y: i,
			})
			if gNode.touched.allTouched() {
				area++
			}
		}
	}

	return area
}

func (g Grid) findTopLeft(loop *Node) *Node {
	for i := 0; i < len(g); i++ {
		for j := 0; j < len(g[i]); j++ {
			if g[i][j].inLoop {
				lastNode := loop
				curNode := loop.attachedNodes[0]
				if loop.coor.x == j && loop.coor.y == i && isCorner(loop.r) {
					return loop
				}

				for lastNode != curNode {
					if curNode.coor.x == j && curNode.coor.y == i && isCorner(curNode.r) {
						return curNode
					}

					nextNode := curNode.nextNode(lastNode)
					lastNode, curNode = curNode, nextNode
				}
			}
		}
	}

	panic("there should be a top left node")
}

type GridNode struct {
	r       rune
	touched *TouchedFromDirections
	inLoop  bool
}

type TouchedFromDirections struct {
	Up    bool
	Down  bool
	Right bool
	Left  bool
}

func (t *TouchedFromDirections) allTouched() bool {
	return t.Up && t.Down && t.Right && t.Left
}

type Coordinate struct {
	x int
	y int
}

func (c Coordinate) findDirection(c2 Coordinate) Direction {
	return Direction{
		x: c2.x - c.x,
		y: c2.y - c.y,
	}
}

func parseGrid(r io.Reader) (Grid, Coordinate) {
	grid := Grid{}
	var coor Coordinate

	y := 0
	utils.ExecutePerLine(r, func(line string) error {
		nodes := make([]GridNode, len(line))
		for i, r := range line {
			nodes[i] = GridNode{
				r:       r,
				touched: &TouchedFromDirections{},
			}
			if r == 'S' {
				coor = Coordinate{
					x: i,
					y: y,
				}
			}
		}
		y++

		grid = append(grid, nodes)
		return nil
	})

	return grid, coor
}
