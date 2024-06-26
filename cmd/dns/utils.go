package dns

func acquireEntryFullName(domain, entry string) string {
	if entry == "@" || entry == "" {
		return domain
	}

	return entry + "." + domain
}
