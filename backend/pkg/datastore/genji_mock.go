package datastore

// MockGenjiDB represents the mocked object
type MockGenjiDB struct {
	callParams []interface{}
}

// Exec implements the GenjiDB interface
func (m *MockGenjiDB) Exec(q string, args ...interface{}) error {
	m.callParams = []interface{}{q}
	m.callParams = append(m.callParams, args...)

	return nil
}

// CalledWith returns the arguments which with the
// Exec function was called.
func (m *MockGenjiDB) CalledWith() []interface{} {
	return m.callParams
}
