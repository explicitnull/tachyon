package applogic

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"tacacs-webconsole/options"
	"tacacs-webconsole/repository"
	"tacacs-webconsole/types"

	"github.com/aerospike/aerospike-client-go"
	"github.com/dchest/uniuri"
	"github.com/sirupsen/logrus"
)

var (
	ErrPasswordsMismatch = errors.New("passwords mismatch")
	ErrBadCharacters     = errors.New("bad characters in password")
	ErrTooShort          = errors.New("too short password")
)

func CreateUserAction(le *logrus.Entry, aClient *aerospike.Client, req types.Account, authenticatedUsername string) (string, error) {
	ctx := context.TODO()
	// normalization
	subdivID, err := repository.GetSubdivisionID(le, aClient, req.Subdivision)
	if err != nil {
		le.WithError(err).Error("getting subdivision ID failed")
		return "", err
	}

	permisID, err := repository.GetPermId(le, aClient, req.Permission)
	if err != nil {
		le.WithError(err).Error("getting permission ID failed")
		return "", err
	}

	cleartext := genPass()
	le.Debug(cleartext)
	hash := MakeHash(le, cleartext)
	le.Debug(hash)

	err = repository.CreateUser(ctx, aClient, req.Name, hash, req.Mail, authenticatedUsername, permisID, subdivID)
	if err != nil {
		le.WithError(err).Errorf("error creating user")
		return "", err
	}

	le.WithField("username", req.Name).Info("user created")

	return cleartext, nil
}

func ChangePasswordAction(le *logrus.Entry, aclient *aerospike.Client, authenticatedUsername, pass, passConfirm string, o *options.Options) error {
	// checking if passwords don't match
	if pass != passConfirm {
		return ErrPasswordsMismatch
	}

	// checking if password has not [[:graph:]] symbols
	ok, _ := regexp.MatchString("^[[:graph:]]+$", pass)
	if !ok {
		return ErrBadCharacters
	}

	// checking password length
	CleanMap := make(map[string]interface{})
	CleanMap["pass"] = pass
	if len(pass) < o.MinPassLen {
		return ErrTooShort
	}

	// changing password
	hash := MakeHash(le, pass)

	ctx := context.TODO()
	err := repository.SetPassword(ctx, aclient, authenticatedUsername, hash)
	if err != nil {
		le.WithError(err).Error("password update failed")
		return err
	}

	le.Info("user password updated")

	/*
		// Set flag "password changed"
		stmt, err2 := db.Prepare("update usr set pass_chd='true' where username=$1")
		checkErr(err2)
		defer stmt.Close()

		res, err3 := stmt.Exec(authUser)
		checkErr(err3)

		affect, err4 := res.RowsAffected()
		checkErr(err4)

		log.Println(affect, `rows changed while setting "password changed" flag for`, authUser)
	*/

	// activating user
	err = repository.SetAccountStatus(ctx, aclient, authenticatedUsername, "active")
	if err != nil {
		le.WithError(err).Error("user activation failed")
		return err
	}

	le.Info("user status switched to active due to password update")

	return nil
}

func EditAccountAction(le *logrus.Entry, aclient *aerospike.Client, fac *types.Account) error {
	// getting account data from DB
	ctx := context.TODO()
	dbac, err := repository.GetAccountByName(ctx, aclient, fac.Name)
	if err != nil {
		le.WithError(err).Error("getting account failed")
		return err
	}

	// applying changes
	if fac.Cleartext != "" {
		hash := MakeHash(le, fac.Cleartext)
		err = repository.SetPassword(ctx, aclient, fac.Name, hash)
		if err != nil {
			return err
		}
	}

	if fac.Subdivision != "--" && fac.Subdivision != dbac.Subdivision {
		subdivID, err := strconv.Atoi(fac.Subdivision)
		if err != nil {
			return err
		}

		err = repository.SetSubdivision(ctx, aclient, fac.Name, subdivID)
		if err != nil {
			return err
		}
	}

	if fac.Permission != "--" && fac.Permission != dbac.Permission {
		permisID, err := strconv.Atoi(fac.Permission)
		if err != nil {
			return err
		}

		err = repository.SetPermission(ctx, aclient, fac.Name, permisID)
		if err != nil {
			return err
		}
	}

	if fac.Mail != "" && fac.Mail != dbac.Mail {
		err = repository.SetMail(ctx, aclient, fac.Name, fac.Mail)
		if err != nil {
			return err
		}
	}

	if fac.Status != dbac.Status {
		err = repository.SetAccountStatus(ctx, aclient, fac.Name, fac.Status)
		if err != nil {
			return err
		}
	}

	if fac.UILevel != dbac.UILevel {
		err = repository.SetUILevel(ctx, aclient, fac.Name, fac.UILevel)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveAccount(le *logrus.Entry, aclient *aerospike.Client, acname, authenticatedUsername string) error {
	ctx := context.TODO()
	err := repository.DeleteAccount(ctx, aclient, acname, authenticatedUsername)
	if err != nil {
		return nil
	}

	return nil
}

func genPass() string {
	return uniuri.NewLen(10)
}

// MakeHash generates SHA hashes for given passwords
func MakeHash(le *logrus.Entry, cleartext string) string {
	hash := sha256.Sum256([]byte(cleartext))
	le.Debug(hash)

	enc := base64.StdEncoding.EncodeToString(hash[:])
	return strings.Replace(enc, "=", "", -1)
}
