package zendesk

import (
	"reflect"
	"testing"
)

func Test_getPriorityLevel(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantPriLvl string
	}{
		{
			name:       "Check for Level One",
			args:       []string{"platinum"},
			wantPriLvl: "LevelOne",
		},
		{
			name:       "Check for Level Two",
			args:       []string{"gold"},
			wantPriLvl: "LevelTwo",
		},
		{
			name:       "Check for Level Three",
			args:       []string{"silver"},
			wantPriLvl: "LevelThree",
		},
		{
			name:       "Check for Level Four",
			args:       []string{"bronze"},
			wantPriLvl: "LevelFour",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPriLvl := getPriorityLevel(tt.args); gotPriLvl != tt.wantPriLvl {
				t.Errorf("getPriorityLevel() = %v, want %v", gotPriLvl, tt.wantPriLvl)
			}
		})
	}
}

func TestCheckSLA(t *testing.T) {
	type args struct {
		tick ZenOutput
	}
	tests := []struct {
		name    string
		args    args
		wantSla []ActiveTicket
	}{
		{
			name: "Check for No SLA - Single Ticket",
			args: args{
				tick: ZenOutput{
					Tickets: Tickets{
						{
							ID:          123,
							Tags:        []string{"bronze"},
							RequesterID: 12345,
							Subject:     "Test Subject",
							Description: "Test ticket",
						},
					},
				},
			},
			wantSla: []ActiveTicket{
				{
					ID:          123,
					Requester:   12345,
					Level:       "LevelFour",
					Tags:        []string{"bronze"},
					Subject:     "Test Subject",
					Description: "Test ticket",
				},
			},
		},
		{
			name: "Check for No SLA - Multiple Ticket",
			args: args{
				tick: ZenOutput{
					Tickets: Tickets{
						{
							ID:          123,
							Tags:        []string{"bronze"},
							RequesterID: 12345,
							Subject:     "Test Subject",
							Description: "Test ticket",
						},
						{
							ID:          234,
							Tags:        []string{"gold"},
							RequesterID: 23456,
							Subject:     "Second Subject",
							Description: "Second ticket",
						},
					},
				},
			},
			wantSla: []ActiveTicket{
				{
					ID:          123,
					Requester:   12345,
					Level:       "LevelFour",
					Tags:        []string{"bronze"},
					Subject:     "Test Subject",
					Description: "Test ticket",
				},
				{
					ID:          234,
					Requester:   23456,
					Level:       "LevelTwo",
					Tags:        []string{"gold"},
					Subject:     "Second Subject",
					Description: "Second ticket",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSla := CheckSLA(tt.args.tick); !reflect.DeepEqual(gotSla, tt.wantSla) {
				t.Errorf("CheckSLA() = %v, want %v", gotSla, tt.wantSla)
			}
		})
	}
}
