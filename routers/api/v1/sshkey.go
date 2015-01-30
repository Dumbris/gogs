// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/Unknwon/com"
	api "github.com/gogits/go-gogs-client"
	"github.com/gogits/gogs/models"
	"github.com/gogits/gogs/modules/log"
	"github.com/gogits/gogs/modules/middleware"
	"strings"
	//"github.com/gogits/gogs/modules/setting"
)

func ListSSHKeys(ctx *middleware.Context) {

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

func GetSSHKey(ctx *middleware.Context) {

	id := com.StrTo(ctx.Query("id")).MustInt64()
	if id <= 0 {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": "id required",
		})
		return
	}
	//var err error
	var publicKey, err = models.GetPublicKeyById(id)

	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, &api.Key{
		ID:    int(publicKey.Id),
		Key:   publicKey.Content,
		Title: publicKey.Name,
	})
}

func PostSSHKey(ctx *middleware.Context, form api.AddSSHKeyOption) {
	var err error
	// Delete SSH key.
	if ctx.Query("_method") == "DELETE" {
		id := com.StrTo(ctx.Query("id")).MustInt64()
		if id <= 0 {
			return
		}

		if err = models.DeletePublicKey(&models.PublicKey{Id: id}); err != nil {
			ctx.JSON(500, map[string]interface{}{
				"ok":    false,
				"error": err.Error(),
			})
		} else {
			log.Trace("SSH key deleted: %s", ctx.User.Name)
			ctx.JSON(204, nil)
		}
		return
	}

	// Add new SSH key.
	if ctx.Req.Method == "POST" {
		if ctx.HasError() {
			ctx.JSON(500, map[string]interface{}{
				"ok":    false,
				"error": "has error",
			})
			return
		}

		// Remove newline characters from form.KeyContent
		cleanContent := strings.Replace(form.Key, "\n", "", -1)

		if ok, err := models.CheckPublicKeyString(cleanContent); !ok {
			if err == models.ErrKeyUnableVerify {
				ctx.JSON(500, map[string]interface{}{
					"ok":    false,
					"error": ctx.Tr("form.unable_verify_ssh_key"),
				})
				return
			} else {
				ctx.JSON(500, map[string]interface{}{
					"ok":    false,
					"error": ctx.Tr("form.invalid_ssh_key", err.Error()),
				})
				return
			}
		}

		k := &models.PublicKey{
			OwnerId: ctx.User.Id,
			Name:    form.Title,
			Content: cleanContent,
		}
		if err := models.AddPublicKey(k); err != nil {
			ctx.JSON(500, map[string]interface{}{
				"ok":    false,
				"error": err.Error(),
			})

			return
		} else {
			log.Trace("SSH key added: %s", ctx.User.Name)
			ctx.JSON(201, map[string]interface{}{
				"ok": true,
			})
			return
		}
	}

	ctx.JSON(500, map[string]interface{}{
		"ok":    false,
		"error": "invalid request",
	})
}
