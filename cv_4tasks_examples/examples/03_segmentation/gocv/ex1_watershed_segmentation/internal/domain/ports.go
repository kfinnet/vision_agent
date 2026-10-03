package domain

// ImageSegmenter는 '이미지 파일 경로를 받아 픽셀 단위로 분할한다'는 핵심 행위만
// 정의하는 작은 인터페이스(ISP) — 유스케이스 계층이 의존하는 "포트"다.
//
// 시각화(색칠)나 파일 저장 같은 책임은 포함하지 않는다(SRP). 어떤 구체 알고리즘
// (Watershed, 향후 다른 분할 알고리즘 등)을 쓰든 이 인터페이스만 구현하면
// 유스케이스 계층의 수정 없이 교체 가능하다(OCP/LSP).
type ImageSegmenter interface {
	// Segment는 입력 이미지(경로)를 분할하여 SegmentationResult를 반환한다.
	Segment(inputPath string) (SegmentationResult, error)
}

// ResultVisualizer는 SegmentationResult를 사람이 보기 좋은 컬러 이미지(파일)로
// 바꾸는 표현(Presenter) 책임만 정의하는 작은 인터페이스(ISP) — 분할 계산
// 책임(ImageSegmenter)과는 분리된 별도의 포트다(SRP).
type ResultVisualizer interface {
	// Colorize는 분할 결과를 색칠하여 outputPath에 저장한다.
	Colorize(result SegmentationResult, outputPath string) error
}
