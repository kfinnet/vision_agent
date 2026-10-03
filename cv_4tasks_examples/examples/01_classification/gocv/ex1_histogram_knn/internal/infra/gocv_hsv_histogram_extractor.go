// Package infra는 ①분류 실습1의 도메인 인터페이스를 gocv(OpenCV)로 구현하는 어댑터들을 담는다.
// 인프라/어댑터 계층(Infrastructure Layer) — 이 패키지만 "gocv.io/x/gocv"를 import하며,
// 실제 OpenCV API 호출은 모두 여기에 있다.
package infra

import (
	"fmt"

	"gocv.io/x/gocv"

	"ex1_histogram_knn/internal/domain"
)

// HSVHistogramExtractor는 domain.FeatureExtractor를 HSV(H·S) 2D 색상 히스토그램으로
// 구현하는 어댑터다. "분류"의 가장 단순한 고전적 접근 — 이미지 전체의 색 분포를
// "지문"처럼 사용한다. (밝기 채널 V는 조명 변화에 민감하므로 제외)
//
// OCP: HOG 등 다른 특징 추출 방식이 필요하면 FeatureExtractor를 구현하는 새 구조체를
//
//	추가하면 되고, 유스케이스 코드는 전혀 바뀌지 않는다.
//
// LSP: Extract()의 입출력 규약이 domain.FeatureExtractor와 완전히 동일하므로
//
//	FeatureExtractor를 기대하는 어떤 코드에도 안전하게 대입할 수 있다.
type HSVHistogramExtractor struct {
	Bins int
}

// NewHSVHistogramExtractor는 히스토그램 bin 개수를 받아 HSVHistogramExtractor를 생성한다.
func NewHSVHistogramExtractor(bins int) *HSVHistogramExtractor {
	return &HSVHistogramExtractor{Bins: bins}
}

// Extract는 이미지 파일을 읽어 HSV의 H(색상)·S(채도) 2채널 2D 히스토그램을 계산하고
// 0~1로 정규화한 뒤, 1차원 domain.Feature로 변환해 반환한다.
func (e *HSVHistogramExtractor) Extract(imagePath string) (domain.Feature, error) {
	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("이미지를 읽을 수 없습니다: %s", imagePath)
	}
	defer img.Close()

	hsv := gocv.NewMat()
	defer hsv.Close()
	gocv.CvtColor(img, &hsv, gocv.ColorBGRToHSV)

	hist := gocv.NewMat()
	defer hist.Close()
	mask := gocv.NewMat()
	defer mask.Close()

	chans := []int{0, 1}
	sizes := []int{e.Bins, e.Bins}
	ranges := []float64{0, 180, 0, 256}
	gocv.CalcHist([]gocv.Mat{hsv}, chans, mask, &hist, sizes, ranges, false)
	gocv.Normalize(hist, &hist, 0, 1, gocv.NormMinMax)

	data, err := hist.DataPtrFloat32()
	if err != nil {
		return nil, err
	}
	feature := make(domain.Feature, len(data))
	copy(feature, data)
	return feature, nil
}

// 컴파일 타임에 HSVHistogramExtractor가 domain.FeatureExtractor 포트를 만족하는지 확인한다.
var _ domain.FeatureExtractor = (*HSVHistogramExtractor)(nil)
