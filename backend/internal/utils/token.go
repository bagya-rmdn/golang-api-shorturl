package utils
import (
"crypto/sha256"
"math/big"
)
var base62 = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
func base62Encode(b []byte) string {
// Convert bytes to big.Int then encode base62
n := new(big.Int).SetBytes(b)
if n.Cmp(big.NewInt(0)) == 0 {
return "0"
}
res := make([]byte, 0, 44)
r := new(big.Int)
q := new(big.Int)
base := big.NewInt(62)
for n.Cmp(big.NewInt(0)) > 0 {
q.QuoRem(n, base, r)
res = append(res, base62[r.Int64()])
n = new(big.Int).Set(q)
}
// reverse
for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
res[i], res[j] = res[j], res[i]
}
return string(res)
}
// Deterministic token from normalized URL; first 8 chars of base62(sha256)
func TokenFromURL(normalized string) string {
sum := sha256.Sum256([]byte(normalized))
enc := base62Encode(sum[:])
if len(enc) < 8 {
return enc
}
return enc[:8]
}