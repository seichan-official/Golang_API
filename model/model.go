package models
type User struct {
    gorm.Model
    SpotifyID    string 
    DisplayName  string
    Email        string
}

type Music struct{
    
}

//DB初期化
func dbInit() {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbInit）")
    }
    db.AutoMigrate(&User{})
    defer db.Close()
}

//DB追加
func dbInsert(SpotifyID string, DisplayName string, Email string), {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        panic("データベース開けず！（dbInsert)")
    }
    db.Create(&User{SpotifyID: spotifyID, DisplayName: displayName, Email: email})
    defer db.Close()
}

//DB更新
// func dbUpdate(id int, token string) {
//     db, err := gorm.Open("sqlite3", "test.sqlite3")
//     if err != nil {
//         panic("データベース開けず！（dbUpdate)")
//     }
//     var user User
//     db.First(&user, id)
//     todo.Token = token
    
//     db.Save(&user)
//     db.Close()
// }

//DB削除
// func dbDelete(id int) {
//     db, err := gorm.Open("sqlite3", "test.sqlite3")
//     if err != nil {
//         panic("データベース開けず！（dbDelete)")
//     }
//     var user User
//     db.First(&user, id)
//     db.Delete(&user)
//     db.Close()
// }

//DB全取得
// func dbGetAll() []Todo {
//     db, err := gorm.Open("sqlite3", "test.sqlite3")
//     if err != nil {
//         panic("データベース開けず！(dbGetAll())")
//     }
//     var todos []Todo
//     db.Order("created_at desc").Find(&todos)
//     db.Close()
//     return todos
// }

//DB一つ取得
// func dbGetOne(id int) Todo {
//     db, err := gorm.Open("sqlite3", "test.sqlite3")
//     if err != nil {
//         panic("データベース開けず！(dbGetOne())")
//     }
//     var todo Todo
//     db.First(&todo, id)
//     db.Close()
//     return todo
// }

func dbGetUserByID(id uint) (User, error) {
    db, err := gorm.Open("sqlite3", "test.sqlite3")
    if err != nil {
        return User{}, err
    }
    var user User
    result := db.First(&user, id)
    defer db.Close()
    return user, result.Error
}
