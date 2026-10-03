// Package infra는 gocv를 직접 사용하는 유일한 계층이다. domain 패키지의
// 인터페이스(MotionDetector, FrameSource, ResultPresenter)를 gocv 기반으로
// 구현한다.
package infra

import (
	"fmt"

	"gocv.io/x/gocv"
)

// GoCVVideoSource는 gocv.VideoCapture로 비디오 파일에서 프레임을 순차적으로
// 읽어오는 domain.FrameSource 구현체다. 같은 패키지의 다른 어댑터
// (GoCVMotionDetector, GoCVWindowPresenter)가 Frame()을 통해 현재 프레임을
// 공유해서 쓴다.
type GoCVVideoSource struct {
	video *gocv.VideoCapture
	frame gocv.Mat
}

// NewGoCVVideoSource는 비디오 파일 경로로부터 GoCVVideoSource를 생성한다.
func NewGoCVVideoSource(videoPath string) (*GoCVVideoSource, error) {
	video, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return nil, fmt.Errorf("비디오를 열 수 없습니다: %w", err)
	}
	return &GoCVVideoSource{video: video, frame: gocv.NewMat()}, nil
}

// Next는 다음 프레임을 읽는다. 더 이상 프레임이 없으면 false를 반환한다.
// domain.FrameSource 인터페이스의 구현이다.
func (s *GoCVVideoSource) Next() bool {
	ok := s.video.Read(&s.frame)
	return ok && !s.frame.Empty()
}

// Frame은 가장 최근에 읽은 프레임을 반환한다.
func (s *GoCVVideoSource) Frame() gocv.Mat {
	return s.frame
}

// Close는 비디오 캡처와 내부 Mat 자원을 해제한다.
func (s *GoCVVideoSource) Close() error {
	s.frame.Close()
	return s.video.Close()
}
