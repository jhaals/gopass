package random

import "math/rand"

func RandomString(l int) string {
    bytes := make([]byte, l)
    for i:=0 ; i<l ; i++ {
      bytes[i] = byte(randInt(65,90))
    }
    return string(bytes)
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}
