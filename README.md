[![CircleCI](https://circleci.com/gh/upvestco/upvest-go.svg?style=svg)](https://circleci.com/gh/upvestco/upvest-go) [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/rpip/upvest-go)

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
tenant = c.NewTenant(apiKey, apiSecret, apiPassphrase)

// create a user
user, err := tenant.User.Create(username, randomString(12))
if err != nil {
    t.Errorf("CREATE User returned error: %v", err)
}

// list users
users, err := tenant.User.List()
if err != nil {
    t.Errorf("List Users returned error: %v", err)
}

// retrieve 20 users
users, err := tenant.User.ListN(20)
if err != nil {
    t.Errorf("List Users returned error: %v", err)
}

// change password
user, err = tenant.User.ChangePassword(username, params)
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
    AssetID: "8fc19cd0-8f50-4626-becb-c9e284d2315b",
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

## Usage

### Tenancy
#### User management

##### Create a user
```go
user, err := tenant.User.Create('username','password')
```

##### Retrieve a user

```go
user, err := tenant.User.Get('username')
```

##### List all users under tenancy

```go
users, err := tenant.User.List()

for _, user := range users.Values {
  //do something with user
}
```

##### List a specific number of users under tenancy

```go
users, err := tenant.User.listN(10)
```

##### Change password of a user

```go
params := &upvest.ChangePasswordParams{
    OldPassword: "current password",
    NewPassword: "new pasword",
}
user, err := tenant.User.ChangePassword(username, params)
```

##### Delete a user

```go
tenant.User.Delete('username')
```

### Clientele

#### Assets

##### List available assets

```go
assets, err := clientele.Asset.List()
for _, asset := range assets.Values {
  //do something with asset
}
```


#### Wallets

##### Create a wallet for a user

```go
wp := &upvest.WalletParams{
    Password: "current user password",
    AssetID:  "asset ID",
    // Type:     "encrypted",
    // Index:    0,
}

// create the wallet
wallet, err := clientele.Wallet.Create(wp)
```

##### Retrieve specific wallet for a user

```go
wallet1, err := clientele.Wallet.Get(walletID)
```

##### List all wallets for a user

```go
wallets, err := clientele.Wallet.List()
for _, wallet := range wallets.Values {
  //do something with wallet
}
```

##### List a specific number of wallets

```go
wallets, err := clientele.Wallet.ListN(40)
```

#### Transactions

##### Create transaction

```go
tp := &upvest.TransactionParams{
    Password:  "current user password",
    AssetID:   "asset ID",
    Quantity:  "quantity, e.g. 10000000000000000",
    Fee:       "fee, e.g. 41180000000000",
    Recipient: "transaction address, e.g. 0xf9b44Ba370CAfc6a7AF77D0BDB0d50106823D91b",
}

// create the transaction
txn, err := clientele.Transaction.Create("wallet ID", tp)
```

#### Retrieve specific transaction

```go
wallet1, err := clientele.Wallet.Get("wallet ID")
```

##### List all transactions of a wallet for a user

```go
transactions, err := clientele.Transaction.List("wallet ID")
for _, txn := range transactions.Values {
  //do something with transaction
}
```

##### List a specific number of transactions of a wallet for a user

```go
transactions, err := clientele.Transaction.ListN("wallet ID", 8)
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

For details on all the functionality in this library, see the [GoDoc documentation](https://godoc.org/github.com/upvestco/upvest-go).
