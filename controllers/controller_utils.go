package controllers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	sh "github.com/codeskyblue/go-sh"
	"github.com/gernest/lora/models"
	"github.com/gernest/lora/utilities/pbkdf2"
)

// Rebuild takes a new saved page object and rebuilds the project with ne updted conent
func Rebuild(p *models.Page) error {
	logThis.Event(" Rebuilding %s.....",p.Title )
	project := new(models.Project)
	db, err := models.Conn()
	if err != nil {
		return err
	}
	err = db.Find(project, p.ProjectId).Error
	if err != nil {
		return err
	}
	err = project.LoadConfigFile()
	if err != nil {
		return err
	}
	for k := range project.Pages {
		pj := &project.Pages[k]
		sections := []models.Section{}
		db.Model(pj).Related(&sections)

		if pj.Id == p.Id {

			if len(sections) > 0 {
				for key := range sections {
					s := &sections[key]
					sub := []models.SubSection{}
					db.Model(s).Related(&sub)
					s.SubSections = sub
				}
				pj.Sections = sections
			}
			pj.Content = p.Content
		}
	}
	err = project.SaveConfigFile()
	if err != nil {
		return err
	}
	err = project.GenContent()
	if err != nil {
		return err
	}
	err = project.Build()
	if err != nil {
		return err
	}
	logThis.Success(" *** done  building %s***", p.Title)
	return nil
}

// Deploy pushes the project to a dokku server
func Deploy(p *models.Project) error {
	sess := sh.NewSession().SetDir(p.ProjectPath)
	remote := fmt.Sprintf("%s-deploy", p.Name)
	remoteURI := fmt.Sprintf("dokku@tushabe.com:%s", p.Name)
	commitMsg := fmt.Sprintf("deloyed at %s", time.Now().Format(time.RFC822))
	beego.Info(" deploying to ", remoteURI)

	err := p.Build()
	if err != nil {
		beego.Debug(err)
		return err
	}
	err = sess.Command("git", "commit", "-m", commitMsg).Run()
	if err != nil {
		beego.Debug(err)
		return err
	}
	err = sess.Command("git", "remote", "add", remote, remoteURI).Run()
	if err != nil {
		beego.Debug(err)
		return err
	}
	err = sess.Command("git", "push", remote, "master").Run()
	if err != nil {
		beego.Debug(err)
		return err
	}
	beego.Info("***deployed successful****")
	return nil

}

func checkUserByName(sess map[string]interface{}) (models.Account, error) {
	usr := models.Account{}
	db, err := models.Conn()
	if err != nil {
		return usr, err
	}

	usr.UserName = sess["username"].(string)
	query := db.Where("user_name= ?", usr.UserName).First(&usr)
	if query.Error != nil {
		return usr, err
	}
	return usr, err
}
func checkUserByEmail(email string) (models.Account, error) {
	usr := models.Account{}
	db, err := models.Conn()
	if err != nil {
		return usr, err
	}

	usr.Email = email
	query := db.Where("email= ?", usr.Email).First(&usr)
	if query.Error != nil {
		return usr, errors.New("Problem querying the database")
	}
	return usr, err
}
func checkUserPassword(usr *models.Account, passwd string) error {
	var x pbkdf2.PasswordHash
	var err error
	x.Hash = make([]byte, 32)
	x.Salt = make([]byte, 16)

	if usr.Password == "" {
		return errors.New("No such User")
	}
	x.Hash, err = hex.DecodeString(usr.Password[:64])
	if err != nil {
		return err
	}
	x.Salt, err = hex.DecodeString(usr.Password[64:])
	if err != nil {
		return err
	}
	if !pbkdf2.MatchPassword(passwd, &x) {
		return errors.New("Wrong Password")
	}
	return nil
}
