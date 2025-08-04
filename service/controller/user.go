package controller

import (
	"encoding/base64"
	"encoding/json"
	"gdp/service/configs"
	"gdp/service/middleware"
	"gdp/utils"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name     string `json:"name"`
	Pass     string `json:"pass"`
	Disabled bool   `json:"disabled"`
	Role     int    `json:"role"`
}

// 获取用户列表
func userList() ([]User, error) {
	userDataFilename := configs.Service.DbDir + "user.json"
	content, err := os.ReadFile(userDataFilename)
	if err != nil {
		utils.FileCreateWithDirs(userDataFilename)
		content = []byte("[]")
		os.WriteFile(userDataFilename, content, 0644)
	}

	users := []User{}
	jsonErr := json.Unmarshal(content, &users)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return users, nil
}

func userAdd(user User) error {
	users, err := userList()
	if err != nil {
		return err
	}
	users = append(users, user)
	content, err := json.Marshal(users)
	if err != nil {
		return err
	}

	userDataFilename := configs.Service.DbDir + "user.json"
	err = os.WriteFile(userDataFilename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Register(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	type Params struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	var params Params
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	users, err := userList()
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	isExist := utils.SliceSome(users, func(user User, index int) bool {
		return user.Name == params.Username
	})
	if isExist {
		ctx.ErrorCustom("User already exists")
		return
	}

	addUserErr := userAdd(User{
		Name:     params.Username,
		Pass:     utils.Md5Encipher(params.Password),
		Disabled: true,
		Role:     1,
	})
	if addUserErr != nil {
		ctx.ErrorCustom(addUserErr.Error())
		return
	}

	ctx.SuccessData("Wait for the administrator's confirmation")
}

func Verification(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	auth := ctx.GetHeader("Authorization")
	if len(auth) < 6 || auth[:6] != "Basic " {
		ctx.ErrorAuth("Verification failed")
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}
	info := strings.Split(string(decoded), ":")
	users, err := userList()
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	query := utils.SliceFind(users, func(user User, index int) bool {
		return user.Name == info[0]
	})
	if query.Name == "" {
		ctx.ErrorCustom("User not found")
		return
	}

	if query.Disabled {
		ctx.ErrorAuth("User not activated")
		return
	}

	passwrod := utils.Md5Encipher(info[1])
	if query.Pass != passwrod {
		ctx.ErrorAuth("Password error")
		return
	}

	ctx.Set("userInfo", query)
	ctx.Next()
}

// 获取当前用户信息，在此之前必须使用 Verification
func GetUserInfo(c *gin.Context) User {
	value, _ := c.Get("userInfo")
	info, _ := value.(User)
	return info
}

func UserInfo(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	info := GetUserInfo(c)
	info.Pass = ""
	ctx.SuccessData(info)
}

func UserList(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	users, err := userList()
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}
	for i := range users {
		users[i].Pass = ""
	}
	ctx.SuccessData(users)
}

// 验证用户角色
func VerificationRole(role int) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := middleware.ContextGet(c)
		info := GetUserInfo(c)
		if info.Role != role {
			ctx.ErrorAuth("Permission denied")
			return
		}
		ctx.Next()
	}
}

// 激活用户
func Activated(c *gin.Context) {
	ctx := middleware.ContextGet(c)

	type Params struct {
		Username string `form:"username" binding:"required"`
		Disabled bool   `form:"disabled"`
	}
	var params Params
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	users, err := userList()
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	for index, user := range users {
		if user.Name == params.Username {
			users[index].Disabled = params.Disabled
			break
		}
	}
	content, err := json.Marshal(users)
	if err != nil {
		ctx.ErrorCustom(err.Error())
	}
	userDataFilename := configs.Service.DbDir + "user.json"
	err = os.WriteFile(userDataFilename, content, 0644)
	if err != nil {
		ctx.ErrorCustom(err.Error())
	}

	ctx.Success()
}

func Login(c *gin.Context) {
	ctx := middleware.ContextGet(c)
	type Params struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	var params Params
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.ErrorParams(err.Error())
		return
	}

	users, err := userList()
	if err != nil {
		ctx.ErrorCustom(err.Error())
		return
	}

	params.Password = utils.Md5Encipher(params.Password)
	query := utils.SliceFind(users, func(user User, index int) bool {
		return user.Name == params.Username && user.Pass == params.Password
	})
	if query.Name == "" {
		ctx.ErrorCustom("Username or password error")
		return
	}

	ctx.SuccessData("Basic " + base64.StdEncoding.EncodeToString([]byte(query.Name+":"+query.Pass)))
}
