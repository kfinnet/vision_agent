package infra

import (
	"fmt"

	"ex2_hog_pedestrian_detection/internal/domain"
	"gocv.io/x/gocv"
)

// HOGPedestrianDetector는 HOG + 선형 SVM 기반 보행자 탐지 어댑터다.
// 얼굴 탐지(Haar Cascade)와 달리 "사람의 전신 실루엣" 패턴을 학습한
// OpenCV 내장 사전학습 모델(DefaultPeopleDetector)을 사용한다.
// domain.ObjectDetector 포트를 구현하므로, 유스케이스 계층 입장에서는
// 다른 어떤 ObjectDetector 구현체와도 동일하게 취급된다(LSP).
type HOGPedestrianDetector struct {
	hog gocv.HOGDescriptor
}

// NewHOGPedestrianDetector는 OpenCV 내장 사전학습 보행자 검출 SVM을 로드해
// HOGPedestrianDetector를 생성한다.
func NewHOGPedestrianDetector() *HOGPedestrianDetector {
	hog := gocv.NewHOGDescriptor()
	hog.SetSVMDetector(gocv.HOGDefaultPeopleDetector())
	return &HOGPedestrianDetector{hog: hog}
}

// Close는 HOG 디스크립터가 점유한 자원을 해제한다.
func (d *HOGPedestrianDetector) Close() error {
	return d.hog.Close()
}

// Detect는 주어진 프레임에서 보행자를 탐지해 domain.Detection 목록으로 반환한다.
func (d *HOGPedestrianDetector) Detect(frame domain.Frame) ([]domain.Detection, error) {
	mf, ok := frame.(MatFrame)
	if !ok {
		return nil, fmt.Errorf("HOGPedestrianDetector는 MatFrame만 지원합니다")
	}
	rects := d.hog.DetectMultiScale(mf.Mat)
	detections := make([]domain.Detection, 0, len(rects))
	for _, r := range rects {
		detections = append(detections, domain.Detection{
			X: r.Min.X, Y: r.Min.Y, W: r.Dx(), H: r.Dy(), Label: "person",
		})
	}
	return detections, nil
}
