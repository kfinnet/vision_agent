// Package infra는 ===== [인프라/어댑터 계층 / Infrastructure Layer] =====.
// gocv 호출은 이 패키지에만 존재한다. GrabCutSegmenter는 분할 "계산"만
// 담당하고, GrabCutPresenter는 결과 마스크를 적용한 파생 이미지(배경 제거)를
// 만들어 파일로 저장하는 표현만 담당한다 — 두 책임을 하나의 구조체로 합치지
// 않음으로써 각자 변경 이유가 분리된다(SRP).
package infra

import (
	"fmt"
	"image"

	"gocv.io/x/gocv"

	"ex2_grabcut_segmentation/internal/domain"
)

// GrabCutSegmenter는 GrabCut 알고리즘으로 사각형 힌트 하나만으로 전경/배경을
// 분리하는 어댑터다.
//
// domain.ImageSegmenter 인터페이스의 구체 구현체(Infrastructure) —
// usecase 패키지는 이 구조체의 존재를 모른 채 domain.ImageSegmenter 타입으로만
// 다룬다(DIP). 다른 전경 분할 알고리즘을 추가하려면 이 인터페이스를 구현하는
// 새 구조체만 만들면 되며, usecase는 전혀 수정하지 않는다(OCP).
type GrabCutSegmenter struct {
	// MarginRatio는 입력 이미지 가로 폭 대비 초기 사각형의 좌우 여백 비율.
	// 세로 여백은 MarginRatio*1.5 비율로 계산된다(기존 Python/Go 구현과 동일).
	MarginRatio float64
	// IterCount는 GrabCut 반복 최적화 횟수.
	IterCount int
}

// NewGrabCutSegmenter는 기본 여백 비율(0.06)로 GrabCutSegmenter를 생성한다.
func NewGrabCutSegmenter() *GrabCutSegmenter {
	return &GrabCutSegmenter{MarginRatio: 0.06, IterCount: 8}
}

// Segment는 domain.ImageSegmenter를 구현한다. 사각형 힌트 하나를 "아마도
// 전경"으로 가정하고 반복적인 그래프 컷 최적화로 전경 마스크를 계산한다.
func (s *GrabCutSegmenter) Segment(inputPath string) (domain.SegmentationResult, error) {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return domain.SegmentationResult{}, fmt.Errorf("이미지를 읽을 수 없습니다: %s", inputPath)
	}
	defer img.Close()

	w, h := img.Cols(), img.Rows()
	mx, my := int(float64(w)*s.MarginRatio), int(float64(h)*s.MarginRatio*1.5)
	rect := image.Rect(mx, my, w-mx, h-my)

	// GrabCut 입출력 버퍼 — 마스크(전경/배경 레이블), 내부용 배경/전경 가우시안 혼합모델
	mask := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8U)
	defer mask.Close()
	bgdModel := gocv.NewMat()
	defer bgdModel.Close()
	fgdModel := gocv.NewMat()
	defer fgdModel.Close()

	// GrabCut 실행: rect 안을 "아마도 전경"으로 가정하고 반복 최적화
	gocv.GrabCut(img, &mask, rect, &bgdModel, &fgdModel, s.IterCount, gocv.GCInitWithRect)

	// 마스크 값: 0=확실한 배경, 1=확실한 전경, 2=아마도 배경, 3=아마도 전경
	// (1, 3)을 최종 전경으로 합쳐 이진 마스크로 변환
	rows, cols := mask.Rows(), mask.Cols()
	flat := make([]uint8, rows*cols)
	fgPixels := 0
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			v := mask.GetUCharAt(y, x)
			if v == 1 || v == 3 {
				flat[y*cols+x] = 255
				fgPixels++
			}
		}
	}
	ratio := float64(fgPixels) / float64(rows*cols)
	fmt.Printf("전경으로 분류된 픽셀 비율: %.1f%%\n", 100*ratio)

	return domain.SegmentationResult{
		ForegroundMask:  flat,
		Rows:            rows,
		Cols:            cols,
		ForegroundRatio: ratio,
		InitRect:        domain.Rect{X: rect.Min.X, Y: rect.Min.Y, W: rect.Dx(), H: rect.Dy()},
	}, nil
}

// GrabCutPresenter는 SegmentationResult의 마스크를 원본 이미지에 적용해
// 배경을 제거한 파생 이미지를 파일로 저장하는 표현(Presenter) 전담 구조체다.
//
// 분할 "계산"(GrabCutSegmenter)과는 완전히 분리된 책임이다(SRP) — 이 구조체를
// 다른 표현 방식(예: 배경 블러)으로 교체해도 분할 알고리즘/유스케이스는 전혀
// 영향받지 않는다.
type GrabCutPresenter struct{}

// NewGrabCutPresenter는 GrabCutPresenter를 생성한다.
func NewGrabCutPresenter() *GrabCutPresenter {
	return &GrabCutPresenter{}
}

// RemoveBackground는 domain.ResultPresenter를 구현한다. 픽셀 단위 마스크로
// 배경을 제거(전경만 남기고 나머지는 검은색)한 이미지를 outputPath에 저장한다.
func (p *GrabCutPresenter) RemoveBackground(inputPath string, result domain.SegmentationResult, outputPath string) error {
	img := gocv.IMRead(inputPath, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("이미지를 읽을 수 없습니다: %s", inputPath)
	}
	defer img.Close()

	rows, cols := result.Rows, result.Cols
	binMask := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8U)
	defer binMask.Close()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			binMask.SetUCharAt(y, x, result.At(y, x))
		}
	}

	foreground := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	defer foreground.Close()
	img.CopyToWithMask(&foreground, binMask)

	if ok := gocv.IMWrite(outputPath, foreground); !ok {
		return fmt.Errorf("결과 저장 실패: %s", outputPath)
	}
	return nil
}
