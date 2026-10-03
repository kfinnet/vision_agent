// Package infra는 ===== [인프라/어댑터 계층 / Infrastructure Layer] =====.
// gocv 호출은 이 패키지에만 존재한다. WatershedSegmenter는 분할 "계산"만
// 담당하고, WatershedVisualizer는 결과를 "색칠해 파일로 저장"하는 표현만
// 담당한다 — 두 책임을 하나의 구조체로 합치지 않음으로써 각자 변경 이유가
// 분리된다(SRP).
package infra

import (
	"fmt"
	"image"
	"math/rand"

	"gocv.io/x/gocv"

	"ex1_watershed_segmentation/internal/domain"
)

// WatershedSegmenter는 Watershed 알고리즘으로 맞닿은 물체(예: 동전)를 픽셀
// 단위로 분리하는 어댑터다.
//
// domain.ImageSegmenter 인터페이스의 구체 구현체(Infrastructure) —
// usecase 패키지는 이 구조체의 존재를 모른 채 domain.ImageSegmenter 타입으로만
// 다룬다(DIP). 다른 분할 알고리즘을 추가하려면 이 인터페이스를 구현하는
// 새 구조체만 만들면 되며, usecase는 전혀 수정하지 않는다(OCP).
type WatershedSegmenter struct {
	// MorphIterations는 모폴로지 Open 연산 반복 횟수.
	MorphIterations int
	// DilateIterations는 확실한 배경(sureBg)을 구하기 위한 팽창 반복 횟수.
	DilateIterations int
}

// NewWatershedSegmenter는 기본 파라미터로 WatershedSegmenter를 생성한다.
func NewWatershedSegmenter() *WatershedSegmenter {
	return &WatershedSegmenter{MorphIterations: 2, DilateIterations: 3}
}

// Segment는 domain.ImageSegmenter를 구현한다. 그레이스케일 → Otsu 이진화 →
// 모폴로지 Open → 거리변환(distance transform) → 마커(marker) 생성 →
// Watershed 순으로 처리한다.
func (s *WatershedSegmenter) Segment(inputPath string) (domain.SegmentationResult, error) {
	// 1) 그레이스케일로 읽기
	gray := gocv.IMRead(inputPath, gocv.IMReadGrayScale)
	if gray.Empty() {
		return domain.SegmentationResult{}, fmt.Errorf("이미지를 읽을 수 없습니다: %s", inputPath)
	}
	defer gray.Close()

	// 2) Otsu 이진화(전경/배경 자동 분리) + 모폴로지 Open(작은 잡음 제거)
	thresh := gocv.NewMat()
	defer thresh.Close()
	gocv.ThresholdWithParams(gray, &thresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	defer kernel.Close()
	opening := gocv.NewMat()
	defer opening.Close()
	gocv.MorphologyEx(thresh, &opening, gocv.MorphOpen, kernel)

	// 3) 확실한 배경(sureBg) / 확실한 전경(sureFg) / 모름(unknown) 영역 계산
	sureBg := gocv.NewMat()
	defer sureBg.Close()
	gocv.Dilate(opening, &sureBg, kernel)

	distTmp := gocv.NewMat()
	defer distTmp.Close()
	dist := gocv.NewMat()
	defer dist.Close()
	gocv.DistanceTransform(opening, &dist, &distTmp, gocv.DistL2, gocv.DistMask5, gocv.DistanceLabelCComp)

	_, maxVal, _, _ := gocv.MinMaxLoc(dist)
	sureFg := gocv.NewMat()
	defer sureFg.Close()
	gocv.Threshold(dist, &sureFg, 0.5*maxVal, 255, gocv.ThresholdBinary)
	sureFg.ConvertTo(&sureFg, gocv.MatTypeCV8U)

	unknown := gocv.NewMat()
	defer unknown.Close()
	gocv.Subtract(sureBg, sureFg, &unknown)

	// 4) 마커(marker) 생성 — 연결 요소(connected components)마다 서로 다른 번호를 매긴다
	markers := gocv.NewMat()
	defer markers.Close()
	numMarkers := gocv.ConnectedComponents(sureFg, &markers)
	fmt.Printf("물체 후보 개수: %d\n", numMarkers-1)

	markers32 := gocv.NewMat()
	defer markers32.Close()
	markers.ConvertTo(&markers32, gocv.MatTypeCV32S)

	rows, cols := markers32.Rows(), markers32.Cols()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			v := markers32.GetIntAt(y, x)
			markers32.SetIntAt(y, x, v+1) // 배경을 0이 아닌 1로 (Watershed 규약)
		}
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			if unknown.GetUCharAt(y, x) == 255 {
				markers32.SetIntAt(y, x, 0) // "모름" 영역 = Watershed가 채워나갈 부분
			}
		}
	}

	// 5) Watershed 실행 (원본은 3채널 BGR이 필요)
	colorImg := gocv.NewMat()
	defer colorImg.Close()
	gocv.CvtColor(gray, &colorImg, gocv.ColorGrayToBGR)
	gocv.Watershed(colorImg, &markers32)

	// 결과를 domain.SegmentationResult(평면 []int32 라벨 맵)로 변환
	flat := make([]int32, rows*cols)
	maxLabel := int32(0)
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			v := markers32.GetIntAt(y, x)
			flat[y*cols+x] = v
			if v > maxLabel {
				maxLabel = v
			}
		}
	}
	numSegments := int(maxLabel)
	fmt.Printf("최종 분할된 물체 개수: %d\n", numSegments-1)

	return domain.SegmentationResult{
		Markers:     flat,
		Rows:        rows,
		Cols:        cols,
		NumSegments: numSegments,
	}, nil
}

