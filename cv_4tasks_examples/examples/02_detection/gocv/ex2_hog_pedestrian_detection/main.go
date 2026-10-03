// CV 4대 태스크 실습 — ②탐지(Detection) 실습2 (GoCV)
// HOG + 선형 SVM 기반 보행자(사람 전신) 탐지
//
// Python/Colab 실습2(ex2_hog_pedestrian_detection.ipynb)와 동일한 OpenCV 내장
// 사전학습 모델(DefaultPeopleDetector)을 순수 Go + GoCV로 구동합니다.
// 얼굴 탐지(실습1, Haar Cascade)와 달리 "사람의 전신 실루엣" 패턴을 학습한 모델입니다.
//
// 이 파일은 클린 아키텍처의 구성 루트(Composition Root)다. CLI 인자를
// 해석해 입력이 이미지인지 비디오인지 판단하고, internal/infra의 구체
// 어댑터(이미지/비디오 I/O, HOG 탐지기, 그리기, 창 렌더러)를 생성해
// internal/usecase에 생성자로 주입(DI)한 뒤, domain의 추상 인터페이스만으로
// 입출력 흐름을 제어한다. "gocv.io/x/gocv"는 이 파일에서 직접 import하지
// 않으며(DIP), 실제 gocv 호출은 모두 infra 패키지 안에서만 일어난다.
//
// 사용법:
//
//	go run main.go <입력이미지 또는 비디오경로> [출력이미지경로]
//	- 이미지 파일이면 정지 이미지 1장을 처리하고 결과를 저장합니다.
//	- 비디오 파일(.avi/.mp4 등)이면 매 프레임 탐지 결과를 창에 표시합니다 (ESC로 종료).
//
// 예시:
//
//	go run main.go street.jpg result.jpg
//	go run main.go vtest.avi
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ex2_hog_pedestrian_detection/internal/infra"
	"ex2_hog_pedestrian_detection/internal/usecase"
)

var videoExts = map[string]bool{".avi": true, ".mp4": true, ".mov": true, ".mkv": true}

// isVideo는 확장자를 보고 비디오 파일 여부를 판단한다. gocv를 쓰지 않는
// 순수 Go 로직이므로 구성 루트(main.go)에 두어도 계층 규칙에 어긋나지 않는다.
func isVideo(path string) bool {
	return videoExts[strings.ToLower(filepath.Ext(path))]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("사용법: go run main.go <입력이미지 또는 비디오경로> [출력이미지경로]")
		os.Exit(1)
	}
	inputPath := os.Args[1]

	// 1) 구성 루트: HOG 디스크립터 + OpenCV 내장 사전학습 보행자 검출 SVM 로드
	detector := infra.NewHOGPedestrianDetector()
	defer detector.Close()
	drawer := infra.NewGoCVFrameDrawer()

	// 2) 구체 어댑터(detector, drawer)를 유스케이스에 생성자 주입(DI)
	pedestrianUseCase := usecase.NewPedestrianDetectionUseCase(detector, drawer)

	if isVideo(inputPath) {
		// 2a) 비디오: 매 프레임 탐지 (실시간 처리 — GUI 환경 필요)
		runVideo(inputPath, pedestrianUseCase)
	} else {
		// 2b) 정지 이미지: 1장 처리 후 결과 저장
		outputPath := "hog_pedestrian_result.jpg"
		if len(os.Args) >= 3 {
			outputPath = os.Args[2]
		}
		runImage(inputPath, outputPath, pedestrianUseCase)
	}
}

// runImage는 정지 이미지 1장을 읽어 탐지 → 시각화 → 저장까지 처리한다.
// domain 추상화(ImageReader/ImageWriter)에만 의존하므로 gocv를 직접
// 호출하지 않는다.
func runImage(inputPath, outputPath string, uc *usecase.PedestrianDetectionUseCase) {
	imageIO := infra.NewGoCVImageIO()

	frame, err := imageIO.Read(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	annotated, detections, err := uc.Execute(frame)
	if err != nil {
		fmt.Println("탐지 오류:", err)
		os.Exit(1)
	}
	fmt.Printf("탐지된 보행자 수: %d\n", len(detections))

	ok, err := imageIO.Write(outputPath, annotated)
	if err != nil || !ok {
		fmt.Println("결과 저장 실패:", outputPath)
		os.Exit(1)
	}
	fmt.Println("결과 저장 완료:", outputPath)
}

// runVideo는 비디오 파일을 매 프레임 읽어 탐지 → 시각화 → 창 출력을
// 반복한다. domain 추상화(VideoSource/WindowRenderer)에만 의존하므로
// gocv를 직접 호출하지 않는다.
func runVideo(inputPath string, uc *usecase.PedestrianDetectionUseCase) {
	video, err := infra.NewFileVideoSource(inputPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer video.Close()

	window := infra.NewGoCVWindowRenderer("GoCV HOG 보행자 탐지 (ESC로 종료)")
	defer window.Close()

	frameCount := 0
	for {
		frame, ok := video.Read()
		if !ok || frame.Empty() {
			fmt.Println("더 이상 읽을 프레임이 없습니다. 종료합니다.")
			break
		}

		annotated, detections, err := uc.Execute(frame)
		if err != nil {
			fmt.Println("탐지 오류:", err)
			break
		}
		fmt.Printf("frame %d: %d명 탐지\n", frameCount, len(detections))

		window.Show(annotated)
		frameCount++
		if window.WaitKey(1) == 27 {
			break
		}
	}
}
