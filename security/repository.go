package security

import (
	"fmt"
	"strings"

	"github.com/nelsonlai-golang/db-util/sqlite"
	"gorm.io/gorm"
)

/*
	This file is the repository methods of the entities in the security package.
	These methods are NOT supposed to expose outside of the package.
*/

//// SecurityPath ////

func connectDB() *gorm.DB {
	return sqlite.Memory()
}

func autoMigrateSecurityPath() {
	db := connectDB()
	db.AutoMigrate(&SecurityPath{})
}

func deleteSecurityPath() {
	db := connectDB()
	err := db.Unscoped().Where("id IS NOT NULL").Delete(&SecurityPath{}).Error
	if err != nil {
		panic(err)
	}
}

func createSecurityPath(path SecurityPath) {
	db := connectDB()
	db.Create(&path)
}

func findPotentialPaths(method string, path string) []SecurityPath {
	fragments := strings.Split(strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/"), "/")

	var build, query string
	for _, fragment := range fragments {
		build += "/" + fragment
		query += "path_regex LIKE '" + build + "%' OR "
	}
	query = strings.TrimSuffix(query, " OR ")

	db := connectDB()
	var paths []SecurityPath
	db.Where(fmt.Sprintf("method = ? AND (%s)", query), method).Find(&paths)
	return paths
}

//// SecuritySession ////

func autoMigrateSecuritySession() {
	db := connectDB()
	db.AutoMigrate(&SecuritySession{})
}

func createSession(session *SecuritySession) {
	db := connectDB()
	db.Create(session)
}

func updateSession(session *SecuritySession) {
	db := connectDB()
	db.Save(session)
}

func deleteSessionById(id string) {
	db := connectDB()
	db.Unscoped().Delete(&SecuritySession{}, "session_id = ?", id)
}

func findSessionById(id string) (*SecuritySession, error) {
	db := connectDB()
	var session SecuritySession
	db.Where("session_id = ?", id).First(&session)
	return &session, nil
}

func findSessionByUserId(userId uint) (*SecuritySession, error) {
	db := connectDB()
	var session SecuritySession
	db.Where("user_id = ?", userId).First(&session)
	return &session, nil
}
