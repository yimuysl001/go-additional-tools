package duckdb

// OrderRandomFunction returns the SQL function for random ordering.
func (d *Driver) OrderRandomFunction() string {
	return "RANDOM()"
}
