package util

// FloatPtr32 float32のポインタを返すヘルパー関数
func Float32Ptr(v float32) *float32 {
	return &v
}

func IntPtr(v int) *int {
	return &v
}

// StringPtr stringのポインタを返すヘルパー関数
func StringPtr(v string) *string {
	return &v
}

// HiraganaToKatakana ひらがなをカタカナに変換する関数
func HiraganaToKatakana(input string) string {
	var result string
	for _, r := range input {
		if r >= 'ぁ' && r <= 'ゖ' {
			// ひらがなからカタカナへの変換
			result += string(r - 'ぁ' + 'ァ')
		} else {
			result += string(r)
		}
	}
	return result
}
