package main

import "testing"

func BenchmarkEffectiveLogger(b *testing.B) {

	eLogger := EffectiveLogger{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eLogger.Info("Effective!")
	}
}

func BenchmarkIneffectiveLogger2(b *testing.B) {

	iLogger := &IneffectiveLogger{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		iLogger.Info("Ineffective!")
	}
}
