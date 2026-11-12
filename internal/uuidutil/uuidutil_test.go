package uuidutil

import (
	"slices"
	"testing"

	"github.com/google/uuid"
)

func TestDecodeNilReturnsNoFields(t *testing.T) {
	t.Parallel()

	if got := Decode(uuid.Nil); got != nil {
		t.Fatalf("Decode(uuid.Nil) = %v, want nil", got)
	}
}

func TestDecodeVersion4(t *testing.T) {
	t.Parallel()

	got := Decode(uuid.MustParse("4326ff5f-774d-4506-a18c-4bc50c761863"))
	want := []Field{
		{Name: "Standard String Format", Value: "4326ff5f-774d-4506-a18c-4bc50c761863"},
		{Name: "Single Integer Value", Value: "89260762576260186387968344354060638307"},
		{Name: "Version", Value: "4"},
		{Name: "Variant", Value: "DCE 1.1, ISO/IEC 11578:1996"},
		{Name: "Contents", Value: "43:26:FF:5F:77:4D:45:06:A1:8C:4B:C5:0C:76:18:63"},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Decode() = %v, want %v", got, want)
	}
}

func TestDecodeVersion1(t *testing.T) {
	t.Parallel()

	got := Decode(uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8"))
	want := []Field{
		{Name: "Standard String Format", Value: "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		{Name: "Single Integer Value", Value: "143098242404177361603877621312831893704"},
		{Name: "Version", Value: "1"},
		{Name: "Variant", Value: "DCE 1.1, ISO/IEC 11578:1996"},
		{Name: "Contents - Time", Value: "1998-02-04 22:13:53.1511824 UTC"},
		{Name: "Contents - Clock", Value: "180"},
		{Name: "Contents - Node", Value: "00:c0:4f:d4:30:c8"},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Decode() = %v, want %v", got, want)
	}
}

func TestFieldsToRows(t *testing.T) {
	t.Parallel()

	got := FieldsToRows([]Field{{Name: "Version", Value: "4"}, {Name: "Variant", Value: "Reserved (future use)"}})
	want := [][]string{{"Version", "4"}, {"Variant", "Reserved (future use)"}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("FieldsToRows() = %v, want %v", got, want)
	}
}

func TestGenerateEveryVersion(t *testing.T) {
	t.Parallel()

	for version := 1; version <= 7; version++ {
		id, err := Generate(version, "https://example.com")
		if err != nil {
			t.Fatalf("Generate(%d) returned error: %v", version, err)
		}
		if got := int(id.Version()); got != version {
			t.Errorf("Generate(%d) version = %d, want %d", version, got, version)
		}
	}
}

func TestGenerateNameBasedVersionsAreStable(t *testing.T) {
	t.Parallel()

	for _, version := range []int{3, 5} {
		first, _ := Generate(version, "https://example.com")
		second, _ := Generate(version, "https://example.com")
		other, _ := Generate(version, "https://example.org")
		if first != second {
			t.Errorf("Generate(%d) = %s and %s for the same namespace", version, first, second)
		}
		if first == other {
			t.Errorf("Generate(%d) = %s for different namespaces", version, first)
		}
	}
}

func TestGenerateUnsupportedVersionReturnsError(t *testing.T) {
	t.Parallel()

	for _, version := range []int{0, 8} {
		id, err := Generate(version, "")
		if err == nil {
			t.Fatalf("Generate(%d) expected error", version)
		}
		if id != uuid.Nil {
			t.Fatalf("Generate(%d) = %s, want nil uuid", version, id)
		}
	}
}
