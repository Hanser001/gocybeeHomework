package api

import (
	"bufio"
	"fmt"
	"ginDemo/api/middleware"
	"ginDemo/dao"
	"ginDemo/model"
	"ginDemo/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"time"
)

func register(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespSuccess(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	//把用户名传入dao/regist.txt
	file, err := os.OpenFile("E:\\jetbrains\\goland\\ginDemo\\dao\\regist.txt", os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		utils.RespFail(c, "file wrong")
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if username == string(line) {
			utils.RespFail(c, "用户已存在")
			return
		}
		if err == io.EOF {
			file.WriteString("\n")
			file.WriteString(username) //传入用户名
			dao.Adduser(username, password)
			break
		}
		if err != nil {
			utils.RespFail(c, "file wrong")
			return
		}
	}

	// 以 JSON 格式返回信息
	utils.RespSuccess(c, "add user successful")
}

// 仅有登录部分有改动
func login(c *gin.Context) {
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}

	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	file, err := os.Open("E:\\jetbrains\\goland\\ginDemo\\dao\\regist.txt")
	if err != nil {
		utils.RespFail(c, "file wrong")
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			utils.RespFail(c, "用户不存在")
			return
		}
		if username == string(line) {
			dao.Adduser(username, password)
			break
		}
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}

	// 正确则登录成功
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "Yxh",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	utils.RespSuccess(c, tokenString)
}

func getUsernameFromToken(c *gin.Context) {
	username, _ := c.Get("username")
	utils.RespSuccess(c, username.(string))
}

// 实现登录留言
func leave(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	message := c.PostForm("message")
	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}
	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}
	dao.LeaveMessage(message)
	//把留言保存到dao库的txt文件
	file, err := os.OpenFile("E:\\jetbrains\\goland\\ginDemo\\dao\\messages.txt", os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		utils.RespFail(c, "file wrong")
		return
	}
	defer file.Close()
	file.WriteString(username)
	file.WriteString(":")
	file.WriteString(message)
	//保证每次留言都会换行
	file.WriteString("\n")
	utils.RespSuccess(c, message)
}

// 添加密保问题
func addQuestion(c *gin.Context) {
	//传入用户名和密码，问题和答案
	username := c.PostForm("username")
	password := c.PostForm("password")
	question := c.PostForm("question")
	answer := c.PostForm("answer")

	// 验证用户名是否存在
	flag := dao.SelectUser(username)
	fmt.Println(flag)
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "用户不存在")
		return
	}
	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username)
	// 若不正确则传出错误,正确则进入下一步
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}
	dao.AddPwdQuestion(username, question)
	dao.AddPwdAnswer(username, answer)
	utils.RespSuccess(c, "设置成功")
}

// 实现找回密码
func findpassword(c *gin.Context) {
	//传入用户名和答案与新密码
	username := c.PostForm("username")
	answer := c.PostForm("answer")
	newPassword := c.PostForm("newPassword")

	// 验证用户名是否存在
	flag1 := dao.SelectUser(username)
	fmt.Println(flag1)
	if !flag1 {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "用户不存在")
		return
	}
	//查找该用户名是否设置密保问题
	flag2 := dao.SelectQuestion(username)
	if !flag2 {
		utils.RespFail(c, "该用户未设置密保问题")
		return
	}
	//查找密保问题答案
	trueAnswer := dao.SelectAnswerFromQuestion(username)
	//答案不正确就报错
	if answer != trueAnswer {
		utils.RespFail(c, "答案错误")
		return
	}
	//正确就修改密码
	dao.FindPassword(username, newPassword)
	utils.RespSuccess(c, "密码修改成功")
}
