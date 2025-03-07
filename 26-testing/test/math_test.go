package test

import (
	"example/testing/utils"
	"fmt"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	a, b := 1, 2
	expected := 3
	result := utils.Add(a, b)

	if result != expected {
		t.Errorf("Add(%d,%d) returned %d, expected %d", a, b, result, expected)
	}

	// 定义测试用例
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"正数相加", 1, 2, 3},
		{"负数相加", -1, -2, -3},
		{"零值相加", 0, 0, 0},
		{"正负数相加", 1, -2, -1},
	}

	// 运行所有测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d,%d) = %d; 期望值 %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}

}

func TestSubtract(t *testing.T) {
	t.Helper()
	a, b := 1, 2
	expected := -1
	result := utils.Subtract(a, b)

	if result != expected {
		t.Errorf("Subtract(%d,%d) returned %d, expected %d", a, b, result, expected)
	}
}

func setup() {
	fmt.Println("-------setup-------")
}

func teardown() {
	fmt.Println("-------teardown-------")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
