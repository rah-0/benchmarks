package meta

import (
	"reflect"
	"testing"
)

func BenchmarkMethodAccess(b *testing.B) {
	p := &Person{}

	b.Run("SetName", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p.SetName("John")
		}
	})

	b.Run("GetName", func(b *testing.B) {
		p.SetName("John")
		for i := 0; i < b.N; i++ {
			_ = p.GetName()
		}
	})
}

func BenchmarkReflectionAccess(b *testing.B) {
	p := &Person{}
	val := reflect.ValueOf(p).Elem()
	nameField := val.FieldByName("Name")

	b.Run("SetName", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			nameField.SetString("John")
		}
	})

	b.Run("GetName", func(b *testing.B) {
		nameField.SetString("John")
		for i := 0; i < b.N; i++ {
			_ = nameField.String()
		}
	})
}
