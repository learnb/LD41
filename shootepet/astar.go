// Copyright 2018 Bryan Learn
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shootepet

import (
	"fmt"
        "math"
)

type Graph struct {
    graph []int
    xMax int
    yMax int
}

func (g *Graph) Indx2Coord(v int) (int, int) {
        if v == -1 {
            return -1, -1
        }
        return v % g.xMax, v / g.xMax
}


func (g *Graph) Coord2Indx(x, y int) (int) {
        if x == -1 || y == -1 {
            return -1
        }
        return (y * g.xMax) + x
}

func (g *Graph) Print() {
    for i, _ := range g.graph {
        sx, sy := g.Indx2Coord(i)
        print(fmt.Sprintf("(%d, %d)", sx, sy))
        if sx == g.xMax-1 {
            print(fmt.Sprintf("\n"))
        }
    }
}

func (g *Graph) Astar(src int, dst int) []int {
    if g.graph == nil {
        return nil
    }
    //sx, sy := gIndx2Coord(i)

    // init open and closed lists
    closedSet := make([]int, 0)
    openSet := make([]int, 0)
    openSet = append(openSet, src)
    cameFrom := map[int]int{}

    gScore := map[int]float64{}
    gScore[src] = 0.0

    fScore := map[int]float64{}

    var q int // current node
    q = openSet[0]

    // while openSet is not empty
    for len(openSet)>0 {
        // pop
        //_, openSet = openSet[len(openSet)-1], openSet[:len(openSet)-1]

        // select q: find node in openSet with lowest fScore
        smallest := 9000.5
        q = openSet[0]
        for _, v := range openSet {
            score, ok := fScore[v]
            if ok && score <= smallest {
                smallest = score
                q = v
            }
        }

        // Check for end condition
        if q == dst {
	    return reconstructPath(cameFrom, q) // return path
        }

        // delete q from openSet
	qPos := lPos(q, openSet)
        openSet = append(openSet[:qPos], openSet[qPos+1:]...)

        // add q to closedSet
        closedSet = append(closedSet, q)

        // for each neighbor of q
        for _, n := range g.getNeighbors(q) {
            //skip if neighbor in closeSet OR doesn't exist (-1)
            if lContains(closedSet, n) || n == -1 {
                continue
            }

            // discover new node
	    if !lContains(openSet, n) {
	        openSet = append(openSet, n)
	    }

	    // calc gScore
            scoreQ, ok := gScore[q]
	    if !ok {
		scoreQ = 9000.5
	    }
	    t_gScore := scoreQ + g.distance_between(q, n)

            scoreN, ok := gScore[n]
	    if !ok {
		scoreN = 9000.5
	    }
	    if t_gScore >= scoreN {
		continue	// This is not a better path
	    }

	    // This path is better for now; record
	    cameFrom[n] = q
	    gScore[n] = t_gScore
	    fScore[n] = gScore[n] + g.heuristic_cost(n, dst)
        }

    } // end while openSet

    return nil
}

func reconstructPath(cameFrom map[int]int, current int) []int {
    var total_path []int
    total_path = append(total_path, current)
    //while current exists in cameFrom
    for {
        if val, ok := cameFrom[current]; ok {
	    current = val
            total_path = append(total_path, current)
        } else {
            break
        }
    }

    return total_path
}

func (g *Graph) getNeighbors(cell int) ([4]int) {
    var neigh [4]int
    x, y := g.Indx2Coord(cell)
    neigh[0] = g.Coord2Indx(x-1, y) //left
    neigh[1] = g.Coord2Indx(x+1, y) //right
    neigh[2] = g.Coord2Indx(x, y-1) //up
    neigh[3] = g.Coord2Indx(x, y+1) //down

    // out of bounds checks
    if x-1 < 0 { //left
        neigh[0] = -1
    }
    if x+1 >= g.xMax { //right
        neigh[1] = -1
    }
    if y-1 < 0 { //up
        neigh[2] = -1
    }
    if y+1 >= g.yMax { //down
        neigh[3] = -1
    }

    // impassable tile check
    if neigh[0] != -1 && g.graph[neigh[0]] != 0 { // non-zero means impassable
        neigh[0] = -1
    }
    if neigh[1] != -1 && g.graph[neigh[1]] != 0 { // non-zero means impassable
        neigh[1] = -1
    }
    if neigh[2] != -1 && g.graph[neigh[2]] != 0 { // non-zero means impassable
        neigh[2] = -1
    }
    if neigh[3] != -1 && g.graph[neigh[3]] != 0 { // non-zero means impassable
        neigh[3] = -1
    }

    return neigh
}

func (g *Graph) distance_between(src, dst int) float64 {
    x1, y1 := g.Indx2Coord(src)
    x2, y2 := g.Indx2Coord(dst)
    return g.dist(x1, y1, x2, y2)
}

func (g *Graph) heuristic_cost(src, dst int) float64 {
    x1, y1 := g.Indx2Coord(src)
    x2, y2 := g.Indx2Coord(dst)
    return g.dist(x1, y1, x2, y2)
}

func (g *Graph) dist(x1, y1, x2, y2 int) float64 {
    return math.Sqrt( math.Pow((float64(x1)-float64(x2)),2) + math.Pow((float64(y1)-float64(y2)),2) )
}

func lPos(value int, s []int) int {
    for p, v := range s {
        if (v == value) {
            return p
        }
    }
    return -1
}

func lContains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

