package integration_test

const (
	UserEmail 				 = "t@t.com"
	UserPass  				 = "password123"
	EmulatorHost 			 = "localhost"
	EmulatorPort 			 = 9099
	ResourceServerPort = 8088
	ResourceServerHost = "localhost"
	ApiKey 						 = "fake-api-key"
	IsSecure    			 = false
)

type ReturnedToken struct {
	Kind         string `json:"kind"`
	IsNewUser    bool   `json:"isNewUser"`
	IdToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpirationIn string `json:"expiresIn"`
}
