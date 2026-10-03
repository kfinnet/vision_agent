package infra

import (
	"bufio"
	"fmt"
	"os"

	"ex2_dnn_classification/internal/domain"
)

// TextFileLabelRepository는 domain.LabelRepository를 줄바꿈으로 구분된 텍스트 파일로부터
// 구현하는 어댑터다 (예: ImageNet 1000개 클래스 레이블 파일).
//
// LSP: Label()의 시그니처와 동작이 domain.LabelRepository 계약을 그대로 만족하므로,
// LabelRepository를 기대하는 어떤 코드에도 안전하게 대입할 수 있다 (범위를 벗어난
// 인덱스에도 panic/에러 없이 "class_<idx>"를 반환해, 인터페이스가 암묵적으로 기대하는
// "항상 문자열을 반환한다"는 사전조건을 그대로 지킨다).
type TextFileLabelRepository struct {
	labels []string
}

// NewTextFileLabelRepository는 레이블 텍스트 파일을 읽어 TextFileLabelRepository를 생성한다.
func NewTextFileLabelRepository(path string) (*TextFileLabelRepository, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var labels []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		labels = append(labels, sc.Text())
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return &TextFileLabelRepository{labels: labels}, nil
}

// Label은 인덱스에 해당하는 레이블 이름을 반환한다. 범위를 벗어나면 "class_<idx>" 형태로 반환한다.
func (r *TextFileLabelRepository) Label(index int) string {
	if index >= 0 && index < len(r.labels) {
		return r.labels[index]
	}
	return fmt.Sprintf("class_%d", index)
}

// 컴파일 타임에 TextFileLabelRepository가 domain.LabelRepository 포트를 만족하는지 확인한다.
var _ domain.LabelRepository = (*TextFileLabelRepository)(nil)
