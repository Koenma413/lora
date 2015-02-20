// Copyright 2015 Geofrey Ernest a.k.a gernest, All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gernest/lora/models"
)

type SubSectionController struct {
	MainController
}

func (s *SubSectionController) Update() {
	flash := beego.NewFlash()

	sess := s.ActivateContent("subsection/edit")
	lora := models.NewLoraObject()
	if sess == nil {
		flash.Error("you need to login inorder to update this page")
		flash.Store(&s.Controller)
		return
	}
	s.LayoutSections["JScripts"] = "jscript/editor.html"

	subSectionID, _ := s.GetInt64(":subSectionID")
	sectionID, _ := s.GetInt64(":sectionID")

	subSection := models.SubSection{}
	page := models.Page{}
	section := models.Section{}

	db, err := models.Conn()
	defer db.Close()
	if err != nil {
		flash.Error("Whacko opening the database")
		flash.Store(&s.Controller)
		return
	}
	err = db.First(&subSection, subSectionID).Error
	if err != nil {
		flash.Error("WHacko ", err)
		flash.Store(&s.Controller)
		return
	}
	if subSection.SectionId != sectionID {
		flash.Error("The page does not belong to this project")
		flash.Store(&s.Controller)
		return
	}
	subSection.Sanitize()
	lora.Add(subSection)
	s.Data["lora"] = lora

	if s.Ctx.Input.Method() == "POST" {
		subSectionContent := s.GetString("content")
		subSection.Body = subSectionContent
		db.Save(&subSection)

		db.First(&section, subSection.SectionId)
		db.First(&page, section.PageId)

		err = Rebuild(&page)
		if err != nil {
			flash.Error(" WHacko ", err)
			flash.Store(&s.Controller)
			return
		}

		s.Redirect("/web/accounts", 302)
	}
}
