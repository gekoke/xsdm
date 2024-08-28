package main

import (
	"fmt"
	gopwd "github.com/Maki-Daisuke/go-pwd"
	pam "github.com/msteinert/pam"
	"log"
)

type conversationHandler struct {
	username string
	password string
	authInfo *authInfo
}

func (conversationHandler conversationHandler) RespondPAM(style pam.Style, str string) (string, error) {
	switch style {
	case pam.PromptEchoOn: // get username
		return conversationHandler.username, nil
	case pam.PromptEchoOff: // get password
		return conversationHandler.password, nil
	case pam.ErrorMsg:
		conversationHandler.authInfo.infos = append(conversationHandler.authInfo.infos, str)
		return "", nil
	case pam.TextInfo:
		conversationHandler.authInfo.errors = append(conversationHandler.authInfo.errors, str)
		return "", nil
	case pam.BinaryPrompt:
		panic("BinaryPrompt unimplemented")
	default:
		panic("unreachable")
	}
}

func login(username string, password string, authInfo *authInfo) error {
	handler := conversationHandler{
		username: username,
		password: password,
		authInfo: authInfo,
	}

	transaction, err := pam.Start("xsdm", username, handler)
	if err != nil {
		log.Printf("user %s: failed to start PAM transaction: %s", username, err)
		return err
	}

	err = transaction.Authenticate(0)
	if err != nil {
		log.Printf("user %s: failed to authenticate user: %s", username, err)
		return err
	}

	err = transaction.AcctMgmt(0)
	if err != nil {
		log.Printf("user %s: failed to determine account validity: %s", username, err)
		return err
	}

	err = transaction.SetCred(0)
	if err != nil {
		log.Printf("user %s: failed to set user's credentials: %s", username, err)
		return err
	}

	err = transaction.OpenSession(0)
	if err != nil {
		log.Printf("user %s: failed to open session for user: %s", username, err)
		credErr := transaction.SetCred(pam.DeleteCred)
		if credErr != nil {
			log.Printf("user %s: failed to delete user's credentials: %s", username, credErr)
		}
		return err
	}

	user := gopwd.Getpwnam(username)

	putPAMEnv("HOME", user.Dir, transaction)
	putPAMEnv("PWD", user.Dir, transaction)

	putPAMEnv("NAME", user.Name, transaction)
	putPAMEnv("LOGNAME", user.Name, transaction)

	putPAMEnv("SHELL", user.Shell, transaction)

	return nil
}

func putPAMEnv(key string, value string, transaction *pam.Transaction) {
	err := transaction.PutEnv(fmt.Sprintf("%s=%s", key, value))
	if err != nil {
		log.Printf("error setting PAM environment variable %s=%s: %s", key, value, err)
	}
}
