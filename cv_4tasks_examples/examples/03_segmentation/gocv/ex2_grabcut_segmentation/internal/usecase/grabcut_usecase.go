// Package usecase는 ===== [유스케이스 계층 / Use Case Layer] =====.
// 오직 domain 패키지의 인터페이스(포트)에만 의존하며(DIP), gocv를 직접
// import하거나 호출하지 않는다. "입력 이미지를 분할해서 배경을 제거한 결과
// 파일로 저장한다"는 애플리케이션 흐름만 책임진다(SRP).
package usecase

import "ex2_grabcut_segmentation/internal/domain"

// GrabCutUseCase는 전경/배경 분할 + 배경 제거 유스케이스를 오케스트레이션한다.
//
// 생성자 주입으로 domain.ImageSegmenter와 domain.ResultPresenter 구현체를
// 받는다. 어떤 구체 알고리즘이 주입되는지는 전혀 알지 못하므로, GrabCut을
// 다른 분할 알고리즘으로 교체해도 이 구조체는 한 줄도 바뀌지 않는다(OCP).
type GrabCutUseCase struct {
	segmenter domain.ImageSegmenter
	presenter domain.ResultPresenter
}

// NewGrabCutUseCase는 분할기와 표현기를 주입받아 GrabCutUseCase를 생성한다.
func NewGrabCutUseCase(segmenter domain.ImageSegmenter, presenter domain.ResultPresenter) *GrabCutUseCase {
	return &GrabCutUseCase{segmenter: segmenter, presenter: presenter}
}

// Execute는 inputPath의 이미지를 전경/배경으로 분할하고, 배경을 제거한
// 결과를 outputPath에 저장한다. 전경 픽셀 비율(0.0~1.0)을 함께 반환한다.
func (uc *GrabCutUseCase) Execute(inputPath, outputPath string) (float64, error) {
	result, err := uc.segmenter.Segment(inputPath)
	if err != nil {
		return 0, err
	}
	if err := uc.presenter.RemoveBackground(inputPath, result, outputPath); err != nil {
		return 0, err
	}
	return result.ForegroundRatio, nil
}
