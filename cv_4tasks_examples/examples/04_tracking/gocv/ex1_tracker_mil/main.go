// CV 4대 태스크 실습 — ④추적(Tracking) 실습1 (GoCV)
// 내장 Tracker(TrackerMIL) — 첫 프레임의 박스 하나를 영상 전체에서 계속 따라간다
//
// Python/Colab 실습1(ex1_meanshift_tracking.ipynb)이 MeanShift(색 히스토그램 기반)를
// 썼다면, 이 GoCV 실습은 OpenCV가 제공하는 범용 "Tracker" API(TrackerMIL)를 사용한다.
// 매 프레임 탐지를 새로 수행하지 않고, 첫 프레임에서 지정한 박스의 외형 모델을 학습해
// 다음 프레임들에서 "가장 비슷해 보이는 위치"를 계속 찾아간다 — 탐지(공간) 결과를
// 시간 축으로 연결하는 추적(Tracking)의 핵심 개념을 가장 단순하게 보여주는 API다.
//
// 이 파일은 구성 루트(Composition Root)다. CLI 인자를 해석하고, gocv 기반
// 구체 어댑터(internal/infra)를 생성해 유스케이스(internal/usecase)에
// 주입한 뒤 실행한다 — 추적 알고리즘 자체의 로직은 여기 없다(internal/infra에
// 있다). Clean Architecture 구조:
//
//	main.go              — 구성 루트: 의존성 조립 + CLI 진입점
//	internal/domain/      — 값 타입(TrackedObject, Trail) + 인터페이스(포트)
//	internal/usecase/      — 추적 루프 오케스트레이션(gocv 비의존)
//	internal/infra/        — gocv 기반 구체 구현(이 패키지만 gocv를 import)
//
// 사용법:
//
//	go run main.go <비디오경로> <x> <y> <w> <h>
//	- x,y,w,h: 첫 프레임에서 추적을 시작할 초기 바운딩 박스
//
// 예시:
//
//	go run main.go vtest.avi 660 255 40 100
package main

import (
	"fmt"
	"image"
	"os"
	"strconv"

	"ex1_tracker_mil/internal/infra"
	"ex1_tracker_mil/internal/usecase"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Println("사용법: go run main.go <비디오경로> <x> <y> <w> <h>")
		os.Exit(1)
	}
	videoPath := os.Args[1]
	x, _ := strconv.Atoi(os.Args[2])
	y, _ := strconv.Atoi(os.Args[3])
	w, _ := strconv.Atoi(os.Args[4])
	h, _ := strconv.Atoi(os.Args[5])
	initRect := image.Rect(x, y, x+w, y+h)

	// --- 구성 루트: 구체 어댑터(infra) 생성 ---
	source, err := infra.NewGoCVVideoSource(videoPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer source.Close()

	tracker := infra.NewGoCVTrackerMIL(source)
	defer tracker.Close()

	sink := infra.NewGoCVWindowSink("GoCV Tracker(MIL) — ESC로 종료", source)
	defer sink.Close()

	// --- 구성 루트: 추상 인터페이스로 유스케이스에 주입(DIP) ---
	trackingUseCase := usecase.NewTrackingUseCase(tracker, source, sink)

	frameCount, err := trackingUseCase.Run(initRect)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("영상이 끝났습니다. 총", frameCount, "프레임 추적")
}
