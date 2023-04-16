package auth

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/fangbc5/gogo/core/db"
	"log"
	"sync"
)

var auth *casbin.Enforcer
var lock sync.Mutex

func MakeEnforcer() *casbin.Enforcer {
	if auth == nil {
		lock.Lock()
		if auth == nil {
			auth = load()
		}
		lock.Unlock()
	}
	return auth
}

func load() *casbin.Enforcer {

	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	a, err := gormadapter.NewAdapterByDB(db.MysqlClient()) // Your driver and data source.
	if err != nil {
		lock.Unlock()
		log.Panicln(err)
	}
	e, err := casbin.NewEnforcer("core/config/rbac_model.conf", a)
	if err != nil {
		lock.Unlock()
		log.Panicln(err)
	}
	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	err = e.LoadPolicy()
	if err != nil {
		lock.Unlock()
		log.Panicln(err)
	}
	// Check the permission.
	//e.Enforce("alice", "data1", "read")
	return e

	// Modify the policy.
	// e.AddPolicy(...)
	// e.RemovePolicy(...)

	// Save the policy back to DB.
	//e.SavePolicy()
}

func AddPolicy(params ...interface{}) {
	auth.AddPolicy(params)
	auth.SavePolicy()
}

func AddPolicies(rules [][]string) {
	auth.AddPolicies(rules)
	auth.SavePolicy()
}

func RemovePolicy(params ...interface{}) {
	auth.RemovePolicy(params)
	auth.SavePolicy()
}

func RemovePolicies(rules [][]string) {
	auth.RemovePolicy(rules)
	auth.SavePolicy()
}
