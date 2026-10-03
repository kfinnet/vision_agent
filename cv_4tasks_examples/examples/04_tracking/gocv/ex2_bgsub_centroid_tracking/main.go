// CV 4대 태스크 실습 — ④추적(Tracking) 실습2 (GoCV)
// 배경 차분(Background Subtraction, MOG2) + 중심점(centroid) 매칭 기반 다중 객체 추적
//
// Python/Colab 실습2(ex2_optical_flow_tracking.ipynb)가 특징점(코너) 단위로 추적했다면,
// 이 GoCV 실습은 강의안에서 설명한 "탐지(공간) + 매칭(시간) = 추적" 공식을 가장 직접적으로
// 구현한다: ① 매 프레임 배경 차분으로 "움직이는 물체"를 탐지(컨투어) → ② 이전 프레임에서
// 추적 중이던 객체들과 중심점 거리가 가장 가까운 것끼리 짝지어 같은 ID를 유지(매칭) →
// ③ 새로 나타난 물체는 새 ID를 부여, 오래 매칭되지 않은 ID는 제거.
//
// 이 파일은 구성 루트(Composition Root)다. CLI 인자를 해석하고, gocv 기반
// 구체 어댑터(internal/infra) + gocv-free 순수 매칭 알고리즘(internal/domain의
// CentroidTracker)을 생성해 유스케이스(internal/usecase)에 주입한 뒤 실행한다.
// Clean Architecture 구조:
//
//	main.go              — 구성 루트: 의존성 조립 + CLI 진입점
//	internal/domain/      — 값 타입(Detection, Track) + 인터페이스(포트) +
//	                        gocv-free 순수 알고리즘(CentroidTracker)
//	internal/usecase/      — 탐지→매칭→표시 루프 오케스트레이션(gocv 비의존)
//	internal/infra/        — gocv 기반 구체 구현(이 패키지만 gocv를 import)
//
// 사용법:
//
//	go run main.go <비디오경로>
//
// 예시:
//
//	go run main.go vtest.avi
package main

import (
	"fmt"
	"os"

	"ex2_bgsub_centroid_tracking/internal/domain"
	"ex2_bgsub_centroid_tracking/internal/infra"
	"ex2_bgsub_centroid_tracking/internal/usecase"
)

const minContourArea = 400.0 // 너무 작은 노이즈 컨투어는 무시

func main() {
	if len(os.Args) < 2 {
		fmt.Println("사용법: go run main.go <비디오경로>")
		os.Exit(1)
	}
	videoPath := os.Args[1]

	// --- 구성 루트: 구체 어댑터(infra) 생성 ---
	source, err := infra.NewGoCVVideoSource(videoPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer source.Close()

	detector := infra.NewGoCVMotionDetector(source, minContourArea)
	defer detector.Close()

	presenter := infra.NewGoCVWindowPresenter("GoCV 배경차분 + 중심점 매칭 추적 (ESC로 종료)", source)
	defer presenter.Close()

	// gocv에 의존하지 않는 순수 매칭 알고리즘 — domain 패키지에 직접 생성
	matcher := domain.NewCentroidTracker(domain.DefaultCentroidTrackerConfig())

	// --- 구성 루트: 추상 인터페이스/순수 타입으로 유스케이스에 주입(DIP) ---
	trackingUseCase := usecase.NewBgSubCentroidTrackingUseCase(source, detector, matcher, presenter)

	frameCount := trackingUseCase.Run()
	fmt.Println("영상이 끝났습니다. 총", frameCount, "프레임 처리")
}
