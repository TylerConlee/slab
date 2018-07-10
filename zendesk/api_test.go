package zendesk

import (
	"reflect"
	"testing"
)

func TestGetAllTickets(t *testing.T) {
	tests := []struct {
		name        string
		wantTickets ZenOutput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTickets := GetAllTickets(); !reflect.DeepEqual(gotTickets, tt.wantTickets) {
				t.Errorf("GetAllTickets() = %v, want %v", gotTickets, tt.wantTickets)
			}
		})
	}
}

func TestGetTicketEvents(t *testing.T) {
	tests := []struct {
		name        string
		wantTickets EventOutput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTickets := GetTicketEvents(); !reflect.DeepEqual(gotTickets, tt.wantTickets) {
				t.Errorf("GetTicketEvents() = %v, want %v", gotTickets, tt.wantTickets)
			}
		})
	}
}

func TestGetTicket(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name       string
		args       args
		wantTicket Ticket
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTicket := GetTicket(tt.args.id); !reflect.DeepEqual(gotTicket, tt.wantTicket) {
				t.Errorf("GetTicket() = %v, want %v", gotTicket, tt.wantTicket)
			}
		})
	}
}

func TestGetTicketRequester(t *testing.T) {
	type args struct {
		user int
	}
	tests := []struct {
		name       string
		args       args
		wantOutput User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := GetTicketRequester(tt.args.user); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("GetTicketRequester() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetOrganization(t *testing.T) {
	type args struct {
		user int
	}
	tests := []struct {
		name    string
		args    args
		wantOrg Orgs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOrg := GetOrganization(tt.args.user); !reflect.DeepEqual(gotOrg, tt.wantOrg) {
				t.Errorf("GetOrganization() = %v, want %v", gotOrg, tt.wantOrg)
			}
		})
	}
}

func TestGetRequestedTickets(t *testing.T) {
	type args struct {
		user int
	}
	tests := []struct {
		name       string
		args       args
		wantOutput ZenOutput
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := GetRequestedTickets(tt.args.user); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("GetRequestedTickets() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
