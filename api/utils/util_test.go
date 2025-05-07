package utils

import (
	"fmt"
	"testing"
)

func TestGenerateSucess(t *testing.T) {
	size := 10

	result := StringWithCharset(size, charset)

	if len(result) != size {
		t.Errorf("generate result length %d != %d", len(result), size)
	}
	fmt.Println(result)
}

func TestGenerateWrongSize(t *testing.T) {
	expectedSize := 10
	actualSize := 9

	result := StringWithCharset(actualSize)

	if len(result) != expectedSize {
		// esse teste "falha de propósito" para mostrar erro
		t.Errorf("Generated result is (%d) it's different to (%d)", len(result), expectedSize)
	} else {
		t.Logf("expected flaws, but the sizes matched")
	}
}
