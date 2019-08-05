[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/rpip/upvest-go) [![Build Status](https://travis-ci.org/rpip/upvest-go.svg?branch=master)](https://travis-ci.org/rpip/upvest-go)

# Go library for the Upvest API.

In order to retrieve your API credentials for using this Go client, you'll need to [sign up with Upvest](https://login.upvest.co/sign-up).

Where possible, the services available on the client groups the API into logical chunks and correspond to the structure of the [Upvest API documentation](https://doc.upvest.co).

First, create an Upvest client and depending on what action to take, you either create tenancy or clientele client. All tenancy related operations must be authenticated using the API Keys Authentication, whereas all actions on a user's behalf need to be authenticated via OAuth. The API calls are built along with those two authentication objects.

``` go
// NewClient creates a new Upvest API client with the given base URL
// and HTTP client, allowing overriding of the HTTP client to use.
// This is useful if you're running in a Google AppEngine environment
// where the http.DefaultClient is not available.

c := NewClient("", nil)

// Configure logging using the Loggingenabled config key
c.Loggingenabled = true
```

### Tenancy API - API Keys Authentication
The Upvest API uses the notion of _tenants_, which represent customers that build their platform upon the Upvest API. The end-users of the tenant (i.e. your customers), are referred to as _clients_. A tenant is able to manage their users directly (CRUD operations for the user instance) and is also able to initiate actions on the user's behalf (create wallets, send transactions).

The authentication via API keys and secret allows you to perform all tenant related operations.
Please create an API key pair within the [Upvest account management](https://login.upvest.co/).

The default `BASE_URL` for both authentication objects is `https://api.playground.upvest.co`, but feel free to adjust it, once you retrieve your live keys. Next, create an `Tenancy` object in order to authenticate your API calls:

```go
tenancy = c.NewTenant(apiKey, apiSecret, apiPassphrase)

// create a user
user, err := tenancy.User.Create(username, randomString(12))
if err != nil {
    t.Errorf("CREATE User returned error: %v", err)
}

// list users
users, err := tenancy.User.List()
if err != nil {
    t.Errorf("List Users returned error: %v", err)
}

// retrieve 20 users
users, err := tenancy.User.ListN(20)
if err != nil {
    t.Errorf("List Users returned error: %v", err)
}

// change password
user, err = tenancy.User.Update(username, params)
```

### Clientele API - OAuth Authentication
The authentication via OAuth allows you to perform operations on behalf of your user.
For more information on the OAuth concept, please refer to our [documentation](https://doc.upvest.co/docs/oauth2-authentication).
Again, please retrieve your client credentials from the [Upvest account management](https://login.upvest.co/).

Next, create an `Clientele` object with these credentials and your user authentication data in order to authenticate your API calls on behalf of a user:

```go

clientele = c.NewClientele(clientID, clientSecret, username, password)
wp := &WalletParams{
    Password: staticUserPW,
    AssetID:  ethWallet.Balances[0].AssetID,
}

// create the wallet
wallet, err := clientele.Wallet.Create(wp)
if err != nil {
    t.Errorf("CREATE Wallet returned error: %v", err)
}

// // retrieve the wallet
wallet1, err := clientele.Wallet.Get(wallet.ID)
if err != nil {
    t.Errorf("GET Wallet returned error: %v", err)
}

```

## Development

1. Code must be `go fmt` compliant: `make fmt`
2. All types, structs and funcs should be documented.
3. Ensure that `make test` succeeds.
4. Set up config settings via environment variables:

    ```shell
    # Set your tenancy API key information here.
    export API_KEY=xxxx
    export API_SECRET=xxxx
    export API_PASSPHRASE=xxxx

    # Set your OAuth2 client information here.
    export OAUTH2_CLIENT_ID=xxxx
    export OAUTH2_CLIENT_SECRET=xxxx
    ```


## Test

Run all tests:

    make test

Run a single test:

    DEBUG=1 go test -run TestChangePassword

## More

For a comprehensive reference, check out the [Upvest documentation](https://doc.upvest.co).

For details on all the functionality in this library, see the [GoDoc documentation](https://godoc.org/github.com/rpip/upvest-go).
