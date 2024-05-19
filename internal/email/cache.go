package email

type Cache struct {
	Email string
	Name  string
}

func NewCache(email *string, name *string) *Cache {
	return &Cache{
		Email: *email,
		Name:  *name,
	}
}
