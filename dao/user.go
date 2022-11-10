package dao

var datebase = make(map[string]string, 100)

var messageBin = make([]string, 100)

var passwordQuestion = make(map[string]string, 100)

var pwdQuetionAnswer = make(map[string]string, 100)

func Adduser(username string, password string) {
	datebase[username] = password
}

func SelectUser(username string) bool {
	if datebase[username] == "" {
		return false
	}
	return true
}

func SelectPasswordFromUsername(username string) string {
	return datebase[username]
}

func SelectQuestion(username string) bool {
	if passwordQuestion[username] == "" {
		return false
	}
	return true
}

func SelectAnswerFromQuestion(username string) string {
	question := passwordQuestion[username]
	return pwdQuetionAnswer[question]
}

func LeaveMessage(message string) {
	//把留言存入切片
	messageBin = append(messageBin, message)
}

// 实现添加密保问题
func AddPwdQuestion(username string, question string) {
	passwordQuestion[username] = question
}

func AddPwdAnswer(username string, answer string) {
	question := passwordQuestion[username]
	pwdQuetionAnswer[question] = answer
}

// 实现找回密码
func FindPassword(username string, newpassword string) {
	datebase[username] = newpassword
}
