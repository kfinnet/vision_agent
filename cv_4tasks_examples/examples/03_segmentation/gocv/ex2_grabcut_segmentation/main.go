// Package main은 ===== [구성 루트 / Composition Root] =====.
// CV 4대 태스크 실습 — ③분할(Segmentation) 실습2 (GoCV)
// GrabCut — 사각형 힌트 하나로 전경/배경을 픽셀 단위로 분리
//
// Python/Colab 실습2(ex2_grabcut_segmentation.ipynb)와 동일한 개념을
// 순수 Go + GoCV로, Clean Architecture 계층(domain/usecase/infra)으로
// 분리해 구현합니다. Watershed(실습1)가 촘촘한 마커가 필요했다면, GrabCut은
// "전경이 대략 여기 있다"는 사각형 하나만 주면 반복적인 그래프 컷 최적화로
// 정교한 픽셀 단위 마스크를 스스로 찾아낸다.
//
// main.go는 CLI 인자 파싱과 의존성 조립(구체 어댑터 생성 → 유스케이스에
// 생성자 주입)만 담당하는 구성 루트이며, gocv를 직접 호출하는 분할 로직은
// 전혀 포함하지 않습니다 — 그런 로직은 모두 internal/infra에 있습니다.
//
// 사용법:
//
//	go run main.go <입력이미지> <출력이미지(배경 제거 결과)> [여백비율(기본 0.06)]
//
// 예시:
//
//	go run main.go cat.jpg cat_foreground.png
package main

import (
	"fmt"
	"os"
	"strconv"

	"ex2_grabcut_segmentation/internal/domain"
	"ex2_grabcut_segmentation/internal/infra"
	"ex2_grabcut_segmentation/internal/usecase"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("사용법: go run main.go <입력이미지> <출력이미지> [여백비율]")
		os.Exit(1)
	}
	inputPath, outputPath := os.Args[1], os.Args[2]
	margin := 0.06
	if len(os.Args) >= 4 {
		if m, err := strconv.ParseFloat(os.Args[3], 64); err == nil {
			margin = m
		}
	}

	// 의존성 조립(Composition) — 구체 어댑터 생성 → 유스케이스에 주입
	segmenterImpl := infra.NewGrabCutSegmenter()
	segmenterImpl.MarginRatio = margin
	var segmenter domain.ImageSegmenter = segmenterImpl
	var presenter domain.ResultPresenter = infra.NewGrabCutPresenter()
	uc := usecase.NewGrabCutUseCase(segmenter, presenter)

	if _, err := uc.Execute(inputPath, outputPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("결과 저장 완료:", outputPath)
}
