package infra

import (
	"image"

	"gocv.io/x/gocv"

	"ex1_tracker_mil/internal/domain"
)

// GoCVTrackerMIL은 OpenCV가 제공하는 범용 Tracker API(TrackerMIL)를 사용해
// domain.Tracker 인터페이스를 구현하는 어댑터다.
//
// Python/Colab 실습1(ex1_meanshift_tracking.ipynb)이 MeanShift(색 히스토그램
// 기반)를 썼다면, 이 GoCV 실습은 TrackerMIL을 사용한다. 매 프레임 탐지를
// 새로 수행하지 않고, 첫 프레임에서 지정한 박스의 외형 모델을 학습해 다음
// 프레임들에서 "가장 비슷해 보이는 위치"를 계속 찾아간다.
//
// OCP: CSRT/KCF 등 다른 내장 트래커로 바꾸려면 이 파일과 같은 패턴으로 새
// 어댑터만 추가하면 되고, usecase 패키지는 전혀 수정할 필요가 없다.
type GoCVTrackerMIL struct {
	tracker gocv.TrackerMIL
	source  *GoCVVideoSource // 같은 패키지 내부에서 현재 프레임을 공유받는 지점
}

// NewGoCVTrackerMIL은 프레임을 공급하는 GoCVVideoSource와 함께
// GoCVTrackerMIL을 생성한다.
func NewGoCVTrackerMIL(source *GoCVVideoSource) *GoCVTrackerMIL {
	return &GoCVTrackerMIL{tracker: gocv.NewTrackerMIL(), source: source}
}

// Init은 현재(첫) 프레임 + 초기 박스로 외형(appearance) 모델을 학습한다.
// domain.Tracker 인터페이스의 구현이다.
func (t *GoCVTrackerMIL) Init(initBox image.Rectangle) error {
	t.tracker.Init(t.source.Frame(), initBox)
	return nil
}

// Update는 "이전 외형과 가장 비슷한 위치"를 현재 프레임에서 찾아
// domain.TrackedObject로 반환한다. domain.Tracker 인터페이스의 구현이다.
func (t *GoCVTrackerMIL) Update() (domain.TrackedObject, error) {
	rect, ok := t.tracker.Update(t.source.Frame())
	return domain.NewTrackedObject(rect, ok), nil
}

// Close는 내부 gocv 추적기 자원을 해제한다.
func (t *GoCVTrackerMIL) Close() error {
	return t.tracker.Close()
}
