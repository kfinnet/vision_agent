package infra

import (
	"fmt"

	"ex2_hog_pedestrian_detection/internal/domain"
	"gocv.io/x/gocv"
)

// GoCVImageIO는 gocv.IMRead/gocv.IMWrite를 사용해 정지 이미지를 읽고 쓰는
// domain.ImageReader / domain.ImageWriter 구현체다. 파일 입출력 책임만
// 가지며(SRP), 탐지나 시각화 로직은 포함하지 않는다.
type GoCVImageIO struct{}

// NewGoCVImageIO는 GoCVImageIO를 생성한다.
func NewGoCVImageIO() *GoCVImageIO {
	return &GoCVImageIO{}
}

// Read는 path의 이미지를 컬러로 읽어 domain.Frame(MatFrame)으로 반환한다.
func (io *GoCVImageIO) Read(path string) (domain.Frame, error) {
	img := gocv.IMRead(path, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("이미지를 읽을 수 없습니다: %s", path)
	}
	return MatFrame{Mat: img}, nil
}

// Write는 frame을 path 경로에 이미지 파일로 저장한다.
func (io *GoCVImageIO) Write(path string, frame domain.Frame) (bool, error) {
	mf, ok := frame.(MatFrame)
	if !ok {
		return false, fmt.Errorf("GoCVImageIO는 MatFrame만 지원합니다")
	}
	return gocv.IMWrite(path, mf.Mat), nil
}
