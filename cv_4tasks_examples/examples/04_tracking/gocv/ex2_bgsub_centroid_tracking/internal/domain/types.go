// Package domain은 ④추적 실습2(배경차분 + 중심점 매칭)의 도메인 계층이다.
// 외부 CV 라이브러리(gocv)에 전혀 의존하지 않는 순수 값 타입/엔티티와,
// gocv 없이도 동작하는 "중심점 매칭" 알고리즘(CentroidTracker)을 담는다.
package domain

import (
	"image"
	"math"
)

// Detection은 한 프레임에서 탐지된 "움직이는 물체" 하나(박스 + 중심점)를
// 나타내는 값 객체다. 배경 차분/컨투어 같은 구현 세부사항을 전혀 담지 않는다.
type Detection struct {
	Box    image.Rectangle
	Center image.Point
}

// NewDetection은 박스로부터 중심점을 계산해 Detection을 만든다.
func NewDetection(box image.Rectangle) Detection {
	return Detection{Box: box, Center: centroid(box)}
}

func centroid(r image.Rectangle) image.Point {
	return image.Pt(r.Min.X+r.Dx()/2, r.Min.Y+r.Dy()/2)
}

func dist(a, b image.Point) float64 {
	dx, dy := float64(a.X-b.X), float64(a.Y-b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Track은 여러 프레임에 걸쳐 유지되는 추적 대상 하나(고유 ID + 현재 중심점)다.
type Track struct {
	ID          int
	Center      image.Point
	FramesSince int // 마지막으로 매칭된 뒤 지난 프레임 수
	Active      bool
}

// CentroidTrackerConfig는 CentroidTracker의 매칭 민감도를 조절하는 파라미터다.
type CentroidTrackerConfig struct {
	MaxMatchDist    float64 // 이 거리(픽셀) 이내여야 "같은 객체"로 매칭
	MaxMissedFrames int     // 이 프레임 수만큼 매칭 안 되면 추적 ID 제거
}

// DefaultCentroidTrackerConfig는 실습 예제에서 쓰는 기본 파라미터다.
func DefaultCentroidTrackerConfig() CentroidTrackerConfig {
	return CentroidTrackerConfig{MaxMatchDist: 60.0, MaxMissedFrames: 10}
}

// CentroidTracker는 "이전 프레임에서 추적 중이던 객체들"과 "이번 프레임의
// 탐지 결과"를 중심점 간 최근접 거리로 짝짓고, ID 수명을 관리하는 순수
// 알고리즘이다.
//
//   - 이 타입은 점(image.Point)과 설정값만 다루며 gocv를 전혀 import하지
//     않는다 — "탐지(공간)"는 infra의 MotionDetector가 책임지고, "매칭(시간)"은
//     이 타입이 책임진다는 SRP 분리를 보여주는 대표적인 예다.
//   - DIP: usecase 계층은 이 구체 타입을 직접 들고 orchestration만 수행하며,
//     gocv 기반 구현은 전혀 모른다. CentroidTracker 자체가 gocv-free이므로
//     domain 계층에 두어도 "인프라에 의존하지 않는 순수 비즈니스 규칙"이라는
//     Clean Architecture의 도메인 정의에 부합한다.
type CentroidTracker struct {
	cfg    CentroidTrackerConfig
	tracks []Track
	nextID int
}

// NewCentroidTracker는 주어진 설정으로 빈 CentroidTracker를 생성한다.
func NewCentroidTracker(cfg CentroidTrackerConfig) *CentroidTracker {
	return &CentroidTracker{cfg: cfg, nextID: 1}
}

// Update는 이번 프레임의 탐지 목록을 기존 추적 목록과 매칭하고, 매칭/신규/소멸
// 규칙을 적용한 뒤 현재 활성 추적 목록을 반환한다.
//
// 2) 매칭(시간): 이전 프레임 추적 목록과 이번 프레임 탐지 결과를 최근접 거리로 매칭
func (c *CentroidTracker) Update(detections []Detection) []Track {
	matchedDet := make([]bool, len(detections))
	var updated []Track

	for _, t := range c.tracks {
		bestIdx, bestDist := -1, c.cfg.MaxMatchDist
		for di, d := range detections {
			if matchedDet[di] {
				continue
			}
			if dd := dist(t.Center, d.Center); dd < bestDist {
				bestDist, bestIdx = dd, di
			}
		}
		if bestIdx >= 0 {
			matchedDet[bestIdx] = true
			updated = append(updated, Track{ID: t.ID, Center: detections[bestIdx].Center, FramesSince: 0, Active: true})
		} else if t.FramesSince+1 < c.cfg.MaxMissedFrames {
			updated = append(updated, Track{ID: t.ID, Center: t.Center, FramesSince: t.FramesSince + 1, Active: false})
		}
		// FramesSince가 임계값을 넘으면 추적 목록에서 자연히 제외(= ID 소멸)
	}

	// 매칭되지 않은 새 탐지 = 새로운 객체 → 새 ID 부여
	for di, matched := range matchedDet {
		if !matched {
			updated = append(updated, Track{ID: c.nextID, Center: detections[di].Center, FramesSince: 0, Active: true})
			c.nextID++
		}
	}

	c.tracks = updated
	return c.tracks
}
