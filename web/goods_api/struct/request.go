/**
    @auther: oreki
    @date: 2022年09月08日 12:15 AM
    @note: 图灵老祖保佑,永无BUG
**/

package _struct

import (
	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId uint
	jwt.StandardClaims
}
