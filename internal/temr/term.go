package term

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

type pepper struct {
	Name   string
	Number int
}
type Data struct {
	Title    string
	List     []pepper
	Selected string
}

func New() *Data {
	return &Data{}
}

func (d *Data) SetList(list []string) {
	items := d.List
	counter := 1
	for _, val := range list {
		items = append(items, pepper{Name: val, Number: counter})
		counter++
	}
	d.List = items
}

func (d *Data) Run() {
	items := d.List

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0000203A {{ .Number | cyan }}. {{ .Name | cyan }}",
		Inactive: "  {{ .Number }}. {{ .Name  }}",
		Selected: "{{ `\U00002714` | green }} {{ .Name | green }}",
		Details:  ``,
	}

	searcher := func(input string, index int) bool {
		pepper := items[index]
		name := strings.Replace(strings.ToLower(pepper.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     d.Title,
		Items:     items,
		Templates: templates,
		Size:      7,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			os.Exit(0)
		}
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	d.Selected = items[i].Name
}

func (d *Data) SetTitle(title string) {
	d.Title = title
}

func (d *Data) GetSelect() string {
	return d.Selected
}

func (d *Data) ClearList() {
	d.List = []pepper{}
}
