package domain

// Frame은 비디오/이미지 한 프레임을 추상화한 인터페이스다.
// main.go와 유스케이스 계층이 gocv.Mat 같은 구체 타입을 직접 참조하지
// 않도록(DIP), "비어 있는지 확인할 수 있는 프레임"이라는 최소 동작만
// 정의한다. 실제 구현(gocv.Mat 래퍼)은 인프라 계층에 있다.
type Frame interface {
	Empty() bool
}

// ObjectDetector는 객체 탐지기의 포트(인터페이스)다.
// ISP에 따라 "탐지"라는 단일 책임만 가지며, 시각화나 캡처 등의 책임은
// 포함하지 않는다. Haar Cascade, HOG, 추후 추가될 DNN 기반 탐지기 등
// 모든 구체 구현은 이 인터페이스를 만족해야 하고(LSP), 유스케이스
// 계층은 오직 이 추상화에만 의존한다(DIP).
type ObjectDetector interface {
	// Detect는 주어진 프레임에서 객체를 탐지해 Detection 목록을 반환한다.
	Detect(frame Frame) ([]Detection, error)
}

// VideoSource는 비디오/웹캠 프레임을 순차적으로 읽어오는 포트다.
// 실제 구현(웹캠 장치, 비디오 파일 등)은 인프라 계층에 있다.
type VideoSource interface {
	// Read는 다음 프레임을 읽어 ok로 성공 여부를 반환한다.
	Read() (frame Frame, ok bool)
	// Close는 사용 중인 자원을 해제한다.
	Close() error
}

// FrameRenderer는 프레임과 탐지 결과를 화면에 표시(오버레이 그리기 +
// 창 출력)하는 포트다. 탐지(ObjectDetector)와 렌더링 책임을 분리해
// ISP/SRP를 만족시킨다 — 렌더링 방식을 바꾸더라도 탐지 로직이나
// 유스케이스는 영향을 받지 않는다(OCP).
type FrameRenderer interface {
	// Render는 frame 위에 detections를 그려 창에 출력한다.
	Render(frame Frame, detections []Detection, frameCount int)
	// WaitKey는 지정한 시간(ms) 동안 키 입력을 기다리고 키 코드를 반환한다.
	WaitKey(delayMs int) int
	// Close는 창 등 사용 중인 자원을 해제한다.
	Close() error
}
