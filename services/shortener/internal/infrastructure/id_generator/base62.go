package id_generator

// Base62Encode
// Alphabet including 0-9, A-Z, a-z.
func Base62Encode(num uint64) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	if num == 0 {
		return string(alphabet[0])
	}
	res := make([]byte, 0, 10)
	for num > 0 {
		res = append(res, alphabet[num%62])
		num /= 62
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return string(res)
}
