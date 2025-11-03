package mqtt

import "testing"

func TestNewClient_InvalidBroker(t *testing.T) {
	logger := &mockLogger{}
	_, err := NewClient(logger, "invalid://broker", "test-client")
	if err == nil {
		t.Error("NewClient() with invalid broker should return error")
	}
}

func TestClient_Publish(t *testing.T) {
	logger := &mockLogger{}

	tests := []struct {
		name     string
		broker   string
		clientID string
		wantErr  bool
	}{
		{
			name:     "invalid broker",
			broker:   "invalid://broker",
			clientID: "test-client",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(logger, tt.broker, tt.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			err = client.Publish("test/topic", "test payload", 1, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			if client != nil {
				client.Close()
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	logger := &mockLogger{}

	client, err := NewClient(logger, "tcp://test.mosquitto.org:1883", "test-client")
	if err != nil {
		t.Skipf("Skipping test: cannot connect to broker: %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}
