package domain

import "image"

// Tracker는 "첫 프레임의 박스 하나를 영상 전체에서 계속 따라간다"는 추적
// 알고리즘의 추상 인터페이스(포트)다.
//
//   - SRP: 이 인터페이스는 오직 "추적"만 책임진다. 영상 입출력이나 화면
//     표시(윈도우)는 포함하지 않는다(ISP).
//   - DIP: 유스케이스 계층은 이 인터페이스에만 의존하고, gocv 기반 구체
//     구현(TrackerMIL 등)은 infra 패키지에만 존재한다.
//   - OCP: CSRT/KCF 등 다른 내장 트래커나 딥러닝 기반 트래커로 교체하려면
//     infra 패키지에 이 인터페이스를 구현하는 새 어댑터만 추가하면 되고,
//     usecase 패키지는 전혀 수정할 필요가 없다.
//
// Frame은 gocv.Mat을 직접 노출하지 않기 위해 빈 인터페이스로 추상화하는
// 대신, 이 예제에서는 실용성을 위해 infra 계층이 실제 프레임 타입(gocv.Mat)을
// 다루고 usecase는 "프레임을 읽어오는 콜백"만 호출하는 방식으로 DIP를
// 적용한다(FrameSource 참고). Tracker 인터페이스 자체는 gocv 타입을
// 직접 노출하지 않도록 image.Rectangle(표준 라이브러리 타입)만 사용한다.
type Tracker interface {
	// Init은 초기 바운딩 박스로 추적기를 초기화한다. 첫 프레임은 어댑터가
	// 생성 시점 또는 FrameSource를 통해 내부적으로 확보한다.
	Init(initBox image.Rectangle) error

	// Update는 다음 프레임에서 추적 결과를 계산해 반환한다.
	Update() (TrackedObject, error)
}

// FrameSource는 비디오로부터 프레임을 순차적으로 읽어오는 포트다.
// usecase 계층은 이 인터페이스를 통해서만 "다음 프레임이 있는가"를 알며,
// 비디오 캡처(gocv.VideoCapture) 구체 타입을 전혀 알지 못한다(DIP).
type FrameSource interface {
	// Next는 다음 프레임을 읽을 수 있으면 true를 반환한다. 더 이상 프레임이
	// 없으면 false를 반환한다.
	Next() bool
}

// ResultSink는 한 프레임의 추적 결과를 소비(시각화/출력)하는 포트다.
// 추적 로직과 화면 표시 로직을 분리하기 위한 작은 인터페이스다(SRP/ISP).
type ResultSink interface {
	// Present는 이번 프레임의 추적 결과와 지금까지의 이동 경로를 받아
	// 화면에 그리거나 콘솔에 출력한다. 반환값 false는 "중단 요청(ESC 등)"을 뜻한다.
	Present(result TrackedObject, trail Trail, frameIdx int) bool
}
