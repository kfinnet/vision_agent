package domain

// Frame은 비디오/이미지 한 프레임을 추상화한 인터페이스다.
// main.go와 유스케이스 계층이 gocv.Mat 같은 구체 타입을 직접 참조하지
// 않도록(DIP), "비어 있는지 확인할 수 있는 프레임"이라는 최소 동작만
// 정의한다. 실제 구현(gocv.Mat 래퍼)은 인프라 계층에 있다.
type Frame interface {
	Empty() bool
}

// ObjectDetector는 객체 탐지기의 포트(인터페이스)다.
// ISP에 따라 "탐지"라는 단일 책임만 가지며, 시각화/저장 등의 책임은
// 포함하지 않는다. HOG 보행자 탐지기, 얼굴 탐지기 등 모든 구체 구현은
// 이 인터페이스를 만족해야 하고(LSP), 유스케이스 계층은 오직 이
// 추상화에만 의존한다(DIP).
type ObjectDetector interface {
	// Detect는 주어진 프레임에서 객체를 탐지해 Detection 목록을 반환한다.
	Detect(frame Frame) ([]Detection, error)
}
