package domain

// ImageSegmenter는 '이미지 파일 경로를 받아 전경/배경으로 분할한다'는 핵심
// 행위만 정의하는 작은 인터페이스(ISP) — 유스케이스 계층이 의존하는 "포트"다.
//
// 시각화(윤곽선, 배경 제거 등)는 포함하지 않는다(SRP). GrabCut 대신 다른
// 전경 분할 알고리즘으로 교체하더라도 이 인터페이스만 구현하면 되므로
// 유스케이스 계층은 전혀 수정할 필요가 없다(OCP/LSP).
type ImageSegmenter interface {
	// Segment는 입력 이미지(경로)를 전경/배경으로 분할하여 SegmentationResult를 반환한다.
	Segment(inputPath string) (SegmentationResult, error)
}

// ResultPresenter는 SegmentationResult를 적용한 파생 이미지(배경 제거 등)를
// 만들어 파일로 저장하는 표현 책임만 정의하는 작은 인터페이스(ISP) — 분할
// 계산 책임(ImageSegmenter)과는 분리된 별도의 포트다(SRP).
type ResultPresenter interface {
	// RemoveBackground는 원본 이미지에서 배경을 제거(검은색 처리)하여 outputPath에 저장한다.
	RemoveBackground(inputPath string, result SegmentationResult, outputPath string) error
}
