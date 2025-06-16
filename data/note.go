package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/john-marinelli/bon/cfg"
)

type Note struct {
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
	DaysLeft int       `json:"days_left"`
	Id       int       `json:"id"`
}

func SaveNote(path string, content string) error {
	fPath := cfg.Config.ArchDir + path + ".md"
	if path != "" {
		dPath := filepath.Dir(fPath)
		if _, err := os.Stat(dPath); os.IsNotExist(err) {
			err := os.MkdirAll(dPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	if path == "" {
		data, err := os.ReadFile(cfg.Config.BonFile)
		if err != nil {
			return err
		}
		notes := []*Note{}

		json.Unmarshal(data, &notes)
		id := 0
		if len(notes) == 0 {
			id = 1
		} else {
			for _, i := range notes {
				if i.Id > id {
					id = i.Id
				}
			}
			id += 1
		}
		notes = append(notes, &Note{
			Content:  content,
			Date:     time.Now(),
			DaysLeft: 7,
			Id:       id,
		})

		f, err := os.Create(cfg.Config.BonFile)
		defer f.Close()
		b, err := json.Marshal(notes)
		if err != nil {
			panic(err)
		}
		f.Write(b)
		return err
	}

	f, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(content)
	return nil
}

func LoadBonNotes() ([]Note, error) {
	data, err := os.ReadFile(cfg.Config.BonFile)
	if err != nil {
		return nil, err
	}
	notes := []Note{}

	json.Unmarshal(data, &notes)

	return notes, nil
}

func DeleteBonNote(id int) ([]Note, error) {
	data, err := os.ReadFile(cfg.Config.BonFile)
	if err != nil {
		return nil, err
	}
	notes := []Note{}

	json.Unmarshal(data, &notes)
	d := []Note{}
	for _, n := range notes {
		if n.Id == id {
			continue
		}
		d = append(d, n)
	}

	f, err := os.Create(cfg.Config.BonFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	f.Write(b)

	return d, nil
}

func EditBonNote(id int, content string) error {
	data, err := os.ReadFile(cfg.Config.BonFile)
	if err != nil {
		return err
	}
	notes := []Note{}

	json.Unmarshal(data, &notes)
	d := []Note{}
	for i := range notes {
		if notes[i].Id == id {
			notes[i].Content = content
		}
	}

	fmt.Println(notes)

	f, err := os.Create(cfg.Config.BonFile)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	f.Write(b)

	return nil

}

func LoadAndClearNotes() ([]Note, error) {
	data, err := os.ReadFile(cfg.Config.BonFile)
	if err != nil {
		return nil, err
	}
	notes := []Note{}

	json.Unmarshal(data, &notes)
	d := []Note{}
	for _, n := range notes {
		n.DaysLeft = int(
			n.Date.AddDate(
				0,
				0,
				cfg.Config.MaxDays,
			).Sub(
				time.Now(),
			).Hours() / 24,
		)
		if n.DaysLeft < 0 {
			continue
		}
		d = append(d, n)
	}

	f, err := os.Create(cfg.Config.BonFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	f.Write(b)

	return d, nil
}
