# ③ 분할 (Segmentation) 실습

사진을 픽셀 하나하나까지 "이 픽셀은 어떤 물체에 속하는가"로 구분해내는,
가장 정밀한 CV 태스크입니다.

## Colab/OpenCV (Python)

### 실습1 — `ex1_watershed_segmentation.ipynb`
서로 맞닿은 동전 사진에 Watershed 알고리즘을 적용해, 탐지(사각 박스)로는
셀 수 없는 맞닿은 물체들을 픽셀 단위로 분리합니다.
(검증 결과: 24개 동전 중 24~26개로 정상 분할 — 결과물 포함)

### 실습2 — `ex2_grabcut_segmentation.ipynb`
사각형 힌트 하나로 GrabCut 알고리즘이 반복적인 그래프 컷 최적화를 통해
전경(고양이)을 픽셀 단위로 추출합니다. 배경 제거·인물사진 모드(배경 블러)까지
함께 실습합니다.
(검증 결과: 전경 34% 추출, 윤곽 육안 확인)

## GoCV (Go)

### 실습1 — `gocv/ex1_watershed_segmentation/`
Python 실습1과 동일한 Watershed 파이프라인을 Go로 구현합니다.

```bash
cd gocv/ex1_watershed_segmentation
go get gocv.io/x/gocv
go run main.go ../../../assets/coins.png watershed_result.png
```

### 실습2 — `gocv/ex2_grabcut_segmentation/`
Python 실습2와 동일한 GrabCut 전경 추출을 Go로 구현합니다.

```bash
cd gocv/ex2_grabcut_segmentation
go get gocv.io/x/gocv
go run main.go ../../../assets/cat.jpg cat_foreground.png
```

> GoCV 2개는 `gofmt`/스텁 기반 `go build`·`go vet`까지 확인했으나, 실제 `gocv`
> 라이브러리로 완전한 컴파일·실행 검증은 하지 못했습니다 — 상위 README의
> "검증 현황"을 참고하세요. Watershed/GrabCut은 GoCV API 중에서도 비교적 복잡한
> 축에 속하므로, 실행 전 본인 환경에서의 `go build` 확인을 특히 권장합니다.

## 구조 — Clean Architecture + SOLID

모든 예제는 단순히 동작만 하는 코드가 아니라, 계층을 분리하고 의존성이
추상화를 향하도록(DIP) 재구성되어 있습니다.

### Python (`colab_opencv/*.py`)

하나의 파일 안에서 주석 배너로 4개 계층을 위에서 아래로 명확히 구분합니다.

1. `[도메인 계층 / Domain Layer]` — `@dataclass` 결과 값 객체(`SegmentationResult`)와
   `abc.ABC` 추상 인터페이스(`ImageSegmenter`). cv2/scipy에 의존하지 않습니다.
2. `[유스케이스 계층 / Use Case Layer]` — `SegmentationUseCase`. 생성자 주입으로
   받은 `ImageSegmenter` 추상화만 사용하며 cv2/scipy를 직접 호출하지 않습니다.
3. `[인프라/어댑터 계층 / Infrastructure Layer]` — `WatershedSegmenter`,
   `GrabCutSegmenter` 등 cv2/scipy를 실제로 호출하는 구체 구현체. 분할 계산과
   시각화(색칠, 배경 제거 등)를 각각 별도 클래스로 분리했습니다.
4. `[구성 루트 / Composition Root]` — Colab 셀에서 구체 어댑터를 생성해
   유스케이스에 주입하고, 실행·시각화·저장까지 수행합니다.

### Go (`gocv/<example>/`)

```
<example>/
  go.mod
  main.go              — 구성 루트: CLI 인자 파싱 + 의존성 조립(DI)만 수행
  internal/
    domain/
      types.go          — 값 타입 (SegmentationResult 등)
      ports.go          — 인터페이스 (ImageSegmenter, ResultVisualizer/Presenter)
    usecase/
      *_usecase.go      — 생성자 주입 + Execute()만 포함, gocv 미사용
    infra/
      gocv_*.go         — gocv를 호출하는 유일한 패키지, 계산과 시각화를 구조체로 분리
```

### SOLID 적용 방식

- **SRP**: 분할 계산 / 결과 시각화(색칠·배경 제거) / 파일 입출력·플로팅을
  각각 별도 클래스·구조체로 분리했습니다.
- **OCP**: Watershed를 GrabCut이나 다른 분할 알고리즘으로 바꾸려면 `ImageSegmenter`를
  구현하는 새 클래스만 추가하면 되며, 유스케이스는 수정하지 않습니다.
- **LSP**: `WatershedSegmenter`, `GrabCutSegmenter` 모두 `ImageSegmenter`의
  완전한 대체물로 동작합니다.
- **ISP**: 분할(`ImageSegmenter`)과 시각화(`ResultVisualizer`/`SegmentationVisualizer`)
  인터페이스를 하나로 합치지 않고 작게 유지했습니다.
- **DIP**: 유스케이스 계층은 cv2/gocv 구체 타입이 아니라 도메인 인터페이스에만
  의존합니다.
