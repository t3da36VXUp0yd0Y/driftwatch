package drift

// Result holds the drift detection outcome for a single service.
type Result struct {
	Service  string
	Expected string
	Actual   string
	Drifted  bool
	Mounts   []string
	Labels   map[string]string
	Ports    []string
	Env      []string
}

// Detect compares expected image tags against actual running images.
func Detect(expected map[string]string, actual map[string]string) []Result {
	var results []Result
	for service, image := range expected {
		actualImage, running := actual[service]
		r := Result{
			Service:  service,
			Expected: image,
			Actual:   actualImage,
			Drifted:  !running || actualImage != image,
		}
		results = append(results, r)
	}
	return results
}

// HasDrift returns true if any result in the slice has drifted.
func HasDrift(results []Result) bool {
	for _, r := range results {
		if r.Drifted {
			return true
		}
	}
	return false
}
