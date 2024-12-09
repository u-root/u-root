package main

// CI is a struct autogenerate from go.yml
type CI struct {
	Name string `yaml:"name"`
	On   struct {
		Push struct {
			Branches []string `yaml:"branches"`
		} `yaml:"push"`
		PullRequest struct {
			Branches []string `yaml:"branches"`
		} `yaml:"pull_request"`
	} `yaml:"on"`
	Jobs struct {
		Linters struct {
			RunsOn string `yaml:"runs-on"`
			Steps  []struct {
				Uses string `yaml:"uses,omitempty"`
				Name string `yaml:"name,omitempty"`
				With struct {
					GoVersion string `yaml:"go-version"`
				} `yaml:"with,omitempty"`
				Run string `yaml:"run,omitempty"`
			} `yaml:"steps"`
		} `yaml:"linters"`
		Build struct {
			RunsOn string `yaml:"runs-on"`
			Steps  []struct {
				Uses string `yaml:"uses,omitempty"`
				Name string `yaml:"name,omitempty"`
				With struct {
					GoVersion string `yaml:"go-version"`
				} `yaml:"with,omitempty"`
				Run string `yaml:"run,omitempty"`
			} `yaml:"steps"`
		} `yaml:"build"`
		Badbuild struct {
			RunsOn string `yaml:"runs-on"`
			Steps  []struct {
				Uses string `yaml:"uses,omitempty"`
				Name string `yaml:"name,omitempty"`
				With struct {
					GoVersion string `yaml:"go-version"`
				} `yaml:"with,omitempty"`
				ID              string `yaml:"id,omitempty"`
				ContinueOnError bool   `yaml:"continue-on-error,omitempty"`
				Run             string `yaml:"run,omitempty"`
			} `yaml:"steps"`
		} `yaml:"badbuild"`
		MultiOsArch struct {
			Strategy struct {
				Matrix struct {
					Os   []string `yaml:"os"`
					Arch []string `yaml:"arch"`
				} `yaml:"matrix"`
			} `yaml:"strategy"`
			RunsOn string `yaml:"runs-on"`
			Steps  []struct {
				Uses string `yaml:"uses,omitempty"`
				Name string `yaml:"name,omitempty"`
				With struct {
					GoVersion string `yaml:"go-version"`
				} `yaml:"with,omitempty"`
				Run string `yaml:"run,omitempty"`
			} `yaml:"steps"`
		} `yaml:"multi-os-arch"`
	} `yaml:"jobs"`
}
