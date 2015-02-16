package packing

import (
	"image"
	"sort"
	"math"
	"errors"
)

type node struct {
	ID string
	Left *node //left node
	Right *node //right node
	Rect image.Rectangle //placed rect
	Used bool //it is ture when node is used
	Path string //path of node from root
}

func createNode(rect image.Rectangle) *node {
	return &node{"", nil, nil, rect, false, "root"}
}

func (n *node) IsUsed() bool {
	if (n.Used) {
		return true
	}
	return false
}

func (n *node) IsPlaceable(rect image.Rectangle) bool {
	if (n != nil && rect.Dx() <= n.Rect.Dx() && rect.Dy() <= n.Rect.Dy()) {
		return true
	}
	
	return false
}

func (n *node) IsInsertable(rect image.Rectangle) bool {
	if (n != nil && !n.IsUsed() && n.IsPlaceable(rect) &&
		n.Left == nil && n.Right == nil) {
		
		return true
	}
	
	return false
}

func (n *node) GetInsertableNode(rect image.Rectangle) *node {
	if (n == nil) {
		return nil
	}
	
	if (!n.IsPlaceable(rect)) {
		return nil
	} else if (n.IsInsertable(rect)) {
		return n
	}
	
	result := n.Left.GetInsertableNode(rect)
	if (result == nil) {
		result = n.Right.GetInsertableNode(rect)
	}
	return result
}

func (n *node) Insert(id string, rect image.Rectangle) bool {
	targetNode := n.GetInsertableNode(rect)
	if (targetNode.IsInsertable(rect)) {
		// vertical divid
		l := createNode(image.Rect(targetNode.Rect.Min.X, targetNode.Rect.Min.Y,
			targetNode.Rect.Min.X + rect.Dx(), targetNode.Rect.Max.Y))
		r := createNode(image.Rect(targetNode.Rect.Min.X + rect.Dx(), targetNode.Rect.Min.Y,
			targetNode.Rect.Max.X, targetNode.Rect.Max.Y))
		l.Path = targetNode.Path + "l"
		targetNode.Left = l
		r.Path = targetNode.Path + "r"
		targetNode.Right = r
		
		// horizonal divid
		ll := createNode(image.Rect(l.Rect.Min.X, l.Rect.Min.Y,
			l.Rect.Min.X + rect.Dx(), l.Rect.Min.Y + rect.Dy()))
		lr := createNode(image.Rect(l.Rect.Min.X, l.Rect.Min.Y + rect.Dy(),
			l.Rect.Min.X + rect.Dx(), l.Rect.Max.Y))
		ll.Path = l.Path + "l"
		l.Left = ll
		lr.Path = l.Path + "r"
		l.Right = lr
		
		// insert
		ll.ID = id
		ll.Used = true
		//log.Println("Node is inserted. id = " + id, " path = " + ll.Path)
		
		return true
	}
	
	return false
}

func (n *node) GetUsedNodes() []node {
	var result []node
	if (n != nil) {
		if (n.IsUsed()) {
			result = append(result, *n)
		} else {
			leftResult := n.Left.GetUsedNodes()
			result = append(result, leftResult...)
			rightResult := n.Right.GetUsedNodes()
			result = append(result, rightResult...)
		}
	}
	
	return result
}

func (n *node) GetUnUsedNodes() []node {
	var result []node
	if (n != nil) {
		if (!n.IsUsed() && n.Left == nil && n.Right == nil) {
			result = append(result, *n)
		} else {
			leftResult := n.Left.GetUsedNodes()
			result = append(result, leftResult...)
			rightResult := n.Right.GetUsedNodes()
			result = append(result, rightResult...)
		}
	}
	
	return result
}

// sortableImages array for sorting.
type sortableImages []packingImage
func (simages sortableImages) Len() int           { return len(simages) }
func (simages sortableImages) Less(i, j int) bool { return (simages[i].Image.Bounds().Dy()) > (simages[j].Image.Bounds().Dy()) }
func (simages sortableImages) Swap(i, j int)      { simages[i], simages[j] = simages[j], simages[i] }

// packingImage object for packing.
type packingImage struct {
	ID string
	Image image.Image
}

func createPackingImage(id string, img image.Image) *packingImage {
	return &packingImage{id, img}
}

// Info information for packing.
type Info struct {
	MaxX, MaxY int
	index map[string]*image.Image
	simages sortableImages
}

// CreatePackingInfo create Info.
// param maxX - Width of rectangle to pack image
// param maxY - Height of rectangle to pack image
// return arg1 - Info
func CreatePackingInfo(maxX int, maxY int) *Info {
	var images []packingImage
	return &Info{maxX, maxY, make(map[string]*image.Image), images}
}

// GetImage get image.
// param id - ID of image
// return arg1 - image
func (pinfo *Info) GetImage(id string) image.Image {
	return *pinfo.index[id]
}

// AddImage add image.
// param id - ID of image
// param img - image
func (pinfo *Info) AddImage(id string, img image.Image) {
	packingImage := createPackingImage(id, img)
	pinfo.simages = append(pinfo.simages, *packingImage)
	pinfo.index[id] = &img
}

// AddImageAt add image at index.
// param index - index
// param pimg - packingImage
func (pinfo *Info) AddImageAt(index int, pimg packingImage) {
	pinfo.simages[index] = pimg
	pinfo.index[pimg.ID] = &pimg.Image
}

