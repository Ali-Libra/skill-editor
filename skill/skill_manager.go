package skill

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Skill struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Times       int
	Times5      int
	Hero        int
	MasterSkill int
	Stages      []int
	Active      int

	Nodes []SkillNode
}

type SkillManager struct {
	Skills   []Skill
	filePath string
}

func NewSkillManager() *SkillManager {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	filePath := filepath.Join(cwd, "skill.json")

	manager := &SkillManager{
		Skills:   []Skill{},
		filePath: filePath,
	}

	manager.load()
	return manager
}

func (m *SkillManager) AddSkill(id, name string) {
	m.Skills = append(m.Skills, Skill{ID: id, Name: name})
}

func (m *SkillManager) RemoveSkill(target *Skill) {
	for i := range m.Skills {
		if &m.Skills[i] == target {
			m.Skills = append(m.Skills[:i], m.Skills[i+1:]...)
			return
		}
	}
}

func (m *SkillManager) Save() error {
	data, err := json.MarshalIndent(m.Skills, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

func (m *SkillManager) load() {
	if _, err := os.Stat(m.filePath); os.IsNotExist(err) {
		return
	}
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &m.Skills)
}
