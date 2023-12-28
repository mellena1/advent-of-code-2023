package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mellena1/advent-of-code-2023/utils"
)

func main() {
	f := utils.ReadFile("input.txt")
	defer f.Close()

	hailstones := parseHailstones(f)
	fmt.Printf("Part one solution: %d\n", hailstones.NumIntersections2D(200000000000000, 400000000000000))
}

type Hailstones []Hailstone

func (h Hailstones) NumIntersections2D(minPos, maxPos int) int {
	type intersection struct {
		h1Idx int
		h2Idx int
		t1    float64
		t2    float64
	}

	possibleIntersections := []intersection{}

	for i, h1 := range h {
		for j, h2 := range h {
			if i <= j {
				continue
			}

			t1, t2, doIntersect := intersectionOf2DVectors([]int{h1.Pos.X, h1.Pos.Y}, []int{h1.Vel.X, h1.Vel.Y}, []int{h2.Pos.X, h2.Pos.Y}, []int{h2.Vel.X, h2.Vel.Y})
			if !doIntersect || t1 < 0 || t2 < 0 {
				continue
			}

			pos := h1.PosAtTime(t1)
			if int(pos.X) < minPos || int(pos.X) > maxPos || int(pos.Y) < minPos || int(pos.Y) > maxPos {
				continue
			}

			possibleIntersections = append(possibleIntersections, intersection{h1Idx: i, h2Idx: j, t1: t1, t2: t2})
		}
	}

	return len(possibleIntersections)
}

type Velocities utils.Coordinate3D[int]

type Hailstone struct {
	Pos utils.Coordinate3D[int]
	Vel Velocities
}

func (h Hailstone) PosAtTime(t float64) utils.Coordinate3D[float64] {
	return utils.NewCoordinate3D(float64(h.Pos.X)+float64(h.Vel.X)*t, float64(h.Pos.Y)+float64(h.Vel.Y)*t, float64(h.Pos.Z)+float64(h.Vel.Z)*t)
}

func parseHailstones(r io.Reader) Hailstones {
	hailstones := []Hailstone{}

	utils.ExecutePerLine(r, func(line string) error {
		pos, vel, _ := strings.Cut(line, " @ ")

		posInts, err := posOrVelToInts(pos)
		if err != nil {
			return err
		}

		velInts, err := posOrVelToInts(vel)
		if err != nil {
			return err
		}

		hailstones = append(hailstones, Hailstone{
			Pos: utils.NewCoordinate3D(posInts[0], posInts[1], posInts[2]),
			Vel: Velocities(utils.NewCoordinate3D(velInts[0], velInts[1], velInts[2])),
		})

		return nil
	})

	return hailstones
}

func posOrVelToInts(s string) ([]int, error) {
	spl := strings.Split(s, ",")
	tSpl := utils.SliceMap(spl, func(v string) string {
		return strings.TrimSpace(v)
	})
	return utils.StrSliceToIntSlice(tSpl)
}

func intersectionOf2DVectors(p1 []int, v1 []int, p2 []int, v2 []int) (float64, float64, bool) {
	c := []int{p2[0] - p1[0], p2[1] - p1[1]}

	dt := float64(c[1]*v2[0] - c[0]*v2[1])
	du := float64(v1[0]*c[1] - v1[1]*c[0])
	d := float64(v1[1]*v2[0] - v1[0]*v2[1])

	return dt / d, du / d, d != 0
}
