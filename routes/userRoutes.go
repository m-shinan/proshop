package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-shinan/project-shop/controllers"
	"github.com/m-shinan/project-shop/middleware"
)

func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")
	{
		user.GET("/signup", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user_signup.html", nil)
		})
		user.POST("/signup", controllers.UserSignUp)

		user.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user_login.html", nil)
		})

		user.POST("/login", controllers.UserLogin)

		user.GET("/uservalidate", middleware.UserAuth, controllers.ValidateUser)

		user.GET("/userHome", middleware.UserAuth, controllers.UserHome)

		user.GET("/logout", middleware.UserAuth, controllers.UserLogout)

		///////////////////  CART   //////////////////////////////

		user.GET("/cart", middleware.UserAuth, controllers.Cart)

		user.POST("/cart/add/:id", middleware.UserAuth, controllers.AddToCart)

		user.POST("/cart/remove/:id", middleware.UserAuth, controllers.RemoveFromCart)

		user.POST("/cart/plus/:id", middleware.UserAuth, controllers.PlusQuantity)

		user.POST("/cart/minus/:id", middleware.UserAuth, controllers.MinusQuantity)

		////////////////// WISHLIST //////////////////////

		user.GET("/wihslist", middleware.UserAuth, controllers.Wishlist)

		user.POST("/wishlist/add/:id", middleware.UserAuth, controllers.AddToWishlist)

		user.POST("/wishlist/remove/:id", middleware.UserAuth, controllers.RemoveFromWishlist)

	}

}
