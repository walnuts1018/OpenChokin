package usecase

import "time"

type MoneySumResponse struct {
	MoneyProviderSum       float64
	ActualMoneyPoolSum     float64
	ForecastedMoneyPoolSum float64
}

// GetMoneyInformation calculates and returns the total balance of MoneyPools and MoneyProviders for a user.
func (u Usecase) GetMoneyInformation(userID string) (MoneySumResponse, error) {
	var response MoneySumResponse
	var err error

	// Retrieve all MoneyPools associated with the user.
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return response, err
	}

	// Calculate the sum of MoneyPool balances.
	for _, pool := range moneyPools {
		balance, err := u.db.GetMoneyPoolBalance(pool.ID, false)
		if err != nil {
			return response, err
		}
		response.ActualMoneyPoolSum += balance
		balance, err = u.db.GetMoneyPoolBalance(pool.ID, true)
		if err != nil {
			return response, err
		}
		response.ForecastedMoneyPoolSum += balance
	}

	// Retrieve all MoneyProviders associated with the user.
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return response, err
	}

	// Calculate the sum of MoneyProvider balances.
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}

	return response, nil
}

// GetMoneyInformationOfDate computes the total balance of MoneyProviders and MoneyPools belonging to a user up to a given date.
func (u Usecase) GetMoneyInformationOfDate(userID string, date time.Time) (MoneySumResponse, error) {
	var response MoneySumResponse

	// Retrieve all MoneyProviders for the user and calculate the sum of their balances.
	moneyProviders, err := u.db.GetMoneyProvidersByUserID(userID)
	if err != nil {
		return response, err
	}
	for _, provider := range moneyProviders {
		response.MoneyProviderSum += provider.Balance
	}

	// Retrieve all MoneyPools for the user and calculate the sum of their balances up to the specified date.
	moneyPools, err := u.db.GetMoneyPoolsByUserID(userID)
	if err != nil {
		return response, err
	}
	for _, pool := range moneyPools {
		balance, err := u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, false)
		if err != nil {
			return response, err
		}
		response.ActualMoneyPoolSum += balance
		balance, err = u.db.GetMoneyPoolBalanceOfDate(pool.ID, date, true)
		if err != nil {
			return response, err
		}
		response.ForecastedMoneyPoolSum += balance
	}

	return response, nil
}
