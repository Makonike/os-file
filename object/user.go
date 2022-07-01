package object

type User struct {
	username string
	password string
}

func (u *User) NewUser(username, password string) *User {
	return &User{
		username: username,
		password: password,
	}
}

func Login(username, password string) *Result {
	if MMemory.CurUser != nil {
		return FailMsg("[user login fail]: plz logout first")
	}

	u, ok := DCache.Disk.UMap[username]
	if !ok || u == nil {
		return FailMsg("[user login fail]: username not found")
	}
	if u.password != password {
		return FailMsg("[user login fail]: wrong password")
	}

	for _, v := range DCache.Disk.Dirs[0].Children {
		if v.Fcb.IsDir() && v.Fcb.Name == username {
			MMemory.CurUser = u
			return SuccessResult()
		}
	}
	// TODO: lazy load to fix error
	return FailMsg("[user login fail]: internal server error, not found the user dir")
}

func Logout() *Result {
	if MMemory.CurUser == nil {
		return FailMsg("[user logout fail]: no user logged now")
	}
	MMemory.CurUser = nil
	// change curDir
	MMemory.CurDir = DCache.Disk.Dirs[0]
	return SuccessResult()
}
