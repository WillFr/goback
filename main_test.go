package main

import (
	"os"
	"strings"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	cli := ".\\goback.exe -maxFile 10 -thre 0.95 -sto 0.88 -mvs 000 -win 1000 -cut 14 -sta 9 -stam 30 -lim 1.045 -vthre 1.03 -com 0 -rs 1 -bat 30000"
	os.Args = strings.Split(cli, " ")

	for i := 0; i < 1; i++ {
		main()
	}
}
