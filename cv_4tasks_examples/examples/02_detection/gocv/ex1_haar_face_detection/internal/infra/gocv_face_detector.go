package infra

import (
	"fmt"

	"ex1_haar_face_detection/internal/domain"
	"gocv.io/x/gocv"
)

// HaarCascadeFaceDetector는 Haar Cascade 분류기 기반 얼굴 탐지 어댑터다.
// domain.ObjectDetector 포트를 구현하므로, 유스케이스 계층 입장에서는
// 다른 어떤 ObjectDetector 구현체와도 동일하게 취급된다(LSP). 이 구조체는
// "탐지"라는 책임만 가지며, 화면 출력이나 캡처 책임은 갖지 않는다(SRP).
type HaarCascadeFaceDetector struct {
	classifier   gocv.CascadeClassifier
	scaleFactor  float64
	minNeighbors int
}

// NewHaarCascadeFaceDetector는 cascadePath의 Haar Cascade XML 파일을 로드해
// HaarCascadeFaceDetector를 생성한다.
func NewHaarCascadeFaceDetector(cascadePath string) (*HaarCascadeFaceDetector, error) {
	classifier := gocv.NewCascadeClassifier()
	if !classifier.Load(cascadePath) {
		return nil, fmt.Errorf("cascade 파일을 불러오지 못했습니다: %s", cascadePath)
	}
	return &HaarCascadeFaceDetector{classifier: classifier}, nil
}

// Close는 분류기가 점유한 자원을 해제한다.
func (d *HaarCascadeFaceDetector) Close() error {
	return d.classifier.Close()
}

// Detect는 주어진 프레임에서 얼굴을 탐지해 domain.Detection 목록으로 반환한다.
// 현재 프레임에서 얼굴 탐지 (매 프레임 반복 호출 -> 실시간 처리).
func (d *HaarCascadeFaceDetector) Detect(frame domain.Frame) ([]domain.Detection, error) {
	mf, ok := frame.(MatFrame)
	if !ok {
		return nil, fmt.Errorf("HaarCascadeFaceDetector는 MatFrame만 지원합니다")
	}
	rects := d.classifier.DetectMultiScale(mf.Mat)
	detections := make([]domain.Detection, 0, len(rects))
	for _, r := range rects {
		detections = append(detections, domain.Detection{
			X: r.Min.X, Y: r.Min.Y, W: r.Dx(), H: r.Dy(),
			Label: "face", Confidence: 1.0,
		})
	}
	return detections, nil
}
