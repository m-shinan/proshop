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

		user.GET("/verifyotp", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "otp.html", nil)
		})

		user.POST("/verify-otp", controllers.VerifyOTP)

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

		user.GET("/wishlist", middleware.UserAuth, controllers.Wishlist)

		user.POST("/wishlist/add/:id", middleware.UserAuth, controllers.AddToWishlist)

		user.POST("/wishlist/remove/:id", middleware.UserAuth, controllers.RemoveFromWishlist)

		/////////////////// CHECKOUT ////////////////////////

		user.GET("/checkout", middleware.UserAuth, controllers.Checkout)

		//////////////////// PROFILE ////////////////////////

		user.GET("/profile", middleware.UserAuth, controllers.UserProfile)

		user.POST("/profile/edit", middleware.UserAuth, controllers.EditUserProfile)

		///////////////// PASSWORD ////////////////////

		user.POST("/change/password", middleware.UserAuth, controllers.UserChangePassword)

		///////////////// ADDRESS /////////////////

		user.POST("/address/add", middleware.UserAuth, controllers.AddAddress)

		user.GET("/address/edit/:id", middleware.UserAuth, controllers.EditAddressPage)

		user.POST("/address/edit/:id", middleware.UserAuth, controllers.EditAddress)

		user.GET("/address/delete/:id", middleware.UserAuth, controllers.DeleteAddress)

		////////////////// PAYMENT ////////////////////////

		user.POST("/razorpay", middleware.UserAuth, controllers.CreateOrderRazorPay)

		user.POST("/payment/razorpay", middleware.UserAuth, controllers.ConfirmRazorpayPayment)

		user.POST("/payment/cod", middleware.UserAuth, controllers.CreateOrderCod)

		///////////////// ORDERS /////////////////////

		user.GET("/orders", middleware.UserAuth, controllers.UserOrders)

		user.POST("/orders/cancel", middleware.UserAuth, controllers.CancelOrder)

	}

}
