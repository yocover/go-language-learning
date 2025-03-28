// AuthMiddleware 权限检查中间件
func (api *DocumentAPI) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.GetHeader("X-User-ID")
		domain := c.GetHeader("X-Domain")
		if user == "" || domain == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证或未指定域"})
			c.Abort()
			return
		}

		// 构建权限检查请求
		path := c.Request.URL.Path
		method := c.Request.Method

		// 获取文档所有者（如果适用）
		var owner string
		if docID := c.Param("id"); docID != "" {
			if doc, exists := api.docs[docID]; exists {
				owner = doc.OwnerID
			}
		}

		// 如果是当前用户的文档，设置 owner 为 "self"
		if owner == user {
			owner = "self"
		}

		// 检查权限（加入域参数和所有者）
		ok, err := api.enforcer.Enforce(user, domain, path, method, owner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限"})
			c.Abort()
			return
		}

		// 将域信息存储在上下文中
		c.Set("domain", domain)
		c.Set("user_id", user)
		c.Next()
	}
}

// CreateDocument 创建新文档
func (api *DocumentAPI) CreateDocument(c *gin.Context) {
	domain := c.MustGet("domain").(string)
	userID := c.MustGet("user_id").(string)

	var doc Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置文档类型和域
	docType := "public"
	if strings.Contains(c.Request.URL.Path, "/private/") {
		docType = "private"
	}
	doc.Type = docType
	doc.Domain = domain
	doc.OwnerID = userID // 设置文档所有者

	api.docs[doc.ID] = &doc
	c.JSON(http.StatusCreated, doc)
}

// UpdateDocument 更新文档
func (api *DocumentAPI) UpdateDocument(c *gin.Context) {
	domain := c.MustGet("domain").(string)
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")

	doc, exists := api.docs[id]
	if !exists || doc.Domain != domain {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		return
	}

	// 检查文档所有权
	owner := "other"
	if doc.OwnerID == userID {
		owner = "self"
	}

	// 验证用户是否有权限更新文档
	ok, err := api.enforcer.Enforce(userID, domain, c.Request.URL.Path, "PUT", owner)
	if err != nil || !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文档"})
		return
	}

	var updatedDoc Document
	if err := c.ShouldBindJSON(&updatedDoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedDoc.ID = id
	updatedDoc.Type = doc.Type       // 保持原有类型
	updatedDoc.Domain = domain       // 保持原有域
	updatedDoc.OwnerID = doc.OwnerID // 保持原有所有者
	api.docs[id] = &updatedDoc

	c.JSON(http.StatusOK, updatedDoc)
}

// DeleteDocument 删除文档
func (api *DocumentAPI) DeleteDocument(c *gin.Context) {
	domain := c.MustGet("domain").(string)
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")

	doc, exists := api.docs[id]
	if !exists || doc.Domain != domain {
		c.JSON(http.StatusNotFound, gin.H{"error": "文档不存在"})
		return
	}

	// 检查文档所有权
	owner := "other"
	if doc.OwnerID == userID {
		owner = "self"
	}

	// 验证用户是否有权限删除文档
	ok, err := api.enforcer.Enforce(userID, domain, c.Request.URL.Path, "DELETE", owner)
	if err != nil || !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文档"})
		return
	}

	delete(api.docs, id)
	c.Status(http.StatusNoContent)
} 