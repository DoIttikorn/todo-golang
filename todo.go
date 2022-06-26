package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*t = append(*t, todo)

}
func (t *Todos) Complete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("index out of range")
	}
	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true
	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("index out of range")
	}
	*t = append(ls[:index-1], ls[index:]...) // understand this

	return nil
}

func (t *Todos) Load(filename string) error {
	file, error := ioutil.ReadFile(filename)
	if error != nil {
		if errors.Is(error, os.ErrNotExist) { // understand errors.Is
			return nil
		}
		return error
	}

	if len(file) == 0 {
		return error
	}

	error = json.Unmarshal(file, t)
	if error != nil {
		return error
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "Created At"},
			{Align: simpletable.AlignRight, Text: "Completed At"},
		},
	}
	var cells [][]*simpletable.Cell

	for i, item := range *t {

		task := blue(item.Task)
		done := blue("NO")
		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			done = green("YES")
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", i+1)},
			{Align: simpletable.AlignLeft, Text: task},
			{Align: simpletable.AlignCenter, Text: done},
			{Align: simpletable.AlignRight, Text: item.CreatedAt.Format(time.RFC822)},
			{Align: simpletable.AlignRight, Text: item.CompletedAt.Format(time.RFC822)},
		})
	}

	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: red(fmt.Sprintf("%d pending tasks", t.CountPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (t *Todos) CountPending() int {
	var count int
	for _, item := range *t {
		if !item.Done {
			count++
		}
	}
	return count
}
