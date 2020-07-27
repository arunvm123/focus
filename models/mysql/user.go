package mysql

import (
	"github.com/arunvm/focus/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile organises user data and returns it
func (*Mysql) GetProfile(user *models.User) (*models.UserProfile, error) {
	return &models.UserProfile{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		ProfilePic: user.ProfilePic,
	}, nil
}

func (db *Mysql) UpdateProfile(args models.UpdateProfileArgs, user *models.User) error {
	if args.Name != nil {
		user.Name = *args.Name
	}
	if args.ProfilePic != nil {
		user.ProfilePic = args.ProfilePic
	}

	err := db.Client.Save(user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdateProfile",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}

// GetUserFromEmail returns user details from the given email id
func (db *Mysql) GetUserFromEmail(email string) (*models.User, error) {
	var user models.User

	err := db.Client.Find(&user, "email = ?", email).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "GetUserFromEmail",
			"info":  "retrieving user info from email",
			"email": email,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

// GetUserFromID returns user details from the given user id
func (db *Mysql) GetUserFromID(userID int) (*models.User, error) {
	var user models.User

	err := db.Client.Find(&user, "id = ?", userID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":   "GetUserFromID",
			"info":   "retrieving user info from id",
			"userID": userID,
		}).Error(err)
		return nil, err
	}

	return &user, nil
}

func (db *Mysql) CheckIfUserExists(email string) bool {
	var count int
	err := db.Client.Table("users").Where("email = ?", email).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":  "CheckIfUserExists",
			"info":  "checking if user with specified email exitst",
			"email": email,
		}).Error(err)
		return true
	}

	if count > 0 {
		return true
	}

	return false
}

func (db *Mysql) UserSignup(args *models.SignUpArgs, googleOauth bool) (*models.User, string, error) {
	var user models.User

	if !googleOauth {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
		if err != nil {
			log.WithFields(log.Fields{
				"func":    "UserSignup",
				"subFunc": "bcrypt.GenerateFromPassword",
				"email":   args.Email,
			}).Error(err)
			return nil, "", err
		}
		user.Password = string(passwordHash)
	}

	user.Email = args.Email
	user.Name = args.Name
	user.Verified = false
	user.GoogleOauth = googleOauth

	err := db.Client.Create(&user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "UserSignup",
			"subFunc":     "user.Create",
			"email":       args.Email,
			"googleOauth": googleOauth,
		}).Error(err)
		return nil, "", err
	}

	org, err := createPersonalOrganisation(db.Client, &user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "UserSignup",
			"subFunc":     "user.createPersonalOrganisation",
			"email":       args.Email,
			"googleOauth": googleOauth,
		}).Error(err)
		return nil, "", err
	}

	err = createPersonalTeam(db.Client, org, &user)
	if err != nil {
		log.WithFields(log.Fields{
			"func":        "UserSignup",
			"subFunc":     "user.createPersonalTeam",
			"email":       args.Email,
			"googleOauth": googleOauth,
		}).Error(err)
		return nil, "", err
	}

	var token string
	if !googleOauth {
		token, err = emailValidateToken(db.Client, &user)
		if err != nil {
			log.WithFields(log.Fields{
				"func":        "UserSignup",
				"subFunc":     "emailValidateToken",
				"email":       args.Email,
				"googleOauth": googleOauth,
			}).Error(err)
			return nil, "", err
		}
	}

	return &user, token, nil
}

func (db *Mysql) UpdatePassword(args *models.UpdatePasswordArgs, user *models.User) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.CurrentPassword))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "bcrypt.CompareHashAndPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "bcrypt.GenerateFromPassword",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	user.Password = string(hashedPassword)
	err = db.Client.Save(user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "UpdatePassword",
			"subFunc": "user.Save",
			"userID":  user.ID,
		}).Error(err)
		return err
	}

	return nil
}
