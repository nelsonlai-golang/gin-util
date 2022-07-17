package security

import (
	"strings"

	"github.com/nelsonlai-golang/db-util/sqlite"
)

/*
	This file is the repository methods of the entities in the security package.
	These methods are NOT supposed to expose outside of the package.
*/

//// SecurityPath ////

func autoMigrateSecurityPath() {
	db := sqlite.Memory()
	db.AutoMigrate(&SecurityPath{})
}

func deleteSecurityPathTable() {
	db := sqlite.Memory()
	db.Unscoped().Delete(&SecurityPath{})
}

func createSecurityPath(path SecurityPath) {
	db := sqlite.Memory()
	db.Create(&path)
}

func findPotentialPaths(method string, path string) []SecurityPath {
	fragments := strings.Split(strings.TrimSuffix(path, "/"), "/")

	var build, query string
	for _, fragment := range fragments {
		build += "/" + fragment
		query += "path LIKE '" + build + "%' OR "
	}
	query = strings.TrimSuffix(query, " OR ")

	db := sqlite.Memory()
	var paths []SecurityPath
	db.Where("method = ? AND (?)", method, query).Find(&paths)
	return paths
}

//// SecuritySession ////

func autoMigrateSecuritySession() {
	db := sqlite.Memory()
	db.AutoMigrate(&SecuritySession{})
}

func createSession(session *SecuritySession) {
	db := sqlite.Memory()
	db.Create(session)
}

func updateSession(session *SecuritySession) {
	db := sqlite.Memory()
	db.Save(session)
}

func deleteSessionById(id string) {
	db := sqlite.Memory()
	db.Unscoped().Delete(&SecuritySession{}, "session_id = ?", id)
}

func findSessionById(id string) (*SecuritySession, error) {
	db := sqlite.Memory()
	var session SecuritySession
	db.Where("session_id = ?", id).First(&session)
	return &session, nil
}

func findSessionByUserId(userId uint) (*SecuritySession, error) {
	db := sqlite.Memory()
	var session SecuritySession
	db.Where("user_id = ?", userId).First(&session)
	return &session, nil
}
