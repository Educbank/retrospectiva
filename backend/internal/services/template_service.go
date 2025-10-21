package services

import (
	"errors"
)

type TemplateService struct{}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

type TemplateDefinition struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Categories   []TemplateCategory `json:"categories"`
	Instructions string             `json:"instructions"`
}

type TemplateCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
}

func (s *TemplateService) GetAvailableTemplates() []TemplateDefinition {
	return []TemplateDefinition{
		{
			ID:           "start_stop_continue",
			Name:         "Start, Stop, Continue",
			Description:  "Identifique o que começar, parar e continuar fazendo",
			Instructions: "Reflita sobre o último período e organize suas ideias em três categorias:",
			Categories: []TemplateCategory{
				{
					ID:          "start",
					Name:        "Start",
					Description: "O que devemos começar a fazer?",
					Color:       "#4CAF50",
					Icon:        "play_circle",
				},
				{
					ID:          "stop",
					Name:        "Stop",
					Description: "O que devemos parar de fazer?",
					Color:       "#F44336",
					Icon:        "stop_circle",
				},
				{
					ID:          "continue",
					Name:        "Continue",
					Description: "O que devemos continuar fazendo?",
					Color:       "#2196F3",
					Icon:        "refresh",
				},
			},
		},
		{
			ID:           "4ls",
			Name:         "4Ls",
			Description:  "Liked, Learned, Lacked, Longed for",
			Instructions: "Avalie o período através de quatro perspectivas:",
			Categories: []TemplateCategory{
				{
					ID:          "liked",
					Name:        "Liked",
					Description: "O que gostamos?",
					Color:       "#4CAF50",
					Icon:        "favorite",
				},
				{
					ID:          "learned",
					Name:        "Learned",
					Description: "O que aprendemos?",
					Color:       "#FF9800",
					Icon:        "school",
				},
				{
					ID:          "lacked",
					Name:        "Lacked",
					Description: "O que faltou?",
					Color:       "#F44336",
					Icon:        "warning",
				},
				{
					ID:          "longed_for",
					Name:        "Longed For",
					Description: "O que desejamos?",
					Color:       "#9C27B0",
					Icon:        "star",
				},
			},
		},
		{
			ID:           "mad_sad_glad",
			Name:         "Mad, Sad, Glad",
			Description:  "Identifique sentimentos e emoções",
			Instructions: "Expresse como se sentiu durante o período:",
			Categories: []TemplateCategory{
				{
					ID:          "mad",
					Name:        "Mad",
					Description: "O que nos deixou irritados?",
					Color:       "#F44336",
					Icon:        "mood_bad",
				},
				{
					ID:          "sad",
					Name:        "Sad",
					Description: "O que nos deixou tristes?",
					Color:       "#2196F3",
					Icon:        "sentiment_dissatisfied",
				},
				{
					ID:          "glad",
					Name:        "Glad",
					Description: "O que nos deixou felizes?",
					Color:       "#4CAF50",
					Icon:        "mood",
				},
			},
		},
		{
			ID:           "sailboat",
			Name:         "Sailboat",
			Description:  "Metáfora do barco a vela",
			Instructions: "Imagine que sua equipe é um barco navegando:",
			Categories: []TemplateCategory{
				{
					ID:          "wind",
					Name:        "Wind",
					Description: "O que nos empurra para frente?",
					Color:       "#4CAF50",
					Icon:        "air",
				},
				{
					ID:          "anchors",
					Name:        "Anchors",
					Description: "O que nos segura?",
					Color:       "#FF9800",
					Icon:        "anchor",
				},
				{
					ID:          "rocks",
					Name:        "Rocks",
					Description: "Riscos e obstáculos",
					Color:       "#F44336",
					Icon:        "warning",
				},
				{
					ID:          "destination",
					Name:        "Destination",
					Description: "Onde queremos chegar?",
					Color:       "#2196F3",
					Icon:        "place",
				},
			},
		},
		{
			ID:           "went_well_to_improve",
			Name:         "Went Well | To Improve",
			Description:  "O que funcionou bem e o que pode ser melhorado",
			Instructions: "Avalie o período focando em dois aspectos principais:",
			Categories: []TemplateCategory{
				{
					ID:          "went_well",
					Name:        "Went Well",
					Description: "O que funcionou bem?",
					Color:       "#4CAF50",
					Icon:        "check_circle",
				},
				{
					ID:          "to_improve",
					Name:        "To Improve",
					Description: "O que pode ser melhorado?",
					Color:       "#FF9800",
					Icon:        "trending_up",
				},
			},
		},
	}
}

func (s *TemplateService) GetTemplate(templateID string) (*TemplateDefinition, error) {
	templates := s.GetAvailableTemplates()

	for _, template := range templates {
		if template.ID == templateID {
			return &template, nil
		}
	}

	return nil, errors.New("template not found")
}

func (s *TemplateService) ValidateTemplate(templateID string) bool {
	templates := s.GetAvailableTemplates()

	for _, template := range templates {
		if template.ID == templateID {
			return true
		}
	}

	return false
}

func (s *TemplateService) GetTemplateCategories(templateID string) ([]TemplateCategory, error) {
	template, err := s.GetTemplate(templateID)
	if err != nil {
		return nil, err
	}

	return template.Categories, nil
}
