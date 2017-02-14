package gorm

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//User : ~
type User struct {
	gorm.Model
	FirstName           string `json:"firstName" form:"firstName"`
	LastName            string `json:"LastName" form:"lastName"`
	Email               string `json:"email" form:"email"`
	Address             string `json:"address" form:"address"`
	Username            string `json:"username" form:"username"`
	PhoneNumber         string `json:"phoneNumber" form:"phoneNumber"`
	HashedPassword      []byte `json:"hashedPassword" form:"hashedPassword"`
	FriendsIDs          []User
	ProblemsPostedIDs   []Problem
	SolutionsIDs        []Solution
	FollowedProblemsIDs []Problem
	CommentIDs          []Comment
	IsDisabled          bool
}

//LoginForm : ~
type LoginForm struct {
	Password string `json:"password" form:"password"`
	Username string `json:"username" form:"username"`
}

// RegistrationForm : A registration struct
type RegistrationForm struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

//PasswordResetForm : ~
type PasswordResetForm struct {
	Email string `json:"email" form:"email"`
}

//API Functions

// CreateUser : check if user is already created,
// if not use RegistrationForm to populate a new one
func CreateUser(form RegistrationForm) {

	//check DB if Username is already taken
	u := User{}
	err := db.Where("username = ?", form.Username).First(&u).Value
	if err == nil {
		glog.Info("error has occured")
	}
	glog.Info(err)
	//If username does not exist
	if u.Username == "" {
		glog.Info("Username not taken...")
		//generate hashpassword
		passwordBytes := []byte(form.Password)
		hashedPassword, hashError := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
		if hashError != nil {
			glog.Error("Bcrypt failed to hash password")
		}

		//Reinitalize user as new user and set it with form details
		u = User{}
		u.HashedPassword = hashedPassword
		u.Email = form.Email
		u.Username = form.Username

		db.Create(&u)
	} else {
		glog.Error("Username is already taken")
	}

}

//GetUserByUsername : get user by name
func (u *User) GetUserByUsername(name string) bool {
	err := db.Where("username = ?", name).First(&u)
	if err == nil {
		glog.Info("There was an error")
	}
	if u.ID == 0 {
		return false
	}
	return true
}

//VerifyUser : Checks db credentials
func (u *User) VerifyUser(username string, password string) bool {
	err := db.Where("username = ? AND hashed_password", username, password).First(&u)
	if err == nil {
		glog.Info("There was an error")
	}
	if u.ID == 0 {
		return false
	}
	return true
}

//LoginAttempt : Logs everytime someone logs on
func (l LoginForm) LoginAttempt() {
	db.Create(l)
}

// PostProblem : User Auth Required> Post Problem
func (u *User) PostProblem(text string, description string) {
	p := Problem{
		OriginalPoster: *u,
		Title:          text,
		Description:    description,
	}
	db.Create(&p)
	u.ProblemsPostedIDs = append(u.ProblemsPostedIDs, p)
}

//PostSolution : User Auth Required> Post Solution
func (u *User) PostSolution(p Problem, text string, description string) {
	s := Solution{
		ProblemID:      p.ID,
		OriginalPoster: *u,
		Text:           text,
		Rank:           0,
	}
	db.Create(s)
}

//FollowProblem : User follows a problem, Add problemID to array
func (u *User) FollowProblem(problemID int) {
	//u.FollowedProblemsIDs = append(u.FollowedProblemsIDs, problemID)
}

// getFollowedProblems : returns problemIDs of all problems followed by the user
//TODO:
//THis doesn't work right
func (u User) getFollowedProblems() []int {
	var followedProblems []int
	err := db.Where("followed_problems = ?").Find(&followedProblems)
	if err == nil {
		glog.Error("Unable to retrieve users followed problems")
	}
	return followedProblems
}

// VoteOnSolution : User votes on a solution to increase it's rank
func (u *User) VoteOnSolution(solutionID int) {
	solution := Solution{}
	solution.GetSolutionByID(solutionID)
	//Check if user has already voted on this problem.
	//if so change vote
	//else add vote
}
