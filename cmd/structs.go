package cmd

type User struct {
	ID       int
	Image    string
	Username string
	Email    string
	Password string
	Role     string
	Ban      int
	Report   string
}

type Post struct {
	ID          int
	Photo       string
	Title       string
	Texte       string
	Hidden      int
	Like        int
	Dislike     int
	Signalement int
	Categorie   string
	Ban         int
	Archived    string
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	SaveInfo string `json:"saveinfo"`
	JwtToken string `json:"jwtToken"`
}

type AdminPanel struct {
	Account         []map[string]interface{} `json:"account"`
	AccountReported []map[string]interface{} `json:"accountReported"`
	Ban             []map[string]interface{} `json:"ban"`
	PostHidden      []map[string]interface{} `json:"postHidden"`
	PostArchived    []map[string]interface{} `json:"postArchived"`
}

type GestionPost struct {
	Post string `json:"post"`
}

type AdminPanelChange struct {
	Key           string `json:"key"`
	UnBanUser     string `json:"unban-user"`
	BanUser       string `json:"ban-user"`
	RoleAdminUser string `json:"role-admin-user"`
	RoleModoUser  string `json:"role-modo-user"`
	DeletePost    string `json:"delete-post"`
}

type Register struct {
	Nom      string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type responseRegister struct {
	Message string `json:"message"`
}
