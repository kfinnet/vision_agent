package infra

import (
	"image"
	"image/color"

	"ex2_hog_pedestrian_detection/internal/domain"
	"gocv.io/x/gocv"
)

// GoCVFrameDrawer는 탐지 박스와 라벨을 프레임 위에 그리는 domain.FrameDrawer
// 구현체다. 탐지 로직과 분리된 책임(SRP)이므로, 박스 색상이나 라벨 표기
// 형식을 바꾸더라도 탐지기나 유스케이스에는 영향이 없다(OCP).
type GoCVFrameDrawer struct {
	boxColor color.RGBA
}

// NewGoCVFrameDrawer는 기본 박스 색상을 갖는 GoCVFrameDrawer를 생성한다.
func NewGoCVFrameDrawer() *GoCVFrameDrawer {
	return &GoCVFrameDrawer{boxColor: color.RGBA{R: 255, G: 140, B: 0, A: 0}}
}

// Draw는 frame을 복제한 뒤 그 위에 detections를 그려 새 Frame(MatFrame)으로 반환한다.
func (d *GoCVFrameDrawer) Draw(frame domain.Frame, detections []domain.Detection) domain.Frame {
	mf, ok := frame.(MatFrame)
	if !ok {
		return frame
	}
	result := mf.Mat.Clone()
	for _, det := range detections {
		rect := image.Rect(det.X, det.Y, det.X+det.W, det.Y+det.H)
		gocv.Rectangle(&result, rect, d.boxColor, 2)
		gocv.PutText(&result, det.Label, image.Pt(det.X, det.Y-8),
			gocv.FontHersheySimplex, 0.6, d.boxColor, 2)
	}
	return MatFrame{Mat: result}
}
