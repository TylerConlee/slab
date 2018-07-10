package zendesk

import (
	"reflect"
	"testing"
	"time"
)

func TestGetTimeRemaining(t *testing.T) {
	type args struct {
		ticket ActiveTicket
	}
	tests := []struct {
		name       string
		args       args
		wantRemain time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRemain := GetTimeRemaining(tt.args.ticket); !reflect.DeepEqual(gotRemain, tt.wantRemain) {
				t.Errorf("GetTimeRemaining() = %v, want %v", gotRemain, tt.wantRemain)
			}
		})
	}
}
