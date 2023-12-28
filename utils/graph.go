package utils

import (
	"math"
	"slices"
)

type ConnectionMap[K comparable] map[K]map[K]int

func (cMap ConnectionMap[K]) Dijkstra(source K) (map[K]int, map[K]K) {
	distances := make(map[K]int, len(cMap))
	pq := NewPriorityQueue[K, int]()
	prev := map[K]K{}

	for vertex := range cMap {
		distances[vertex] = math.MaxInt
	}
	distances[source] = 0

	pq.Init(distances)

	for pq.Len() > 0 {
		curNode, curNodeDist := pq.Pop()

		for neighbor, dist := range cMap[curNode] {
			alt := curNodeDist + dist
			if alt < distances[neighbor] {
				distances[neighbor] = alt
				pq.Update(neighbor, alt)
				prev[neighbor] = curNode
			}
		}
	}

	return distances, prev
}

func (cMap ConnectionMap[K]) DijkstraWithDest(source K, destination K) int {
	distances := make(map[K]int, len(cMap))
	pq := NewPriorityQueue[K, int]()

	for vertex := range cMap {
		distances[vertex] = math.MaxInt
	}
	distances[source] = 0

	for vertex, distance := range distances {
		pq.Push(vertex, distance)
	}

	for pq.Len() > 0 {
		curNode, curNodeDist := pq.Pop()

		if curNode == destination {
			break
		}

		for neighbor, dist := range cMap[curNode] {
			alt := curNodeDist + dist
			if alt < distances[neighbor] {
				distances[neighbor] = alt
				pq.Update(neighbor, alt)
			}
		}
	}

	return distances[destination]
}

func (cMap ConnectionMap[K]) LongestDijkstraWithDest(source K, destination K) (int, []K) {
	distances := make(map[K]int, len(cMap))
	pq := NewPriorityQueue[K, int]()
	pq.SetPriorityOrder(true)
	prev := map[K]K{}

	for vertex := range cMap {
		distances[vertex] = math.MinInt
	}
	distances[source] = 0

	for vertex, distance := range distances {
		pq.Push(vertex, distance)
	}

	for pq.Len() > 0 {
		curNode, curNodeDist := pq.Pop()

		for neighbor, dist := range cMap[curNode] {
			if nodeIsAlreadyInPath(prev, curNode, neighbor) {
				continue
			}

			alt := curNodeDist + dist
			if alt > distances[neighbor] {
				distances[neighbor] = alt
				pq.Update(neighbor, alt)
				prev[neighbor] = curNode
			}
		}
	}

	return distances[destination], reconstructPath(prev, source, destination)
}

func nodeIsAlreadyInPath[K comparable](prev map[K]K, source, node K) bool {
	curNode := source
	ok := true

	for ok {
		curNode, ok = prev[curNode]

		if ok && curNode == node {
			return true
		}
	}

	return false
}

func reconstructPath[K comparable](prev map[K]K, source, dest K) []K {
	path := []K{dest}

	curNode := dest
	for curNode != source {
		curNode = prev[curNode]
		path = append(path, curNode)
	}

	slices.Reverse(path)
	return path
}
