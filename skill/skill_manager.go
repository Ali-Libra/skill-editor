package skill

type Skill struct {
	ID   int
	Name string
}

type SkillManager struct {
	Skills []Skill
}

func NewSkillManager() *SkillManager {
	return &SkillManager{
		Skills: []Skill{
			{ID: 1, Name: "火球术"},
			{ID: 2, Name: "冰冻术"},
		},
	}
}
