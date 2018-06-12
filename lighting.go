package main

import (
	// "fmt"
	"math"
)

const ()

type Light struct {
	color    Color
	location []float64
}

func CalcLighting(p0, p1, p2, Ka, Kd, Ks, view []float64, Ia Color, lights map[string]Light) Color {
	c := Color{
		r: 0,
		g: 0,
		b: 0,
	}
	ambient := calcAmbientLighting(Ia, Ka)
	c.r += ambient.r
	c.g += ambient.g
	c.b += ambient.b
	for _, l := range lights {
		diffuse := calcDiffuseLighting(p0, p1, p2, Kd, l)
		spec := calcSpecularLighting(p0, p1, p2, view, Ks, l)
		c.r += diffuse.r + spec.r
		c.g += diffuse.g + spec.g
		c.b += diffuse.b + spec.b
	}
	c.Limit()
	return c
}

func calcAmbientLighting(Ia Color, Ka []float64) Color {
	return Color{
		r: int(float64(Ia.r) * Ka[0]),
		g: int(float64(Ia.g) * Ka[1]),
		b: int(float64(Ia.b) * Ka[2]),
	}
}

func calcDiffuseLighting(p0, p1, p2, Kd []float64, l Light) Color {
	// TODO: Error checking lol
	n, _ := Normal(p0, p1, p2)
	n = Normalize(n)
	light := Normalize(l.location)
	diffuse, _ := DotProduct(n, light)
	r := float64(l.color.r) * Kd[0] * diffuse
	g := float64(l.color.g) * Kd[1] * diffuse
	b := float64(l.color.b) * Kd[2] * diffuse
	rnew := int(math.Max(r, 0))
	gnew := int(math.Max(g, 0))
	bnew := int(math.Max(b, 0))
	return Color{
		r: rnew,
		g: gnew,
		b: bnew,
	}
}

func calcSpecularLighting(p0, p1, p2, view, Ks []float64, l Light) Color {
	view = Normalize(view)
	n, _ := Normal(p0, p1, p2)
	n = Normalize(n)
	light := Normalize(l.location)
	dot, _ := DotProduct(n, light)
	dot *= 2
	for i := range n {
		n[i] = n[i] * dot
	}
	diff := make([]float64, 3)
	for i := range diff {
		diff[i] = n[i] - l.location[i]
	}
	spec, _ := DotProduct(diff, view)
	spec = math.Pow(spec, 22)
	r := int(math.Max(float64(l.color.r)*Ks[0]*spec, 0))
	g := int(math.Max(float64(l.color.g)*Ks[1]*spec, 0))
	b := int(math.Max(float64(l.color.b)*Ks[2]*spec, 0))
	c := Color{
		r: r,
		g: g,
		b: b,
	}
	return c
}
