package ursa

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGroup(id, name string) Group {
	return Group{ID: id, Name: name}
}

type Groups []Group

// GetIDs returns a list of unique group IDs.
func (g Groups) GetIDs() []string {
	idExists := map[string]struct{}{}

	var ids []string
	for _, group := range g {
		if _, ok := idExists[group.ID]; !ok {
			idExists[group.ID] = struct{}{}
			ids = append(ids, group.ID)
		}
	}

	return ids
}

// GetNames returns a list of unique group names.
func (g Groups) GetNames() []string {
	nameExists := map[string]struct{}{}

	var names []string
	for _, group := range g {
		if _, ok := nameExists[group.Name]; !ok {
			nameExists[group.Name] = struct{}{}
			names = append(names, group.Name)
		}
	}

	return names
}
