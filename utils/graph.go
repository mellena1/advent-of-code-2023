package utils

import (
	"math"
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
	prev := map[K]K{}

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
				prev[neighbor] = curNode
			}
		}
	}

	return distances[destination]
}
