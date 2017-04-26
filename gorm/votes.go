package gorm

import "github.com/jinzhu/gorm"
import "github.com/golang/glog"
import "work/xprincipia/backend/util"

//Vote : ~
type Vote struct {
	gorm.Model
	Type     int
	TypeID   int
	Username string
}

//VoteForm : ~
type VoteForm struct {
	Type     int
	TypeID   int
	Username string
}

//CreateVote : ~
func CreateVote(form VoteForm) bool {
	//Create Vote
	v := Vote{}
	v.Type = form.Type
	v.TypeID = form.TypeID
	v.Username = form.Username

	foundVotes := []Vote{}
	db.Where("username = ? AND type_id = ? AND type = ?", v.Username, v.TypeID, v.Type).Find(&foundVotes)
	glog.Info(foundVotes)
	if len(foundVotes) > 0 {
		return false
	}
	db.Create(&v)

	//Change solution rank
	if v.Type == util.SOLUTION {
		s := Solution{}
		s.VoteSolution(v.TypeID)
	} else {
		if v.Type == util.PROBLEM {
			p := Problem{}
			p.VoteProblem(v.TypeID)
		} else {
			if v.Type == util.QUESTION {
				q := Question{}
				q.VoteQuestion(v.TypeID)
			} else {
				if v.Type == util.SUGGESTION {
					s := Suggestion{}
					s.VoteSuggestion(v.TypeID)
				} else {
					if v.Type == util.ANSWER {
						a := Answer{}
						a.VoteAnswer(v.TypeID)
					} else {
						if v.Type == util.COMMENT {
							c := Comment{}
							c.VoteComment(v.TypeID)

						}
					}
				}
			}
		}
	}
	return true

}

//GetNumberOfVotesByTypeID : !
func GetNumberOfVotesByTypeID(id int) int {
	v := []Vote{}
	db.Where("type_id = ?", id).Find(&v)
	return len(v)
}
