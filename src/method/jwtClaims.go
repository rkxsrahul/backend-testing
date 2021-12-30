package method

import (
	"fmt"

	"git.xenonstack.com/util/continuous-security-backend/config"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

//ExtractClaims function for extract claims
func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {

	// parsing token and checking its validity
	rtoken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Conf.JWT.PrivateKey), nil
	})
	// if any err return nil claims
	if err != nil {
		return nil, err
	}
	claims := rtoken.Claims.(jwt.MapClaims)
	return claims, nil

}
