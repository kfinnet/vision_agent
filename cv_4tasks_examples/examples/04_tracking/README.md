# ④ 추적 (Tracking) 실습

영상의 여러 프레임에 걸쳐 "같은 물체"가 시간에 따라 어떻게 이동하는지
계속 따라가는 태스크입니다. 탐지(공간) + 매칭(시간)의 결합입니다.

## Colab/OpenCV (Python)

### 실습1 — `ex1_meanshift_tracking.ipynb`
첫 프레임에서 지정한 보행자의 색 분포(히스토그램)를 기준으로, MeanShift가
매 프레임 "색이 가장 비슷한 위치"로 추적 윈도우를 이동시킵니다.
(검증 결과: 40프레임에 걸친 이동 경로 정상 생성 — 결과물 포함.
단, MeanShift는 대상이 가려지거나 프레임을 벗어나면 추적을 놓칠 수 있다는
한계도 함께 확인했습니다. 실무에서는 CSRT/KCF/딥러닝 기반 트래커로 보완합니다.)

### 실습2 — `ex2_optical_flow_tracking.ipynb`
MeanShift와 달리, 영상 전체에서 추적하기 좋은 특징점(코너) 여러 개를 뽑아
Lucas-Kanade Optical Flow로 각각의 이동을 동시에 추적합니다.
(검증 결과: 다중 특징점의 이동 궤적 정상 생성 — 결과물 포함)

## GoCV (Go)

### 실습1 — `gocv/ex1_tracker_mil/`
OpenCV가 제공하는 범용 Tracker API(TrackerMIL)로 박스 하나를 영상 전체에서
계속 따라갑니다. 매 프레임 재탐지 없이 "이전 외형과 가장 비슷한 위치"를 찾는
방식으로, 추적 API를 가장 단순하게 보여줍니다.

```bash
cd gocv/ex1_tracker_mil
go get gocv.io/x/gocv
go run main.go ../../../assets/vtest.avi 660 255 40 100
```

### 실습2 — `gocv/ex2_bgsub_centroid_tracking/`
강의안의 "탐지(공간) + 매칭(시간) = 추적" 공식을 가장 직접적으로 구현한 예제입니다.
배경 차분(MOG2)으로 움직이는 물체를 매 프레임 탐지하고, 이전 프레임 추적 목록과
중심점 거리로 매칭해 같은 ID를 유지합니다 — SORT류 추적기의 핵심 아이디어를
간단한 형태로 보여줍니다.

```bash
cd gocv/ex2_bgsub_centroid_tracking
go get gocv.io/x/gocv
go run main.go ../../../assets/vtest.avi
```

> GoCV 2개는 `gofmt`/스텁 기반 `go build`·`go vet`까지 확인했으나, 실제 `gocv`
> 라이브러리로 완전한 컴파일·실행 검증은 하지 못했습니다 — 상위 README의
> "검증 현황"을 참고하세요. 비디오 창을 띄우는 예제이므로 GUI가 있는 로컬
> PC(Windows/macOS/Linux 데스크톱)에서 실행해야 합니다.

## 아키텍처 — Clean Architecture + SOLID 적용

이 폴더의 모든 예제는 "같은 핵심 로직을 다른 알고리즘/언어로 바꿔 끼울 수
있어야 한다"는 추적(Tracking) 태스크의 특성에 맞게, 계층을 분리하고
SOLID 원칙을 실제 코드 구조에 적용했습니다.

### Python(Colab) — 파일 내부를 4개 계층으로 구분

하나의 `.py` 파일(노트북 셀 마커 `# Colab 셀 N`은 그대로 유지) 안에서
주석 배너로 아래 4개 계층을 top-to-bottom으로 구분합니다.

1. `[도메인 계층]` — `@dataclass` 값 객체(`TrackedRegion`/`TrackedPoint`,
   `TrackState`)와 `abc.ABC` 추상 인터페이스(`ObjectTracker`). cv2를
   전혀 쓰지 않습니다.
2. `[유스케이스 계층]` — `TrackingUseCase`/`OpticalFlowTrackingUseCase`가
   생성자로 `ObjectTracker`를 주입받아 프레임 반복·결과 누적만 담당합니다.
   cv2를 import하지 않습니다.
3. `[인프라/어댑터 계층]` — `MeanShiftTracker`, `LucasKanadeOpticalFlowTracker`가
   `ObjectTracker`를 구현하며, cv2를 호출하는 유일한 계층입니다.
4. `[구성 루트]` — 실제 Colab 실행 코드: 영상을 열고, 구체 어댑터를 만들어
   유스케이스에 주입하고, matplotlib로 시각화·저장합니다.

### Go(GoCV) — 폴더 구조로 계층 분리

```
<example>/
  main.go              구성 루트(의존성 조립 + CLI, 디스플레이 루프)
  internal/domain/      값 타입 + 인터페이스(포트) — gocv 비의존
  internal/usecase/      오케스트레이션 — gocv 비의존
  internal/infra/        gocv 기반 구체 구현 — gocv를 import하는 유일한 패키지
```

실습2(`ex2_bgsub_centroid_tracking`)는 기존에 하나로 섞여 있던 "배경 차분
탐지", "중심점 거리 매칭/ID 수명 관리", "그리기"를 각각 `infra.GoCVMotionDetector`,
`domain.CentroidTracker`, `infra.GoCVWindowPresenter`로 분리했습니다. 특히
`CentroidTracker`는 점(`image.Point`)과 거리 계산만 다루는 순수 알고리즘이라
gocv를 전혀 import하지 않으며, 그대로 `domain` 패키지에 두어 "인프라에
의존하지 않는 순수 규칙"이라는 Clean Architecture 정의를 그대로 보여줍니다.

### SOLID 적용 방식

- **SRP**: 추적 알고리즘, 궤적/매칭 누적, 시각화를 각각 별도 클래스/타입으로
  분리했습니다(예: MeanShift 추적 ≠ trajectory 누적 ≠ matplotlib 그리기).
- **OCP**: 새 추적 알고리즘(CSRT/KCF/딥러닝 트래커 등)을 추가하려면 해당
  추상 인터페이스를 구현하는 어댑터만 추가하면 되고, 유스케이스는 수정할
  필요가 없습니다.
- **LSP / ISP 설계 판단**: MeanShift는 "영역 하나"를, Optical Flow는
  "다중 점"을 추적한다는 본질적 차이가 있습니다. 이를 하나의 인터페이스에
  억지로 맞추면(예: 항상 점 1개만 반환하도록 강제) 한쪽이 거짓 구현이 되어
  LSP를 위반하게 됩니다. 그래서 각 예제는 "자신에게 자연스러운 자료형의
  리스트를 반환한다"는 더 일반적인 계약을 따르는, 필요한 만큼만 작은
  `ObjectTracker`/`Tracker` 인터페이스를 각각 정의했습니다(ISP) — 불필요하게
  큰 공용 인터페이스를 강제하지 않으면서도 각 구현이 정직하게 계약을
  지키도록 했습니다.
- **DIP**: 유스케이스/오케스트레이션 계층은 추상 인터페이스에만 의존하고,
  cv2/gocv 호출은 인프라 계층에만 존재합니다. 구체 구현 생성과 주입은
  구성 루트(Colab 셀 또는 Go `main.go`)에서만 일어납니다.
