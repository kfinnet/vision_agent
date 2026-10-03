package infra

import (
	"fmt"
	"image"
	"image/color"

	"ex1_haar_face_detection/internal/domain"
	"gocv.io/x/gocv"
)

// GoCVWindowRenderer는 gocv.Window를 사용해 탐지 결과를 화면에 표시하는
// domain.FrameRenderer 구현체다. 박스/라벨을 그리는 시각화 책임과 탐지
// 책임을 분리해(SRP) 탐지 알고리즘 교체가 렌더링에 영향을 주지 않게 한다.
type GoCVWindowRenderer struct {
	window   *gocv.Window
	boxColor color.RGBA
}

// NewGoCVWindowRenderer는 title 제목의 창을 생성해 GoCVWindowRenderer를 반환한다.
func NewGoCVWindowRenderer(title string) *GoCVWindowRenderer {
	return &GoCVWindowRenderer{
		window:   gocv.NewWindow(title),
		boxColor: color.RGBA{R: 255, G: 140, B: 0, A: 0},
	}
}

// Render는 frame 위에 탐지 박스와 상태 텍스트를 그려 창에 출력한다.
func (r *GoCVWindowRenderer) Render(frame domain.Frame, detections []domain.Detection, frameCount int) {
	mf, ok := frame.(MatFrame)
	if !ok {
		return
	}
	for _, d := range detections {
		rect := image.Rect(d.X, d.Y, d.X+d.W, d.Y+d.H)
		gocv.Rectangle(&mf.Mat, rect, r.boxColor, 3)
		label := fmt.Sprintf("%s (%d)", d.Label, frameCount)
		gocv.PutText(&mf.Mat, label, image.Pt(d.X, d.Y-10),
			gocv.FontHersheySimplex, 0.8, r.boxColor, 2)
	}
	statusText := fmt.Sprintf("frame=%d faces=%d", frameCount, len(detections))
	gocv.PutText(&mf.Mat, statusText, image.Pt(10, 25),
		gocv.FontHersheySimplex, 0.7, color.RGBA{G: 255, A: 0}, 2)

	r.window.IMShow(mf.Mat)
}

// WaitKey는 지정한 시간(ms) 동안 키 입력을 기다린다.
func (r *GoCVWindowRenderer) WaitKey(delayMs int) int {
	return r.window.WaitKey(delayMs)
}

// Close는 창 자원을 해제한다.
func (r *GoCVWindowRenderer) Close() error {
	return r.window.Close()
}
