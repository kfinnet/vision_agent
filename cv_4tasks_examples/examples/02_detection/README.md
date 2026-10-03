# ② 탐지 (Detection) 실습

사진 속 여러 물체 각각의 "무엇"과 "어디(위치)"를 함께 찾아내는 태스크입니다.

## Colab/OpenCV (Python)

### 실습1 — `ex1_haar_face_detection.ipynb`
OpenCV에 내장된 사전학습 Haar Cascade 분류기로 얼굴을 탐지합니다.
(검증 결과: 샘플 이미지에서 얼굴 정상 탐지)

### 실습2 — `ex2_hog_pedestrian_detection.ipynb`
OpenCV 내장 HOG+SVM 보행자 검출기로 CCTV 영상 속 사람을 탐지합니다.
얼굴(부분) 탐지와 달리 "사람 전신 실루엣"을 찾는 모델입니다.
(검증 결과: vtest.avi 프레임에서 보행자 4명 탐지, NMS로 중복 제거)

## GoCV (Go)

### 실습1 — `gocv/ex1_haar_face_detection/`
웹캠 또는 비디오 파일에서 매 프레임 얼굴을 탐지하는 실시간 버전입니다.

```bash
cd gocv/ex1_haar_face_detection
go get gocv.io/x/gocv
curl -o haarcascade_frontalface_default.xml \
  https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml
go run main.go haarcascade_frontalface_default.xml ../../../assets/vtest.avi
```

### 실습2 — `gocv/ex2_hog_pedestrian_detection/`
Python 실습2와 동일한 HOG 보행자 검출기를 Go로 구동합니다. 정지 이미지와
비디오 파일을 모두 지원합니다.

```bash
cd gocv/ex2_hog_pedestrian_detection
go get gocv.io/x/gocv
go run main.go ../../../assets/vtest.avi   # 비디오: 실시간 표시(ESC로 종료)
```

> GoCV 2개는 `gofmt`/스텁 기반 `go build`·`go vet`까지 확인했으나, 실제 `gocv`
> 라이브러리로 완전한 컴파일·실행 검증은 하지 못했습니다 — 상위 README의
> "검증 현황"을 참고하세요. 비디오 창을 띄우는 예제이므로 GUI가 있는 로컬
> PC(Windows/macOS/Linux 데스크톱)에서 실행해야 합니다.

## 아키텍처 — 클린 아키텍처 + SOLID 적용

Python/Go 4개 예제 모두 "탐지"라는 핵심 로직을 구체 구현(cv2, gocv)에서
분리하는 클린 아키텍처 계층 구조로 재구성했습니다.

### Python(Colab) — 한 파일 내 4계층

각 `_code.py` 파일은 하나의 파일이지만, 주석 배너로 구분된 4개 계층으로
위에서 아래로 구성됩니다.

1. `[도메인 계층]` — `DetectionBox`(값 객체), `ObjectDetector`(추상 인터페이스).
   cv2 의존성이 전혀 없는 순수 타입만 정의합니다.
2. `[유스케이스 계층]` — `DetectionUseCase`가 생성자로 `ObjectDetector`를
   주입받아 탐지(및 ex2의 경우 NMS) 흐름만 조율합니다. cv2를 직접 호출하지 않습니다.
3. `[인프라/어댑터 계층]` — `HaarCascadeFaceDetector`, `HOGPedestrianDetector`,
   `NMSBoxSuppressor` 등 cv2를 실제로 호출하는 구체 클래스들입니다.
4. `[구성 루트]` — 어댑터를 생성해 유스케이스에 주입(DI)하고, 샘플 이미지/영상에
   대해 실행·시각화·저장하는 실제 Colab 실행 코드입니다.

### Go(GoCV) — 패키지 단위 4계층

각 예제 폴더를 `internal/domain`(값 타입 + 포트 인터페이스),
`internal/usecase`(생성자 주입 기반 오케스트레이션),
`internal/infra`(gocv를 호출하는 유일한 패키지), `main.go`(구성 루트) 로
분리했습니다. `gocv.io/x/gocv`는 `internal/infra` 패키지에서만 import합니다.

### SOLID 적용 방식

- **SRP**: 탐지, 후처리(NMS), 시각화(그리기/창 출력)를 각각 별도 클래스·구조체로 분리했습니다.
- **OCP**: Haar Cascade→HOG(또는 향후 DNN) 교체는 `ObjectDetector`를 구현하는
  새 인프라 클래스 추가만으로 가능하며, 유스케이스 코드는 수정하지 않습니다.
- **LSP**: 모든 구체 탐지기는 동일한 `DetectionBox`/`Detection` 반환 형태를
  지키며, 인터페이스가 허용하는 입력에 대해 예외를 던지지 않습니다.
- **ISP**: 탐지(`ObjectDetector`)와 시각화·후처리(`BoxSuppressor`, `FrameDrawer` 등)
  인터페이스를 하나로 합치지 않고 작게 분리했습니다.
- **DIP**: 유스케이스는 `ObjectDetector` 등 추상화에만 의존하며, cv2/gocv
  구체 타입을 직접 참조하지 않습니다.
