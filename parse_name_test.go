package main

import (
	"testing"
)

func TestGetProgramName(t *testing.T) {
	for _, basename := range testcase {
		GetProgramName(basename)
	}
}

var testcase []string = []string{
	"トリリオンゲーム　♯４「ワルい男」_20241018.ts.program.txt",
	"株式会社マジルミエ　ＦＲＩＤＡＹ　ＡＮＩＭＥ　ＮＩＧＨＴ_20241018.ts",
	"２．５次元の誘惑（リリサ）　＃１４「あなたと一緒に」_20241004.ts.program.txt",
	"アイドルマスター　シャイニーカラーズ　2nd　第2話「Straylight.run()／／playback」_20241012.ts",
}
