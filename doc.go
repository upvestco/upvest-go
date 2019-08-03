/*
Package upvest provides the binding for Upvest REST APIs.
Where possible, the services available on the client groups the API into
logical chunks and correspond to the structure of the Upvest API documentation
at https://doc.upvest.co/reference.

First, create an Upvest client and depending on what action to take,
you either create tenancy or clientele client. All tenancy related operations
must be authenticated using the API Keys Authentication, whereas all actions
on a user's behalf need to be authenticated via OAuth. The API calls are
built along with those two authentication objects.

Usage:

	  // NewClient creates a new Upvest API client with the given base URL
	  // and HTTP client, allowing overriding of the HTTP client to use.
	  // This is useful if you're running in a Google AppEngine environment
	  // where the http.DefaultClient is not available.

	  c := NewClient("", nil)

	  // Tenant API - API Keys Authentication
	  // The Upvest API uses the notion of _tenants_, which represent customers that build their platform upon the Upvest API.
	  // The authentication via API keys and secret allows you to perform all tenant related operations.

	  tenancyClient = c.NewTenant(apiKey, apiSecret, apiPassphrase)

	  // create a user
	  user, err := tenancyClient.User.Create(username, randomString(12))
	  if err != nil {
		  t.Errorf("CREATE User returned error: %v", err)
	  }

	  // list users
	  users, err := tenancyClient.User.List()
	  if err != nil {
		  t.Errorf("List Users returned error: %v", err)
	  }

	  // retrieve 20 users
	  users, err := tenancyClient.User.ListN(20)
	  if err != nil {
		  t.Errorf("List Users returned error: %v", err)
	  }

	  // change password
	  user, err = tenancyClient.User.Update(username, params)

	  // Clinetele API - OAuth Authentication
	  // The authentication via OAuth allows you to perform operations on behalf of your user.
	  // For more information on the OAuth concept, please refer to our documentation at https://doc.upvest.co/docs/oauth2-authentication
	  // Next, create an `Clientele` object with these credentials
      // and your user authentication data in order to authenticate your API calls on behalf of a user:

	  clienteleClient = c.NewClientele(clientID, clientSecret, username, password)
	  wp := &WalletParams{
		  Password: staticUserPW,
		  AssetID:  ethWallet.Balances[0].AssetID,
	  }

	  // create the wallet
	  wallet, err := clienteleTestClient.Wallet.Create(wp)
	  if err != nil {
		  t.Errorf("CREATE Wallet returned error: %v", err)
	  }

	  // // retrieve the wallet
	  wallet1, err := clienteleTestClient.Wallet.Get(wallet.ID)
	  if err != nil {
		  t.Errorf("GET Wallet returned error: %v", err)
	  }
*/
package upvest
