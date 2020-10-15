package hd

// func TestKeyring(t *testing.T) {
// 	dir, cleanup := tests.NewTestCaseDir(t)
// 	mockIn := strings.NewReader("")
// 	t.Cleanup(cleanup)

// 	kr, err := keyring.New("ethermint", keyring.BackendTest, dir, mockIn, EthSecp256k1Option()...)
// 	require.NoError(t, err)

// 	// fail in retrieving key
// 	info, err := kr.Get("foo")
// 	require.Error(t, err)
// 	require.Nil(t, info)

// 	mockIn.Reset("password\npassword\n")
// 	info, mnemonic, err := kr.CreateMnemonic("foo", keyring.English, ethermint.BIP44HDPath, EthSecp256k1)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, mnemonic)
// 	require.Equal(t, "foo", info.GetName())
// 	require.Equal(t, "local", info.GetType().String())
// 	require.Equal(t, EthSecp256k1, info.GetAlgo())

// 	params := *hd.NewFundraiserParams(0, ethermint.Bip44CoinType, 0)
// 	hdPath := params.String()

// 	bz, err := DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, keyring.Secp256k1)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, bz)

// 	bz, err = DeriveSecp256k1(mnemonic, keyring.DefaultBIP39Passphrase, hdPath)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, bz)

// 	bz, err = DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, keyring.SigningAlgo(""))
// 	require.Error(t, err)
// 	require.Empty(t, bz)

// 	bz, err = DeriveSecp256k1(mnemonic, keyring.DefaultBIP39Passphrase, "/wrong/hdPath")
// 	require.Error(t, err)
// 	require.Empty(t, bz)

// 	bz, err = DeriveKey(mnemonic, keyring.DefaultBIP39Passphrase, hdPath, EthSecp256k1)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, bz)

// 	privkey := &ethsecp256k1.PrivKey{Key: bz}
// 	addr := common.BytesToAddress(privkey.PubKey().Address().Bytes())

// 	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
// 	require.NoError(t, err)

// 	path := hdwallet.MustParseDerivationPath(hdPath)

// 	account, err := wallet.Derive(path, false)
// 	require.NoError(t, err)
// 	require.Equal(t, addr.String(), account.Address.String())
// }
