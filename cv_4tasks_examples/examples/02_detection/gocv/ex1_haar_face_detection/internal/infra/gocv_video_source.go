package infra

import (
	"fmt"

	"ex1_haar_face_detection/internal/domain"
	"gocv.io/x/gocv"
)

// GoCVVideoSource는 gocv.VideoCapture를 사용해 웹캠 또는 비디오 파일에서
// 프레임을 읽어오는 domain.VideoSource 구현체다. 비디오 소스가 웹캠인지
// 파일인지는 생성자에서 결정되며, 이후 Read() 호출부는 이를 알 필요가
// 없다(SRP: "프레임 읽기"라는 책임만 가진다).
type GoCVVideoSource struct {
	capture *gocv.VideoCapture
	mat     gocv.Mat
}

// NewWebcamVideoSource는 지정한 장치 번호의 웹캠을 여는 GoCVVideoSource를 생성한다.
func NewWebcamVideoSource(deviceID int) (*GoCVVideoSource, error) {
	cap, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		return nil, fmt.Errorf("웹캠(%d)을 열 수 없습니다: %w", deviceID, err)
	}
	return &GoCVVideoSource{capture: cap, mat: gocv.NewMat()}, nil
}

// NewFileVideoSource는 지정한 경로의 비디오 파일을 여는 GoCVVideoSource를 생성한다.
func NewFileVideoSource(path string) (*GoCVVideoSource, error) {
	cap, err := gocv.VideoCaptureFile(path)
	if err != nil {
		return nil, fmt.Errorf("비디오 파일(%s)을 열 수 없습니다: %w", path, err)
	}
	return &GoCVVideoSource{capture: cap, mat: gocv.NewMat()}, nil
}

// Read는 다음 프레임을 읽어 domain.Frame(MatFrame)으로 감싸 반환한다.
func (s *GoCVVideoSource) Read() (frame domain.Frame, ok bool) {
	if !s.capture.Read(&s.mat) || s.mat.Empty() {
		return MatFrame{Mat: s.mat}, false
	}
	return MatFrame{Mat: s.mat}, true
}

// Close는 비디오 캡처와 내부 Mat 자원을 해제한다.
func (s *GoCVVideoSource) Close() error {
	s.mat.Close()
	return s.capture.Close()
}
