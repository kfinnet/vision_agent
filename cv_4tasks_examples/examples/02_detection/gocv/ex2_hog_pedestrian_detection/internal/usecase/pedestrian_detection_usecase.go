// Package usecase는 클린 아키텍처의 유스케이스 계층이다.
// 애플리케이션의 흐름(탐지 실행 → 시각화)을 조율하며, gocv를 비롯한 어떤
// 구체 구현에도 의존하지 않고 오직 domain 패키지의 추상 인터페이스(포트)만
// 사용한다(DIP). 생성자로 의존성을 주입받는다(Constructor Injection).
package usecase

import "ex2_hog_pedestrian_detection/internal/domain"

// PedestrianDetectionUseCase는 보행자 탐지 + 시각화 유스케이스다.
// detector와 drawer는 둘 다 추상 인터페이스 타입으로 주입받으므로, 구체
// 구현(HOG, 다른 그리기 스타일 등)을 교체해도 이 구조체는 수정되지
// 않는다(OCP, DIP).
type PedestrianDetectionUseCase struct {
	detector domain.ObjectDetector
	drawer   domain.FrameDrawer
}

// NewPedestrianDetectionUseCase는 ObjectDetector와 FrameDrawer를 주입받아
// PedestrianDetectionUseCase를 생성한다.
func NewPedestrianDetectionUseCase(detector domain.ObjectDetector, drawer domain.FrameDrawer) *PedestrianDetectionUseCase {
	return &PedestrianDetectionUseCase{detector: detector, drawer: drawer}
}

// Execute는 한 프레임에 대해 탐지를 수행하고, 결과를 그린 새 프레임과
// 탐지 목록을 함께 반환한다. gocv 관련 호출은 전혀 포함하지 않는다(DIP).
func (uc *PedestrianDetectionUseCase) Execute(frame domain.Frame) (domain.Frame, []domain.Detection, error) {
	detections, err := uc.detector.Detect(frame)
	if err != nil {
		return nil, nil, err
	}
	annotated := uc.drawer.Draw(frame, detections)
	return annotated, detections, nil
}
