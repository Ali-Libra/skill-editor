package skill

import "skill-editor/node"

type Skill struct {
	ID   int
	Name string
	Node node.Node
}

type SkillManager struct {
	Skills []Skill
}

func NewSkillManager() *SkillManager {
	return &SkillManager{
		Skills: []Skill{
			{ID: 1, Name: "单骑"},
			{ID: 2, Name: "武圣"},
		},
	}
}
