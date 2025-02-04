package handlebars

import (
	"io/ioutil"
	"path"
	"sync"

	find "github.com/1800alex/go-find"
)

type Templater struct {
	handlebars *Handlebars

	templatesMu sync.Mutex
	templates   map[string]*Template
}

// NewTemplater creates a new Templater instance
func NewTemplater() *Templater {
	return &Templater{
		handlebars: New(),
		templates:  make(map[string]*Template),
	}
}

// LoadTemplates loads all files in a directory and parses them as handlebars templates
func (t *Templater) LoadTemplates(dir string) error {
	t.templatesMu.Lock()
	defer t.templatesMu.Unlock()

	files, err := find.Find(dir, find.Options{
		RegularFilesOnly: true,
		Recursive:        true,
	})

	if err != nil {
		return err
	}

	for _, file := range files {
		tpl, err := t.handlebars.ParseFile(file.Path)
		if err != nil {
			return err
		}

		t.templates[path.Base(file.Path)] = tpl
	}

	return nil
}

// LoadPartials loads all files in a directory and parses them as handlebars partials
func (t *Templater) LoadPartials(dir string) error {
	files, err := find.Find(dir, find.Options{
		RegularFilesOnly: true,
		Recursive:        true,
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		contents, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return err
		}

		t.handlebars.RegisterPartial(path.Base(file.Path), string(contents))
	}

	return nil
}

// Get returns the previously loaded template from the map
func (t *Templater) Get(name string) (*Template, bool) {
	t.templatesMu.Lock()
	defer t.templatesMu.Unlock()

	res, ok := t.templates[name]
	return res, ok
}
