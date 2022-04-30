package snowid

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewNode(t *testing.T) {

	_, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	_, err = NewNode(5000)
	if err == nil {
		t.Fatalf("no error creating NewNode, %s", err)
	}

}

// lazy check if Generate will create duplicate IDs
// would be good to later enhance this with more smarts
func TestGenerateDuplicateID(t *testing.T) {

	node, _ := NewNode(1)

	var x, y ID
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			t.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}

// I feel like there's probably a better way
func TestRace(t *testing.T) {

	node, _ := NewNode(1)

	go func() {
		for i := 0; i < 1000000000; i++ {

			NewNode(1)
		}
	}()

	for i := 0; i < 4000; i++ {

		node.Generate()
	}

}

func TestBase58(t *testing.T) {
	node, err := NewNode(0)
	if err != nil {
		t.Fatalf("error creating NewNode, %s", err)
	}

	for i := 0; i < 10; i++ {

		sf := node.Generate()
		b58 := sf.String()
		psf, err := Parse([]byte(b58))
		if err != nil {
			t.Fatal(err)
		}
		if sf != psf {
			t.Fatal("Parsed does not match String.")
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	id := ID(13587)
	expected := "\"13587\""

	bytes, err := id.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error during MarshalJSON")
	}

	if string(bytes) != expected {
		t.Fatalf("Got %s, expected %s", string(bytes), expected)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		json        string
		expectedID  ID
		expectedErr string
	}{
		{`"13587"`, 13587, ""},
		{`1`, 0, `invalid base58 ID "1"`},
		{`"invalid`, 0, `invalid base58 ID "\"invalid"`},
	}

	for _, tc := range testCases {
		t.Run(tc.json, func(t *testing.T) {
			var id ID
			err := id.UnmarshalJSON([]byte(tc.json))
			if tc.expectedErr == "" {
				assert.NilError(t, err)
				assert.Equal(t, id, tc.expectedID)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr)
		})
	}
}

func BenchmarkParseBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()
	b58 := sf.String()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Parse([]byte(b58))
	}
}
func BenchmarkBase58(b *testing.B) {

	node, _ := NewNode(1)
	sf := node.Generate()

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sf.String()
	}
}
func BenchmarkGenerate(b *testing.B) {

	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {

	NodeBits = 1
	StepBits = 21
	node, _ := NewNode(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	// Generate the ID to unmarshal
	node, _ := NewNode(1)
	id := node.Generate()
	bytes, _ := id.MarshalJSON()

	var id2 ID

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = id2.UnmarshalJSON(bytes)
	}
}

func BenchmarkMarshal(b *testing.B) {
	// Generate the ID to marshal
	node, _ := NewNode(1)
	id := node.Generate()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = id.MarshalJSON()
	}
}

func TestParseBase58(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    ID
		wantErr bool
	}{
		{
			name:    "ok",
			arg:     "4jgmnx8Js8A",
			want:    1428076403798048768,
			wantErr: false,
		},
		{
			name:    "0 not allowed",
			arg:     "0jgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "I not allowed",
			arg:     "Ijgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "O not allowed",
			arg:     "Ojgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
		{
			name:    "l not allowed",
			arg:     "ljgmnx8Js8A",
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
