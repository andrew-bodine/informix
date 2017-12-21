package policy_test

type MockDownstreamer struct {
    Payloads   map[string]map[string]interface{}
}

func (m *MockDownstreamer) Connect() error {
    return nil
}

func (m *MockDownstreamer) Publish(t string, payload map[string]interface{}) error {
    m.Payloads[t] = payload

    return nil
}
