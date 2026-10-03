package domain

// MotionDetector는 한 프레임에서 "움직이는 물체"를 탐지하는 추상 인터페이스
// (포트)다.
//
//   - SRP/ISP: 오직 "탐지"만 책임진다. 매칭(시간 축 ID 유지)은 CentroidTracker가,
//     시각화는 ResultPresenter가 각각 별도로 책임진다 — 배경 차분, 중심점
//     매칭, 그리기 세 책임을 하나로 묶지 않는다.
//   - DIP: usecase 계층은 이 인터페이스에만 의존한다. 실제 배경 차분
//     알고리즘(MOG2 등)의 gocv 기반 구현은 infra 패키지에만 존재한다.
//   - OCP: MOG2 대신 다른 배경 차분 알고리즘(KNN 등)이나 딥러닝 기반 객체
//     탐지기로 바꾸려면 이 인터페이스를 구현하는 새 infra 어댑터만 추가하면
//     되고, usecase/CentroidTracker는 전혀 수정할 필요가 없다.
type MotionDetector interface {
	// Detect는 현재 프레임에서 탐지된 움직이는 물체 목록을 반환한다.
	Detect() []Detection
}

// FrameSource는 비디오로부터 프레임을 순차적으로 읽어오는 포트다.
type FrameSource interface {
	// Next는 다음 프레임을 읽을 수 있으면 true를 반환한다.
	Next() bool
}

// ResultPresenter는 한 프레임의 탐지/추적 결과를 소비(시각화/출력)하는
// 포트다. 탐지·매칭 로직과 화면 표시 로직을 분리하기 위한 작은
// 인터페이스다(SRP/ISP).
type ResultPresenter interface {
	// Present는 이번 프레임의 탐지 박스와 유지 중인 추적 목록을 받아 화면에
	// 그리거나 콘솔에 출력한다. 반환값 false는 "중단 요청(ESC 등)"을 뜻한다.
	Present(detections []Detection, tracks []Track, frameIdx int) bool
}
