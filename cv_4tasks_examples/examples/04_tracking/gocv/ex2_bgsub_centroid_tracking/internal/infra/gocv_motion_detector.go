package infra

import (
	"image"

	"gocv.io/x/gocv"

	"ex2_bgsub_centroid_tracking/internal/domain"
)

// GoCVMotionDetector는 배경 차분(MOG2) + 모폴로지 정리 + 컨투어로 "움직이는
// 물체"를 탐지하는 domain.MotionDetector 구현체다.
//
// MOG2: 최근 프레임들을 학습해 "배경"을 추정하고, 거기서 벗어나는 픽셀
// (= 움직이는 물체)을 찾는다.
//
// SRP: 오직 "탐지"만 책임진다. 탐지된 박스들을 가지고 "이전 프레임과 같은
// 객체인지" 짝짓는 매칭(시간) 책임은 이 타입에 없다 — domain.CentroidTracker가
// 전담한다. OCP: 다른 배경 차분 알고리즘으로 바꾸려면 이 구조체와 같은
// 패턴의 새 어댑터만 추가하면 된다.
type GoCVMotionDetector struct {
	source         *GoCVVideoSource
	mog2           gocv.BackgroundSubtractorMOG2
	fgMask         gocv.Mat
	kernel         gocv.Mat
	minContourArea float64 // 너무 작은 노이즈 컨투어는 무시
}

// NewGoCVMotionDetector는 주어진 프레임 소스에서 움직임을 탐지하는
// GoCVMotionDetector를 생성한다.
func NewGoCVMotionDetector(source *GoCVVideoSource, minContourArea float64) *GoCVMotionDetector {
	return &GoCVMotionDetector{
		source:         source,
		mog2:           gocv.NewBackgroundSubtractorMOG2(),
		fgMask:         gocv.NewMat(),
		kernel:         gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5)),
		minContourArea: minContourArea,
	}
}

// Detect는 현재 프레임에서 배경 차분 → 모폴로지 정리 → 컨투어 탐색을 거쳐
// 움직이는 물체들의 박스 목록을 반환한다. domain.MotionDetector 인터페이스의
// 구현이다.
//
// 1) 탐지(공간): 배경 차분 → 모폴로지 정리 → 컨투어 찾기 → 중심점 계산
func (d *GoCVMotionDetector) Detect() []domain.Detection {
	frame := d.source.Frame()
	d.mog2.Apply(frame, &d.fgMask)
	gocv.MorphologyEx(d.fgMask, &d.fgMask, gocv.MorphOpen, d.kernel)
	contours := gocv.FindContours(d.fgMask, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	defer contours.Close()

	var detections []domain.Detection
	for i := 0; i < contours.Size(); i++ {
		c := contours.At(i)
		if gocv.ContourArea(c) < d.minContourArea {
			continue
		}
		r := gocv.BoundingRect(c)
		detections = append(detections, domain.NewDetection(r))
	}
	return detections
}

// Close는 내부 gocv 자원(MOG2, Mat)을 해제한다.
func (d *GoCVMotionDetector) Close() error {
	d.kernel.Close()
	d.fgMask.Close()
	return d.mog2.Close()
}
