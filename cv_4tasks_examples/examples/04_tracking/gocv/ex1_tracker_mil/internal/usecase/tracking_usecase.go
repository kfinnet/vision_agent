// Package usecase는 ④추적 실습1의 유스케이스 계층이다. "비디오에서 프레임을
// 순차적으로 읽으며 추적기를 갱신하고 이동 경로를 누적한 뒤 결과를 전달한다"는
// 오케스트레이션 로직만 담으며, gocv를 직접 import하지 않는다(DIP).
package usecase

import (
	"fmt"
	"image"

	"ex1_tracker_mil/internal/domain"
)

// TrackingUseCase는 단일 대상 추적의 전체 흐름을 조율한다.
//
//   - SRP: "프레임 읽기 → 추적기 Init/Update 호출 → 경로 누적 → 결과 전달"
//     흐름만 책임지며, 실제 추적 알고리즘이나 화면 표시 방식은 알지 못한다.
//   - DIP: 생성자로 Tracker/FrameSource/ResultSink 세 개의 추상 인터페이스만
//     주입받는다. 이들의 gocv 기반 구체 구현(infra 패키지)은 전혀 모른다.
type TrackingUseCase struct {
	tracker domain.Tracker
	source  domain.FrameSource
	sink    domain.ResultSink
}

// NewTrackingUseCase는 추적기, 프레임 소스, 결과 싱크를 주입받아
// TrackingUseCase를 생성한다(생성자 주입 — DIP의 실제 적용 지점).
func NewTrackingUseCase(tracker domain.Tracker, source domain.FrameSource, sink domain.ResultSink) *TrackingUseCase {
	return &TrackingUseCase{tracker: tracker, source: source, sink: sink}
}

// Run은 첫 프레임으로 추적기를 초기화한 뒤, 영상이 끝나거나 ResultSink가
// 중단을 요청할 때까지 매 프레임 추적을 반복한다. 반환값은 처리된 프레임 수다.
func (u *TrackingUseCase) Run(initBox image.Rectangle) (int, error) {
	// 1) 추적기 초기화 — 첫 프레임 + 초기 박스로 외형(appearance) 모델을 학습
	if !u.source.Next() {
		return 0, fmt.Errorf("첫 프레임을 읽을 수 없습니다")
	}
	if err := u.tracker.Init(initBox); err != nil {
		return 0, fmt.Errorf("추적기 초기화 실패: %w", err)
	}

	trail := domain.Trail{}
	center := image.Pt(initBox.Min.X+initBox.Dx()/2, initBox.Min.Y+initBox.Dy()/2)
	trail.Add(center)

	frameIdx := 0
	for u.source.Next() {
		frameIdx++

		// 2) 매 프레임 Update 호출 — "이전 외형과 가장 비슷한 위치"를 현재 프레임에서 찾는다
		result, err := u.tracker.Update()
		if err != nil || !result.Found {
			fmt.Println("추적 실패(대상을 놓침) — frame", frameIdx)
			continue
		}
		trail.Add(result.Center)

		if !u.sink.Present(result, trail, frameIdx) {
			break
		}
	}
	return frameIdx, nil
}
