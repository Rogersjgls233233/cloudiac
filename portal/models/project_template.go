package models

import "cloudiac/portal/libs/db"

type ProjectTemplate struct {
	BaseModel
	ProjectId	string  `json:"projectId" gorm:"not null"`
	TemplateId  string  `json:"templateId" gorm:"not null"`
}

func (ProjectTemplate) TableName() string {
	return "iac_project_template"
}


func (u ProjectTemplate) Migrate(sess *db.Session) error {
	return u.AddUniqueIndex(sess, "unique__project__template", "project_id", "template_id")
}