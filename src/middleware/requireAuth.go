package middleware

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"fmt"
	"time"
	"go-jwt/initializers"
	"go-jwt/model"
)

func RequireAuth(c *gin.Context){
	// Get the cookie off req

	tokenString,err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Decode/validate it


// Parse takes the token string and a function for looking up the key. The latter is especially
// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
// head of the token to identify which key to use, but the parsed token (head and claims) is provided
// to the callback, providing flexibility.
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte(os.Getenv("SECRET")),nil
})

if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	
	// check the exp

	if float64(time.Now().Unix()) > claims["exp"].(float64){
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// find the user with token sub

	var user model.User
	initializers.DB.First(&user,claims["sub"])

	if user.ID == 0{
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Attach to req
     c.Set("user",user)

	//Continue

	c.Next()
} else {
	c.AbortWithStatus(http.StatusUnauthorized)
}


}