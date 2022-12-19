package config

type DummyConfig struct {
}

const dummy = "dummy"

func (d DummyConfig) GetBillingToken() string {
	return dummy
}

func (d DummyConfig) GetDBConn() string {
	return dummy
}

func (d DummyConfig) GetServerConn() string {
	return dummy
}

func (d DummyConfig) GetParallelGoroutines() int {
	return 1
}

func (d DummyConfig) GetEndpoint() string {
	return dummy
}

func (d DummyConfig) JobIsEnabled() bool {
	return true
}
