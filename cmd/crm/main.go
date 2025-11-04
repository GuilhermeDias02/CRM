package main

import (
	"github.com/GuilhermeDias02/CRM/internal/app"
	"github.com/GuilhermeDias02/CRM/internal/contact"
)

func main() {
	store := contact.NewMemoryStore()
	app.App(store)
}
