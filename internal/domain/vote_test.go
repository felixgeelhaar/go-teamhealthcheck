package domain

import "testing"

func TestParseVoteColor(t *testing.T) {
	tests := []struct {
		input   string
		want    VoteColor
		wantErr bool
	}{
		{"green", VoteGreen, false},
		{"yellow", VoteYellow, false},
		{"red", VoteRed, false},
		{"blue", "", true},
		{"", "", true},
		{"GREEN", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseVoteColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVoteColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseVoteColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestVoteColorScore(t *testing.T) {
	tests := []struct {
		color VoteColor
		want  float64
	}{
		{VoteGreen, 3},
		{VoteYellow, 2},
		{VoteRed, 1},
		{VoteColor("invalid"), 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.color), func(t *testing.T) {
			if got := tt.color.Score(); got != tt.want {
				t.Errorf("%s.Score() = %v, want %v", tt.color, got, tt.want)
			}
		})
	}
}
