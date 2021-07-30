package main

import (
	"TGPL-exercise-solutions/chapter6/geometry"
	"fmt"
)

func main() {
	p := geometry.Point{X: 0, Y: 0}
	q := geometry.Point{X: 1, Y: 0}
	fmt.Printf("The distance between %v and %v is %v\n", p, q, p.Distance(q))
	fmt.Printf("The distance between %v and %v is %v\n", q, p, geometry.Distance(q, p))

	var path geometry.Path
	r := geometry.Point{X: 0, Y: 1}
	path = append(path, p, q, r)
	fmt.Printf("The distance on path %v is %v\n", path, path.Distance())
}
