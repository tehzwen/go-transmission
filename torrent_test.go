package go_transmission

import "testing"

func Test_cleanMagnet(t *testing.T) {
	type args struct {
		original string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "parses magnet properly",
			args: args{
				original: "magnet:?xt=urn:btih:test-torrent-magnet-name&tr=faketracker1&tr=udpfaketracker2&tr=udpfaketracker3",
			},
			want: "magnet:?xt=urn:btih:test-torrent-magnet-name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanMagnet(tt.args.original); got != tt.want {
				t.Errorf("cleanMagnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
