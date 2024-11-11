package main

import (
	"testing"
)

func TestGetProgramName(t *testing.T) {
	for idx, basename := range testcase {
		name, _ := GetProgramName(basename)
		if correctcase[idx] != name {
			t.Errorf("excepted name: <%v>, but <%v>", correctcase[idx], name)
		}
	}

}

var testcase []string = []string{
	"トリリオンゲーム　♯４「ワルい男」_20241018.ts.program.txt",
	"株式会社マジルミエ　ＦＲＩＤＡＹ　ＡＮＩＭＥ　ＮＩＧＨＴ_20241018.ts",
	"２．５次元の誘惑（リリサ）　＃１４「あなたと一緒に」_20241004.ts.program.txt",
	"アイドルマスター　シャイニーカラーズ　2nd　第2話「Straylight.run()／／playback」_20241012.ts",
	"魔王様、リトライ！Ｒ　＃５「リマインド －ＲｅＭｉｎｄ－」_20241102.ts.program.txt",
	"アニメ ラブライブ!スーパースター!! 3期(6)_20241110.ts.program",
}

var correctcase []string = []string{
	"トリリオンゲーム",
	"株式会社マジルミエ FRIDAY ANIME NIGHT",
	"2.5次元の誘惑(リリサ)",
	"アイドルマスター シャイニーカラーズ 2nd",
	"魔王様、リトライ!R",
	"アニメ ラブライブ!スーパースター!! 3期",
}
