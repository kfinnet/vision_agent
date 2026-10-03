// Package usecase는 ④추적 실습2의 유스케이스 계층이다. "프레임을 읽으며
// 탐지기를 호출하고, 중심점 매칭기로 ID를 유지한 뒤 결과를 전달한다"는
// 오케스트레이션 로직만 담으며, gocv를 직접 import하지 않는다(DIP).
package usecase

import (
	"fmt"

	"ex2_bgsub_centroid_tracking/internal/domain"
)

// BgSubCentroidTrackingUseCase는 "탐지(공간) + 매칭(시간) = 추적" 흐름 전체를
// 조율한다.
//
//   - SRP: "프레임 읽기 → 탐지 호출 → 매칭 호출 → 결과 전달" 흐름만
//     책임진다. 배경 차분이나 컨투어 처리(탐지), 그리기(표시)는 전혀 모른다.
//   - DIP: 생성자로 MotionDetector / FrameSource / ResultPresenter 세 개의
//     추상 인터페이스와, gocv에 의존하지 않는 순수 도메인 타입인
//     *domain.CentroidTracker를 주입받는다.
type BgSubCentroidTrackingUseCase struct {
	source   domain.FrameSource
	detector domain.MotionDetector
	matcher  *domain.CentroidTracker
	sink     domain.ResultPresenter
}

// NewBgSubCentroidTrackingUseCase는 프레임 소스, 탐지기, 중심점 매칭기,
// 결과 표시기를 주입받아 유스케이스를 생성한다(생성자 주입 — DIP).
func NewBgSubCentroidTrackingUseCase(
	source domain.FrameSource,
	detector domain.MotionDetector,
	matcher *domain.CentroidTracker,
	sink domain.ResultPresenter,
) *BgSubCentroidTrackingUseCase {
	return &BgSubCentroidTrackingUseCase{source: source, detector: detector, matcher: matcher, sink: sink}
}

// Run은 영상이 끝나거나 ResultPresenter가 중단을 요청할 때까지 매 프레임
// 탐지 → 매칭 → 표시를 반복한다. 반환값은 처리된 프레임 수다.
func (u *BgSubCentroidTrackingUseCase) Run() int {
	frameIdx := 0
	for u.source.Next() {
		frameIdx++

		// 1) 탐지(공간): infra의 MotionDetector에 위임(배경 차분 → 컨투어 → 박스)
		detections := u.detector.Detect()

		// 2) 매칭(시간): gocv-free 순수 알고리즘인 CentroidTracker에 위임
		tracks := u.matcher.Update(detections)

		fmt.Printf("frame %d: 탐지 %d개 / 유지 중인 추적 ID %d개\n", frameIdx, len(detections), len(tracks))

		if !u.sink.Present(detections, tracks, frameIdx) {
			break
		}
	}
	return frameIdx
}
