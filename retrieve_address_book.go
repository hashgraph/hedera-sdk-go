package hedera

func GetMainnetAddressBook() (*NodeAddressBook, error) {
	client := ClientForMainnet()

	result, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetPreviewnetAddressBook() (*NodeAddressBook, error) {
	client := ClientForPreviewnet()

	result, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func GetTestnetAddressBook() (*NodeAddressBook, error) {
	client := ClientForTestnet()

	result, err := NewAddressBookQuery().
		SetFileID(FileIDForAddressBook()).
		Execute(client)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
