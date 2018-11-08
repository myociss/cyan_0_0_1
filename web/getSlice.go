package web

import (
	"../algorithm"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"../plane3d"
)

/* contains code to return all information the web browswer uses to 
 * draw an svg from the point groups generated by the algorithm.
 */

 type SliceRes struct {
	SvgShapes 	[]SvgShape
	MinX 		float64
	MinY 		float64
	XRange		float64
	YRange		float64
}

type SvgShape struct {
	Points 		[][2]float64
	TetId		string
	CentroidX	string
	CentroidY	string
	TissueId	string
	Weight		float64
}

type ZRequestRes struct {
	MinZ float64
	MaxZ float64
}

func GetZRangeRequestHandler(w http.ResponseWriter, r *http.Request){
	//log.Println("here")
	minZ, maxZ := algorithm.GetZRange()
	a, err := json.Marshal(&ZRequestRes{minZ, maxZ})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
    w.Write(a)
}

/* this file handles the initial request for all 2d slices of mesh. these are
 * created when the server starts and are stored in an array. this file
 * returns the entire array of 2d slices as a json object.
 */

func SelectionSliceRequestHandler(w http.ResponseWriter, r *http.Request){
	//log.Println(r)

	zValRequest, _ := r.URL.Query()["zVal"]
	log.Println(zValRequest)
	
	zVal, err := strconv.ParseFloat(strings.Join(zValRequest, ","), 64)
	
	if err != nil {
		log.Fatal(err)
	}

	//res := algorithm.GetSelectionSlice(zVal)
	//log.Println(len(slice))
	slice := algorithm.GetSelectionSlice(zVal)
	//log.Println(slice)
	a, err := json.Marshal(newSliceRes(slice))
    if err != nil {
		log.Println(err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
    w.Write(a)
}

func newSliceRes(groups []*plane3d.VertexGroup)(*SliceRes){
	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	var svgShapes []SvgShape 

	for _, group := range groups {
		for _, v := range group.FoundVertices {
			if v.Vec[0] < minX {
				minX = v.Vec[0]
			}

			if v.Vec[0] > maxX {
				maxX = v.Vec[0]
			}

			if v.Vec[1] < minY {
				minY = v.Vec[1]
			}

			if v.Vec[1] > maxY {
				maxY = v.Vec[1]
			}
		}
		svgShapes = append(svgShapes, newSvgShape(group))
	}

	xRange, yRange := maxX - minX, maxY - minY
	return &SliceRes{svgShapes, minX, minY, xRange, yRange}
}

//need to implement calculation of shape center
//which entails figuring out why the points are out of order

//reverses x and y for browser and orders points correctly
func newSvgShape(group *plane3d.VertexGroup)(SvgShape){
	points := make([][2]float64, len(group.FoundVertices))
	for idx, point := range group.FoundVertices {
		points[idx] = [2]float64{point.Vec[0], point.Vec[1]}
	}

	tetId := strconv.Itoa(group.TetId)
	tissueId := strconv.Itoa(group.TissueId)
	centerY, centerX := calcCentroid(points)
	return SvgShape{points, tetId,centerX, centerY, tissueId, group.Weight}
}

func calcCentroid(group [][2]float64)(string, string){
	t1Y, t1X := 0.0, 0.0
	for i := 0; i < 3; i++ {
		t1Y += group[i][1]
		t1X += group[i][0]
	}

	t1Y, t1X = t1Y / 3, t1X / 3
	
	if len(group) == 3 {
		y := strconv.FormatFloat(t1Y, 'f', 2, 64)
		x := strconv.FormatFloat(t1X, 'f', 2, 64)
		return y, x
	} else {
		t2Y, t2X := 0.0, 0.0

		for i := 2; i <= 4; i++ {
			t2Y += group[i % 4][1]
			t2X += group[i % 4][0]
		}

		t2Y, t2X = t2Y / 3, t2X / 3

		y, x := (t1Y + t2Y) / 2, (t1X + t2X) / 2

		yStr := strconv.FormatFloat(y, 'f', 2, 64)
		xStr := strconv.FormatFloat(x, 'f', 2, 64)

		return yStr, xStr
	}
}