// WatershedVisualizer는 SegmentationResult를 물체마다 다른 색으로 칠해
// 이미지 파일로 저장하는 표현(Presenter) 전담 구조체다.
//
// 분할 "계산"(WatershedSegmenter)과는 완전히 분리된 책임이다(SRP) — 이
// 구조체를 다른 시각화 방식으로 교체해도 분할 알고리즘/유스케이스는 전혀
// 영향받지 않는다.
type WatershedVisualizer struct {
	// Seed는 세그먼트별 색상을 결정하는 난수 시드 (재현 가능한 색상 배정을 위함).
	Seed int64
}

// NewWatershedVisualizer는 기본 시드로 WatershedVisualizer를 생성한다.
func NewWatershedVisualizer() *WatershedVisualizer {
	return &WatershedVisualizer{Seed: 42}
}

// Colorize는 domain.ResultVisualizer를 구현한다. 경계선은 흰색, 배경은
// 연회색, 각 물체는 임의의 색으로 칠해 outputPath에 저장한다.
func (v *WatershedVisualizer) Colorize(result domain.SegmentationResult, outputPath string) error {
	rows, cols := result.Rows, result.Cols
	out := gocv.NewMatWithSize(rows, cols, gocv.MatTypeCV8UC3)
	defer out.Close()
	out.SetTo(gocv.NewScalar(235, 235, 235, 0))

	rng := rand.New(rand.NewSource(v.Seed))
	colorMap := map[int32][3]uint8{}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			label := result.At(y, x)
			switch {
			case label == -1:
				out.SetUCharAt3(y, x, 0, 255)
				out.SetUCharAt3(y, x, 1, 255)
				out.SetUCharAt3(y, x, 2, 255)
			case label <= 1:
				// 배경 — 미리 채운 연회색 유지
			default:
				c, ok := colorMap[label]
				if !ok {
					c = [3]uint8{uint8(rng.Intn(195) + 60), uint8(rng.Intn(195) + 60), uint8(rng.Intn(195) + 60)}
					colorMap[label] = c
				}
				out.SetUCharAt3(y, x, 0, c[0])
				out.SetUCharAt3(y, x, 1, c[1])
				out.SetUCharAt3(y, x, 2, c[2])
			}
		}
	}

	if ok := gocv.IMWrite(outputPath, out); !ok {
		return fmt.Errorf("결과 저장 실패: %s", outputPath)
	}
	return nil
}
