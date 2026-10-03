// CV 실습 2 (GoCV) — 실시간 얼굴 탐지 (Haar Cascade Detection)
//
// Python/OpenCV 실습 2(02_face_detection.ipynb)가 정지 이미지 1장을 처리했다면,
// 이 예제는 웹캠 또는 비디오 파일에서 "실시간으로" 매 프레임마다 얼굴을 탐지합니다.
// GoCV가 강점을 갖는 실시간 비디오 스트림 처리를 보여주는 예제입니다.
//
// 이 파일은 클린 아키텍처의 구성 루트(Composition Root)다. CLI 인자를
// 해석하고, internal/infra의 구체 어댑터(비디오 소스, Haar Cascade 탐지기,
// 창 렌더러)를 생성해 internal/usecase에 생성자로 주입(DI)한 뒤, domain의
// 추상 인터페이스만으로 캡처→탐지→렌더링 루프를 제어한다. "gocv.io/x/gocv"는
// 이 파일에서 직접 import하지 않으며(DIP), 실제 gocv 호출은 모두 infra
// 패키지 안에서만 일어난다.
//
// 사용법:
//
//	go run main.go <cascade.xml 경로> [비디오소스]
//	- 비디오소스가 숫자면 웹캠 장치 번호(기본 0), 파일 경로면 해당 비디오 파일을 사용합니다.
//	- 비디오소스를 생략하면 웹캠(장치 0)을 사용합니다.
//
// 예시:
//
//	go run main.go haarcascade_frontalface_default.xml           (웹캠)
//	go run main.go haarcascade_frontalface_default.xml sample.mp4 (비디오 파일)
//
// haarcascade_frontalface_default.xml 다운로드:
//
//	https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml
package main

import (
	"fmt"
	"os"

	"ex1_haar_face_detection/internal/domain"
	"ex1_haar_face_detection/internal/infra"
	"ex1_haar_face_detection/internal/usecase"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("사용법: go run main.go <cascade.xml 경로> [비디오소스(웹캠번호 또는 파일경로)]")
		os.Exit(1)
	}
	cascadePath := os.Args[1]

	// 비디오 소스 결정: 생략 시 웹캠(0), 인자가 있으면 파일/장치
	videoSource := "0"
	if len(os.Args) >= 3 {
		videoSource = os.Args[2]
	}

	// 1) 구성 루트: 구체 어댑터(비디오 소스) 생성 — 웹캠 번호 또는 비디오 파일 경로
	var source domain.VideoSource
	var err error
	if videoSource == "0" || videoSource == "1" || videoSource == "2" {
		deviceID := 0
		fmt.Sscanf(videoSource, "%d", &deviceID)
		source, err = infra.NewWebcamVideoSource(deviceID)
	} else {
		source, err = infra.NewFileVideoSource(videoSource)
	}
	if err != nil {
		fmt.Println("비디오 소스를 열 수 없습니다:", err)
		os.Exit(1)
	}
	defer source.Close()

	// 2) 구성 루트: Haar Cascade 탐지기 어댑터 생성
	detector, err := infra.NewHaarCascadeFaceDetector(cascadePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer detector.Close()

	// 3) 구성 루트: 디스플레이 창 렌더러 생성 (GUI 환경 필요; 헤드리스 서버에서는 실행되지 않음)
	renderer := infra.NewGoCVWindowRenderer("GoCV 실시간 얼굴 탐지 (ESC로 종료)")
	defer renderer.Close()

	// 4) 구체 어댑터(detector)를 유스케이스에 생성자 주입(DI)
	faceUseCase := usecase.NewFaceDetectionUseCase(detector)

	frameCount := 0
	fmt.Println("실시간 얼굴 탐지를 시작합니다. 창에서 ESC 키를 누르면 종료됩니다.")

	// 캡처 → 탐지 → 렌더링 루프. domain 추상화(VideoSource/ObjectDetector/
	// FrameRenderer)에만 의존하므로 main.go는 gocv를 직접 호출하지 않는다.
	for {
		frame, ok := source.Read()
		if !ok || frame.Empty() {
			fmt.Println("더 이상 읽을 프레임이 없습니다. 종료합니다.")
			break
		}

		// 현재 프레임에서 얼굴 탐지 (매 프레임 반복 -> 실시간 처리)
		faces, err := faceUseCase.Execute(frame)
		if err != nil {
			fmt.Println("탐지 오류:", err)
			break
		}

		renderer.Render(frame, faces, frameCount)
		frameCount++

		// 5) ESC(27) 키 입력 시 종료
		if renderer.WaitKey(1) == 27 {
			break
		}
	}

	fmt.Printf("종료: 총 %d 프레임 처리\n", frameCount)
}
