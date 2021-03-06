package DB

//the package should use a database to provide a function, which should also be able to work if simultanoues query happens
//in this version it is just a single file

import "SSClusterManager/executer/UserManager/UserType"
import "os"
import "io/ioutil"
import "encoding/json"
import "log"
import "sync"

type user UserType.User

const filename = "user.db"

var users = map[uint16]user{}
var mutex sync.RWMutex

func init() {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatalln("Reading user database file:", err)
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatalln("Decoding user database file:", err)
	}
}

func Add(u UserType.User) {
	mutex.Lock()
	defer mutex.Unlock()
	defer save()

	users[u.Port] = user(u)
}

func GetAll() []UserType.User {
	mutex.RLock()
	defer mutex.RUnlock()

	v := make([]UserType.User, len(users))
	for _, value := range users {
		v = append(v, UserType.User(value))
	}
	return v
}

func Get(UserPort uint16) (UserType.User, bool) {
	mutex.RLock()
	defer mutex.RUnlock()

	u, exist := users[UserPort]
	return UserType.User(u), exist
}

func Del(UserPort uint16) {
	mutex.Lock()
	defer mutex.Unlock()
	defer save()

	delete(users, UserPort)
}

func save() {
	data, err := json.Marshal(users)
	if err != nil {
		log.Fatalln("Encoding user database:", err)
	}
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		log.Fatalln("Writing to user database:", err)
	}
}
