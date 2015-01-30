// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v1

import (
	api "github.com/Dumbris/go-gogs-client"

	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/middleware"
	//"github.com/gogits/gogs/modules/setting"
)

func SettingsSSHKeys(ctx *middleware.Context) {

	//var err error
	var publicKeys, err = models.ListPublicKeys(ctx.User.Id)

	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	keys := make([]*api.Key, len(publicKeys))

	for i := range publicKeys {
		k := &api.Key{
			ID:    int(publicKeys[i].Id),
			Key:   publicKeys[i].Content,
			Title: publicKeys[i].Name,
		}
		keys[i] = k
	}

	ctx.JSON(200, &keys)
}

//func SettingsSSHKeysPost(ctx *middleware.Context, form auth.AddSSHKeyForm) {
//	ctx.Data["Title"] = ctx.Tr("settings")
//	ctx.Data["PageIsUserSettings"] = true
//	ctx.Data["PageIsSettingsSSHKeys"] = true

//	var err error
//	ctx.Data["Keys"], err = models.ListPublicKeys(ctx.User.Id)
//	if err != nil {
//		ctx.Handle(500, "ssh.ListPublicKey", err)
//		return
//	}

//	// Delete SSH key.
//	if ctx.Query("_method") == "DELETE" {
//		id := com.StrTo(ctx.Query("id")).MustInt64()
//		if id <= 0 {
//			return
//		}

//		if err = models.DeletePublicKey(&models.PublicKey{Id: id}); err != nil {
//			ctx.Handle(500, "DeletePublicKey", err)
//		} else {
//			log.Trace("SSH key deleted: %s", ctx.User.Name)
//			ctx.Redirect(setting.AppSubUrl + "/user/settings/ssh")
//		}
//		return
//	}

//	// Add new SSH key.
//	if ctx.Req.Method == "POST" {
//		if ctx.HasError() {
//			ctx.HTML(200, SETTINGS_SSH_KEYS)
//			return
//		}

//		// Remove newline characters from form.KeyContent
//		cleanContent := strings.Replace(form.Content, "\n", "", -1)

//		if ok, err := models.CheckPublicKeyString(cleanContent); !ok {
//			if err == models.ErrKeyUnableVerify {
//				ctx.Flash.Info(ctx.Tr("form.unable_verify_ssh_key"))
//			} else {
//				ctx.Flash.Error(ctx.Tr("form.invalid_ssh_key", err.Error()))
//				ctx.Redirect(setting.AppSubUrl + "/user/settings/ssh")
//				return
//			}
//		}

//		k := &models.PublicKey{
//			OwnerId: ctx.User.Id,
//			Name:    form.SSHTitle,
//			Content: cleanContent,
//		}
//		if err := models.AddPublicKey(k); err != nil {
//			if err == models.ErrKeyAlreadyExist {
//				ctx.RenderWithErr(ctx.Tr("form.ssh_key_been_used"), SETTINGS_SSH_KEYS, &form)
//				return
//			}
//			ctx.Handle(500, "ssh.AddPublicKey", err)
//			return
//		} else {
//			log.Trace("SSH key added: %s", ctx.User.Name)
//			ctx.Flash.Success(ctx.Tr("settings.add_key_success"))
//			ctx.Redirect(setting.AppSubUrl + "/user/settings/ssh")
//			return
//		}
//	}

//	ctx.HTML(200, SETTINGS_SSH_KEYS)
//}
