package config

//go:generate mockgen -destination mocks/mock_composite.go . Composite
type Composite interface {
	CloudflareAPIKey() string
}

type composite struct {
	entries []Composite
}

func NewComposite() (Composite, error) {
	localConfig, err := LoadLocalConfig()
	if err != nil {
		return nil, err
	}

	return &composite{
		entries: []Composite{
			localConfig,
			envConfig{},
		},
	}, nil
}

func (c *composite) CloudflareAPIKey() string {
	for _, entry := range c.entries {
		if key := entry.CloudflareAPIKey(); key != "" {
			return key
		}
	}

	return ""
}
