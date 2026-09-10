package sessiontoken

type AccessToken struct {
	/**
	 * 访问token
	 */
	accessToken string

	/**
	 * 刷新token
	 */
	refreshToken string

	/**
	 * 访问token的生效时间(秒)
	 */
	accessTokenValidity int64

	/**
	 * 刷新token的生效时间(秒)
	 */
	refreshTokenValidity int64
}
