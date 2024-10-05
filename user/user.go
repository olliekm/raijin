package user

type User struct {
	Username    string `bson:"username" json:"username"`
	DisplayName string `bson:"display_name" json:"display_name"`
	Password    string `bson:"password" json:"password"`
	ProfilePic  string `bson:"profile_pic" json:"profile_pic"` // URL or path to the profile picture
}
