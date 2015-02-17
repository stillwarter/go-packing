package packing

import (
	"testing"
	"math/rand"
	"io/ioutil"
	"strconv"
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"encoding/json"
)

func createImage(width int, height int, color color.RGBA) *image.RGBA {
	r := image.Rect(0, 0, width, height);
	eleImg := image.NewRGBA(r);
	for y := eleImg.Rect.Min.Y; y < eleImg.Rect.Max.Y; y++ {
		for x := eleImg.Rect.Min.X; x < eleImg.Rect.Max.X; x++ {
			eleImg.Set(x, y, color)
		}
	}

	return eleImg;
}

func createMeta(presults []*Result, names []string) ([]byte, error) {
	type frame struct {
		FileID string `json:"fileId"`
		X      int    `json:"x"`
		Y      int    `json:"y"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Rotate bool   `json:"rotate"`
	}
	obj := make(map[string][]*frame)
	for i, presult := range presults {
		var frames []*frame
		for fileID, rect := range presult.Rects {
			f := &frame{}
			f.FileID = fileID
			f.X = rect.Min.X
			f.Y = rect.Min.Y
			f.Height = rect.Max.Y - rect.Min.Y
			f.Width = rect.Max.X - rect.Min.X
			f.Rotate = presult.IsRotated(fileID)
			frames = append(frames, f)
		}
		obj[names[i]] = frames
	}
	jsonb, err := json.MarshalIndent(&obj, "", " ")
	if err != nil {
		return nil, err
	}
	return jsonb, nil
}

func TestOK_Pack(t *testing.T) {
	// create packing infomation
	pinfo := CreatePackingInfo(512, 512)
	for i := 0; i < 10; i++ {
		fileID := "img_" + strconv.FormatUint(uint64(i), 10) + ".png"
		randR := uint8(rand.Intn(255));
		randG := uint8(rand.Intn(255));
		randB := uint8(rand.Intn(255));
		pinfo.AddImage(fileID, createImage(rand.Intn(255), rand.Intn(255), color.RGBA{randR,randG,randB,0x88}))
	}
	
	// pack
	presults, err := Pack(*pinfo)
	if err != nil {
		t.Error(err)
	}
	
	// draw images into txatlas
	imgNum := len(presults)
	pimgs := make([][]byte, imgNum, imgNum)
	for i := 0; i < imgNum; i++ {
		presult := presults[i]
		img := image.NewRGBA(presult.BaseRect)
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				img.Set(x, y, color.RGBA{0x00, 0x00, 0x00, 0x00})
			}
		}
		for id, rect := range presult.Rects {
			eleImg := pinfo.GetImage(id)
			draw.Draw(img, rect, eleImg, image.Pt(0, 0), draw.Over)
		}

		buff := &bytes.Buffer{} //empty buffer
		png.Encode(buff, img)

		pimgs[i] = buff.Bytes()
	}
	
	// output image file
	var names []string
	for i, v := range pimgs {
		var name = "txatlas_" + strconv.FormatUint(uint64(i), 10) + ".png"
		err = ioutil.WriteFile("output/" + name, v, 0777)
		if err != nil {
			t.Error(err)
		}
		names = append(names, name)
	}
	
	// output meta file
	jsonb, err := createMeta(presults, names)
	err = ioutil.WriteFile("output/meta.json", jsonb, 0777)
	if err != nil {
		t.Error(err)
	}
	
	
}

func TestOK_PadToPow2_Up(t *testing.T) {
	// test data
	x := 999
	expected := 1024
	
	result := PadToPow2(x)
	
	if result != expected {
		t.Errorf("result(%d) is invalid. expected = %d", result, expected)
	}
}

func TestOK_PadToPow2_Same(t *testing.T) {
	// test data
	x := 512
	expected := 512
	
	result := PadToPow2(x)
	
	if result != expected {
		t.Errorf("result(%d) is invalid. expected = %d", result, expected)
	}
}