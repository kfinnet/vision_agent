// Package main은 ===== [구성 루트 / Composition Root] =====.
// CV 4대 태스크 실습 — ③분할(Segmentation) 실습1 (GoCV)
// Watershed 알고리즘 — 서로 맞닿은 물체(동전 등)를 픽셀 단위로 분리
//
// Python/Colab 실습1(ex1_watershed_segmentation.ipynb)과 동일한 파이프라인을
// 순수 Go + GoCV로, Clean Architecture 계층(domain/usecase/infra)으로
// 분리해 구현합니다.
//
// main.go는 CLI 인자 파싱과 의존성 조립(구체 어댑터 생성 → 유스케이스에
// 생성자 주입)만 담당하는 구성 루트이며, gocv를 직접 호출하는 분할 로직은
// 전혀 포함하지 않습니다 — 그런 로직은 모두 internal/infra에 있습니다.
//
// 처리 흐름:
//
//	그레이스케일 → Otsu 이진화 → 모폴로지 Open → 거리변환(distance transform)
//	→ 마커(marker) 생성 → Watershed → 결과를 색칠하여 저장
//
// 사용법:
//
//	go run main.go <입력이미지(물체 사진)> <출력이미지>
//
// 예시:
//
//	go run main.go coins.png watershed_result.png
package main

import (
	"fmt"
	"os"

	"ex1_watershed_segmentation/internal/domain"
	"ex1_watershed_segmentation/internal/infra"
	"ex1_watershed_segmentation/internal/usecase"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("사용법: go run main.go <입력이미지> <출력이미지>")
		os.Exit(1)
	}
	inputPath, outputPath := os.Args[1], os.Args[2]

	// 의존성 조립(Composition) — 구체 어댑터 생성 → 유스케이스에 주입
	var segmenter domain.ImageSegmenter = infra.NewWatershedSegmenter()
	var visualizer domain.ResultVisualizer = infra.NewWatershedVisualizer()
	uc := usecase.NewWatershedUseCase(segmenter, visualizer)

	if _, err := uc.Execute(inputPath, outputPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("결과 저장 완료:", outputPath)
}
