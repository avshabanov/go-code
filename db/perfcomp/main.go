package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/avshabanov/go-code/db/perfcomp/logic"
	"github.com/avshabanov/go-code/fixture"
)

const sqliteDaoType = "sqlite"

var (
	dbPath      = flag.String("db-path", "", "Path to identity service database, e.g. /tmp/perfcomp-sqlite.db")
	dbType      = flag.String("db-type", sqliteDaoType, "Type of the database to test")
	initSize    = flag.Int("init-size", 10, "Size of initial data sample, applicable to initialization mode only")
	offsetToken = flag.String("ot", "", "Offset token, applicable to select mode only")
	mode        = flag.String("mode", "select", "App launch mode, possible values: select, reinit")
)

func main() {
	flag.Parse()

	if len(*dbPath) == 0 {
		log.Printf("db path is empty")
		flag.Usage()
		return
	}

	if *mode == "reinit" {
		deleteFileIfExists(*dbPath)
	}

	var dao logic.Dao
	var err error
	switch *dbType {
	case sqliteDaoType:
		dao, err = logic.NewSqliteDao(*dbPath)
	case "bolt":
		dao, err = logic.NewBoltDao(*dbPath)
	default:
		log.Fatalf("unable to create new DAO of type %s", *dbType)
	}

	if err != nil {
		log.Fatalf("cannot create dao: %v", err)
	}
	defer dao.Close()

	switch *mode {
	case "select":
		selectUsers(dao)
	case "reinit":
		reinit(dao)
	default:
		log.Fatalf("unknown mode=%s", *mode)
	}
}

//
// Private
//

func deleteFileIfExists(filePath string) {
	// file exists, try to delete it
	if err := os.Remove(filePath); err != nil {
		log.Printf("reinit db file=%s, remove operation failed: %v", filePath, err)
	}
}

func iterate(dao logic.Dao, limits []int, iterations int) {
	offsetToken := ""
	for n := 0; n < iterations; n++ {
		limit := n % len(limits)

		page, err := dao.QueryUsers(offsetToken, limit)
		if err != nil {
			log.Printf("unexpected error while querying users: %v", err)
		}

		offsetToken = page.OffsetToken
	}
}

func parallelSelects(dao logic.Dao) {

}

func reinit(dao logic.Dao) {
	// insert fixture
	if err := dao.Add(getUserFixture(*initSize, 1)); err != nil {
		log.Fatalf("cannot add user profile: %v", err)
	}
}

func selectUsers(dao logic.Dao) {
	userPage, err := dao.QueryUsers(*offsetToken, 8)
	if err != nil {
		log.Fatalf("cannot get user profiles: %v", err)
	}

	fmt.Println("users:")
	for _, p := range userPage.Profiles {
		fmt.Printf("# %s\n", p)
	}

	if len(userPage.OffsetToken) > 0 {
		fmt.Printf("# offsetToken: %s\n", userPage.OffsetToken)
	} else {
		fmt.Println("# <last page>")
	}

}

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
