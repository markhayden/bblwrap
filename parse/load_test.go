package parse

// import (
// 	"appengine/aetest"
// 	"fmt"
// 	"github.com/bmizerany/assert"
// 	"github.com/lytics/gaeservices/lib/models"
// 	"testing"
// )

// func prepareLyticsClient() models.Account {
// 	return models.Account{
// 		Name:   "Lytics",
// 		Apikey: "EYBDi6MectBTMbm8ewX7rsIb",
// 		Id:     "4f7b525bdba058fc096757a6",
// 		Aid:    12,
// 	}
// }

// func prepareFakeClient() models.Account {
// 	return models.Account{
// 		Name:   "Broken",
// 		Apikey: "humptydumpty",
// 		Id:     "2352525",
// 		Aid:    8675309,
// 	}
// }

// func TestLoadAudiences(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// audiences returns success
// 	client1 := prepareLyticsClient()
// 	audience1 := loadAudiences(client1, ctx)
// 	assert.Equal(t, len(audience1) > 0, true)
// 	fmt.Println("+++ Get audience list pass verified")

// 	// audiences returns fail
// 	client2 := prepareFakeClient()
// 	audience2 := loadAudiences(client2, ctx)
// 	assert.Equal(t, len(audience2) == 0, true)
// 	fmt.Println("+++ Get audience list pass verified")
// }

// func TestLoadAction(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	client := prepareLyticsClient()

// 	// action returns success
// 	action1 := loadAction("b7ebf6780c8c44e29d8fe1c570f9110b", client, ctx)
// 	assert.Equal(t, action1.Name, "any_email_imported")
// 	fmt.Println("+++ Get action pass verified")

// 	// action returns fail
// 	action2 := loadAction("imnotarealaction", client, ctx)
// 	assert.Equal(t, action2.Id, "")
// 	fmt.Println("+++ Get action fail verified")
// }

// func TestLoadProvider(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	client := prepareLyticsClient()

// 	// privder returns success
// 	provider1 := loadProvider("d88940e9e727440ba12b4a0ba9582cc8", client, ctx)
// 	assert.Equal(t, provider1.Slug, "campaignmonitor")
// 	fmt.Println("+++ Get provider pass verified")

// 	// provider returns fail
// 	provider2 := loadProvider("imnotarealprovider", client, ctx)
// 	assert.Equal(t, provider2.Slug, "")
// 	fmt.Println("+++ Get provider fail verified")
// }

// func TestLoadSuggestion(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// suggestion returns success
// 	client1 := prepareLyticsClient()
// 	suggestion1 := loadSuggestion("4bf1ac6d03bb4952b40def409dfc120e", client1, ctx)
// 	assert.Equal(t, suggestion1.Id, "4bf1ac6d03bb4952b40def409dfc120e")
// 	fmt.Println("+++ Get suggestion pass verified")

// 	// suggestion returns fail
// 	client2 := prepareFakeClient()
// 	suggestion2 := loadSuggestion("imnotarealsuggestion", client2, ctx)
// 	assert.Equal(t, suggestion2.Id, "")
// 	fmt.Println("+++ Get suggestion fail verified")
// }

// func TestLoadSuggestions(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// suggestions returns success
// 	client1 := prepareLyticsClient()
// 	suggestion1 := loadSuggestions(client1, ctx)
// 	assert.Equal(t, len(suggestion1) > 1, true)
// 	fmt.Println("+++ Get suggestion pass verified")

// 	// suggestions returns fail
// 	client2 := prepareFakeClient()
// 	suggestion2 := loadSuggestions(client2, ctx)
// 	assert.Equal(t, len(suggestion2) > 1, false)
// 	fmt.Println("+++ Get suggestion pass verified")
// }

// func TestLoadAccounts(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// accounts returns success
// 	accounts1 := loadAccounts(ctx, LYTICS_API, ADMIN_KEY)
// 	assert.Equal(t, len(accounts1) > 0, true)
// 	fmt.Println("+++ Get account list pass verified")

// 	// accounts returns fail
// 	accounts2 := loadAccounts(ctx, "fake", "key")
// 	assert.Equal(t, len(accounts2) == 0, true)
// 	fmt.Println("+++ Get account list fail verified")
// }

// func TestLoadAuths(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// suggestions returns success
// 	client1 := prepareLyticsClient()

// 	// auths returns success
// 	auth1 := loadAuths(client1, ctx)
// 	assert.Equal(t, len(auth1) > 1, true)
// 	fmt.Println("+++ Get account auths verified")
// }

// func TestLoadFromContentApi(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	client := prepareLyticsClient()

// 	// content returns success
// 	content1 := loadFromContentApi("f3cae01b40914aecbb7d49767b0aabff", client, ctx)
// 	assert.Equal(t, len(content1) > 0, true)
// 	fmt.Println("+++ Get content fail verified")

// 	// content returns fail
// 	content2 := loadFromContentApi("imnotrealcontent", client, ctx)
// 	assert.Equal(t, len(content2) == 0, true)
// 	fmt.Println("+++ Get content fail verified")
// }

// func TestLoadLocalFile(t *testing.T) {
// 	// load file success
// 	file1, _ := loadLocalFile("templates/test_mobile.html")
// 	assert.Equal(t, file1 != "", true)
// 	fmt.Println("+++ Get local file pass verified")

// 	// load file fail
// 	_, err := loadLocalFile("templates/notarealfile.html")
// 	assert.Equal(t, err != nil, true)
// 	fmt.Println("+++ Get local file fail verified")
// }

// func TestSean(t *testing.T) {
// 	ctx, _ := aetest.NewContext(nil)
// 	defer ctx.Close()

// 	// auths returns success
// 	seanTest(ctx)
// }
