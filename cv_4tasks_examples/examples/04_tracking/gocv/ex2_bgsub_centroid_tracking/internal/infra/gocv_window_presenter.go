package infra

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"

	"ex2_bgsub_centroid_tracking/internal/domain"
)

// GoCVWindowPresenter는 탐지 박스 + 유지 중인 추적 ID를 gocv 윈도우에 그려
// 보여주는 domain.ResultPresenter 구현체다.
//
// SRP/ISP: 탐지(GoCVMotionDetector)나 매칭(domain.CentroidTracker)과 완전히
// 분리된, 오직 "그리기/표시" 책임만 지는 어댑터다.
type GoCVWindowPresenter struct {
	window *gocv.Window
	source *GoCVVideoSource
}

// NewGoCVWindowPresenter는 지정한 제목의 윈도우를 띄우는 GoCVWindowPresenter를
// 생성한다.
func NewGoCVWindowPresenter(title string, source *GoCVVideoSource) *GoCVWindowPresenter {
	return &GoCVWindowPresenter{window: gocv.NewWindow(title), source: source}
}

// Present는 현재 프레임을 복제해 그 위에 탐지 박스와, 이번 프레임에 매칭된
// 추적 ID 라벨을 그려 화면에 표시한다. ESC 키 입력 시 false를 반환한다.
// domain.ResultPresenter 인터페이스의 구현이다.
//
// 3) 시각화 — 박스 + 유지되고 있는 ID
func (p *GoCVWindowPresenter) Present(detections []domain.Detection, tracks []domain.Track, frameIdx int) bool {
	frame := p.source.Frame()
	vis := frame.Clone()
	defer vis.Close()

	for _, d := range detections {
		gocv.Rectangle(&vis, d.Box, color.RGBA{R: 0, G: 140, B: 255, A: 0}, 2)
	}
	for _, t := range tracks {
		if t.Active {
			label := fmt.Sprintf("ID:%02d", t.ID)
			gocv.PutText(&vis, label, image.Pt(t.Center.X-20, t.Center.Y-10),
				gocv.FontHersheySimplex, 0.6, color.RGBA{R: 122, G: 92, B: 158, A: 0}, 2)
		}
	}

	p.window.IMShow(vis)
	return p.window.WaitKey(30) != 27
}

// Close는 윈도우 자원을 해제한다.
func (p *GoCVWindowPresenter) Close() error {
	return p.window.Close()
}
