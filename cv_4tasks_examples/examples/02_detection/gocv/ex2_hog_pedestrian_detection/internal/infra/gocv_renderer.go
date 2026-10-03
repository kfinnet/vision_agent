package infra

import (
	"ex2_hog_pedestrian_detection/internal/domain"
	"gocv.io/x/gocv"
)

// GoCVWindowRenderer는 gocv.Window를 사용해 프레임을 화면에 표시하는
// domain.WindowRenderer 구현체다. 박스를 그리는 책임은 GoCVFrameDrawer가
// 맡으므로, 이 구조체는 순수하게 "보여주기"만 담당한다(SRP).
type GoCVWindowRenderer struct {
	window *gocv.Window
}

// NewGoCVWindowRenderer는 title 제목의 창을 생성해 GoCVWindowRenderer를 반환한다.
func NewGoCVWindowRenderer(title string) *GoCVWindowRenderer {
	return &GoCVWindowRenderer{window: gocv.NewWindow(title)}
}

// Show는 frame을 창에 출력한다.
func (r *GoCVWindowRenderer) Show(frame domain.Frame) {
	mf, ok := frame.(MatFrame)
	if !ok {
		return
	}
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
