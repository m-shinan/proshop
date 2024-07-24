package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/controllers"
	"github.com/m-shinan/project-shop/middleware"
)

func AdminRoutes(c *gin.Engine) {
	admin := c.Group("/admin")
	{
		admin.GET("/signup", func(c *gin.Context) {
			c.HTML(http.StatusOK, "adminSignup.html", nil)
		})
		admin.POST("/signup", controllers.AdminSignUp)

		admin.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_login.html", nil)
		})

		admin.POST("/login", controllers.AdminLogin)

		admin.GET("/adminvalidate", middleware.AdminAuth, controllers.ValidateAdmin)

		admin.GET("/adminHome", middleware.AdminAuth, func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_home.html", nil)
		})

		admin.GET("/logout", middleware.AdminAuth, controllers.AdminLogout)

		//////////////////////////////product management

		admin.GET("/products", middleware.AdminAuth, controllers.AdminProducts)

		admin.GET("/products/add", middleware.AdminAuth, controllers.AdminAddProductPage)

		admin.POST("/products/add", middleware.AdminAuth, controllers.AdminAddProduct)

		admin.GET("/products/edit/:id", middleware.AdminAuth, controllers.AdminEditProductPage)

		admin.POST("/products/edit/:id", middleware.AdminAuth, controllers.AdminEditProduct)

		admin.POST("/products/delete/:id", middleware.AdminAuth, controllers.AdminDeleteProduct)

		//////////////////////////////category management

		admin.GET("/categories/add", middleware.AdminAuth, func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_add_cat.html", nil)
		})

		admin.POST("/categories/add", middleware.AdminAuth, controllers.AdminAddCategory)

		admin.GET("/categories/edit/:id", middleware.AdminAuth, controllers.AdminEditcatPage)

		admin.POST("/categories/edit/:id", middleware.AdminAuth, controllers.AdminEditCat)

		admin.POST("/categories/delete/:id", middleware.AdminAuth, controllers.AdminDeleteCat)

		/////////////////////////////User management

		admin.GET("/users", middleware.AdminAuth, controllers.AdminViewUsers)

		admin.POST("/users/block/:id", middleware.AdminAuth, controllers.AdminBlockUsers)

		admin.POST("/users/unblock/:id", middleware.AdminAuth, controllers.AdminUnblockUsers)

		admin.POST("/users/delete/:id", middleware.AdminAuth, controllers.AdminDeleteUsers)

	}
}

// admin.POST("/login", controllers.AdminLogin())
// admin.GET("/logout", controllers.AdminLogout())
// admin.GET("/dashboard", controllers.AdminDashboard())

// admin user controlls

// admin.GET("/viewUsers", controllers.ViewUsers())
// admin.GET("/searchUsers", controllers.searchUsers())
// admin.GET("/deleteUser", controllers.DeleteUser())
// admin.GET("/blockUser", controllers.blockUser())
// admin.GET("/unblockUser", controllers.unblockUser())
