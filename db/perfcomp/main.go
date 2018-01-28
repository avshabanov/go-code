package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/avshabanov/go-code/db/perfcomp/logic"
	"github.com/avshabanov/go-code/fixture"
)

type newDao func(openParams string) (logic.Dao, error)

const sqliteDaoType = "sqlite"

var (
	dbPath = flag.String("db-path", "/tmp/perftest-sqlite.db", "Path to identity service database")
	dbType = flag.String("db-type", sqliteDaoType, "Type of the database to test")

	daoFactoryMap = map[string]newDao{
		sqliteDaoType: logic.NewSqliteDao,
	}
)

func main() {
	flag.Parse()

	newDaoFunc := daoFactoryMap[*dbType]
	if newDaoFunc == nil {
		log.Fatalf("unknown dao type: %s", *dbType)
		return
	}

	dao, err := newDaoFunc(*dbPath)
	if err != nil {
		log.Fatalf("cannot create dao: %v", err)
		return
	}
	defer dao.Close()

	userPage, err := dao.QueryUsers("", 1)
	if err != nil {
		log.Fatalf("cannot get user profiles: %v", err)
		return
	}
	if len(userPage.Profiles) == 0 {
		// insert fixture
		if err = dao.Add(getUserFixture(10, 1)); err != nil {
			log.Fatalf("cannot add user profile: %v", err)
			return
		}
	}

	userPage, err = dao.QueryUsers("", 10)
	if err != nil {
		log.Fatalf("cannot get user profiles: %v", err)
		return
	}
	for _, p := range userPage.Profiles {
		fmt.Printf("User: %s\n", p)
	}
}

//
// Private
//

func getUserFixture(count int, startID int) []*logic.UserProfile {
	result := []*logic.UserProfile{}
	r := rand.New(rand.NewSource(1))
	from := time.Date(2000, time.January, 01, 0, 0, 0, 0, time.UTC)
	now := time.Now()

	for i := 0; i < count; i++ {
		created := fixture.GetRandomDateBetween(r, from, now)
		result = append(result, &logic.UserProfile{
			ID:       startID + i,
			Name:     fixture.GetRandomStr(r, fixture.PersonFirstNames) + " " + fixture.GetRandomStr(r, fixture.PersonLastNames),
			Created:  created,
			Accounts: getRandomOauthAccounts(r, created, now),
			Roles:    getRandomRoles(r),
		})
	}

	//log.Printf("Prepared users: %s\n", result)
	return result
}

var oauthAccountDistribution = []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}

var oauthProviders = []string{
	"VK", "Facebook", "Google", "Twitter",
}

func getRandomOauthAccount(r *rand.Rand, from time.Time, to time.Time) *logic.OauthAccount {
	tokenBytes := [16]byte{}
	r.Read(tokenBytes[:])
	return &logic.OauthAccount{
		Provider: fixture.GetRandomStr(r, oauthProviders),
		Token:    hex.EncodeToString(tokenBytes[:]),
		Created:  fixture.GetRandomDateBetween(r, from, to),
	}
}

func getRandomOauthAccounts(r *rand.Rand, from time.Time, to time.Time) []*logic.OauthAccount {
	result := []*logic.OauthAccount{}
	resultCount := fixture.GetRandomInt(r, oauthAccountDistribution)
	for i := 0; i < resultCount; i++ {
		result = append(result, getRandomOauthAccount(r, from, to))
	}

	return result
}

var roleCountDistribution = []int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}

func getRandomRoles(r *rand.Rand) []string {
	roleMap := map[string]bool{}
	roleCount := fixture.GetRandomInt(r, roleCountDistribution)
	for j := 0; j < roleCount; j++ {
		roleMap[fixture.GetRandomStr(r, logic.Roles[:])] = true
	}

	result := []string{}
	for k := range roleMap {
		result = append(result, k)
	}

	return result
}
