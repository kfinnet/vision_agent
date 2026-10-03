package infra

import (
	"image/color"

	"gocv.io/x/gocv"

	"ex1_tracker_mil/internal/domain"
)

// GoCVWindowSink는 추적 결과(박스 + 이동 경로)를 gocv 윈도우에 그려 보여주는
// domain.ResultSink 구현체다.
//
// SRP/ISP: 추적 알고리즘(GoCVTrackerMIL)과 완전히 분리된, 오직 "그리기/표시"
// 책임만 지는 작은 어댑터다. 추적 로직을 전혀 모른 채 TrackedObject와 Trail
// 값만 받아 그린다.
type GoCVWindowSink struct {
	window *gocv.Window
	source *GoCVVideoSource // 그릴 대상 프레임(Mat)을 공유받는 지점
}

// NewGoCVWindowSink는 지정한 제목의 윈도우를 띄우는 GoCVWindowSink를 생성한다.
func NewGoCVWindowSink(title string, source *GoCVVideoSource) *GoCVWindowSink {
	return &GoCVWindowSink{window: gocv.NewWindow(title), source: source}
}

// Present는 현재 프레임 위에 추적 박스와 지금까지의 이동 경로를 그려 화면에
// 표시한다. ESC 키 입력 시 false를 반환해 상위 유스케이스에 중단을 알린다.
// domain.ResultSink 인터페이스의 구현이다.
func (s *GoCVWindowSink) Present(result domain.TrackedObject, trail domain.Trail, frameIdx int) bool {
	frame := s.source.Frame()

	gocv.Rectangle(&frame, result.Box, color.RGBA{R: 122, G: 92, B: 158, A: 0}, 2)
	for i := 1; i < len(trail.Points); i++ {
		gocv.Line(&frame, trail.Points[i-1], trail.Points[i], color.RGBA{R: 0, G: 140, B: 255, A: 0}, 2)
	}

	s.window.IMShow(frame)
	return s.window.WaitKey(30) != 27
}

// Close는 윈도우 자원을 해제한다.
func (s *GoCVWindowSink) Close() error {
	return s.window.Close()
}
