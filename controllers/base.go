package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gernest/lora/utilities/logs"
)

var logThis = logs.NewLoraLog()

// MainController provides base methods for all lora controllers
type MainController struct {
	beego.Controller
}

// ActivateContent makes it easy to add layout to templates, it also checks
// Session cookie if is set and do the initializing stuffs
func (c *MainController) ActivateContent(view string) map[string]interface{} {
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHeader"] = "header.html"
	c.LayoutSections["HtmlFooter"] = "footer.html"
	c.TplNames = view + ".html"
	c.Layout = "base.html"

	logThis.Info("Checking session")
	sess := c.GetSession("xshabe")
	if sess != nil {
		c.Data["InSession"] = 1
		m := sess.(map[string]interface{})
		c.Data["Username"] = m["username"]
		logThis.Success("Session found *%v*", m["username"])
		return m

	}
	logThis.Warning("No session found")
	return nil
}

// Get takes you home baby
func (c *MainController) Get() {
	_ = c.ActivateContent("index")
	c.SetNotice()
}

// Notice this is an old school notice page
func (c *MainController) Notice() {
	_ = c.ActivateContent("notice")

	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		c.Data["notice"] = n
	}

}

// SetNotice makes it easier to set flash notices
func (c *MainController) SetNotice() {
	flash := beego.ReadFromRequest(&c.Controller)
	if n, ok := flash.Data["notice"]; ok {
		c.Data["notice"] = n
	}
}

func (c *MainController) ActivateView(view string) {
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["HtmlHeader"] = "header.html"
	c.LayoutSections["HtmlFooter"] = "footer.html"
	c.TplNames = view + ".html"
	c.Layout = "base.html"
}