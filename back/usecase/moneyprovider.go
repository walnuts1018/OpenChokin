package usecase

type MoneyProviderSummary struct {
	ID      string
	Name    string
	Balance float64
}
type MoneyProvidersSummaryResponse struct {
	Providers []MoneyProviderSummary
}

func (u Usecase) GetMoneyProvidersSummary(userID string) (MoneyProvidersSummaryResponse, error) {
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return MoneyProvidersSummaryResponse{}, err
	}

	var providersSummary []MoneyProviderSummary
	for _, provider := range moneyProviders {
		balance, err := u.db.GetMoneyPoolBalance(provider.ID, false)
		if err != nil {
			// handle the error according to your error policy
			// for example, you could log it and continue with the next provider
			continue
		}
		providersSummary = append(providersSummary, MoneyProviderSummary{
			ID:      provider.ID,
			Name:    provider.Name,
			Balance: balance,
		})
	}

	return MoneyProvidersSummaryResponse{Providers: providersSummary}, nil
}
