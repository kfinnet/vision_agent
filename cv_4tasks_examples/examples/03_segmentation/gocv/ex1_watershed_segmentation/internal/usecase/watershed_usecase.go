// Package usecase는 ===== [유스케이스 계층 / Use Case Layer] =====.
// 오직 domain 패키지의 인터페이스(포트)에만 의존하며(DIP), gocv를 직접
// import하거나 호출하지 않는다. "입력 이미지를 분할해서 색칠된 결과 파일로
// 저장한다"는 애플리케이션 흐름만 책임진다(SRP).
package usecase

import "ex1_watershed_segmentation/internal/domain"

// WatershedUseCase는 이미지 분할 + 결과 시각화 유스케이스를 오케스트레이션한다.
//
// 생성자 주입(Constructor Injection)으로 domain.ImageSegmenter와
// domain.ResultVisualizer 구현체를 받는다. 어떤 구체 알고리즘이 주입되는지는
// 전혀 알지 못하므로, Watershed를 다른 분할 알고리즘으로 교체해도 이 구조체는
// 한 줄도 바뀌지 않는다(OCP).
type WatershedUseCase struct {
	segmenter  domain.ImageSegmenter
	visualizer domain.ResultVisualizer
}

// NewWatershedUseCase는 분할기와 시각화기를 주입받아 WatershedUseCase를 생성한다.
func NewWatershedUseCase(segmenter domain.ImageSegmenter, visualizer domain.ResultVisualizer) *WatershedUseCase {
	return &WatershedUseCase{segmenter: segmenter, visualizer: visualizer}
}

// Execute는 inputPath의 이미지를 분할하고, 그 결과를 outputPath에 색칠하여
// 저장한다. 분할된 물체 개수를 함께 반환한다.
func (uc *WatershedUseCase) Execute(inputPath, outputPath string) (int, error) {
	result, err := uc.segmenter.Segment(inputPath)
	if err != nil {
		return 0, err
	}
	if err := uc.visualizer.Colorize(result, outputPath); err != nil {
		return 0, err
	}
	return result.NumSegments - 1, nil
}
