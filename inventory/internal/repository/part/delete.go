package part

// DeletePart удаляет деталь по UUID. Потокобезопасно.
func (r *repository) DeletePart(uuid string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, existed := r.parts[uuid]
	delete(r.parts, uuid)
	return existed
}