// Result result of packing.
type Result struct {
	BaseRect image.Rectangle
	Rects map[string]image.Rectangle //key:ID, value:image.Rectangle
	rotate map[string]bool //key:ID, value:isRotated
}

// CreatePackingResult create Result.
// return arg1 - Result
func CreatePackingResult() *Result {
	return &Result{image.Rect(0,0,0,0), make(map[string]image.Rectangle), make(map[string]bool)}
}

// GetRect get rectangle of image.
// param id - ID of image
// return arg1 - rectangle of image
func (packingResult *Result) GetRect(id string) image.Rectangle {
	return packingResult.Rects[id]
}

// IsRotated get whether an image rotated.
// param  id - ID of image
// return arg1 - If image is rotated, return true.
func (packingResult *Result) IsRotated(id string) bool {
	return packingResult.rotate[id]
}

// SetRotated set whether an image rotated.
// param id - ID of image
// param rotated - whether an image rotated
func (packingResult *Result) SetRotated(id string, rotated bool) {
	packingResult.rotate[id] = rotated
}

// Pack pack an image within a rectangle.
// param info - images and information for packing
// return arg1 - result of packing
// return arg2 - Error
func Pack(info Info) ([]*Result, error) {
	//println(len(info.simages))
	if len(info.simages) <= 0 {
		return nil, errors.New("number of images is 0")
	}

	//rotate
	rotate := make(map[string]bool)
	for index, pimg := range info.simages {
		img := pimg.Image
		if (img.Bounds().Dx() > img.Bounds().Dy()) {
			pimg.Image = Rotate90(pimg.Image)
			info.AddImageAt(index, pimg)
			rotate[pimg.ID] = true
//debug
//			log.Println("id = " + pimg.ID + " x = " + strconv.FormatUint(uint64(pimg.Image.Bounds().Dx()), 10) + " y = " + strconv.FormatUint(uint64(pimg.Image.Bounds().Dy()), 10))
		}
	}
	
//debug
//	for _,pimg := range info.simages {
//		log.Println("id = " + pimg.ID + " x = " + strconv.FormatUint(uint64(pimg.Image.Bounds().Dx()), 10) + " y = " + strconv.FormatUint(uint64(pimg.Image.Bounds().Dy()), 10))
//	}
	
	sort.Sort(info.simages)
	
	// Check largest Image is placeable to BaseNode
	n := createNode(image.Rect(0,0,info.MaxX,info.MaxY))
	placeable := n.IsPlaceable(info.simages[0].Image.Bounds())
	if !placeable {
		return nil, errors.New("not placeable image is included")
	}
	
	
	// Pack
	var nodes []*node
	nodes = append(nodes, n)
	for _, pimg := range info.simages {
		img := pimg.Image
		inserted := false
		for _, node := range nodes {
			inserted = node.Insert(pimg.ID, img.Bounds())
			if inserted {
				//log.Println(pimg.ID + " is inserted.")
				break
			}
		}
		
		// next node
		if !inserted {
			node := createNode(image.Rect(0,0,info.MaxX,info.MaxY))
			inserted = node.Insert(pimg.ID, img.Bounds())
			if !inserted {
				return nil, errors.New(pimg.ID + "is not insertable.")
			}
			nodes = append(nodes, node)
			//log.Println(pimg.ID + " is inserted to next node.")
		}
	}
	
	var results []*Result
	for _, node := range nodes {
		result := CreatePackingResult()
		maxX := 0
		maxY := 0
		for _, v := range node.GetUsedNodes() {
			if (maxX < v.Rect.Max.X) {
				maxX = v.Rect.Max.X
			}
			if (maxY < v.Rect.Max.Y) {
				maxY = v.Rect.Max.Y
			}
			result.Rects[v.ID] = v.Rect
			result.SetRotated(v.ID, rotate[v.ID])
		}
		result.BaseRect.Max.X = maxX
		result.BaseRect.Max.Y = maxY
		results = append(results, result)
	}
	
	return results, nil
}

// PadToPow2 calculate maximum of power of 2 within x.
// param x - 上限値
// return arg1 : maximum of power of 2 within x
func PadToPow2(x int) int {
	result := 1
	for i := x; i > 0; i /= 2 {
		result = result * 2
	}
	if x * 2 == result {
		return x
	}
	return result
}

// Rotate90 rotates 90 degrees clockwise and returns an image. 
// param src - source image
// return arg1 - rotated image
func Rotate90(src image.Image) image.Image {
	radian := math.Pi * 2 * float64(90) / 360
	//base img
	x := src.Bounds().Max.X
	y := src.Bounds().Max.Y
	dst := image.NewRGBA(image.Rect(0, 0, y, x))
	for i := dst.Rect.Min.Y; i < dst.Rect.Max.Y; i++ {
		for j := dst.Rect.Min.X; j < dst.Rect.Max.X; j++ {
        	srcX := int(float64(j) * math.Cos(-radian) - float64(i) * math.Sin(-radian)) //rot
        	srcY := int(float64(j) * math.Sin(-radian) + float64(i) * math.Cos(-radian)) + y - 1 //rot and trans y
			dst.Set(j, i, src.At(srcX, srcY))
		}
	}
	return dst
}