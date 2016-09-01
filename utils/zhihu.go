package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	. "github.com/rexchen123/lovelydata/extensions"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
)

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type User struct {
	ZhihuId          string
	Name             string
	Address          string
	Avatar           string
	Business         string
	Gender           string
	Education        string
	Major            string
	Description      string
	FolloweesCount   int
	FollowersCount   int
	SpecialCount     int
	FollowTopicCount int
	ApprovalCount    int
	ThankCount       int
	AnswerCount      int
	ArticleCount     int
	PvCount          int
	StartedCount     int
	PublicEditCount  int
	AskCount         int
}

func FetchUserInfo(name string) *User {
	client := &http.Client{}
	url := fmt.Sprintf("http://www.zhihu.com/people/%s", name)
	request, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return nil
	}

	addCookies(request)
	beego.Info(request.Cookies())
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
		beego.Info(response)
		body, _ := ioutil.ReadAll(response.Body)
		return getUserInfo(string(body))
	}
	if err == nil && response.StatusCode == http.StatusOK {
		beego.Info(response.Body)
		return nil
	}
	beego.Info(err)

	return nil
}

func GetFollowerList(name, userType string) {
	client := &http.Client{}
	url := fmt.Sprintf("https://www.zhihu.com/people/%s/%s", name, userType)
	request, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return
	}

	addCookies(request)
	beego.Info(request.Cookies())
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
		beego.Info(response)
		body, _ := ioutil.ReadAll(response.Body)
		handleUserList(string(body))
		return
	}
	if err == nil && response.StatusCode == http.StatusOK {
		beego.Info(response.Body)
		return
	}
	beego.Info(err)

	return
}

func handleUserList(body string) {
	reg := regexp.MustCompile(`<h2 class="zm-list-content-title"><a data-hovercard=".*?" href="https://www.zhihu.com/people/(.*?)" class="zg-link author-link" title="(.*?)"\s>`)
	result := reg.FindAllStringSubmatch(body, 20)

	for _, item := range result {
		if len(item) >= 1 {
			beego.Warn(item[1])
			Q.Queue.Put(item[1])
		}
	}
}

func addCookies(request *http.Request) {
	conf, err := config.NewConfig("json", "conf/zhihucookies.json")
	if err != nil {
		beego.Error(err)
	} else {
		data, err := conf.DIY("cookies")
		if err != nil {
			beego.Error(err)
			panic(err)
		}
		beego.Info(reflect.TypeOf(data))
		beego.Info(data)
		cookies, ok := data.([]interface{})
		if ok {
			for _, cookie := range cookies {
				cookie, ok := cookie.(map[string]interface{})
				if ok {
					name, _ := cookie["name"].(string)
					value, _ := cookie["value"].(string)
					request.AddCookie(&http.Cookie{
						Name:  name,
						Value: value,
					})
				} else {
					beego.Error(fmt.Errorf("add cookies failed"))
				}
			}
		} else {
			beego.Error(fmt.Errorf("add cookies failed"))
		}
	}
}

type Html struct {
	Content string
}

func getUserInfo(body string) *User {
	user := User{}
	html := Html{
		Content: body,
	}
	reg := regexp.MustCompile(`<span class="name">(.*)</span><a class="icon-badge-wrapper" href='#'>`)
	user.Name = html.getItem(reg)

	reg = regexp.MustCompile(`<a class="item home first active"\nhref="/people/(.*)">`)
	user.ZhihuId = html.getItem(reg)

	reg = regexp.MustCompile(`<span class="info-wrap">\n\n<span class="location item" title="(.*)">`)
	user.Address = html.getItem(reg)

	reg = regexp.MustCompile(`<span class="business item" title=["|\'](.*?)["|\']>`)
	user.Business = html.getItem(reg)

	reg = regexp.MustCompile(`<i class="icon icon-profile-(.*?)male"></i>`)
	user.Gender = html.getItem(reg) + "male"

	reg = regexp.MustCompile(`<span class="education item" title=["|\'](.*?)["|\']>`)
	user.Education = html.getItem(reg)

	reg = regexp.MustCompile(`<span class="education-extra item" title=["|\'](.*?)["|\']>`)
	user.Major = html.getItem(reg)

	reg = regexp.MustCompile(`<span class="content">\s(.*?)\s</span>s`)
	user.Description = html.getItem(reg)

	reg = regexp.MustCompile(`<span class="zg-gray-normal">关注了</span><br />\s<strong>(.*?)</strong><label> 人</label>`)
	user.FolloweesCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`<span class="zg-gray-normal">关注者</span><br />\s<strong>(.*?)</strong><label> 人</label>`)
	user.FollowersCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`<strong>(.*?) 个专栏</strong>`)
	user.SpecialCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`<strong>(.*?) 个话题</strong>`)
	user.FollowTopicCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`<span class="zm-profile-header-user-agree"><span class="zm-profile-header-icon"></span><strong>(.*?)</strong>赞同</span>`)
	user.ApprovalCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`<span class="zm-profile-header-user-thanks"><span class="zm-profile-header-icon"></span><strong>(.*?)</strong>感谢</span>`)
	user.ThankCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`提问\s<span class="num">(.*?)</span>`)
	user.AskCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`回答\s<span class="num">(.*?)</span>`)
	user.AnswerCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`文章\s<span class="num">(.*?)</span>`)
	user.ArticleCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`个人主页被 <strong>(.*?)</strong> 人浏览`)
	user.PvCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`收藏\s<span class="num">(.*?)</span>`)
	user.StartedCount = cast.ToInt(html.getItem(reg))

	reg = regexp.MustCompile(`公共编辑\s<span class="num">(.*?)</span>`)
	user.PublicEditCount = cast.ToInt(html.getItem(reg))

	beego.Info(user)
	return &user
}

func (html Html) getItem(reg *regexp.Regexp) string {
	result := reg.FindStringSubmatch(html.Content)
	if len(result) <= 1 {
		return ""
	} else {
		return result[1]
	}
}
