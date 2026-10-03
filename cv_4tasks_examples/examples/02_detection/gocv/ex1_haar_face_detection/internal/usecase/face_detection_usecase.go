// Package usecase는 클린 아키텍처의 유스케이스 계층이다.
// 애플리케이션의 흐름(탐지 실행)을 조율하며, gocv를 비롯한 어떤 구체
// 구현에도 의존하지 않고 오직 domain 패키지의 추상 인터페이스(포트)만
// 사용한다(DIP). 생성자로 의존성을 주입받는 방식(Constructor Injection)을
// 사용한다.
package usecase

import "ex1_haar_face_detection/internal/domain"

// FaceDetectionUseCase는 얼굴 탐지 유스케이스다.
// 생성자 주입으로 domain.ObjectDetector 추상화를 전달받아 Execute()를
// 호출할 때마다 탐지를 수행한다. Haar Cascade를 다른 탐지기(HOG, DNN 등)로
// 교체하더라도 이 구조체는 전혀 수정할 필요가 없다(OCP).
type FaceDetectionUseCase struct {
	detector domain.ObjectDetector
}

// NewFaceDetectionUseCase는 주어진 ObjectDetector를 주입받아
// FaceDetectionUseCase를 생성한다.
func NewFaceDetectionUseCase(detector domain.ObjectDetector) *FaceDetectionUseCase {
	return &FaceDetectionUseCase{detector: detector}
}

// Execute는 한 프레임에 대해 탐지를 실행하고 결과를 반환한다.
// 주입된 domain.ObjectDetector 외에는 아무것도 알지 못하므로(DIP),
// gocv 관련 호출은 전혀 포함하지 않는다.
func (uc *FaceDetectionUseCase) Execute(frame domain.Frame) ([]domain.Detection, error) {
	return uc.detector.Detect(frame)
}
